package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/phnallamothu/msbotbuilder-go/connector/auth"
	"github.com/phnallamothu/msbotbuilder-go/connector/client"
	"github.com/phnallamothu/msbotbuilder-go/core"
	"github.com/phnallamothu/msbotbuilder-go/core/activity"
	"github.com/phnallamothu/msbotbuilder-go/schema"
)

func putRequest(u string, data []byte) error {
	client := &http.Client{}
	dec, err := url.QueryUnescape(u)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, dec, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	size := fmt.Sprintf("%d", len(data))
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Length", size)
	req.Header.Set("Content-Range", fmt.Sprintf("bytes 0-%d/%d", len(data)-1, len(data)))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return fmt.Errorf("failed to upload file with status %d", resp.StatusCode)
	}
	return nil
}

// HTTPHandler handles the HTTP requests from then connector service
type HTTPHandler struct {
	Adapter *core.BotFrameworkAdapter
}

func (ht *HTTPHandler) cleanupConsents(ID string, ref schema.ConversationReference) {
	fmt.Printf("Deleting activity %s\n", ID)
	ctx := context.Background()
	if err := ht.Adapter.DeleteActivity(ctx, ID, ref); err != nil {
		log.Printf("Failed to delete activity. %s", err.Error())
	}
}

func (ht *HTTPHandler) processMessage(w http.ResponseWriter, req *http.Request) {
	// Process the activity directly using the request
	// Create a custom handler that processes both message and invoke activities
	handler := func(turn *activity.TurnContext) error {
		// Handle different activity types
		switch turn.Activity.Type {
		case schema.Message:
			// Handle message activity
			fi, err := os.Stat("data.txt")
			if err != nil {
				return fmt.Errorf("failed to read file: %s", err.Error())
			}

			// send file upload request
			attachments := []schema.Attachment{
				{
					ContentType: "application/vnd.microsoft.teams.card.file.consent",
					Name:        "data.txt",
					Content: map[string]interface{}{
						"description": "Sample data",
						"sizeInBytes": fi.Size(),
					},
				},
			}
			_, err = turn.SendActivity(activity.WithText("Echo: "+turn.Activity.Text), activity.WithAttachments(attachments))
			return err

		case schema.Invoke:
			// Handle invoke activity
			ht.cleanupConsents(turn.Activity.ReplyToID, activity.GetCoversationReference(turn.Activity))
			data, err := ioutil.ReadFile("data.txt")
			if err != nil {
				return fmt.Errorf("failed to read file: %s", err.Error())
			}
			if turn.Activity.Value["type"] != "fileUpload" {
				return nil
			}
			if turn.Activity.Value["action"] != "accept" {
				return nil
			}

			// parse upload info from invoke accept response
			uploadInfo := schema.UploadInfo{}
			infoJSON, err := json.Marshal(turn.Activity.Value["uploadInfo"])
			if err != nil {
				return err
			}
			err = json.Unmarshal(infoJSON, &uploadInfo)
			if err != nil {
				return err
			}

			// upload file
			err = putRequest(uploadInfo.UploadURL, data)
			if err != nil {
				return fmt.Errorf("failed to upload file: %s", err.Error())
			}

			// notify user about uploaded file
			fileAttach := []schema.Attachment{
				{
					ContentType: "application/vnd.microsoft.teams.card.file.info",
					ContentURL:  uploadInfo.ContentURL,
					Name:        uploadInfo.Name,
					Content: map[string]interface{}{
						"uniqueId": uploadInfo.UniqueID,
						"fileType": uploadInfo.FileType,
					},
				},
			}

			_, err = turn.SendActivity(activity.WithAttachments(fileAttach))
			return err

		default:
			return nil
		}
	}

	err := ht.Adapter.ProcessActivity(w, req, handler)
	if err != nil {
		fmt.Println("Failed to process request.", err)
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
