package activity

// This file is kept for backward compatibility
// It is recommended to use the new approach with activity handlers
// that implement the ActivityHandler interface

import (
	"fmt"

	"github.com/phnallamothu/msbotbuilder-go/schema"
)

// Handler acts as the interface for the client program to define actions on various events from connector service.
type Handler interface {
	OnMessage(context *TurnContext) (schema.Activity, error)
	OnInvoke(context *TurnContext) (schema.Activity, error)
	OnConversationUpdate(context *TurnContext) (schema.Activity, error)
}

// PrepareActivityContext routes the received Activity to respective handler function.
// Returns the result of the handler function.
func PrepareActivityContext(handler Handler, context *TurnContext) (schema.Activity, error) {
	switch context.Activity.Type {
	case schema.Message:
		return handler.OnMessage(context)
	case schema.Invoke:
		return handler.OnInvoke(context)
	case schema.ConversationUpdate:
		return handler.OnConversationUpdate(context)
	}
	return schema.Activity{}, fmt.Errorf("Activity type %s not supported yet", context.Activity.Type)
}
