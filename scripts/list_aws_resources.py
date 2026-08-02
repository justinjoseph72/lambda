#!/usr/bin/env python3
"""
List AWS resources running in an account.

Uses the Resource Groups Tagging API to enumerate most resource types in one
or more regions, plus a handful of explicit calls for services that API
doesn't cover well (S3 buckets, IAM users/roles, CloudFormation stacks are
global or need dedicated list calls).

Usage:
    python3 list_aws_resources.py                 # default region (eu-west-2)
    python3 list_aws_resources.py --region us-east-1
    python3 list_aws_resources.py --all-regions
    python3 list_aws_resources.py --profile myprofile

Requires AWS credentials configured (env vars, ~/.aws/credentials, or --profile)
with at least read-only permissions (e.g. the ReadOnlyAccess managed policy).
"""

import argparse
import sys
from collections import defaultdict

import boto3
from botocore.exceptions import ClientError, NoCredentialsError

DEFAULT_REGION = "eu-west-2"


def get_session(profile: str | None) -> boto3.Session:
    return boto3.Session(profile_name=profile) if profile else boto3.Session()


def list_regions(session: boto3.Session) -> list[str]:
    ec2 = session.client("ec2", region_name=DEFAULT_REGION)
    resp = ec2.describe_regions(AllRegions=False)
    return sorted(r["RegionName"] for r in resp["Regions"])


def list_tagged_resources(session: boto3.Session, region: str) -> dict[str, list[str]]:
    """Enumerate resources in a region via the Resource Groups Tagging API."""
    by_type: dict[str, list[str]] = defaultdict(list)
    client = session.client("resourcegroupstaggingapi", region_name=region)
    paginator = client.get_paginator("get_resources")
    try:
        for page in paginator.paginate():
            for resource in page.get("ResourceTagMappingList", []):
                arn = resource["ResourceARN"]
                # arn:partition:service:region:account:resource-type/id
                parts = arn.split(":", 5)
                service = parts[2] if len(parts) > 2 else "unknown"
                by_type[service].append(arn)
    except ClientError as e:
        print(f"  [warn] resourcegroupstaggingapi in {region}: {e}", file=sys.stderr)
    return by_type


def list_global_resources(session: boto3.Session) -> dict[str, list[str]]:
    """Resources that are global or not well covered by the tagging API."""
    results: dict[str, list[str]] = defaultdict(list)

    try:
        s3 = session.client("s3")
        for bucket in s3.list_buckets().get("Buckets", []):
            results["s3-buckets"].append(bucket["Name"])
    except ClientError as e:
        print(f"  [warn] s3 list_buckets: {e}", file=sys.stderr)

    try:
        iam = session.client("iam")
        paginator = iam.get_paginator("list_users")
        for page in paginator.paginate():
            for user in page.get("Users", []):
                results["iam-users"].append(user["UserName"])
    except ClientError as e:
        print(f"  [warn] iam list_users: {e}", file=sys.stderr)

    try:
        iam = session.client("iam")
        paginator = iam.get_paginator("list_roles")
        for page in paginator.paginate():
            for role in page.get("Roles", []):
                if not role["RoleName"].startswith("AWSServiceRoleFor"):
                    results["iam-roles"].append(role["RoleName"])
    except ClientError as e:
        print(f"  [warn] iam list_roles: {e}", file=sys.stderr)

    return results


def list_cloudformation_stacks(session: boto3.Session, region: str) -> list[str]:
    cf = session.client("cloudformation", region_name=region)
    stacks = []
    try:
        paginator = cf.get_paginator("list_stacks")
        for page in paginator.paginate(
            StackStatusFilter=[
                "CREATE_COMPLETE", "UPDATE_COMPLETE", "UPDATE_ROLLBACK_COMPLETE",
                "ROLLBACK_COMPLETE", "IMPORT_COMPLETE",
            ]
        ):
            for stack in page.get("StackSummaries", []):
                stacks.append(stack["StackName"])
    except ClientError as e:
        print(f"  [warn] cloudformation in {region}: {e}", file=sys.stderr)
    return stacks


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--region", default=DEFAULT_REGION, help=f"AWS region to scan (default: {DEFAULT_REGION})")
    parser.add_argument("--all-regions", action="store_true", help="Scan every enabled region instead of just --region")
    parser.add_argument("--profile", default=None, help="AWS CLI profile to use")
    args = parser.parse_args()

    try:
        session = get_session(args.profile)
        account_id = session.client("sts").get_caller_identity()["Account"]
    except NoCredentialsError:
        print("No AWS credentials found. Configure them via env vars, ~/.aws/credentials, or --profile.", file=sys.stderr)
        return 1
    except ClientError as e:
        print(f"Failed to get caller identity: {e}", file=sys.stderr)
        return 1

    print(f"AWS Account: {account_id}\n")

    regions = list_regions(session) if args.all_regions else [args.region]

    print("=== Global resources ===")
    global_resources = list_global_resources(session)
    if global_resources:
        for rtype, items in sorted(global_resources.items()):
            print(f"\n{rtype} ({len(items)}):")
            for item in items:
                print(f"  - {item}")
    else:
        print("  (none found or no permissions)")

    for region in regions:
        print(f"\n=== Region: {region} ===")

        stacks = list_cloudformation_stacks(session, region)
        if stacks:
            print(f"\ncloudformation-stacks ({len(stacks)}):")
            for s in stacks:
                print(f"  - {s}")

        by_type = list_tagged_resources(session, region)
        if not by_type:
            print("  (no resources found via resourcegroupstaggingapi, or no permissions)")
            continue

        for rtype, arns in sorted(by_type.items()):
            print(f"\n{rtype} ({len(arns)}):")
            for arn in arns:
                print(f"  - {arn}")

    return 0


if __name__ == "__main__":
    sys.exit(main())
