package payment

import (
	"bikeshop/v3/cmd/server/config"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type CreatePaymentRequest struct {
	Key     string `json:"idempotency_key"`
	Amount  Amount `json:"amount_money"`
	Source  string `json:"source_id"`
	OrderId string `json:"order_id"`
	Details struct {
		Amount Amount `json:"buyer_supplied_money"`
	} `json:"cash_details"`
}

type Amount struct {
	Value    float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func CreatePayment(orderId string, cost float64) (string, error) {
	client := &http.Client{}
	req, err := config.Config.NewRequest("POST", config.PAYMENT_URL)
	if err != nil {
		return "", err
	}
	body, err := json.Marshal(CreatePaymentRequest{
		Key:     uuid.New().String(),
		Amount:  Amount{Value: cost, Currency: "CAD"},
		Source:  "CASH",
		OrderId: orderId,
		Details: struct {
			Amount Amount "json:\"buyer_supplied_money\""
		}{Amount: Amount{Value: cost, Currency: "CAD"}},
	})
	if err != nil {
		return "", err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	resByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	json.Unmarshal(resByte, &result)
	return result["payment"].(map[string]interface{})["id"].(string), nil
}
