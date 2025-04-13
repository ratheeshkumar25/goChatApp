package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ratheeshkumar25/chatApp/pkg/chat"
)

// Handler manages HTTP request handling
type Handler struct {
	chatRoom *chat.ChatRoom
}

// NewHandler creates a new handler with the provided chat room
func NewHandler(chatRoom *chat.ChatRoom) *Handler {
	return &Handler{
		chatRoom: chatRoom,
	}
}

// JoinHandler handles client join requests
func (h *Handler) JoinHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}

	_, err := h.chatRoom.Join(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// log.Print(client)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Joined chat successfully",
		"id":      id,
	})
}

// LeaveHandler handles client leave requests
func (h *Handler) LeaveHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}

	success := h.chatRoom.Leave(id)
	if !success {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Left chat successfully",
	})
}

// SendHandler handles message sending
func (h *Handler) SendHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	message := r.URL.Query().Get("message")

	if id == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}

	if message == "" {
		http.Error(w, "Missing message", http.StatusBadRequest)
		return
	}

	// Check if client exists
	_, exists := h.chatRoom.GetClient(id)
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	// Broadcast the message
	h.chatRoom.Broadcast(id, message)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Message sent",
	})
}

// MessagesHandler handles retrieving messages
func (h *Handler) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}

	// Check if client exists
	client, exists := h.chatRoom.GetClient(id)
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	// Get messages with timeout (10 seconds)
	messages, ok := client.GetMessages(10 * time.Second)
	if !ok {
		http.Error(w, "Client disconnected", http.StatusGone)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"messages": messages,
	})
}
