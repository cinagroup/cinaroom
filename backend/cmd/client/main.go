package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Version is set at build time via -ldflags
var Version = "dev"

const (
	defaultRelayURL   = "ws://api.cinaseek.ai/ws/register"
	heartbeatInterval = 30 * time.Second
	reconnectDelay    = 5 * time.Second
)

// Command represents a command received from the relay server
type Command struct {
	ID     string          `json:"id"`
	Action string          `json:"action"`
	Params json.RawMessage `json:"params,omitempty"`
}

// Response is sent back to the relay server
type Response struct {
	ID     string      `json:"id"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func printBanner() {
	banner := "\n" +
		"   ██████╗██╗███╗   ██╗███████╗ ██████╗ █████╗ ██████╗ ███████╗\n" +
		"  ██╔════╝██║████╗  ██║██╔════╝██╔════╝██╔══██╗██╔══██╗██╔════╝\n" +
		"  ██║     ██║██╔██╗ ██║███████╗██║     ███████║██████╔╝█████╗  \n" +
		"  ██║     ██║██║╚██╗██║╚════██║██║     ██╔══██║██╔═══╝ ██╔══╝  \n" +
		"  ╚██████╗██║██║ ╚████║███████║╚██████╗██║  ██║██║     ███████╗\n" +
		"   ╚═════╝╚═╝╚═╝  ╚═══╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝     ╚══════╝\n" +
		"\n" +
		"  CinaSeek Client v%s — Cloud Dev Workspace Connector\n" +
		"  https://cinaseek.ai\n\n"
	fmt.Printf(banner, Version)
}

func main() {
	relayURL := flag.String("relay", defaultRelayURL, "WebSocket relay server URL")
	token := flag.String("token", "", "Authentication token")
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("CinaSeek Client v%s\n", Version)
		os.Exit(0)
	}

	printBanner()
	log.Printf("[INFO] CinaSeek Client v%s starting...", Version)
	log.Printf("[INFO] Relay: %s", *relayURL)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Connect to relay and run main loop
	go connectRelay(ctx, *relayURL, *token)

	// Wait for shutdown signal
	select {
	case sig := <-sigCh:
		log.Printf("[INFO] Received signal %v, shutting down...", sig)
		cancel()
	case <-ctx.Done():
	}

	log.Println("[INFO] CinaSeek Client stopped.")
}

func connectRelay(ctx context.Context, relayURL, token string) {
	header := http.Header{}
	if token != "" {
		header.Set("Authorization", "Bearer "+token)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		log.Printf("[INFO] Connecting to relay: %s", relayURL)

		// Use net/http to establish WebSocket upgrade
		req, err := http.NewRequestWithContext(ctx, "GET", relayURL, nil)
		if err != nil {
			log.Printf("[ERROR] Invalid relay URL: %v", err)
			time.Sleep(reconnectDelay)
			continue
		}
		for k, v := range header {
			req.Header[k] = v
		}
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Sec-WebSocket-Version", "13")
		req.Header.Set("Sec-WebSocket-Key", "cinaseek-client-key")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[ERROR] Relay connection failed: %v", err)
			log.Printf("[INFO] Retrying in %v...", reconnectDelay)
			time.Sleep(reconnectDelay)
			continue
		}
		resp.Body.Close()

		log.Println("[INFO] Connected to relay server (polling mode)")

		// Simple polling loop as placeholder for full WebSocket implementation
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				log.Println("[DEBUG] Heartbeat sent")
			}
		}
	}
}

func handleMessage(message []byte) {
	var cmd Command
	if err := json.Unmarshal(message, &cmd); err != nil {
		log.Printf("[WARN] Invalid message: %s", strings.TrimSpace(string(message)))
		return
	}

	log.Printf("[INFO] Received command: %s (id=%s)", cmd.Action, cmd.ID)

	switch cmd.Action {
	case "ping":
		// Pong handled implicitly
	case "status":
		log.Println("[INFO] Status requested — responding with client info")
	case "exec":
		var params struct {
			Command string `json:"command"`
		}
		if err := json.Unmarshal(cmd.Params, &params); err != nil {
			log.Printf("[WARN] Invalid exec params: %v", err)
			break
		}
		log.Printf("[INFO] Exec command forwarded: %s", params.Command)
	default:
		log.Printf("[WARN] Unknown action: %s", cmd.Action)
	}
}
