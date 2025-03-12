package wsclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/coder/websocket"
	_ "github.com/gorilla/websocket"
	websocket "github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/api"
	"github.com/pubgo/funk/log"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

// ClientConfig contains configuration options for the websocket client
type ClientConfig struct {
	// PingInterval defines how often to send ping messages
	PingInterval time.Duration

	// PongTimeout defines how long to wait for a pong response
	PongTimeout time.Duration

	// WriteTimeout defines how long to wait for a write to complete
	WriteTimeout time.Duration

	// ReconnectInterval defines how long to wait between reconnection attempts
	ReconnectInterval time.Duration

	// MaxReconnectAttempts defines maximum number of reconnect attempts (0 for unlimited)
	MaxReconnectAttempts int

	// HealthCheckInterval defines how often to check connection health
	HealthCheckInterval time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		PingInterval:         10 * time.Second,
		PongTimeout:          15 * time.Second,
		WriteTimeout:         10 * time.Second,
		ReconnectInterval:    5 * time.Second,
		MaxReconnectAttempts: 0, // unlimited
		HealthCheckInterval:  30 * time.Second,
	}
}

// ClientOption is a functional option type for configuring the client
type ClientOption func(*ClientConfig)

// WithPingInterval configures the ping interval
func WithPingInterval(d time.Duration) ClientOption {
	return func(c *ClientConfig) {
		c.PingInterval = d
	}
}

// WithPongTimeout configures the pong timeout
func WithPongTimeout(d time.Duration) ClientOption {
	return func(c *ClientConfig) {
		c.PongTimeout = d
	}
}

// WithReconnectInterval configures the reconnect interval
func WithReconnectInterval(d time.Duration) ClientOption {
	return func(c *ClientConfig) {
		c.ReconnectInterval = d
	}
}

// WithMaxReconnectAttempts configures the maximum reconnect attempts
func WithMaxReconnectAttempts(n int) ClientOption {
	return func(c *ClientConfig) {
		c.MaxReconnectAttempts = n
	}
}

// ConnectionState represents the current state of the websocket connection
type ConnectionState int32

const (
	// StateDisconnected indicates the client is disconnected
	StateDisconnected ConnectionState = iota
	// StateConnecting indicates the client is attempting to connect
	StateConnecting
	// StateConnected indicates the client is connected
	StateConnected
)

// wsClient implements the Client interface
type wsClient struct {
	// Configuration
	config      *ClientConfig
	wsURL       string
	cloudClient api.Client
	actions     map[string]func(ctx context.Context, payload Message) Message

	// Connection state
	conn      *websocket.Conn
	connMutex sync.Mutex
	state     atomic.Int32

	// Communication channels
	sendChan    chan Message
	reconnectCh chan struct{}

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc

	// Monitoring
	lastPong               atomic.Value
	reconnectAttempts      atomic.Int32
	totalReconnectAttempts atomic.Int32

	// Logging
	log log.Logger
}

func (c *wsClient) RegisterAction(act string, handler func(ctx context.Context, payload Message) Message) error {
	if c.actions[act] != nil {
		return fmt.Errorf("action %s already registered", act)
	}
	c.actions[act] = handler
	return nil
}

// FormatWSURL formats a websocket URL with session information
func FormatWSURL(wsServer string, boxID string) string {
	return fmt.Sprintf("%s?session_id=%s:%s", wsServer, boxID, lo.RandomString(32, lo.LowerCaseLettersCharset))
}

