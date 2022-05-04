package webhook

import (
	"bikeshop/v3/cmd/server/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type ListSubscriptionResp struct {
	Subscriptions []Subscription
}
type CreateSubscriptionReq struct {
	Key string       `json:"idempotency_key"`
	Sub Subscription `json:"subscription"`
}

type Subscription struct {
	Id              string   `json:"id"`
	ApplicationID   string   `json:"application_id"`
	Name            string   `json:"name"`
	Enabled         bool     `json:"enabled"`
	EventTypes      []string `json:"event_types"`
	NotificationURL string   `json:"notification_url"`
	APIVersion      string   `json:"api_version"`
	SignatureKey    string   `json:"signature_key"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

type subscriptionConfig struct {
	name            string
	eventTypes      []string
	notificationURL string
}

var Subscriptions = map[string]subscriptionConfig{
	"bikestoreEmail": {
		name:            "bikestoreEmail",
		eventTypes:      []string{"inventory.count.updated"},
		notificationURL: "https://webhook.site/c443bc91-7d2a-42cf-a179-ad5f1013fea5",
	},
}

func ConfigureWebhooks() error {
	// first delete the existing webhooks which requires getting the webhooks and then deleting them. as long as they are in the subscrions list.
	ids, err := getSubscriptions()
	if err != nil {
		return err
	}
	for _, id := range ids {
		err = deleteSubscription(id)
		if err != nil {
			return err
		}
	}

	// now make the subscriptions we need with the events we want. For now we use the same the same notifaction url
	return createSubscriptions()
}

func getSubscriptions() ([]string, error) {
	client := &http.Client{}
	req, err := config.Config.NewRequest("GET", config.WEBHOOK_URL+"/subscriptions")
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to get subscriptions %v", err)
	}

	resByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ListSubscriptionResp
	json.Unmarshal(resByte, &result)
	ids := make([]string, 0)
	for _, sub := range result.Subscriptions {
		if _, found := Subscriptions[sub.Name]; found {
			ids = append(ids, sub.Id)
		}
	}
	return ids, nil
}

func deleteSubscription(id string) error {
	client := &http.Client{}
	req, err := config.Config.NewRequest("DELETE", fmt.Sprintf("%s/subscriptions/%s", config.WEBHOOK_URL, id))
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to delete subscription %s with err %v", id, err)
	}
	return nil
}

func createSubscriptions() error {
	client := &http.Client{}
	req, err := config.Config.NewRequest("POST", config.WEBHOOK_URL+"/subscriptions")
	if err != nil {
		return err
	}

	for _, conf := range Subscriptions {
		body, err := json.Marshal(CreateSubscriptionReq{
			Key: uuid.New().String(),
			Sub: Subscription{
				Name:            conf.name,
				EventTypes:      conf.eventTypes,
				NotificationURL: conf.notificationURL,
				APIVersion:      config.Config.GetAPIVersion(),
				Enabled:         true,
			},
		})
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewReader(body))
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unable to create subscription %swith err %v", conf.name, err)
		}
	}
	return nil
}
