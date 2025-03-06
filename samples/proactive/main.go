package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/phnallamothu/msbotbuilder-go/connector/auth"
	"github.com/phnallamothu/msbotbuilder-go/connector/client"
	"github.com/phnallamothu/msbotbuilder-go/core"
	"github.com/phnallamothu/msbotbuilder-go/core/activity"
	"github.com/phnallamothu/msbotbuilder-go/schema"
)

// Card content
// Visit: https://adaptivecards.io/explorer to build your own card format
var cardJSON = []byte(`{
  "$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
  "type": "AdaptiveCard",
  "version": "1.0",
  "body": [
    {
      "type": "TextBlock",
      "text": "This is some text",
      "size": "large"
    },
    {
      "type": "TextBlock",
      "text": "It doesn't wrap by default",
      "weight": "bolder"
    },
    {
      "type": "TextBlock",
      "text": "So set **wrap** to true if you plan on showing a paragraph of text",
      "wrap": true
    },
    {
      "type": "TextBlock",
      "text": "You can also use **maxLines** to prevent it from getting out of hand",
      "wrap": true,
      "maxLines": 2
    },
    {
      "type": "TextBlock",
      "text": "You can even draw attention to certain text with color",
      "wrap": true,
      "color": "attention"
    }
  ]
}`)

// conversationRef to store conversation reference to which proactive messages will be sent
var conversationRef schema.ConversationReference

// Define handler functions that match the expected signatures
var welcomeHandler = func(turn *activity.TurnContext) error {
	_, err := turn.SendActivity(activity.WithText("Welcome"))
	return err
}

var attachHandler = func(turn *activity.TurnContext) error {
	var obj map[string]interface{}
	err := json.Unmarshal(cardJSON, &obj)
	if err != nil {
		return err
	}
	attachments := []schema.Attachment{
		{
			ContentType: "application/vnd.microsoft.card.adaptive",
			Content:     obj,
		},
	}
	_, err = turn.SendActivity(activity.WithText("Sample attachment"), activity.WithAttachments(attachments))
	return err
}

// HTTPHandler handles the HTTP requests from then connector service
type HTTPHandler struct {
	Adapter *core.BotFrameworkAdapter
}

func (ht *HTTPHandler) processMessage(w http.ResponseWriter, req *http.Request) {
	// Process the activity directly using the request
	// First, parse the request to get the conversation reference
	ctx := context.Background()
	act, err := ht.Adapter.ParseRequest(ctx, req)
	if err != nil {
		fmt.Println("Failed to parse request.", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set conversation reference
	conversationRef = activity.GetCoversationReference(act)

	// Process the activity
	err = ht.Adapter.ProcessActivity(w, req, welcomeHandler)
	if err != nil {
		fmt.Println("Failed to process request.", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Request processed successfully.")
	// Send proactive message
	ht.welcome()
}

func (ht *HTTPHandler) welcome() {
	ctx := context.Background()
	err := ht.Adapter.ProactiveMessage(ctx, conversationRef, attachHandler)
	if err != nil {
		fmt.Println("Failed to send proactive message.", err)
		return
	}
	fmt.Println("Proactive message sent successfully.")
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
