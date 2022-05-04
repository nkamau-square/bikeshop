package catalog

import (
	"bikeshop/v3/cmd/server/config"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type CatalogueRequest struct {
	Key    string `json:"idempotency_key"`
	Object Item   `json:"object"`
}

type Item struct {
	Id       string   `json:"id"`
	ItemType string   `json:"type"`
	ItemData ItemData `json:"item_data"`
}

type ItemData struct {
	Abbreviation string      `json:"abbreviation"`
	Name         string      `json:"name"`
	Variations   []Variation `json:"variations"`
}

type Variation struct {
	Id            string        `json:"id"`
	VariationType string        `json:"type"`
	VariationData VariationData `json:"item_variation_data"`
}

type VariationData struct {
	Name           string    `json:"name"`
	ItemId         string    `json:"item_id"`
	PricingType    string    `json:"pricing_type"`
	TrackCatalogue bool      `json:"track_inventory"`
	Stockable      bool      `json:"stockable"`
	Sellable       bool      `json:"sellable"`
	Price          PriceData `json:"price_money"`
}

//🥹 lol 😅

type PriceData struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type BatchDeleteCatalogueRequest struct {
	Ids []string `json:"object_ids"`
}

//TODO
func GetCatalogue() ([]byte, error) {
	req, err := config.Config.NewRequest("GET", fmt.Sprintf("%s/list?types=ITEM_VARIATION", config.CATALOGUE_URL))
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func DropCatalog() error {
	req, err := config.Config.NewRequest("GET", fmt.Sprintf("%s/list?types=ITEM", config.CATALOGUE_URL))
	if err != nil {
		return err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	json.Unmarshal(resByte, &result)
	if len(result) == 0 {
		return nil
	}
	objects := result["objects"].([]interface{})
	itemIDs := make([]string, len(objects))
	for i, item := range objects {
		itemIDs[i] = item.(map[string]interface{})["id"].(string)
	}
	body, _ := json.Marshal(BatchDeleteCatalogueRequest{
		Ids: itemIDs,
	})
	req, err = config.Config.NewRequest("POST", fmt.Sprintf("%s/batch-delete", config.CATALOGUE_URL))
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	_, err = client.Do(req)
	return err
}

func CreateCatalogue(filepath string) (map[string]string, error) {
	data, err := readCatalogueFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to read file at %s with error %v", filepath, err)
	}
	cat := convertToItemStruct(data)
	catalogue, err := pushCatalogue(cat)
	if err != nil {
		return nil, fmt.Errorf("unable to push catalogue with error %v", err)
	}
	return catalogue, nil
}

func pushCatalogue(items []*Item) (map[string]string, error) {
	client := &http.Client{}
	req, err := config.Config.NewRequest("POST", fmt.Sprintf("%s/object", config.CATALOGUE_URL))
	if err != nil {
		return nil, err
	}
	catalogIds := make(map[string]string)
	for _, item := range items {
		body, err := json.Marshal(CatalogueRequest{
			Key:    uuid.New().String(),
			Object: *item,
		})
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(body))
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		resByte, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		var result map[string]interface{}
		json.Unmarshal(resByte, &result)
		idMap := result["id_mappings"].([]interface{})
		var variant string
		for i, val := range idMap {
			valMap := val.(map[string]interface{})
			if i == 0 {
				variant = valMap["client_object_id"].(string)
				continue
			}
			catalogIds[variant+valMap["client_object_id"].(string)] = valMap["object_id"].(string)
		}
	}
	return catalogIds, nil
}

func readCatalogueFile(filepath string) ([][]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func convertToItemStruct(data [][]string) (items []*Item) {
	itemMap := make(map[string]*Item)
	for i := 1; i < len(data); i++ {
		line := data[i]
		var item *Item
		var found bool
		if item, found = itemMap[line[1]]; !found {
			item = &Item{
				ItemType: line[0],
				Id:       "#" + line[1],
			}
			itemMap[line[1]] = item
			items = append(items, item)
		}
		item.ItemData.Abbreviation = line[2]
		item.ItemData.Name = line[3]
		variation := Variation{}
		variation.Id = "#" + line[4]
		variation.VariationType = line[5]
		variation.VariationData.Name = line[6]
		variation.VariationData.ItemId = "#" + line[1]
		variation.VariationData.PricingType = line[7]
		variation.VariationData.TrackCatalogue = strings.ToLower(line[8]) == "true"
		variation.VariationData.Stockable = strings.ToLower(line[9]) == "true"
		variation.VariationData.Sellable = strings.ToLower(line[10]) == "true"
		if s, err := strconv.ParseFloat(line[11], 64); err == nil {
			variation.VariationData.Price.Amount = s
		}
		variation.VariationData.Price.Currency = line[12]
		item.ItemData.Variations = append(item.ItemData.Variations, variation)
	}
	return
}
