package main

import (
	"bikeshop/v3/cmd/server/app"
	"bikeshop/v3/cmd/server/catalog"
	"bikeshop/v3/cmd/server/inventory"
	"bikeshop/v3/cmd/server/webhook"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func init() {
	// handle setup of webhooks
	err := webhook.ConfigureWebhooks()
	if err != nil {
		log.Fatal(err)
	}
	// setup the catalog for usage
	err = catalog.DropCatalog()
	if err != nil {
		log.Fatal(err)
	}
	catalogue, err := catalog.CreateCatalogue("../../inventory.csv")
	if err != nil {
		log.Fatal(err)
	}
	// add inventory to the catalogue
	inventory.GenerateInventory(catalogue)
	// set up the required webhooks for usage.
	//first drop existing webbhooks, could also check if they exist
	// create the required webhooks
}

func main() {

	server := app.NewServer(":8081")
	go func() {
		log.Printf("Starting up server")
		if err := server.Start(); err != nil && !strings.Contains(err.Error(), "Server closed") {
			log.Fatalf("problem starting server %v", err)
		}
	}()

	time.Sleep(5 * time.Second)
	log.Print("Server Running ...")
	initShutdownHook(server)
	log.Printf("server shutdown successfully")
}

func initShutdownHook(server *app.Server) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	if err := server.Shutdown(); err != nil {
		log.Fatalf("problem shutting down webserver %v", err)
	}
}
