package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/phnallamothu/msbotbuilder-go/connector/auth"
	"github.com/phnallamothu/msbotbuilder-go/schema"
	"github.com/stretchr/testify/assert"
)

// TestMockClient is a mock implementation of the client.Client interface
type TestMockClient struct {
	PostError   error
	DeleteError error
	GetError    error
	PutError    error
}

// Post is a mock implementation
func (m *TestMockClient) Post(ctx context.Context, url url.URL, activity schema.Activity) error {
	return m.PostError
}

// Delete is a mock implementation
func (m *TestMockClient) Delete(ctx context.Context, url url.URL) error {
	return m.DeleteError
}

// Get is a mock implementation
func (m *TestMockClient) Get(ctx context.Context, url url.URL) (json.RawMessage, error) {
	return json.RawMessage{}, m.GetError
}

// Put is a mock implementation
func (m *TestMockClient) Put(ctx context.Context, url url.URL, activity schema.Activity) error {
	return m.PutError
}

// createTestAdapter creates a test adapter with mock components
func createTestAdapter() *BotFrameworkAdapter {
	settings := AdapterSettings{
		AppID:       "test-app-id",
		AppPassword: "test-app-password",
	}

	settings.CredentialProvider = auth.SimpleCredentialProvider{
		AppID:    settings.AppID,
		Password: settings.AppPassword,
	}

	adapter, _ := NewBotFrameworkAdapter(settings)
	// Replace the token validator with our mock
	adapter.TokenValidator = &MockTokenValidator{}
	// Replace the client with our mock
	adapter.Client = &TestMockClient{}

	return adapter
}

// createTestActivity creates a test activity
func createTestActivity() schema.Activity {
	return schema.Activity{
		Type: schema.Message,
		From: schema.ChannelAccount{
			ID:   "user1",
			Name: "User One",
		},
		Recipient: schema.ChannelAccount{
			ID:   "bot1",
			Name: "Bot One",
		},
		Conversation: schema.ConversationAccount{
			ID:   "conversation1",
			Name: "Conversation One",
		},
		ChannelID:   "test-channel",
		ServiceURL:  "https://test.com",
		Text:        "Hello, Bot!",
		InputHint:   schema.AcceptingInput,
		Attachments: []schema.Attachment{},
	}
}

// createTestRequest creates a test HTTP request with the given activity
func createTestRequest(t *testing.T, activity schema.Activity) *http.Request {
	activityJSON, err := json.Marshal(activity)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/messages", bytes.NewBuffer(activityJSON))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	return req
}

// TestNewBotFrameworkAdapter tests the creation of a new adapter
func TestNewBotFrameworkAdapter(t *testing.T) {
	settings := AdapterSettings{
		AppID:       "test-app-id",
		AppPassword: "test-app-password",
	}

	adapter, err := NewBotFrameworkAdapter(settings)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)
	assert.Equal(t, settings.AppID, adapter.AppID)
	assert.Equal(t, settings.AppPassword, adapter.AppPassword)
	assert.NotNil(t, adapter.CredentialProvider)
	assert.NotNil(t, adapter.TokenValidator)
	assert.NotNil(t, adapter.Client)
}

// TestNewBotFrameworkAdapterWithCustomClients tests the creation of a new adapter with custom clients
func TestNewBotFrameworkAdapterWithCustomClients(t *testing.T) {
	customAuthClient := &http.Client{}
	customReplyClient := &http.Client{}

	settings := AdapterSettings{
		AppID:       "test-app-id",
		AppPassword: "test-app-password",
		AuthClient:  customAuthClient,
		ReplyClient: customReplyClient,
	}

	adapter, err := NewBotFrameworkAdapter(settings)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)
}

// TestSendActivity tests the SendActivity method
func TestSendActivity(t *testing.T) {
	adapter := createTestAdapter()
	activity := createTestActivity()

	err := adapter.SendActivity(context.Background(), activity)
	assert.NoError(t, err)
}

// TestSendActivityError tests the SendActivity method when the client returns an error
func TestSendActivityError(t *testing.T) {
	adapter := createTestAdapter()
	// Override the client with one that returns an error
	adapter.Client = &TestMockClient{PostError: errors.New("post error")}
	activity := createTestActivity()

	err := adapter.SendActivity(context.Background(), activity)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "post error")
}

// TestSendActivities tests the SendActivities method
func TestSendActivities(t *testing.T) {
	adapter := createTestAdapter()
	activities := []schema.Activity{createTestActivity(), createTestActivity()}

	err := adapter.SendActivities(context.Background(), activities)
	assert.NoError(t, err)
}

// TestSendActivitiesEmpty tests the SendActivities method with an empty slice
func TestSendActivitiesEmpty(t *testing.T) {
	adapter := createTestAdapter()

	err := adapter.SendActivities(context.Background(), []schema.Activity{})
	assert.NoError(t, err)
}

// TestSendActivitiesError tests the SendActivities method when the client returns an error
func TestSendActivitiesError(t *testing.T) {
	adapter := createTestAdapter()
	// Override the client with one that returns an error
	adapter.Client = &TestMockClient{PostError: errors.New("post error")}
	activities := []schema.Activity{createTestActivity()}

	err := adapter.SendActivities(context.Background(), activities)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "post error")
}

// TestUpdateActivity tests the UpdateActivity method
func TestUpdateActivity(t *testing.T) {
	adapter := createTestAdapter()
	activity := createTestActivity()
	activity.ID = "activity1"

	err := adapter.UpdateActivity(context.Background(), activity)
	assert.NoError(t, err)
}

// TestUpdateActivityError tests the UpdateActivity method when the client returns an error
func TestUpdateActivityError(t *testing.T) {
	adapter := createTestAdapter()
	// Override the client with one that returns an error
	adapter.Client = &TestMockClient{PutError: errors.New("put error")}
	activity := createTestActivity()
	activity.ID = "activity1"

	err := adapter.UpdateActivity(context.Background(), activity)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "put error")
}

// TestDeleteActivity tests the DeleteActivity method
func TestDeleteActivity(t *testing.T) {
	adapter := createTestAdapter()
	conversationRef := schema.ConversationReference{
		ActivityID:   "activity1",
		Conversation: schema.ConversationAccount{ID: "conversation1"},
		ServiceURL:   "https://test.com",
	}

	err := adapter.DeleteActivity(context.Background(), "activity1", conversationRef)
	assert.NoError(t, err)
}

// TestDeleteActivityError tests the DeleteActivity method when the client returns an error
func TestDeleteActivityError(t *testing.T) {
	adapter := createTestAdapter()
	// Override the client with one that returns an error
	adapter.Client = &TestMockClient{DeleteError: errors.New("delete error")}
	conversationRef := schema.ConversationReference{
		ActivityID:   "activity1",
		Conversation: schema.ConversationAccount{ID: "conversation1"},
		ServiceURL:   "https://test.com",
	}

	err := adapter.DeleteActivity(context.Background(), "activity1", conversationRef)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete error")
}
