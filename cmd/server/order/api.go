package order

import (
	"bikeshop/v3/cmd/server/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	Order Order  `json:"order"`
	Key   string `json:"idempotency_key"`
}

type PayOrderRequest struct {
	Key      string   `json:"idempotency_key"`
	Payments []string `json:"payment_ids"`
}

type Order struct {
	Location string `json:"location_id"`
	LineItem []struct {
		Quantity string `json:"quantity"`
		ItemId   string `json:"catalog_object_id"`
		ItemType string `json:"item_type"`
	} `json:"line_items"`
}

func CreateOrder(quantity, itemId string) (string, float64, error) {
	client := &http.Client{}
	req, err := config.Config.NewRequest("POST", "https://connect.squareupsandbox.com/v2/orders")
	if err != nil {
		return "", 0, err
	}
	body, err := json.Marshal(CreateOrderRequest{
		Key: uuid.New().String(),
		Order: Order{
			Location: config.Config.GetLocation(),
			LineItem: []struct {
				Quantity string "json:\"quantity\""
				ItemId   string "json:\"catalog_object_id\""
				ItemType string "json:\"item_type\""
			}{
				{
					Quantity: quantity,
					ItemId:   itemId,
					ItemType: "ITEM",
				},
			},
		},
	})
	if err != nil {
		return "nil", 0, err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}

	resByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf(string(resByte))
	}

	var result map[string]interface{}
	json.Unmarshal(resByte, &result)

	return result["order"].(map[string]interface{})["id"].(string), result["order"].(map[string]interface{})["total_money"].(map[string]interface{})["amount"].(float64), nil
}

func PayOrder(orderID, paymentId string) error {
	client := &http.Client{}
	req, err := config.Config.NewRequest("POST", fmt.Sprintf("%s/%s/pay", config.ORDER_URL, orderID))
	if err != nil {
		return err
	}

	body, err := json.Marshal(PayOrderRequest{
		Key:      uuid.New().String(),
		Payments: []string{paymentId},
	})
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(body))
	_, err = client.Do(req)
	return err
}
