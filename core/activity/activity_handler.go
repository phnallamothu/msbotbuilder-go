package activity

import (
	"errors"
	"net/http"

	"github.com/phnallamothu/msbotbuilder-go/schema"
)

// ActivityHandler is the base class for handling activities from the connector service.
// It provides default implementations for all activity types and can be extended to provide
// custom behavior for specific activity types.
type ActivityHandler interface {
	// OnTurn is the main method that handles all activity types.
	OnTurn(context *TurnContext) error

	// Message Activities
	OnMessageActivity(context *TurnContext) (schema.Activity, error)
	OnMessageUpdateActivity(context *TurnContext) (schema.Activity, error)
	OnMessageDeleteActivity(context *TurnContext) (schema.Activity, error)

	// Conversation Update Activities
	OnConversationUpdateActivity(context *TurnContext) (schema.Activity, error)
	OnMembersAddedActivity(membersAdded []schema.ChannelAccount, context *TurnContext) (schema.Activity, error)
	OnMembersRemovedActivity(membersRemoved []schema.ChannelAccount, context *TurnContext) (schema.Activity, error)

	// Message Reaction Activities
	OnMessageReactionActivity(context *TurnContext) (schema.Activity, error)
	OnReactionsAddedActivity(messageReactions []schema.MessageReaction, context *TurnContext) (schema.Activity, error)
	OnReactionsRemovedActivity(messageReactions []schema.MessageReaction, context *TurnContext) (schema.Activity, error)

	// Event Activities
	OnEventActivity(context *TurnContext) (schema.Activity, error)
	OnTokenResponseEvent(context *TurnContext) (schema.Activity, error)

	// Invoke Activities
	OnInvokeActivity(context *TurnContext) (*InvokeResponse, error)
	OnSignInInvokeActivity(context *TurnContext) (schema.Activity, error)
	OnAdaptiveCardInvokeActivity(context *TurnContext, invokeValue interface{}) (*AdaptiveCardInvokeResponse, error)

	// Other Activities
	OnEndOfConversationActivity(context *TurnContext) (schema.Activity, error)
	OnTypingActivity(context *TurnContext) (schema.Activity, error)
	OnInstallationUpdateActivity(context *TurnContext) (schema.Activity, error)
	OnInstallationUpdateAddActivity(context *TurnContext) (schema.Activity, error)
	OnInstallationUpdateRemoveActivity(context *TurnContext) (schema.Activity, error)
	OnUnrecognizedActivityType(context *TurnContext) (schema.Activity, error)
}

// InvokeResponse represents a response to an invoke activity
type InvokeResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body,omitempty"`
}

// AdaptiveCardInvokeResponse represents a response to an adaptive card invoke activity
type AdaptiveCardInvokeResponse struct {
	StatusCode int         `json:"statusCode"`
	Type       string      `json:"type"`
	Value      interface{} `json:"value"`
}

// BaseActivityHandler provides default implementations for all activity types.
// It can be embedded in a struct to provide a base implementation of the ActivityHandler interface.
type BaseActivityHandler struct{}

// OnTurn handles the incoming activity
func (h *BaseActivityHandler) OnTurn(context *TurnContext) error {
	var err error

	// Handle different activity types
	switch context.Activity.Type {
	case schema.Message:
		// Use the receiver directly to ensure dynamic dispatch
		response, err := h.OnMessageActivity(context)
		if err == nil && response.Type != "" {
			_, err = context.SendActivity(WithActivity(response))
		}
	case schema.MessageUpdate:
		response, err := h.OnMessageUpdateActivity(context)
		if err == nil && response.Type != "" {
			_, err = context.SendActivity(WithActivity(response))
		}
	case schema.MessageDelete:
		response, err := h.OnMessageDeleteActivity(context)
		if err == nil && response.Type != "" {
			_, err = context.SendActivity(WithActivity(response))
		}
	case schema.ConversationUpdate:
		response, err := h.OnConversationUpdateActivity(context)
		if err == nil && response.Type != "" {
			_, err = context.SendActivity(WithActivity(response))
		}
	case schema.MessageReactionType:
		response, err := h.OnMessageReactionActivity(context)
		if err == nil && response.Type != "" {
			_, err = context.SendActivity(WithActivity(response))
		}
	case schema.Event:
		response, err := h.OnEventActivity(context)
		if err == nil && response.Type != "" {
			_, err = context.SendActivity(WithActivity(response))
		}
	case schema.Invoke:
		// Special handling for invoke activities
		_, err = h.OnInvokeActivity(context)
	case schema.EndOfConversation:
		response, err := h.OnEndOfConversationActivity(context)
		if err == nil && response.Type != "" {
			_, err = context.SendActivity(WithActivity(response))
		}
	case schema.Typing:
		response, err := h.OnTypingActivity(context)
		if err == nil && response.Type != "" {
			_, err = context.SendActivity(WithActivity(response))
		}
	case schema.InstallationUpdate:
		response, err := h.OnInstallationUpdateActivity(context)
		if err == nil && response.Type != "" {
			_, err = context.SendActivity(WithActivity(response))
		}
	default:
		response, err := h.OnUnrecognizedActivityType(context)
		if err == nil && response.Type != "" {
			_, err = context.SendActivity(WithActivity(response))
		}
	}

	return err
}

