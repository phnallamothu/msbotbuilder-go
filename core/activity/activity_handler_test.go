package activity_test

import (
	"testing"
	"time"

	"github.com/phnallamothu/msbotbuilder-go/core/activity"
	"github.com/phnallamothu/msbotbuilder-go/schema"
	"github.com/stretchr/testify/assert"
)

// TestBot is a test bot that implements the ActivityHandler interface
type TestBot struct {
	activity.BaseActivityHandler
	messageActivityCalled           bool
	membersAddedActivityCalled      bool
	messageReactionActivityCalled   bool
	typingActivityCalled            bool
	endOfConversationActivityCalled bool
	eventActivityCalled             bool
}

// OnTurn handles the incoming activity
func (bot *TestBot) OnTurn(context *activity.TurnContext) error {
	// Directly call the appropriate handler based on activity type
	switch context.Activity.Type {
	case schema.Message:
		_, err := bot.OnMessageActivity(context)
		return err
	case schema.ConversationUpdate:
		_, err := bot.OnConversationUpdateActivity(context)
		return err
	case schema.MessageReactionType:
		_, err := bot.OnMessageReactionActivity(context)
		return err
	case schema.Typing:
		_, err := bot.OnTypingActivity(context)
		return err
	case schema.EndOfConversation:
		_, err := bot.OnEndOfConversationActivity(context)
		return err
	case schema.Event:
		_, err := bot.OnEventActivity(context)
		return err
	default:
		return nil
	}
}

// OnMessageActivity handles message activities
func (bot *TestBot) OnMessageActivity(context *activity.TurnContext) (schema.Activity, error) {
	bot.messageActivityCalled = true
	return schema.Activity{}, nil
}

// OnConversationUpdateActivity handles conversation update activities
func (bot *TestBot) OnConversationUpdateActivity(context *activity.TurnContext) (schema.Activity, error) {
	bot.membersAddedActivityCalled = true
	return schema.Activity{}, nil
}

// OnMessageReactionActivity handles message reaction activities
func (bot *TestBot) OnMessageReactionActivity(context *activity.TurnContext) (schema.Activity, error) {
	bot.messageReactionActivityCalled = true
	return schema.Activity{}, nil
}

// OnTypingActivity handles typing activities
func (bot *TestBot) OnTypingActivity(context *activity.TurnContext) (schema.Activity, error) {
	bot.typingActivityCalled = true
	return schema.Activity{}, nil
}

// OnEndOfConversationActivity handles end of conversation activities
func (bot *TestBot) OnEndOfConversationActivity(context *activity.TurnContext) (schema.Activity, error) {
	bot.endOfConversationActivityCalled = true
	return schema.Activity{}, nil
}

// OnEventActivity handles event activities
func (bot *TestBot) OnEventActivity(context *activity.TurnContext) (schema.Activity, error) {
	bot.eventActivityCalled = true
	return schema.Activity{}, nil
}

func TestOnTurn_MessageActivity(t *testing.T) {
	// Arrange
	bot := &TestBot{}
	turnContext := activity.NewTurnContext(schema.Activity{
		Type: schema.Message,
		Text: "Hello",
	})

	// Act
	err := bot.OnTurn(turnContext)

	// Assert
	assert.NoError(t, err)
	assert.True(t, bot.messageActivityCalled)
}

func TestOnTurn_ConversationUpdateActivity(t *testing.T) {
	// Arrange
	bot := &TestBot{}
	turnContext := activity.NewTurnContext(schema.Activity{
		Type: schema.ConversationUpdate,
		MembersAdded: []schema.ChannelAccount{
			{ID: "user1", Name: "User 1"},
		},
		Recipient: schema.ChannelAccount{ID: "bot1"},
	})

	// Act
	err := bot.OnTurn(turnContext)

	// Assert
	assert.NoError(t, err)
	assert.True(t, bot.membersAddedActivityCalled)
}

func TestOnTurn_MessageReactionActivity(t *testing.T) {
	// Arrange
	bot := &TestBot{}
	turnContext := activity.NewTurnContext(schema.Activity{
		Type: schema.MessageReactionType,
		ReactionsAdded: []schema.MessageReaction{
			{Type: schema.Like},
		},
	})

	// Act
	err := bot.OnTurn(turnContext)

	// Assert
	assert.NoError(t, err)
	assert.True(t, bot.messageReactionActivityCalled)
}

