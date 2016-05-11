package main

import (
	"fmt"
	"os"

	"github.com/m2mtu/facebookbot/endpoints"
	"github.com/m2mtu/facebookbot/reply"
)

func main() {
	os.Setenv("HTTP_PROXY", os.Getenv("FIXIE_URL"))
	os.Setenv("HTTPS_PROXY", os.Getenv("FIXIE_URL"))
	fmt.Println("starting...")
	endpoints.Listen(handleReceiveMessage)
}

func handleReceiveMessage(receivedEvent endpoints.Event) {
	sendEvent := endpoints.Event{}
	sendEvent.SenderID = receivedEvent.RecepientID
	sendEvent.RecepientID = receivedEvent.SenderID
	switch content := receivedEvent.Content.(type) {
	case endpoints.TextContent:
		sendEvent.Content = endpoints.TextContent{Text: reply.Get(content.Text)}
	}
	endpoints.Send(sendEvent)
}