// newClient creates a new websocket client
func newClient(ctx context.Context, wsURL string, cloudClient api.Client, options ...ClientOption) (Client, error) {
	if wsURL == "" {
		return nil, errors.New("websocket URL cannot be empty")
	}

	if cloudClient == nil {
		return nil, errors.New("cloud client cannot be nil")
	}

	// Parse URL to validate it
	_, err := url.Parse(wsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid websocket URL: %w", err)
	}

	// Apply configuration options
	config := DefaultConfig()
	for _, option := range options {
		option(config)
	}

	logs := log.GetLogger("websocket")
	logs.Debug().Msg(.Msgf("config: %+v", config)

	// Create client context with cancellation
	clientCtx, cancel := context.WithCancel(ctx)

	client := &wsClient{
		config:      config,
		wsURL:       wsURL,
		cloudClient: cloudClient,
		actions:     make(map[string]func(ctx context.Context, payload Message) Message),
		sendChan:    make(chan Message, 100),
		reconnectCh: make(chan struct{}, 1),
		ctx:         clientCtx,
		cancel:      cancel,
		log:         logs,
	}

	// Initialize state
	client.state.Store(int32(StateDisconnected))
	client.lastPong.Store(time.Now())
	client.reconnectAttempts.Store(0)
	client.totalReconnectAttempts.Store(0)

	return client, nil
}

func (c *wsClient) Start() {
	// Start the client
	c.log.Info().Msg("Starting client")

	// Initial connection
	c.connect()

	// Start management goroutines
	go c.reconnectLoop()

	// Start connection health checker
	go c.healthCheckLoop()
}

// start begins the client operation
func (c *wsClient) start() {
	c.log.Info().Msg("Starting client")

	// Initial connection
	c.connect()

	// Start management goroutines
	go c.reconnectLoop()
}

// healthCheckLoop periodically checks connection status and initiates reconnection if needed
func (c *wsClient) healthCheckLoop() {
	ticker := time.NewTicker(c.config.HealthCheckInterval)
	defer ticker.Stop()

	c.log.Info().Msg(.Msg("Starting connection health check loop")

	for {
		select {
		case <-c.ctx.Done():
			c.log.Info().Msg(.Msg("Client shutdown, stopping health check loop")
			return

		case <-ticker.C:
			if c.getState() != StateConnected && c.getState() != StateConnecting {
				c.log.Warn("Health check detected disconnected state, triggering reconnect")
				c.triggerReconnect()
			} else {
				c.log.Debug().Msg("Health check: connection is active")
			}
		}
	}
}

// reconnectLoop manages reconnection attempts
func (c *wsClient) reconnectLoop() {
	c.log.Info().Msg("Starting reconnect loop")

	for {
		select {
		case <-c.ctx.Done():
			c.log.Info().Msg("Client shutdown, stopping reconnect loop")
			return

		case <-c.reconnectCh:
			c.log.Info().Msg("Reconnect signal received")

			// If already connecting, don't start another attempt
			if !c.compareAndSetState(StateDisconnected, StateConnecting) {
				c.log.Info().Msg("Already connecting, skipping")
				continue
			}

			currentAttempts := c.reconnectAttempts.Add(1)
			totalAttempts := c.totalReconnectAttempts.Load()
			c.log.Infof("Attempting to reconnect to WebSocket server (attempt: %d, total: %d)",
				currentAttempts, totalAttempts)

			// Check if maximum reconnect attempts reached
			if c.config.MaxReconnectAttempts > 0 && int(currentAttempts) > c.config.MaxReconnectAttempts {
				c.log.Warn().Msgf("Maximum reconnection attempts reached (%d), but connection is essential - continuing to retry", currentAttempts)
				// Reset counter to avoid integer overflow during long-running operations
				c.reconnectAttempts.Store(1)
				// Don't cancel, keep trying
			}

			// Wait before reconnecting
			reconnectDelay := c.config.ReconnectInterval
			c.log.Info().Msg("Reconnecting to WebSocket server after delay")

			select {
			case <-time.After(reconnectDelay):
				c.log.Info().Msg("Reconnecting to WebSocket server")
				c.connect()
			case <-c.ctx.Done():
				c.log.Info().Msg("Reconnect loop context done, exiting")
				return
			}
		}
	}
}

// connect establishes a connection to the WebSocket server
func (c *wsClient) connect() {
	reconnectCount := c.reconnectAttempts.Load()
	totalReconnects := c.totalReconnectAttempts.Load()
	c.log.Infof("Connecting to WebSocket server (attempt: %d, total: %d)",
		reconnectCount, totalReconnects)
	c.setState(StateConnecting)

	// Close existing connection if present
	c.closeConnection()

	// Prepare headers with authorization
	headers := http.Header{}
	token := c.cloudClient.GetJWTToken(c.ctx)
	if token == "" {
		c.log.Error("Failed to get authentication token")
		c.triggerReconnect()
		return
	}
	headers.Add("Authorization", token)

	// Attempt to connect
	conn, resp, err := websocket.DefaultDialer.Dial(c.wsURL, headers)
	if resp != nil {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		c.log.Debug().Msgf("Connection response: %s", string(respBody))
	}

	if err != nil {
		c.log.Error().Msgf("Failed to connect to %s: %v (attempt: %d, total: %d)",
			c.wsURL, err, reconnectCount, totalReconnects)
		c.triggerReconnect()
		return
	}

	// Configure the connection
	c.connMutex.Lock()
	c.conn = conn
	c.connMutex.Unlock()

	// Reset reconnect attempts on successful connection
	c.reconnectAttempts.Store(0)

	// Set up handlers
	conn.SetPingHandler(func(data string) error {
		c.log.Debug().Msg("Received PING message")
		deadline := time.Now().Add(c.config.WriteTimeout)
		return c.writeControlMessage(websocket.PongMessage, []byte(data), deadline)
	})

	conn.SetPongHandler(func(data string) error {
		c.log.Debug().Msg("Received PONG message")
		c.lastPong.Store(time.Now())
		return nil
	})

	conn.SetCloseHandler(func(code int, text string) error {
		c.log.Infof("Connection closed with code %d: %s", code, text)
		c.triggerReconnect()
		return nil
	})

	// Mark as connected
	c.setState(StateConnected)
	c.lastPong.Store(time.Now())
	c.log.Infof("Successfully connected to WebSocket server (local: %s, remote: %s, total reconnects: %d)",
		conn.LocalAddr(), conn.RemoteAddr(), totalReconnects)

	// Start communication goroutines
	go c.readLoop()
	go c.writeLoop()
}

// readLoop handles incoming messages
func (c *wsClient) readLoop() {
	c.log.Info().Msg("Starting read loop")

	defer func() {
		if r := recover(); r != nil {
			c.log.Error().Msgf("Panic in read loop: %v", r)
			c.triggerReconnect()
		}
	}()

	for {
		if c.getState() != StateConnected {
			c.log.Info().Msg("Not connected, skipping read loop")
			return
		}

		// Get connection (with mutex protection)
		conn := c.getConnection()
		if conn == nil {
			c.log.Info().Msg("Connection is nil, skipping read loop")
			return
		}

		// Read message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			c.log.Error().Msgf("Error reading message: %v", err)
			c.triggerReconnect()
			return
		}

		// Process different message types
		switch messageType {
		case websocket.TextMessage:
			if c.log.Level == logrus.DebugLevel {
				c.log.Debug().Msgf("Received message: %s", string(message))
			}

			// Process message in a separate goroutine
			go func(msg []byte) {
				// Try to send with retries
				reqMsg, err := ToResultMessage(msg)
				if err != nil || reqMsg.GetAction() == "" {
					c.log.Error().Msgf("failed to parse websocket request message: %s, %v", msg, err)
				} else {
					c.sendWithRetry(c.invoke(c.ctx, reqMsg))
				}
			}(message)

		case websocket.BinaryMessage:
			c.log.Error().Msgf("Received binary message of length %d", len(message))

		case websocket.CloseMessage:
			c.log.Info().Msg("Received close message")
			c.triggerReconnect()
			return
		}
	}
}

func unknownAction(ctx context.Context, msg Message) Message {
	err := fmt.Errorf("do not know how to handle msg of action %v: %s", msg.GetAction(), string(msg.Marshal()))
	//return msg.ReplyMessage(codec.NewError(codec.Err_Not_Support_Act, err.Error()))
	return msg.ReplyMessage(codec.NewError(codec.Err_Unknown, err.Error()))
}

func (c *wsClient) invoke(ctx context.Context, payload Message) Message {
	defer func() {
		if r := recover(); r != nil {
			c.log.WithContext(ctx).WithFields(logrus.Fields{
				"action": payload.GetAction(),
				"id":     payload.GetID(),
			}).Errorf("Panic in invoke: %v, stack: %s", r, debug.Stack())
		}
	}()

	actionFunc, ok := c.actions[payload.GetAction()]
	if !ok {
		return unknownAction(ctx, payload)
	}

	return actionFunc(ctx, payload)
}

// sendWithRetry attempts to send a message with unlimited retries
func (c *wsClient) sendWithRetry(payload Message) {
	if payload == nil {
		return
	}

	const initialBackoff = 100 * time.Millisecond
	const maxBackoff = 5 * time.Second

	backoff := initialBackoff
	attempt := 0

	for {
		// Create a timeout context for each send attempt
		ctx, cancel := context.WithTimeout(c.ctx, 2*time.Second)

		err := c.Send(ctx, payload)
		cancel() // Cancel the context to avoid leaks

		if err == nil {
			// Send successful
			if attempt > 0 {
				c.log.Debug().Msgf("Successfully sent reply after %d retries", attempt)
			}
			return
		}

		// Log the failure and prepare to retry
		c.log.Debug().Msgf("Failed to send reply (attempt %d): %v, retrying in %v", attempt+1, err, backoff)

		// Wait before retrying, but respect client context
		select {
		case <-c.ctx.Done():
			c.log.Info().Msg("Client context cancelled during send retry")
			return
		case <-time.After(backoff):
			// Continue with retry
		}

		// Increase backoff for next attempt (with a cap)
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}

		attempt++
	}
}

// writeLoop handles outgoing messages and ping maintenance
func (c *wsClient) writeLoop() {
	c.log.Info().Msg("Starting write loop")

	pingTicker := time.NewTicker(c.config.PingInterval)
	defer pingTicker.Stop()

	defer func() {
		if r := recover(); r != nil {
			c.log.Error().Msgf("Panic in write loop: %v", r)
			c.triggerReconnect()
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			c.log.Info().Msg("Write loop context done, exiting")
			return

		case message := <-c.sendChan:
			if c.getState() != StateConnected {
				c.log.Warn("Not connected, using sendWithRetry for reliable delivery")
				// Use sendWithRetry instead of buffering
				go c.sendWithRetry(message)
				continue
			}

			err := c.writeMessage(websocket.TextMessage, message.Marshal())
			if err != nil {
				c.log.Error().Msgf("Error sending message: %v", err)
				// Use sendWithRetry to handle the failed message
				go c.sendWithRetry(message)
				c.triggerReconnect()
				return
			}

		case <-pingTicker.C:
			// Check if we've received a pong recently
			lastPong, ok := c.lastPong.Load().(time.Time)
			if !ok || time.Since(lastPong) > c.config.PongTimeout {
				c.log.Warn().Msgf("Pong timeout detected, reconnecting, duration=%s", time.Since(lastPong))
				c.triggerReconnect()
				return
			}

			// Send ping
			c.log.Debug().Msg("Sending PING message")
			deadline := time.Now().Add(c.config.WriteTimeout)
			err := c.writeControlMessage(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				c.log.Error().Msgf("Error sending ping: %v", err)
				c.triggerReconnect()
				return
			}
		}
	}
}

// Send implements the Client interface
func (c *wsClient) Send(ctx context.Context, payload Message) error {
	// Check connection state first
	if c.getState() != StateConnected {
		return errors.New("not connected to server")
	}

	// Non-blocking send - returns error if channel is full
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.sendChan <- payload:
		return nil
	default:
		return errors.New("send buffer full, message dropped")
	}
}

