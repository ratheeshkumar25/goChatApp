package chat

import (
	"time"
)

// Client represents a chat client
type Client struct {
	ID      string
	MsgChan chan Message
}

// NewClient creates a new client with the specified ID
func NewClient(id string) *Client {
	return &Client{
		ID:      id,
		MsgChan: make(chan Message, 100), // Buffer for 100 messages
	}
}

// GetMessages retrieves messages for the client with timeout
func (c *Client) GetMessages(timeout time.Duration) ([]Message, bool) {
	var messages []Message

	// Set timeout
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// Try to get at least one message
	select {
	case msg, ok := <-c.MsgChan:
		if !ok {
			return nil, false // Channel closed
		}
		messages = append(messages, msg)
	case <-timer.C:
		return messages, true // Timeout with no messages
	}

	// Collect any remaining messages without blocking
drainLoop:
	for {
		select {
		case msg, ok := <-c.MsgChan:
			if !ok {
				break drainLoop // Channel closed
			}
			messages = append(messages, msg)
		default:
			break drainLoop // No more messages available
		}
	}

	return messages, true
}
