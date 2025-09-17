package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"irpl.com/kanban-commons/utils"
)

/*

to run as client
 /RUBBER/scripts/service NatsClient  0.0.0.0:4222 CMDLINE.BASH.COMMAND &

*/

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please provide a valid command, ex: utils status")
		os.Exit(1)
	}

	command := os.Args[1]

	var err error
	switch command {
	case "encrypt":

		if len(os.Args) < 4 {
			fmt.Println("encrypt: Please provide an input filename, and an output filename")
			os.Exit(1)
		} else {
			err = utils.KeylessFileEncrypt(os.Args[2], os.Args[3])
		}

	case "decrypt":

		if len(os.Args) < 4 {
			fmt.Println("decrypt: Please provide an input filename, and an output filename")
			os.Exit(1)
		} else {
			err = utils.KeylessFileDecrypt(os.Args[2], os.Args[3])
		}

	case "signTarFile":

		if len(os.Args) < 3 {
			fmt.Println("signtar: Please provide the tar file full path")
			os.Exit(1)
		} else {
			utils.ProcessTarFile(os.Args[2])
		}

	case "validateTarFile":

		if len(os.Args) < 3 {
			fmt.Println("signtar: Please provide the tar file full path")
			os.Exit(1)
		} else {
			utils.ReverseTarFile(os.Args[2])
		}

	case "changeDate":

		if len(os.Args) < 4 {
			fmt.Println("changeDate: Please provide the new date and key")
			os.Exit(1)
		} else {
			utils.ChangeDate(os.Args[2])
		}

	case "StartSettings":

		if len(os.Args) < 3 {
			fmt.Println("StartSettings: Please provide the key")
			os.Exit(1)
		} else {
			utils.StartSettings(os.Args[2])
			exec.Command("sudo", "gnome-terminal", "--", "bash", "-c", "su - varco; exec bash").Run()

		}

	case "Post":

		if len(os.Args) < 3 {
			fmt.Println("Post: Please provide URL, Authorization, payload (string) ", os.Args)
			os.Exit(1)
		} else {

			data, err := Post(os.Args[2])
			if err == nil {
				fmt.Println(string(data))
				os.Exit(0)
			} else {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

	case "status":

		fmt.Println("OK")

	case "reboot":

		reboot()

	case "terminal":

		exec.Command("sudo", "gnome-terminal", "--", "bash", "-c", "su - varco; exec bash").Run()

	case "Test":

		exec.Command("bash", "-c", `for arg in "$@"; do echo "$arg"; done`, "--").Run()

	case "NatsClient":

		if len(os.Args) < 4 {
			fmt.Println("NatsClient: Please provide the broker ip:port and Topic")
			os.Exit(1)
		} else {
			natsCli(os.Args[2], os.Args[3])
		}

	default:

		fmt.Println("Invalid command, please provide a valid command")
		os.Exit(1)

	}

	if err != nil {
		fmt.Printf("Failed to %s file: %v\n", command, err)
		os.Exit(1)
	}

}

func reboot() {

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("shutdown", "-r", "-t", "0")
	case "linux":
		cmd = exec.Command("sudo", "shutdown", "-r", "now")
	case "darwin":
		cmd = exec.Command("sudo", "shutdown", "-r", "now")
	default:
		log.Fatalf("Unsupported operating system: %s", runtime.GOOS)
	}

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to restart: %v", err)
	}
}

func Post(compress string) ([]byte, error) {

	// Create a new POST request with the JSON data
	decodedBytes, _ := base64.StdEncoding.DecodeString(compress)

	parts := strings.SplitN(string(decodedBytes), ";", 3)
	if len(parts) != 3 {
		fmt.Println("Input does not match format")
	}

	url := parts[0]
	invite := parts[1]
	data := parts[2]

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", invite)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-200 status: %d", resp.StatusCode)
	}

	return bodyBytes, nil

}

// run bash commands from nats messages
func natsCli(broker, subject string) {

	// Setup options with retries and a timeout
	opts := []nats.Option{
		nats.MaxReconnects(-1),              // Infinite retries
		nats.ReconnectWait(2 * time.Second), // Wait time between retries
		nats.Timeout(5 * time.Minute),       // Connection timeout
		nats.RetryOnFailedConnect(true),     // Enable retry on initial connect
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Println("Disconnected from NATS server.", err)
		}),
	}

	// Connect to NATS
	nc, err := nats.Connect(fmt.Sprintf("nats://%s", broker), opts...)
	if err != nil {
		fmt.Println("Failed to connect to NATS:", err)
		return
	}
	defer nc.Close()

	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {

		go func() { // looks the best way for long commands or scripts

			command, _ := base64.StdEncoding.DecodeString(string(msg.Data))

			cmd := exec.Command("bash", "-c", string(command))

			// Create buffers to capture standard output and standard error.
			var stdoutBuf, stderrBuf bytes.Buffer
			cmd.Stdout = &stdoutBuf
			cmd.Stderr = &stderrBuf

			// Execute the command.
			err = cmd.Run()
			output := ""

			if err != nil {
				output = "error:" + stderrBuf.String() + err.Error()
			} else {
				output = stdoutBuf.String()
			}

			// Reply with the command output
			if err := nc.Publish(msg.Reply, []byte(output)); err != nil {
				log.Println("problem sending the answer:", err.Error())
			}
		}()

	})

	if err != nil {
		fmt.Println("Failed to subscribe to subject:", err)
		return
	}

	select {}

}
