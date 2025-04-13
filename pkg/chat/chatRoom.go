package chat

import (
	"fmt"
	"sync"
	"time"
)

// Message represents a chat message
type Message struct {
	SenderID string
	Content  string
	Time     time.Time
}

// ChatRoom manages the chat clients and messages
type ChatRoom struct {
	clients       map[string]*Client
	joinChan      chan *Client
	leaveChan     chan string
	leaveAckChan  chan bool
	broadcastChan chan Message
	mutex         sync.RWMutex
}

// NewChatRoom creates and starts a new chat room
func NewChatRoom() *ChatRoom {
	chat := &ChatRoom{
		clients:       make(map[string]*Client),
		joinChan:      make(chan *Client),
		leaveChan:     make(chan string),
		leaveAckChan:  make(chan bool),
		broadcastChan: make(chan Message),
	}
	go chat.Run()
	return chat
}

// Run handles chat room operations in a separate goroutine
func (c *ChatRoom) Run() {
	for {
		select {
		case client := <-c.joinChan:
			c.mutex.Lock()
			c.clients[client.ID] = client
			c.mutex.Unlock()
			fmt.Printf("Client joined: %s\n", client.ID)

			// Notify everyone about the new client
			joinMsg := Message{
				SenderID: "System",
				Content:  fmt.Sprintf("User %s has joined the chat", client.ID),
				Time:     time.Now(),
			}
			c.broadcastMessage(joinMsg)

		case id := <-c.leaveChan:
			c.mutex.Lock()
			client, exists := c.clients[id]

			// Store existence for acknowledgment
			clientExists := exists

			if exists {
				delete(c.clients, id)
				close(client.MsgChan) // Close the client's message channel
				fmt.Printf("Client left: %s\n", id)
			}
			c.mutex.Unlock()

			// Send acknowledgment
			c.leaveAckChan <- clientExists

			// If client existed, broadcast after releasing the lock
			if clientExists {
				leaveMsg := Message{
					SenderID: "System",
					Content:  fmt.Sprintf("User %s has left the chat", id),
					Time:     time.Now(),
				}
				c.broadcastMessage(leaveMsg)
			}

		case msg := <-c.broadcastChan:
			c.broadcastMessage(msg)
		}
	}
}

// broadcastMessage sends a message to all connected clients
func (c *ChatRoom) broadcastMessage(msg Message) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, client := range c.clients {
		// Use a non-blocking send to prevent one slow client from blocking others
		select {
		case client.MsgChan <- msg:
			// Message sent successfully
		default:
			// Client's channel buffer is full, log and continue
			fmt.Printf("Failed to send message to client %s: channel full\n", client.ID)
		}
	}
}

// Join adds a client to the chat room
func (c *ChatRoom) Join(id string) (*Client, error) {
	// Check if client already exists
	c.mutex.RLock()
	_, exists := c.clients[id]
	c.mutex.RUnlock()

	if exists {
		return nil, fmt.Errorf("client with ID %s already exists", id)
	}

	// Create a new client
	client := NewClient(id)

	// Add client to chat room
	c.joinChan <- client

	return client, nil
}

// Leave removes a client from the chat room
func (c *ChatRoom) Leave(id string) bool {
	c.leaveChan <- id
	return <-c.leaveAckChan
}

// Broadcast sends a message to all clients
func (c *ChatRoom) Broadcast(senderID, content string) {
	msg := Message{
		SenderID: senderID,
		Content:  content,
		Time:     time.Now(),
	}
	c.broadcastChan <- msg
}

// GetClient returns a client by ID
func (c *ChatRoom) GetClient(id string) (*Client, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	client, exists := c.clients[id]
	return client, exists
}
