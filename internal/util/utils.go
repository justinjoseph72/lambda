package util

import (
	"bytes"
	"encoding/json"

	"fmt"
	"log"
)

type Event struct {
	HttpMethod            string                 `json:"httpMethod"`
	Path                  string                 `json:"path"`
	Body                  string                 `json:"body"`
	Headers               map[string]string      `json:"headers,omitempty"`
	RequestContext        map[string]interface{} `json:"requestContext,omitempty"`
	QueryStringParameters map[string]string      `json:"queryStringParameters,omitempty"`
}

type Order struct {
	OrderId string  `json:"orderId"`
	Amount  float64 `json:"amount"`
	Item    string  `json:"item"`
}

func UnmarshalEvent(event json.RawMessage) (*Event, error) {
	log.Printf("Received event: %s", string(event))

	decoder := json.NewDecoder(bytes.NewReader(event))

	var eventData Event
	err := decoder.Decode(&eventData)
	if err != nil {
		log.Printf("Error parsing event: %v", err)
		return nil, fmt.Errorf("failed to parse event: %v", err)
	}
	return &eventData, nil
}

func UnmarshalOrder(orderData string) (*Order, error) {
	var order Order
	err := json.Unmarshal([]byte(orderData), &order)
	if err != nil {
		log.Printf("Error parsing order: %v", err)
		return nil, fmt.Errorf("failed to parse order: %v", err)
	}
	return &order, nil
}
