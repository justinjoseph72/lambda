package svc

import (
	"encoding/json"
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
	// IsBase64Encoded bool              `json:"isBase64Encoded"`
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

func ProcessEvent(event json.RawMessage) (*Result, error) {
	log.Printf("Received event: %s", string(event))

	eventData, err := utils.UnmarshalEvent(event)
	if err != nil {
		log.Printf("Error processing event: %v", err)
		return nil, fmt.Errorf("failed to process event: %v", err)
	}

	if path.Clean(eventData.Path) == "/dummy" {

		if eventData.HttpMethod == "GET" {
			log.Printf("Received GET request for /dummy")

			orderResponse := OrderResponse{
				OrderId: "12345",
				Amount:  100.00,
				Item:    "Sample Item",
			}
			orderJSON, _ := json.Marshal(orderResponse)
			return &Result{
				Success: true,
				Message: string(orderJSON),
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
	if result.Success == false {
		code = 400
	}

	response := EventResponse{
		// IsBase64Encoded: false,
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
