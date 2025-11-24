# Zinx Connection Property Check Example

This example demonstrates how to implement client connection property validation in Zinx framework.

## Overview

The server checks if the client has set a valid `name` property with value `test` when the connection is established.
If the property is missing or has an invalid value, the server disconnects the client immediately.

## Files

- `server/server.go`: Server implementation with connection property validation
- `client/client.go`: Client implementation with valid connection property
- `client/client_invalid.go`: Client implementation with invalid connection property

## How to Run

### 1. Start the Server

```bash
cd server
go run server.go
```

The server will start listening on `127.0.0.1:8999`.

### 2. Run the Valid Client

In a new terminal:

```bash
cd client
go run client.go
```

This client sets the `name` property to `test` (which is valid), so the connection will be accepted by the server.

### 3. Run the Invalid Client

In another terminal:

```bash
cd client
go run client_invalid.go
```

This client sets the `name` property to `invalid_name` (which is invalid), so the server will disconnect it immediately.

## Expected Output

### Server Output

When valid client connects:
```
[INFO] 2023/10/05 14:30:00 Server starting
[INFO] 2023/10/05 14:30:05 Server connection started
[INFO] 2023/10/05 14:30:05 Client connected with valid name: test
[DEBUG] 2023/10/05 14:30:05 Call PingRouter Handle
[DEBUG] 2023/10/05 14:30:05 recv from client : msgId=1, data=Hello server, I'm client with valid name
```

When invalid client connects:
```
[INFO] 2023/10/05 14:30:10 Server connection started
[ERROR] 2023/10/05 14:30:10 Invalid client name: invalid_name
[INFO] 2023/10/05 14:30:10 Server connection lost
```

### Valid Client Output

```
[INFO] 2023/10/05 14:30:05 Client starting
[INFO] 2023/10/05 14:30:05 Client connected to server
[DEBUG] 2023/10/05 14:30:05 Call PingRouter Handle
[DEBUG] 2023/10/05 14:30:05 recv from server : msgId=2, data=pong-server
```

### Invalid Client Output

```
[INFO] 2023/10/05 14:30:10 Client (with invalid name) starting
[INFO] 2023/10/05 14:30:10 Client connected to server
[INFO] 2023/10/05 14:30:10 Client disconnected from server (expected, since we used invalid name)
```

## Implementation Details

### Server Side

1. Set a connection start callback function using `s.SetOnConnStart(DoConnectionBegin)`
2. In the callback function:
   - Get the `name` property from the connection using `conn.GetProperty("name")`
   - Check if the property exists and has the value `test`
   - If validation fails, disconnect the client using `conn.Stop()`

### Client Side

1. Set a connection start callback function using `client.SetOnConnStart(DoClientConnectedBegin)`
2. In the callback function:
   - Set the `name` property using `conn.SetProperty("name", "test")` (for valid client)
   - Set an invalid name for the invalid client

This example shows how to use Zinx's built-in connection property mechanism to implement client validation without modifying the framework's core code.