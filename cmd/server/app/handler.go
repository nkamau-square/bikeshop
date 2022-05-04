package app

import (
	"bikeshop/v3/cmd/server/catalog"
	"bikeshop/v3/cmd/server/inventory"
	"bikeshop/v3/cmd/server/order"
	"bikeshop/v3/cmd/server/payment"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (s *Server) test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func (s *Server) getCatalogue(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	cat, err := catalog.GetCatalogue()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(cat)
}

func (s *Server) getInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var result []string
	json.Unmarshal(reqBody, &result)
	inventory, err := inventory.GetInventory(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(inventory)
}

func (s *Server) purchase(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var result map[string]string
	json.Unmarshal(reqBody, &result)
	// first check if there is remaining quantity
	inventory, err := inventory.GetInventoryCounts([]string{result["id"]})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	icount, err := strconv.Atoi(inventory[result["id"]])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	reqQuant, err := strconv.Atoi(result["quantity"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if reqQuant > icount {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "requested quantity greater than that available"}`))
		return
	}

	orderId, cost, err := order.CreateOrder(result["quantity"], result["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = payment.CreatePayment(orderId, cost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
