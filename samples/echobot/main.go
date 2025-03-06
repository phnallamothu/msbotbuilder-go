package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/phnallamothu/msbotbuilder-go/connector/auth"
	"github.com/phnallamothu/msbotbuilder-go/connector/client"
	"github.com/phnallamothu/msbotbuilder-go/core"
	"github.com/phnallamothu/msbotbuilder-go/core/activity"
)

// Define a simple handler function that matches the expected signature
var customHandler = func(turn *activity.TurnContext) error {
	_, err := turn.SendActivity(activity.WithText("Echo: " + turn.Activity.Text))
	return err
}

// HTTPHandler handles the HTTP requests from then connector service
type HTTPHandler struct {
	Adapter *core.BotFrameworkAdapter
}

func (ht *HTTPHandler) processMessage(w http.ResponseWriter, req *http.Request) {
	// Process the activity directly using the request
	err := ht.Adapter.ProcessActivity(w, req, customHandler)
	if err != nil {
		fmt.Println("Failed to process request", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Request processed successfully.")
}

func main() {
	// Create adapter settings
	setting := core.AdapterSettings{
		AppID:       os.Getenv("APP_ID"),
		AppPassword: os.Getenv("APP_PASSWORD"),
	}

	// Create credential provider
	setting.CredentialProvider = auth.SimpleCredentialProvider{
		AppID:    setting.AppID,
		Password: setting.AppPassword,
	}

	// Create client config
	clientConfig, err := client.NewClientConfig(setting.CredentialProvider, auth.ToChannelFromBotLoginURL[0])
	if err != nil {
		log.Fatal("Error creating client config: ", err)
	}

	// Create connector client
	connectorClient, err := client.NewClient(clientConfig)
	if err != nil {
		log.Fatal("Error creating connector client: ", err)
	}

	// Create bot framework adapter
	adapter := &core.BotFrameworkAdapter{
		AdapterSettings: setting,
		TokenValidator:  &core.MockTokenValidator{},
		Client:          connectorClient,
	}

	httpHandler := &HTTPHandler{adapter}

	http.HandleFunc("/api/messages", httpHandler.processMessage)
	fmt.Println("Starting server on port:3978...")
	http.ListenAndServe(":3978", nil)
}
