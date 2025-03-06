package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/phnallamothu/msbotbuilder-go/connector/auth"
	"github.com/phnallamothu/msbotbuilder-go/connector/client"
	"github.com/phnallamothu/msbotbuilder-go/core/activity"
	"github.com/phnallamothu/msbotbuilder-go/schema"
	"github.com/pkg/errors"
)

// Adapter is the primary interface for the user program to perform operations with
// the connector service.
type Adapter interface {
	ParseRequest(ctx context.Context, req *http.Request) (schema.Activity, error)
	ProcessActivity(w http.ResponseWriter, r *http.Request, handler func(*activity.TurnContext) error) error
	ProactiveMessage(ctx context.Context, ref schema.ConversationReference, handler func(*activity.TurnContext) error) error
	DeleteActivity(ctx context.Context, activityID string, ref schema.ConversationReference) error
	UpdateActivity(ctx context.Context, activity schema.Activity) error
	SendActivity(ctx context.Context, activity schema.Activity) error
	SendActivities(ctx context.Context, activities []schema.Activity) error
}

// AdapterSettings is the configuration for the Adapter.
type AdapterSettings struct {
	AppID              string
	AppPassword        string
	ChannelAuthTenant  string
	OauthEndpoint      string
	OpenIDMetadata     string
	ChannelService     string
	CredentialProvider auth.CredentialProvider
	AuthClient         *http.Client
	ReplyClient        *http.Client
}

// BotFrameworkAdapter implements Adapter and is currently the only implementation returned to the user program.
type BotFrameworkAdapter struct {
	AdapterSettings
	auth.TokenValidator
	client.Client
}

// NewBotFrameworkAdapter creates and returns a new BotFrameworkAdapter with the specified AdapterSettings.
func NewBotFrameworkAdapter(settings AdapterSettings) (*BotFrameworkAdapter, error) {
	// If no credential provider is specified, create a simple one
	if settings.CredentialProvider == nil {
		settings.CredentialProvider = auth.SimpleCredentialProvider{
			AppID:    settings.AppID,
			Password: settings.AppPassword,
		}
	}

	if settings.ChannelService == "" {
		settings.ChannelService = auth.ChannelService
	}

	// Prepare new config and Client
	clientConfig, err := client.NewClientConfig(settings.CredentialProvider, auth.ToChannelFromBotLoginURL[0])
	if err != nil {
		return nil, err
	}

	if settings.AuthClient != nil {
		clientConfig.AuthClient = settings.AuthClient
	}

	if settings.ReplyClient != nil {
		clientConfig.ReplyClient = settings.ReplyClient
	}

	connectorClient, err := client.NewClient(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Connector Client.")
	}

	return &BotFrameworkAdapter{settings, auth.NewJwtTokenValidator(), connectorClient}, nil
}

// ProcessActivity receives an activity, processes it as specified in by the 'handler' and
// sends it to the connector service.
func (bf *BotFrameworkAdapter) ProcessActivity(w http.ResponseWriter, r *http.Request, handler func(*activity.TurnContext) error) error {
	ctx := r.Context()

	// Parse the request
	activityRequest, err := bf.ParseRequest(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return errors.Wrap(err, "Failed to parse request.")
	}

	// Create a new turn context
	turnContext := activity.NewTurnContext(activityRequest)

	// Process the activity
	err = handler(turnContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return errors.Wrap(err, "Handler returned an error.")
	}

	return nil
}

// ProactiveMessage sends activity to a conversation.
// This methods is used for Bot initiated conversation.
func (bf *BotFrameworkAdapter) ProactiveMessage(ctx context.Context, ref schema.ConversationReference, handler func(*activity.TurnContext) error) error {
	// Prepare activity with conversation reference
	activityMsg := activity.ApplyConversationReference(schema.Activity{Type: schema.Message}, ref, true)
	return bf.ProcessActivity(nil, nil, func(turnContext *activity.TurnContext) error {
		turnContext.Activity = activityMsg
		return handler(turnContext)
	})
}

// DeleteActivity Deletes an existing activity by Activity ID
func (bf *BotFrameworkAdapter) DeleteActivity(ctx context.Context, activityID string, ref schema.ConversationReference) error {
	// Prepare activity with conversation reference
	activityMsg := activity.ApplyConversationReference(schema.Activity{Type: schema.Message}, ref, true)
	activityMsg.ID = activityID

	// Create a response object to handle the deletion
	response, err := activity.NewActivityResponse(bf.Client)
	if err != nil {
		return errors.Wrap(err, "Failed to create response object.")
	}

	return response.DeleteActivity(ctx, activityMsg)
}

// ParseRequest parses the received activity in a HTTP reuqest to:
//
// 1. Validate the structure.
//
// 2. Authenticate the request (using authenticateRequest())
//
// Returns an Activity value on successfull parsing.
func (bf *BotFrameworkAdapter) ParseRequest(ctx context.Context, req *http.Request) (schema.Activity, error) {
	activity := schema.Activity{}
	// Find auth headers
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return activity, errors.New("Authentication headers are missing in the request")
	}

	// Parse request body
	err := json.NewDecoder(req.Body).Decode(&activity)
	if err != nil {
		return activity, errors.Wrap(err, "Error while parsing Bot request")
	}
	return activity, bf.authenticateRequest(ctx, activity, authHeader)
}

func (bf *BotFrameworkAdapter) authenticateRequest(ctx context.Context, req schema.Activity, headers string) error {

	_, err := bf.TokenValidator.AuthenticateRequest(ctx, req, headers, bf.CredentialProvider, bf.ChannelService)

	return errors.Wrap(err, "Authentication failed.")
}

// UpdateActivity updates an existing activity.
func (bf *BotFrameworkAdapter) UpdateActivity(ctx context.Context, activityToUpdate schema.Activity) error {
	response, err := activity.NewActivityResponse(bf.Client)

	if err != nil {
		return errors.Wrap(err, "Failed to create response object.")
	}
	return response.UpdateActivity(ctx, activityToUpdate)
}

// SendActivity sends an activity to the conversation referenced in the activity.
func (bf *BotFrameworkAdapter) SendActivity(ctx context.Context, outgoingActivity schema.Activity) error {
	response, err := activity.NewActivityResponse(bf.Client)
	if err != nil {
		return errors.Wrap(err, "Failed to create response object.")
	}

	return response.SendActivity(ctx, outgoingActivity)
}

// SendActivities sends multiple activities to the conversation referenced in the activities.
func (bf *BotFrameworkAdapter) SendActivities(ctx context.Context, outgoingActivities []schema.Activity) error {
	if len(outgoingActivities) == 0 {
		return nil
	}

	response, err := activity.NewActivityResponse(bf.Client)
	if err != nil {
		return errors.Wrap(err, "Failed to create response object.")
	}

	for _, outAct := range outgoingActivities {
		err := response.SendActivity(ctx, outAct)
		if err != nil {
			return errors.Wrap(err, "Failed to send activity.")
		}
	}

	return nil
}
