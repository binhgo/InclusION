package main

import (
	"github.com/centrifugal/centrifuge"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/InclusION/chat"
)

//**********************************************************************************//
// centrifuge
func handleLog(e centrifuge.LogEntry) {
	log.Printf("%s: %v", e.Message, e.Fields)
}

// Wait until program interrupted. When interrupted gracefully shutdown Node.
func waitExitSignal(n *centrifuge.Node) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		n.Shutdown(ctx)
		done <- true
	}()
	<-done
}

func initCentrifuge() *centrifuge.Node {

	cfg := centrifuge.DefaultConfig
	cfg.ClientInsecure = true
	cfg.Publish = true
	node, _ := centrifuge.New(cfg)

	node.On().Connect(func(ctx context.Context, client *centrifuge.Client, e centrifuge.ConnectEvent) centrifuge.ConnectReply {

		client.On().Subscribe(func(e centrifuge.SubscribeEvent) centrifuge.SubscribeReply {

			client.ID()
			log.Printf("user id %s", client.UserID())
			log.Printf("id %s", client.ID())
			log.Printf("client %x", client)

			ok := chat.ValidateUserJoinRoom(client.UserID())

			if ok {
				log.Printf("client subscribes on channel %s", e.Channel)
				return centrifuge.SubscribeReply{}

			} else {

				err1 := centrifuge.Error {
					Code:    109,
					Message: "!permission",
				}

				return centrifuge.SubscribeReply{Error: &err1}
			}
		})

		client.On().Publish(func(e centrifuge.PublishEvent) centrifuge.PublishReply {
			log.Printf("client publishes into channel %s: %s", e.Channel, string(e.Data))
			return centrifuge.PublishReply{}
		})

		// Set Disconnect Handler to react on client disconnect events.
		client.On().Disconnect(func(e centrifuge.DisconnectEvent) centrifuge.DisconnectReply {
			log.Printf("client disconnected")
			return centrifuge.DisconnectReply{}
		})

		// In our example transport will always be Websocket but it can also be SockJS.
		transportName := client.Transport().Name()
		// In our example clients connect with JSON protocol but it can also be Protobuf.
		transportEncoding := client.Transport().Encoding()

		log.Printf("client connected via %s (%s)", transportName, transportEncoding)
		return centrifuge.ConnectReply{}
	})

	node.SetLogHandler(centrifuge.LogLevelDebug, handleLog)

	if err := node.Run(); err != nil {
		panic(err)
	}

	return node
}
// centrifuge
//**********************************************************************************//
