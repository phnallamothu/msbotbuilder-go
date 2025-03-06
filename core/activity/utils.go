package activity

import (
	"github.com/phnallamothu/msbotbuilder-go/schema"
)

// GetCoversationReference returns conversation reference from the activity
func GetCoversationReference(activity schema.Activity) schema.ConversationReference {
	return schema.ConversationReference{
		ActivityID:   activity.ID,
		User:         activity.From,
		Bot:          activity.Recipient,
		Conversation: activity.Conversation,
		ChannelID:    activity.ChannelID,
		ServiceURL:   activity.ServiceURL,
	}
}

// ApplyConversationReference sets delivery information to the activity from conversation reference
func ApplyConversationReference(activity schema.Activity, reference schema.ConversationReference, isIncoming bool) schema.Activity {
	activity.ChannelID = reference.ChannelID
	activity.ServiceURL = reference.ServiceURL
	activity.Conversation = reference.Conversation
	if isIncoming {
		activity.From = reference.User
		activity.Recipient = reference.Bot
		if reference.ActivityID != "" {
			activity.ID = reference.ActivityID
		}
		return activity
	}
	activity.From = reference.Bot
	activity.Recipient = reference.User
	if reference.ActivityID != "" {
		activity.ReplyToID = reference.ActivityID
	}
	return activity
}
