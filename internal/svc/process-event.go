package svc

import (
	"context"
	"encoding/json"
	db "first-excercise/internal/db"
	utils "first-excercise/internal/util"
	"fmt"
	"log"
	"path"
	"time"
)

type Result struct {
	Success bool
	Message string
}

type EventResponse struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

type OrderResponse struct {
	OrderId string  `json:"orderId"`
	Amount  float64 `json:"amount"`
	Item    string  `json:"item"`
	CreatedAt string  `json:"createdAt,omitempty"`
}

func ProcessEvent(ctx context.Context, event json.RawMessage) (*Result, error) {
	log.Printf("Received event: %s", string(event))

	eventData, err := utils.UnmarshalEvent(event)
	if err != nil {
		log.Printf("Error processing event: %v", err)
		return nil, fmt.Errorf("failed to process event: %v", err)
	}

	if path.Clean(eventData.Path) == "/dummy" {

		if eventData.HttpMethod == "GET" {
			log.Printf("Received GET request for /dummy")

			// Fetching based on orderId query parameter if present, otherwise scan all orders
			if eventData.QueryStringParameters != nil {
				log.Printf("Query parameters: %v", eventData.QueryStringParameters)
				orderId := eventData.QueryStringParameters["orderId"]
				if orderId != "" {
					order, err := db.GetOrder(ctx, orderId)
					if err != nil {
						log.Printf("Error retrieving order: %v", err)
						return &Result{
							Success: false,
							Message: fmt.Sprintf("failed to retrieve order: %v", err),
						}, nil
					}

					orderJSON, _ := json.Marshal(order)
					return &Result{
						Success: true,
						Message: string(orderJSON),
					}, nil
				}
			}
			

			orders, err := db.ScanOrders(ctx)
			if err != nil {
				log.Printf("Error scanning DynamoDB: %v", err)
				return &Result{
					Success: false,
					Message: fmt.Sprintf("failed to retrieve orders: %v", err),
				}, nil
			}

			ordersJSON, _ := json.Marshal(orders)
			return &Result{
				Success: true,
				Message: string(ordersJSON),
			}, nil
		}

		if eventData.HttpMethod == "POST" {
			order, err := processOrder(eventData.Body)
			if err != nil {
				log.Printf("Error processing order: %v", err)
				return &Result{
					Success: false,
					Message: fmt.Sprintf("failed to process order: %v", err),
				}, nil
			}

			log.Printf("Received order: %+v", order)

			if err := db.PutOrder(ctx, order); err != nil {
				log.Printf("Error saving order to DynamoDB: %v", err)
				return &Result{
					Success: false,
					Message: fmt.Sprintf("failed to save order: %v", err),
				}, nil
			}

			orderResponse := OrderResponse{
				OrderId: order.OrderId,
				Amount:  order.Amount,
				Item:    order.Item,
				CreatedAt: time.Now().Format("2006-01-02T15:04:05Z"),
			}
			orderJSON, _ := json.Marshal(orderResponse)
			return &Result{
				Success: true,
				Message: string(orderJSON),
			}, nil
		}
	}

	return &Result{
		Success: false,
		Message: fmt.Sprintf("Unhandled request for %s %s", eventData.HttpMethod, eventData.Path),
	}, nil
}

func processOrder(orderData string) (*utils.Order, error) {
	order, err := utils.UnmarshalOrder(orderData)
	if err != nil {
		log.Printf("Error processing order: %v", err)
		return nil, fmt.Errorf("failed to process order: %v", err)
	}
	return order, nil
}

func BuildResponse(result Result) EventResponse {
	log.Printf("Building response for result: %+v", result)
	code := 200
	if !result.Success {
		code = 400
	}

	response := EventResponse{
		StatusCode: code,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       result.Message,
	}

	return response
}

func BuildErrorResponse(cause string) EventResponse {
	log.Printf("Building error response for cause: %s", cause)

	response := EventResponse{
		StatusCode: 500,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       cause,
	}

	return response
}
