package activity

import (
	"github.com/phnallamothu/msbotbuilder-go/schema"
	"github.com/pkg/errors"
)

// TurnContext wraps the Activity received and provides operations for the user
// program of this SDK.
//
// The return value is Activity as provided by the client program, to be send to the connector service.
type TurnContext struct {
	Activity  schema.Activity
	Responded bool
	TurnState map[string]interface{}
}

// NewTurnContext creates a new TurnContext with the given activity
func NewTurnContext(activity schema.Activity) *TurnContext {
	return &TurnContext{
		Activity:  activity,
		Responded: false,
		TurnState: make(map[string]interface{}),
	}
}

// SendActivity sends an activity to user.
func (t *TurnContext) SendActivity(options ...MsgOption) (schema.Activity, error) {
	activity, err := applyMsgOptions(schema.Activity{Type: schema.Message}, options...)
	if err != nil {
		return activity, errors.Wrap(err, "Failed to apply MsgOptions.")
	}

	// Mark that we've responded to this activity
	t.Responded = true

	return ApplyConversationReference(activity, GetCoversationReference(t.Activity), false), nil
}

// SendActivities sends multiple activities to user.
func (t *TurnContext) SendActivities(activities []schema.Activity) ([]schema.Activity, error) {
	if len(activities) == 0 {
		return []schema.Activity{}, nil
	}

	resultActivities := make([]schema.Activity, len(activities))

	for i, activity := range activities {
		// Apply conversation reference
		resultActivities[i] = ApplyConversationReference(activity, GetCoversationReference(t.Activity), false)
	}

	// Mark that we've responded to this activity
	t.Responded = true

	return resultActivities, nil
}

// UpdateActivity updates an existing activity.
func (t *TurnContext) UpdateActivity(activity schema.Activity) (schema.Activity, error) {
	// Apply conversation reference
	updatedActivity := ApplyConversationReference(activity, GetCoversationReference(t.Activity), false)

	return updatedActivity, nil
}

// DeleteActivity deletes an existing activity.
func (t *TurnContext) DeleteActivity(activityID string) error {
	// This is just a stub - actual implementation would depend on the connector
	return nil
}

// GetConversationReference gets the conversation reference from the current activity.
func (t *TurnContext) GetConversationReference() schema.ConversationReference {
	return GetCoversationReference(t.Activity)
}

// SetTurnState sets a value in the turn state.
func (t *TurnContext) SetTurnState(key string, value interface{}) {
	t.TurnState[key] = value
}

// GetTurnState gets a value from the turn state.
func (t *TurnContext) GetTurnState(key string) interface{} {
	return t.TurnState[key]
}

func applyMsgOptions(activity schema.Activity, options ...MsgOption) (schema.Activity, error) {
	for _, opt := range options {
		if err := opt(&activity); err != nil {
			return activity, err
		}
	}
	return activity, nil
}
