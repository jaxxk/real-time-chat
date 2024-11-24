# Real-Time Chat Application

## Project Goals
### Primary Skills
- WebSocket handling and communication.
- Go concurrency with goroutines and channels.
- State management for multiple clients.
- Data persistence with a database (e.g., PostgreSQL).
- Real-time messaging between connected users.

### Stretch Goals (Optional)
- Authentication (e.g., JWT).
- Offline message handling (e.g., saving undelivered messages).
- Scaling with Redis for pub/sub.

---

## Project Architecture
### High-Level Components
1. **WebSocket Server**:
   - Handles client connections.
   - Routes messages between users.
   - Manages connection state.

2. **Message Delivery Service**:
   - Receives messages from WebSocket clients.
   - Sends messages to the appropriate recipient.
   - Handles offline messages and persistence.

3. **Database**:
   - Stores chat history and user data.

4. **Frontend**:
   - A minimal web-based client (e.g., HTML + JavaScript) to test the WebSocket functionality.

---

## Features and Milestones

### Feature 1: User Connections
- **Objective**:
  - Set up a WebSocket server to accept client connections.
  - Each client sends a `username` during the connection handshake.
- **Steps**:
  - Use the `gorilla/websocket` library.
  - Maintain a map of active connections (e.g., `map[string]*websocket.Conn`).
- **Deliverables**:
  - Clients can connect to the server.
  - Server logs all connected users.

### Feature 2: Message Broadcasting
- **Objective**:
  - Allow clients to send messages that the server broadcasts to all other connected users.
- **Steps**:
  - Parse messages from WebSocket clients (e.g., JSON format).
  - Broadcast messages to all active WebSocket connections.
- **Deliverables**:
  - Messages sent by one client appear on all other clients.

### Feature 3: Private Messaging
- **Objective**:
  - Add support for direct, private messages between users.
- **Steps**:
  - Include a `recipient` field in the message payload.
  - Use the connection map to send messages only to the intended recipient.
- **Deliverables**:
  - Users can send direct messages to each other.

### Feature 4: Message Persistence
- **Objective**:
  - Store messages in a database for history and offline delivery.
- **Steps**:
  - Define a `messages` table with fields: `id`, `sender`, `recipient`, `content`, `timestamp`.
  - Write messages to the database as they are sent.
- **Deliverables**:
  - All messages are stored in the database.
  - A simple REST endpoint to retrieve message history.

### Feature 5: Handling Offline Users
- **Objective**:
  - Queue messages for offline users and deliver them when they reconnect.
- **Steps**:
  - Track user presence with an `online_users` map.
  - Flag messages in the database as `delivered` or `undelivered`.
  - Check for undelivered messages during user reconnection.
- **Deliverables**:
  - Messages sent to offline users are stored and delivered when they reconnect.

### Stretch Feature: Authentication
- **Objective**:
  - Authenticate users using JWT or session tokens.
- **Steps**:
  - Validate tokens during the WebSocket handshake.
  - Reject unauthenticated connections.
- **Deliverables**:
  - Only authenticated users can connect and send messages.

---

## Tech Stack
- **Backend**: 
  - Golang with `gorilla/websocket`.
- **Database**:
  - PostgreSQL (or SQLite for simplicity during development).
- **Frontend**:
  - Minimal client with HTML, JavaScript, and WebSocket APIs.
- **Testing**:
  - Unit testing with `testing` package.
  - Load testing with tools like `k6` or `wrk`.

---

## Milestone Plan

| **Milestone**          | **Time Estimate** | **Deliverable**                                                   |
|-------------------------|-------------------|--------------------------------------------------------------------|
| User Connections        | 1-2 days         | WebSocket server accepting and tracking client connections.        |
| Message Broadcasting    | 2 days           | Real-time broadcast messages to all connected clients.             |
| Private Messaging       | 2-3 days         | Direct messaging between specific users.                           |
| Message Persistence     | 3-4 days         | Store messages in a database; retrieve chat history.               |
| Offline Handling        | 3 days           | Queue and deliver messages for offline users.                      |
| Authentication (Stretch)| 2-3 days         | Authenticate users using JWT.                                      |

