package activity

import (
	"github.com/phnallamothu/msbotbuilder-go/schema"
)

// MsgOption is a functional option for schema.Activity
type MsgOption func(*schema.Activity) error

// WithText sets the text property of an activity
func WithText(text string) MsgOption {
	return func(a *schema.Activity) error {
		a.Text = text
		return nil
	}
}

// WithAttachment adds an attachment to an activity
func WithAttachment(attachment schema.Attachment) MsgOption {
	return func(a *schema.Activity) error {
		if a.Attachments == nil {
			a.Attachments = []schema.Attachment{}
		}
		a.Attachments = append(a.Attachments, attachment)
		return nil
	}
}

// WithAttachments sets the attachments property of an activity
func WithAttachments(attachments []schema.Attachment) MsgOption {
	return func(a *schema.Activity) error {
		a.Attachments = attachments
		return nil
	}
}

// WithSuggestedActions sets the suggestedActions property of an activity
func WithSuggestedActions(suggestedActions schema.SuggestedActions) MsgOption {
	return func(a *schema.Activity) error {
		a.SuggestedActions = suggestedActions
		return nil
	}
}

// WithInputHint sets the inputHint property of an activity
func WithInputHint(inputHint schema.InputHints) MsgOption {
	return func(a *schema.Activity) error {
		a.InputHint = inputHint
		return nil
	}
}

// WithSpeak sets the speak property of an activity
func WithSpeak(speak string) MsgOption {
	return func(a *schema.Activity) error {
		a.Speak = speak
		return nil
	}
}

// WithValue sets the value property of an activity
func WithValue(value map[string]interface{}) MsgOption {
	return func(a *schema.Activity) error {
		a.Value = value
		return nil
	}
}

// WithName sets the name property of an activity
func WithName(name string) MsgOption {
	return func(a *schema.Activity) error {
		a.Name = name
		return nil
	}
}

// WithType sets the type property of an activity
func WithType(activityType schema.ActivityTypes) MsgOption {
	return func(a *schema.Activity) error {
		a.Type = activityType
		return nil
	}
}

// WithEntities sets the entities property of an activity
func WithEntities(entities []schema.Entity) MsgOption {
	return func(a *schema.Activity) error {
		a.Entities = entities
		return nil
	}
}

// WithChannelData sets the channelData property of an activity
func WithChannelData(channelData map[string]interface{}) MsgOption {
	return func(a *schema.Activity) error {
		a.ChannelData = channelData
		return nil
	}
}

// WithImportance sets the importance property of an activity
func WithImportance(importance schema.ActivityImportance) MsgOption {
	return func(a *schema.Activity) error {
		a.Importance = importance
		return nil
	}
}

// WithDeliveryMode sets the deliveryMode property of an activity
func WithDeliveryMode(deliveryMode schema.DeliveryModes) MsgOption {
	return func(a *schema.Activity) error {
		a.DeliveryMode = deliveryMode
		return nil
	}
}

// WithSummary sets the summary property of an activity
func WithSummary(summary string) MsgOption {
	return func(a *schema.Activity) error {
		a.Summary = summary
		return nil
	}
}

// WithActivity sets the entire activity
func WithActivity(activity schema.Activity) MsgOption {
	return func(a *schema.Activity) error {
		*a = activity
		return nil
	}
}
