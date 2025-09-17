package bash

import (
	"encoding/base64"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"irpl.com/kanban-commons/utils"
)

var onceNats sync.Once
var instanceBashClient *NatsClientEvents

type NatsClientEvents struct {
	nc      *nats.Conn
	subject string
}

const (
	DefaultNatsSubject = "CMDLINE.BASH.COMMAND"
)

// TODO- ###
func NatsCli() *NatsClientEvents {

	Broker := os.Getenv("BROKER")
	if strings.TrimSpace(Broker) == "" {
		Broker = "0.0.0.0:4222"
	}

	onceNats.Do(func() {

		path := []string{"nats://", Broker}
		// Setup options with retries and a timeout
		opts := []nats.Option{
			nats.MaxReconnects(-1),              // Infinite retries
			nats.ReconnectWait(2 * time.Second), // Wait time between retries
			nats.Timeout(5 * time.Minute),       // Connection timeout
			nats.RetryOnFailedConnect(true),     // Enable retry on initial connect
			nats.Name("bash.nats.client-events"),
			nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
				log.Println("Disconnected from NATS server.", err)
			}),
		}

		var error error

		instanceBashClient = &NatsClientEvents{
			nc:      &nats.Conn{},
			subject: DefaultNatsSubject,
		}
		if instanceBashClient.nc, error = nats.Connect(strings.Join(path, ""), opts...); error != nil {
			log.Fatal(error.Error())
		}

		if len(instanceBashClient.subject) == 0 {
			instanceBashClient.subject = DefaultNatsSubject // Set the default NATS subject for Bash NATS client
		}

	})

	return instanceBashClient
}

func (d *NatsClientEvents) Command(command ...string) (string, error) {
	return d.SendCommand(utils.JoinStr(command...))
}

func (d *NatsClientEvents) SendCommand(command string) (string, error) {

	return d.SlowCommand(2*time.Second, command)

}

func (d *NatsClientEvents) SlowCommand(timeout time.Duration, command string) (string, error) {

	if d.nc == nil {
		return "", errors.New("nats connection is not established")
	}

	encodedBytes := base64.StdEncoding.EncodeToString([]byte(command))

	// Publish the script name to the specified subject and request a reply
	msg, err := d.nc.Request(d.subject, []byte(encodedBytes), timeout)
	if err != nil {
		return "", errors.New(utils.JoinStr("error: sending request to NATS: ", err.Error()))
	}

	result := string(msg.Data)

	if strings.HasPrefix(result, "error:") {

		return "", errors.New(strings.Replace(result, "error:", "", 1))

	}

	return result, nil
}
