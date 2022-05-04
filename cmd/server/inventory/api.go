package inventory

import (
	"bikeshop/v3/cmd/server/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CreateInventoryRequest struct {
	Key             string            `json:"idempotency_key"`
	Changes         []InventoryChange `json:"changes"`
	IgnoreUnchanged bool              `json:"ignore_unchanged_counts"`
}

type GetInventoryRequest struct {
	CatalogueIds []string `json:"catalog_object_ids"`
}

type InventoryChange struct {
	Count struct {
		Id            string `json:"catalog_object_id"`
		Location      string `json:"location_id"`
		Quantity      string `json:"quantity"`
		State         string `json:"state"`
		OccurenceTime string `json:"occurred_at"`
	} `json:"physical_count"`
	ChangeType string `json:"type"`
}

func GetInventory(catalogue []string) ([]byte, error) {
	req, err := config.Config.NewRequest("POST", config.INVENTORY_URL+"/counts/batch-retrieve")
	if err != nil {
		return nil, err
	}
	client := http.Client{}

	body, _ := json.Marshal(GetInventoryRequest{
		CatalogueIds: catalogue,
	})

	req.Body = io.NopCloser(bytes.NewReader(body))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func GetInventoryCounts(ids []string) (map[string]string, error) {
	inventory, err := GetInventory(ids)
	if err != nil {
		return nil, err
	}

	counts := make(map[string]string)
	var result map[string]interface{}
	json.Unmarshal(inventory, &result)
	for _, count := range result["counts"].([]interface{}) {
		counts[count.(map[string]interface{})["catalog_object_id"].(string)] = count.(map[string]interface{})["quantity"].(string)
	}
	return counts, nil
}

func GenerateInventory(catalogIDs map[string]string) error {
	changes := getInventoryChanges(catalogIDs)
	return pushInventoryChanges(changes)
}

func getInventoryChanges(catalogIDs map[string]string) []InventoryChange {
	changes := make([]InventoryChange, len(catalogIDs))
	i := 0
	for _, val := range catalogIDs {
		changes[i].ChangeType = "PHYSICAL_COUNT"
		changes[i].Count.Id = val
		changes[i].Count.Location = config.Config.GetLocation()
		changes[i].Count.Quantity = fmt.Sprintf("%d", rand.Intn(20-10)+10)
		changes[i].Count.State = "IN_STOCK"
		changes[i].Count.OccurenceTime = time.Now().Format(time.RFC3339)
		i++
	}
	return changes
}

func pushInventoryChanges(inventory []InventoryChange) error {
	client := &http.Client{}
	req, err := config.Config.NewRequest("POST", config.INVENTORY_URL+"/changes/batch-create")
	if err != nil {
		return err
	}
	body, err := json.Marshal(CreateInventoryRequest{
		Key:             uuid.New().String(),
		Changes:         inventory,
		IgnoreUnchanged: true,
	})
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(string(resByte))
	}

	return nil
}