// Close implements the Client interface
func (c *wsClient) Close() {
	c.log.Info().Msg("Closing WebSocket client")
	c.cancel()
	c.closeConnection()
}

// Helper methods

// getConnection safely gets the current connection
func (c *wsClient) getConnection() *websocket.Conn {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	return c.conn
}

// closeConnection safely closes the current connection
func (c *wsClient) closeConnection() {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	conn := c.conn
	if conn == nil {
		c.log.Debug().Msg("Connection is nil, skipping close")
		return
	}

	c.log.Debug().Msgf("Closing connection: local:%s remote:%s", c.conn.LocalAddr(), c.conn.RemoteAddr())
	// Send close message
	deadline := time.Now().Add(c.config.WriteTimeout)
	_ = conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		deadline,
	)

	c.log.Debug().Msg("Closed connection")
	// Close connection
	_ = conn.Close()
	c.conn = nil
	c.log.Debug().Msg("Connection closed")
}

// writeMessage safely writes a message to the connection
func (c *wsClient) writeMessage(messageType int, data []byte) error {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	conn := c.conn
	if conn == nil {
		c.log.Info().Msg("Connection is nil, skipping write")
		return errors.New("connection is nil")
	}

	c.log.Debug().Msgf("Writing message: %s", string(data))
	return conn.WriteMessage(messageType, data)
}

