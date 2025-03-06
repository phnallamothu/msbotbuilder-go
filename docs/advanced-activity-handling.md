# Advanced Activity Handling in Bot Builder SDK for Go

## Overview

The Bot Builder SDK for Go now includes advanced activity handling capabilities similar to those found in the Python and .NET SDKs. This document explains how to use these new features to build more sophisticated bots.

## ActivityHandler

The `ActivityHandler` interface provides a more comprehensive way to handle different activity types. It includes methods for handling various activity types such as messages, conversation updates, message reactions, typing, and events.

### Key Components

1. **ActivityHandler Interface**: Defines methods for handling different activity types.
2. **BaseActivityHandler Struct**: Provides default implementations for all activity handler methods.
3. **TurnContext**: Enhanced to support state management and multiple activity operations.

## Using the ActivityHandler

To use the ActivityHandler, create a struct that embeds the `BaseActivityHandler` and override the methods you want to customize:

```go
type MyBot struct {
    activity.BaseActivityHandler
}

// Override OnMessageActivity to handle message activities
func (bot *MyBot) OnMessageActivity(context *activity.TurnContext) (schema.Activity, error) {
    return context.SendActivity(
        activity.WithText(fmt.Sprintf("Echo: %s", context.Activity.Text)),
    )
}
```

## Activity Types

The ActivityHandler supports the following activity types:

### Message Activities
- `OnMessageActivity`: Handles regular message activities.
- `OnMessageUpdateActivity`: Handles message update activities.
- `OnMessageDeleteActivity`: Handles message delete activities.

### Conversation Update Activities
- `OnConversationUpdateActivity`: Handles conversation update activities.
- `OnMembersAddedActivity`: Handles members added activities.
- `OnMembersRemovedActivity`: Handles members removed activities.

### Message Reaction Activities
- `OnMessageReactionActivity`: Handles message reaction activities.
- `OnReactionsAddedActivity`: Handles reactions added activities.
- `OnReactionsRemovedActivity`: Handles reactions removed activities.

### Event Activities
- `OnEventActivity`: Handles event activities.
- `OnTokenResponseEvent`: Handles token response events.

### Invoke Activities
- `OnInvokeActivity`: Handles invoke activities.
- `OnSignInInvokeActivity`: Handles sign-in invoke activities.
- `OnAdaptiveCardInvokeActivity`: Handles adaptive card invoke activities.

### Other Activities
- `OnEndOfConversationActivity`: Handles end of conversation activities.
- `OnTypingActivity`: Handles typing activities.
- `OnInstallationUpdateActivity`: Handles installation update activities.
- `OnUnrecognizedActivityType`: Handles unrecognized activity types.

## Enhanced TurnContext

The `TurnContext` has been enhanced to support state management and multiple activity operations:

- `SendActivity`: Sends a single activity.
- `SendActivities`: Sends multiple activities.
- `UpdateActivity`: Updates an existing activity.
- `DeleteActivity`: Deletes an existing activity.
- `SetTurnState`: Sets a value in the turn state.
- `GetTurnState`: Gets a value from the turn state.

## Example

Here's a complete example of a bot that handles multiple activity types:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/phnallamothu/msbotbuilder-go/connector/auth"
	"github.com/phnallamothu/msbotbuilder-go/core"
	"github.com/phnallamothu/msbotbuilder-go/core/activity"
	"github.com/phnallamothu/msbotbuilder-go/schema"
)

// AdvancedBot demonstrates advanced activity handling capabilities
type AdvancedBot struct {
	activity.BaseActivityHandler
}

// OnMessageActivity handles message activities
func (bot *AdvancedBot) OnMessageActivity(context *activity.TurnContext) (schema.Activity, error) {
	log.Printf("Received message: %s", context.Activity.Text)

	return context.SendActivity(
		activity.WithText(fmt.Sprintf("Echo: %s", context.Activity.Text)),
	)
}

// OnMembersAddedActivity handles members added activities
func (bot *AdvancedBot) OnMembersAddedActivity(membersAdded []schema.ChannelAccount, context *activity.TurnContext) (schema.Activity, error) {
	for _, member := range membersAdded {
		// Skip if the member added is the bot itself
		if member.ID != context.Activity.Recipient.ID {
			return context.SendActivity(
				activity.WithText(fmt.Sprintf("Welcome to Advanced Bot, %s!", member.Name)),
			)
		}
	}

	return schema.Activity{}, nil
}

func main() {
	settings := auth.CredentialProviderSettings{
		AppID:       os.Getenv("APP_ID"),
		AppPassword: os.Getenv("APP_PASSWORD"),
	}

	credentialProvider, err := auth.NewCredentialProvider(settings)
	if err != nil {
		log.Fatal("Error creating credential provider: ", err)
	}

	adapter, err := core.NewBotFrameworkAdapter(core.AdapterSettings{
		CredentialProvider: credentialProvider,
	})
	if err != nil {
		log.Fatal("Error creating adapter: ", err)
	}

	// Create the bot
	bot := &AdvancedBot{}

	// Setup HTTP server
	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		adapter.ProcessActivity(w, r, bot.OnTurn)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3978"
	}

	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
```

## Conclusion

The advanced activity handling capabilities in the Bot Builder SDK for Go provide a more comprehensive way to build sophisticated bots. By leveraging the ActivityHandler interface and the enhanced TurnContext, you can handle various activity types and manage state more effectively.
