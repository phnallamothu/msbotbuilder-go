/*
Bot Framework echo bot sample.
This bot uses msbotbuilder-go: https://github.com/phnallamothu/msbotbuilder-go. It shows
how to create a simple bot that accepts input from the user and echoes it back.

# Run the example

Bring up a terminal and run. Set two variables for the session as APP_ID and APP_PASSWORD to the values of
your BotFramework app_id and password. Then, run:

	go run main.go

# This will start a server which will listen on port 3978

# Understanding the example

The program starts by creating a handler function that will be called when a message is received.
This function takes a TurnContext parameter and returns an error.

	var customHandler = func(turn *activity.TurnContext) error {
		_, err := turn.SendActivity(activity.WithText("Echo: " + turn.Activity.Text))
		return err
	}

A webserver is started with a handler that passes the received HTTP request directly to `adapter.ProcessActivity`
This method authenticates the payload, parses the request, creates a TurnContext, and calls the handler function.

	err := ht.Adapter.ProcessActivity(w, req, customHandler)

# In case of no error, this web responds with a 200 status

To expose this local IP outside your local network, a tool like ngrok can be used.

	ngrok http 3978

The server is then available on a IP similar to http://92832de0.ngrok.io
*/
package main
