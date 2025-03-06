package activity

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/phnallamothu/msbotbuilder-go/connector/client"
	"github.com/phnallamothu/msbotbuilder-go/schema"
	"github.com/pkg/errors"
)

// Response provides functionalities to send activity to the connector service.
type Response interface {
	SendActivity(ctx context.Context, activity schema.Activity) error
	DeleteActivity(ctx context.Context, activity schema.Activity) error
	UpdateActivity(ctx context.Context, activity schema.Activity) error
}

const (
	// APIVersion for response URLs
	APIVersion = "v3"

	sendToConversationURL = "/%s/conversations/%s/activities"
	activityResourceURL   = "/%s/conversations/%s/activities/%s"
)

// DefaultResponse is the default implementation of Response.
type DefaultResponse struct {
	Client client.Client
}

// DeleteActivity sends a Delete activity method to the BOT connector service.
func (response *DefaultResponse) DeleteActivity(ctx context.Context, activity schema.Activity) error {
	u, err := url.Parse(activity.ServiceURL)
	if err != nil {
		return errors.Wrapf(err, "Failed to parse ServiceURL %s.", activity.ServiceURL)
	}

	respPath := fmt.Sprintf(activityResourceURL, APIVersion, activity.Conversation.ID, activity.ID)

	// Send activity to client
	u.Path = path.Join(u.Path, respPath)
	err = response.Client.Delete(ctx, *u)
	return errors.Wrap(err, "Failed to delete response.")
}

// SendActivity sends an activity to the BOT connector service.
func (response *DefaultResponse) SendActivity(ctx context.Context, activity schema.Activity) error {
	u, err := url.Parse(activity.ServiceURL)
	if err != nil {
		return errors.Wrapf(err, "Failed to parse ServiceURL %s.", activity.ServiceURL)
	}

	respPath := fmt.Sprintf(sendToConversationURL, APIVersion, activity.Conversation.ID)

	// if ReplyToID is set in the activity, we send reply to that particular activity
	if activity.ReplyToID != "" {
		respPath = fmt.Sprintf(activityResourceURL, APIVersion, activity.Conversation.ID, activity.ID)
	}

	// Send activity to client
	u.Path = path.Join(u.Path, respPath)
	err = response.Client.Post(ctx, *u, activity)
	return errors.Wrap(err, "Failed to send response.")
}

// UpdateActivity sends a Put activity method to the BOT connector service.
func (response *DefaultResponse) UpdateActivity(ctx context.Context, activity schema.Activity) error {
	u, err := url.Parse(activity.ServiceURL)
	if err != nil {
		return errors.Wrapf(err, "Failed to parse ServiceURL %s.", activity.ServiceURL)
	}

	respPath := fmt.Sprintf(activityResourceURL, APIVersion, activity.Conversation.ID, activity.ID)

	// Send activity to client
	u.Path = path.Join(u.Path, respPath)
	err = response.Client.Put(ctx, *u, activity)
	return errors.Wrap(err, "Failed to update response.")
}

// NewActivityResponse provides a DefaultResponse implementaton of Response.
func NewActivityResponse(connectorClient client.Client) (Response, error) {
	if connectorClient == nil {
		return nil, errors.New("Invalid connector client for ActivityResponse")
	}

	return &DefaultResponse{connectorClient}, nil
}