// OnMessageActivity handles message activities
func (h *BaseActivityHandler) OnMessageActivity(context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnMessageUpdateActivity handles message update activities
func (h *BaseActivityHandler) OnMessageUpdateActivity(context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnMessageDeleteActivity handles message delete activities
func (h *BaseActivityHandler) OnMessageDeleteActivity(context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnConversationUpdateActivity handles conversation update activities
func (h *BaseActivityHandler) OnConversationUpdateActivity(context *TurnContext) (schema.Activity, error) {
	var err error
	var response schema.Activity

	if context.Activity.MembersAdded != nil && len(context.Activity.MembersAdded) > 0 {
		response, err = h.OnMembersAddedActivity(context.Activity.MembersAdded, context)
	} else if context.Activity.MembersRemoved != nil && len(context.Activity.MembersRemoved) > 0 {
		response, err = h.OnMembersRemovedActivity(context.Activity.MembersRemoved, context)
	}

	return response, err
}

// OnMembersAddedActivity handles members added activities
func (h *BaseActivityHandler) OnMembersAddedActivity(membersAdded []schema.ChannelAccount, context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnMembersRemovedActivity handles members removed activities
func (h *BaseActivityHandler) OnMembersRemovedActivity(membersRemoved []schema.ChannelAccount, context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnMessageReactionActivity handles message reaction activities
func (h *BaseActivityHandler) OnMessageReactionActivity(context *TurnContext) (schema.Activity, error) {
	var err error
	var response schema.Activity

	if context.Activity.ReactionsAdded != nil && len(context.Activity.ReactionsAdded) > 0 {
		response, err = h.OnReactionsAddedActivity(context.Activity.ReactionsAdded, context)
	} else if context.Activity.ReactionsRemoved != nil && len(context.Activity.ReactionsRemoved) > 0 {
		response, err = h.OnReactionsRemovedActivity(context.Activity.ReactionsRemoved, context)
	}

	return response, err
}

// OnReactionsAddedActivity handles reactions added activities
func (h *BaseActivityHandler) OnReactionsAddedActivity(messageReactions []schema.MessageReaction, context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnReactionsRemovedActivity handles reactions removed activities
func (h *BaseActivityHandler) OnReactionsRemovedActivity(messageReactions []schema.MessageReaction, context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnEventActivity handles event activities
func (h *BaseActivityHandler) OnEventActivity(context *TurnContext) (schema.Activity, error) {
	if context.Activity.Name == "tokens/response" {
		return h.OnTokenResponseEvent(context)
	}

	return schema.Activity{}, nil
}

// OnTokenResponseEvent handles token response events
func (h *BaseActivityHandler) OnTokenResponseEvent(context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnInvokeActivity handles invoke activities
func (h *BaseActivityHandler) OnInvokeActivity(context *TurnContext) (*InvokeResponse, error) {
	var err error

	if context.Activity.Name == "signin/verifyState" || context.Activity.Name == "signin/tokenExchange" {
		_, err = h.OnSignInInvokeActivity(context)
		if err != nil {
			return &InvokeResponse{Status: http.StatusNotImplemented}, err
		}
		return &InvokeResponse{Status: http.StatusOK}, nil
	}

	if context.Activity.Name == "adaptiveCard/action" {
		// TODO: Parse invoke value
		invokeValue := context.Activity.Value
		response, err := h.OnAdaptiveCardInvokeActivity(context, invokeValue)
		if err != nil {
			return &InvokeResponse{Status: http.StatusBadRequest, Body: response}, err
		}
		return &InvokeResponse{Status: http.StatusOK, Body: response}, nil
	}

	return &InvokeResponse{Status: http.StatusNotImplemented}, errors.New("not implemented")
}

// OnSignInInvokeActivity handles sign-in invoke activities
func (h *BaseActivityHandler) OnSignInInvokeActivity(context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, errors.New("not implemented")
}

// OnAdaptiveCardInvokeActivity handles adaptive card invoke activities
func (h *BaseActivityHandler) OnAdaptiveCardInvokeActivity(context *TurnContext, invokeValue interface{}) (*AdaptiveCardInvokeResponse, error) {
	return nil, errors.New("not implemented")
}

// OnEndOfConversationActivity handles end of conversation activities
func (h *BaseActivityHandler) OnEndOfConversationActivity(context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnTypingActivity handles typing activities
func (h *BaseActivityHandler) OnTypingActivity(context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnInstallationUpdateActivity handles installation update activities
func (h *BaseActivityHandler) OnInstallationUpdateActivity(context *TurnContext) (schema.Activity, error) {
	if context.Activity.Action == "add" || context.Activity.Action == "add-upgrade" {
		return h.OnInstallationUpdateAddActivity(context)
	} else if context.Activity.Action == "remove" || context.Activity.Action == "remove-upgrade" {
		return h.OnInstallationUpdateRemoveActivity(context)
	}

	return schema.Activity{}, nil
}

// OnInstallationUpdateAddActivity handles installation update add activities
func (h *BaseActivityHandler) OnInstallationUpdateAddActivity(context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnInstallationUpdateRemoveActivity handles installation update remove activities
func (h *BaseActivityHandler) OnInstallationUpdateRemoveActivity(context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}

// OnUnrecognizedActivityType handles unrecognized activity types
func (h *BaseActivityHandler) OnUnrecognizedActivityType(context *TurnContext) (schema.Activity, error) {
	return schema.Activity{}, nil
}