func TestOnTurn_TypingActivity(t *testing.T) {
	// Arrange
	bot := &TestBot{}
	turnContext := activity.NewTurnContext(schema.Activity{
		Type: schema.Typing,
	})

	// Act
	err := bot.OnTurn(turnContext)

	// Assert
	assert.NoError(t, err)
	assert.True(t, bot.typingActivityCalled)
}

func TestOnTurn_EndOfConversationActivity(t *testing.T) {
	// Arrange
	bot := &TestBot{}
	turnContext := activity.NewTurnContext(schema.Activity{
		Type: schema.EndOfConversation,
	})

	// Act
	err := bot.OnTurn(turnContext)

	// Assert
	assert.NoError(t, err)
	assert.True(t, bot.endOfConversationActivityCalled)
}

func TestOnTurn_EventActivity(t *testing.T) {
	// Arrange
	bot := &TestBot{}
	turnContext := activity.NewTurnContext(schema.Activity{
		Type: schema.Event,
		Name: "testEvent",
	})

	// Act
	err := bot.OnTurn(turnContext)

	// Assert
	assert.NoError(t, err)
	assert.True(t, bot.eventActivityCalled)
}

func TestTurnContext_SendActivity(t *testing.T) {
	// Arrange
	turnContext := activity.NewTurnContext(schema.Activity{
		Type:      schema.Message,
		ID:        "activity1",
		Timestamp: time.Now(),
		From: schema.ChannelAccount{
			ID:   "user1",
			Name: "User 1",
		},
		Recipient: schema.ChannelAccount{
			ID:   "bot1",
			Name: "Bot 1",
		},
		Conversation: schema.ConversationAccount{
			ID:   "conversation1",
			Name: "Conversation 1",
		},
		ChannelID:  "test",
		ServiceURL: "https://test.com",
	})

	// Act
	activity, err := turnContext.SendActivity(
		activity.WithText("Hello"),
	)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, schema.Message, activity.Type)
	assert.Equal(t, "Hello", activity.Text)
	assert.Equal(t, "bot1", activity.From.ID)
	assert.Equal(t, "user1", activity.Recipient.ID)
	assert.Equal(t, "conversation1", activity.Conversation.ID)
	assert.Equal(t, "test", activity.ChannelID)
	assert.Equal(t, "https://test.com", activity.ServiceURL)
	assert.True(t, turnContext.Responded)
}

func TestTurnContext_SendActivities(t *testing.T) {
	// Arrange
	turnContext := activity.NewTurnContext(schema.Activity{
		Type:      schema.Message,
		ID:        "activity1",
		Timestamp: time.Now(),
		From: schema.ChannelAccount{
			ID:   "user1",
			Name: "User 1",
		},
		Recipient: schema.ChannelAccount{
			ID:   "bot1",
			Name: "Bot 1",
		},
		Conversation: schema.ConversationAccount{
			ID:   "conversation1",
			Name: "Conversation 1",
		},
		ChannelID:  "test",
		ServiceURL: "https://test.com",
	})

	// Create activities to send
	activity1 := schema.Activity{
		Type: schema.Message,
		Text: "Hello 1",
	}

	activity2 := schema.Activity{
		Type: schema.Message,
		Text: "Hello 2",
	}

	// Act
	activities, err := turnContext.SendActivities([]schema.Activity{activity1, activity2})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(activities))
	assert.Equal(t, "Hello 1", activities[0].Text)
	assert.Equal(t, "Hello 2", activities[1].Text)
	assert.True(t, turnContext.Responded)
}

func TestTurnContext_TurnState(t *testing.T) {
	// Arrange
	turnContext := activity.NewTurnContext(schema.Activity{
		Type: schema.Message,
		Text: "Hello",
	})

	// Act
	turnContext.SetTurnState("key1", "value1")
	turnContext.SetTurnState("key2", 42)

	// Assert
	assert.Equal(t, "value1", turnContext.GetTurnState("key1"))
	assert.Equal(t, 42, turnContext.GetTurnState("key2"))
}
