# Go Chat Application

A dynamic chat application built in Go using goroutines and channels for concurrency with a RESTful API interface.

## Overview

This chat application allows clients to dynamically join, send messages, and leave a chat room using HTTP APIs. The core of the application leverages goroutines and channels for handling concurrent client connections.

## Project Directory

ratheeshkumar25-gochatapp/
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   └── main.go
└── pkg/
    ├── api/
    │   └── handlers.go
    ├── chat/
    │   ├── chatRoom.go
    │   └── client.go
    ├── di/
    │   └── di.go
    └── server/
        └── server.go

## Features

- **Chat Room Management**: A central chat room where multiple clients can join, leave, and send messages.
- **Concurrent Client Handling**: Uses goroutines and channels to manage multiple clients concurrently.
- **Thread-safe Operations**: All operations on shared data structures are thread-safe.
- **RESTful API**: Simple HTTP endpoints for chat interactions.
- **Message Broadcasting**: All messages are broadcast to all connected clients.
- **Timeouts**: Implements timeouts for long-polling message retrieval.


## API Endpoints

### Join Chat
**GET** `/join?id=<client_id>`

- Allows a client to join the chat room.
- Returns client information on success.

### Send Message
**GET** `/send?id=<client_id>&message=<message>`

- Allows a client to send a message to the chat room.
- The message is broadcast to all connected clients.

### Get Messages
**GET** `/messages?id=<client_id>`

- Allows a client to receive broadcast messages from the chat room.
- Implements a timeout to prevent indefinite blocking.

### Leave Chat
**GET** `/leave?id=<client_id>`

- Allows a client to leave the chat room.
- Cleans up resources associated with the client.

## Installation

### Clone the repository
```bash
git clone https://github.com/ratheeshkumar25/goChatApp.git
cd chat-app
```

### Install dependencies
```bash
go mod tidy
```

### Build the application
```bash
go build -o chat-cmd ./cmd
```

### Run the server
```bash
./chat-cmd
```

## Testing with Postman

1. Install Postman.
2. Create a new request for each endpoint:
    - **Join**: `GET http://localhost:3000/join?id=user1`
    - **Send**: `GET http://localhost:3000/send?id=user1&message=Hello%20World`
    - **Messages**: `GET http://localhost:3000/messages?id=user1`
    - **Leave**: `GET http://localhost:3000/leave?id=user1`

## Implementation Details

### Concurrency Model

The application uses Go's concurrency primitives:

- **Goroutines**: Each client is managed by a dedicated goroutine.
- **Channels**: Used for communication between components.
- **Mutex**: Used for thread-safe access to shared data structures.

### Message Handling

Messages are broadcast to all clients using a non-blocking approach to ensure that slow clients don't block the entire system.

### Error Handling

The API provides appropriate error responses:

- **400 Bad Request**: For invalid input parameters.
- **404 Not Found**: When a client ID doesn't exist.
- **410 Gone**: When a client has disconnected.

## Future Improvements

- Add WebSocket support for real-time communication.
- Implement message persistence.
- Add user authentication.
- Create a web frontend.
- Support direct messaging between clients.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