// writeControlMessage safely writes a control message to the connection
func (c *wsClient) writeControlMessage(messageType int, data []byte, deadline time.Time) error {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	conn := c.conn
	if conn == nil {
		c.log.Info().Msg("Connection is nil, skipping write")
		return errors.New("connection is nil")
	}

	c.log.Debug().Msgf("Writing control message: %s", string(data))
	return conn.WriteControl(messageType, data, deadline)
}

// getState gets the current connection state
func (c *wsClient) getState() ConnectionState {
	return ConnectionState(c.state.Load())
}

// setState sets the current connection state
func (c *wsClient) setState(state ConnectionState) {
	c.state.Store(int32(state))
}

// compareAndSetState atomically updates the state if the current value matches the expected value
func (c *wsClient) compareAndSetState(expected, new ConnectionState) bool {
	c.log.Debug().Msgf("Comparing and setting state: %d -> %d", expected, new)
	return c.state.CompareAndSwap(int32(expected), int32(new))
}

// triggerReconnect safely triggers a reconnection attempt
func (c *wsClient) triggerReconnect() {
	c.totalReconnectAttempts.Add(1)
	c.log.Debug().Msg("Triggering reconnect")
	if c.getState() == StateDisconnected {
		c.log.Debug().Msg("Already disconnected, skipping reconnect")
		return
	}

	c.log.Debug().Msg("Setting state to disconnected")
	c.setState(StateDisconnected)

	c.log.Debug().Msg("Sending reconnect signal")
	select {
	case c.reconnectCh <- struct{}{}:
		c.log.Info().Msg("Reconnect signal sent")
	default:
		c.log.Info().Msg("Channel already has a pending signal")
	}
}
