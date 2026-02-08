# Axis API Documentation

This document provides a comprehensive overview of the RESTful API endpoints for the Axis application, designed to mimic core functionalities of a Slack-like messaging platform. All API endpoints are prefixed with `/api`.

## Base URL

`/api`

## Authentication

Authentication mechanisms (e.g., JWT, session tokens) are expected for most protected routes, typically handled via an `Authorization` header. Specific details for authentication are not covered in this document but are assumed to be implemented.

## Error Responses

In case of an error, the API will typically return a JSON object with an `error` key and a descriptive message, along with an appropriate HTTP status code.

Example error response:

```json
{
  "error": "User not found"
}
```

Common error status codes:
*   `400 Bad Request`: Invalid input, missing required parameters.
*   `401 Unauthorized`: Authentication required or failed.
*   `403 Forbidden`: User does not have permission to access the resource.
*   `404 Not Found`: The requested resource does not exist.
*   `500 Internal Server Error`: An unexpected server-side error occurred.

## API Endpoints

---

### User Management

**`POST /api/register`**

*   **Description:** Registers a new user.
*   **Request Body Example:**
    ```json
    {
      "name": "John Doe",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "password": "securepassword",
      "timezone": "America/New_York",
      "locale": "en-US"
    }
    ```
*   **Response Body Example (201 Created):**
    ```json
    {
      "id": 1,
      "name": "John Doe",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "status": "active",
      "timezone": "America/New_York",
      "locale": "en-US",
      "is_verified": false,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "last_login_at": null
    }
    ```

**`POST /api/login`**

*   **Description:** Authenticates a user and returns user details.
*   **Request Body Example:**
    ```json
    {
      "email": "john.doe@example.com",
      "password": "securepassword"
    }
    ```
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "name": "John Doe",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "status": "active",
      "timezone": "America/New_York",
      "locale": "en-US",
      "is_verified": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "last_login_at": "2024-01-05T10:30:00Z"
      // Potentially includes a JWT or session token in a header or response body
    }
    ```

**`GET /api/users/:userID`**

*   **Description:** Retrieves a user by their ID.
*   **Path Parameters:**
    *   `userID`: The ID of the user.
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "name": "John Doe",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "status": "active",
      "timezone": "America/New_York",
      "locale": "en-US",
      "is_verified": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "last_login_at": "2024-01-05T10:30:00Z"
    }
    ```

**`GET /api/users/by-email?email={email}`**

*   **Description:** Retrieves a user by their email address.
*   **Query Parameters:**
    *   `email`: The email address of the user.
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "name": "John Doe",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "status": "active",
      "timezone": "America/New_York",
      "locale": "en-US",
      "is_verified": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "last_login_at": "2024-01-05T10:30:00Z"
    }
    ```

**`GET /api/users/by-username?username={username}`**

*   **Description:** Retrieves a user by their username.
*   **Query Parameters:**
    *   `username`: The username of the user.
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "name": "John Doe",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "status": "active",
      "timezone": "America/New_York",
      "locale": "en-US",
      "is_verified": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "last_login_at": "2024-01-05T10:30:00Z"
    }
    ```

**`PUT /api/users/:userID`**

*   **Description:** Updates an existing user's information.
*   **Path Parameters:**
    *   `userID`: The ID of the user to update.
*   **Request Body Example:**
    ```json
    {
      "name": "Jonathan Doe",
      "timezone": "Europe/London"
    }
    ```
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "name": "Jonathan Doe",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "status": "active",
      "timezone": "Europe/London",
      "locale": "en-US",
      "is_verified": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-06T15:00:00Z",
      "last_login_at": "2024-01-05T10:30:00Z"
    }
    ```

**`DELETE /api/users/:userID`**

*   **Description:** Deletes a user by their ID.
*   **Path Parameters:**
    *   `id`: The ID of the user to delete.
*   **Response:** `204 No Content` on successful deletion.

---

### Workspace Management

**`POST /api/workspaces`**

*   **Description:** Creates a new workspace.
*   **Request Body Example:**
    ```json
    {
      "name": "My Team Workspace",
      "description": "A place for my team to collaborate",
      "creator_id": 1
    }
    ```
*   **Response Body Example (201 Created):**
    ```json
    {
      "id": 1,
      "name": "My Team Workspace",
      "description": "A place for my team to collaborate",
      "creator_id": 1,
      "created_at": "2024-01-07T08:00:00Z",
      "updated_at": "2024-01-07T08:00:00Z"
    }
    ```

**`GET /api/workspaces/:workspaceID`**

*   **Description:** Retrieves a workspace by its ID.
*   **Path Parameters:**
    *   `workspaceID`: The ID of the workspace.
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "name": "My Team Workspace",
      "description": "A place for my team to collaborate",
      "creator_id": 1,
      "created_at": "2024-01-07T08:00:00Z",
      "updated_at": "2024-01-07T08:00:00Z"
    }
    ```

**`PUT /api/workspaces/:workspaceID`**

*   **Description:** Updates an existing workspace's information.
*   **Path Parameters:**
    *   `workspaceID`: The ID of the workspace to update.
*   **Request Body Example:**
    ```json
    {
      "name": "My Awesome Team Workspace",
      "description": "An even better place for my team"
    }
    ```
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "name": "My Awesome Team Workspace",
      "description": "An even better place for my team",
      "creator_id": 1,
      "created_at": "2024-01-07T08:00:00Z",
      "updated_at": "2024-01-07T09:30:00Z"
    }
    ```

**`DELETE /api/workspaces/:workspaceID`**

*   **Description:** Deletes a workspace by its ID.
*   **Path Parameters:**
    *   `id`: The ID of the workspace to delete.
*   **Response:** `204 No Content` on successful deletion.

**`GET /api/users/:userID/workspaces`**

*   **Description:** Retrieves all workspaces a specific user is a member of.
*   **Path Parameters:**
    *   `userID`: The ID of the user.
*   **Response Body Example (200 OK):**
    ```json
    [
      {
        "id": 1,
        "name": "My Team Workspace",
        "description": "A place for my team to collaborate",
        "creator_id": 1,
        "created_at": "2024-01-07T08:00:00Z",
        "updated_at": "2024-01-07T08:00:00Z"
      },
      {
        "id": 2,
        "name": "Project Alpha",
        "description": "Workspace for Project Alpha",
        "creator_id": 2,
        "created_at": "2024-01-08T10:00:00Z",
        "updated_at": "2024-01-08T10:00:00Z"
      }
    ]
    ```

---

### Workspace Member Management

**`POST /api/workspaces/:workspaceID/members`**

*   **Description:** Adds a user as a member to a specific workspace.
*   **Path Parameters:**
    *   `workspaceID`: The ID of the workspace.
*   **Request Body Example:**
    ```json
    {
      "user_id": 2,
      "role": "member"
    }
    ```
*   **Response Body Example (201 Created):**
    ```json
    {
      "id": 1,
      "workspace_id": 1,
      "user_id": 2,
      "role": "member",
      "joined_at": "2024-01-07T08:05:00Z"
    }
    ```

**`DELETE /api/workspaces/:workspaceID/members/:userID`**

*   **Description:** Removes a user from a specific workspace.
*   **Path Parameters:**
    *   `workspaceID`: The ID of the workspace.
    *   `userID`: The ID of the user to remove.
*   **Response:** `204 No Content` on successful deletion.

**`GET /api/workspaces/:workspaceID/members`**

*   **Description:** Retrieves all members of a specific workspace.
*   **Path Parameters:**
    *   `workspaceID`: The ID of the workspace.
*   **Response Body Example (200 OK):**
    ```json
    [
      {
        "id": 1,
        "workspace_id": 1,
        "user_id": 1,
        "role": "owner",
        "joined_at": "2024-01-07T08:00:00Z"
      },
      {
        "id": 2,
        "workspace_id": 1,
        "user_id": 2,
        "role": "member",
        "joined_at": "2024-01-07T08:05:00Z"
      }
    ]
    ```

---

### Channel Management

**`POST /api/channels`**

*   **Description:** Creates a new channel within a workspace.
*   **Request Body Example:**
    ```json
    {
      "workspace_id": 1,
      "name": "general",
      "description": "General discussions",
      "type": "public",
      "creator_id": 1
    }
    ```
*   **Response Body Example (201 Created):**
    ```json
    {
      "id": 1,
      "workspace_id": 1,
      "name": "general",
      "description": "General discussions",
      "type": "public",
      "created_by": 1,
      "created_at": "2024-01-07T08:10:00Z",
      "updated_at": "2024-01-07T08:10:00Z"
    }
    ```

**`GET /api/channels/:channelID`**

*   **Description:** Retrieves a channel by its ID.
*   **Path Parameters:**
    *   `channelID`: The ID of the channel.
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "workspace_id": 1,
      "name": "general",
      "description": "General discussions",
      "type": "public",
      "created_by": 1,
      "created_at": "2024-01-07T08:10:00Z",
      "updated_at": "2024-01-07T08:10:00Z"
    }
    ```

**`PUT /api/channels/:channelID`**

*   **Description:** Updates an existing channel's information.
*   **Path Parameters:**
    *   `channelID`: The ID of the channel to update.
*   **Request Body Example:**
    ```json
    {
      "name": "general-chat",
      "description": "General discussions for the team"
    }
    ```
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "workspace_id": 1,
      "name": "general-chat",
      "description": "General discussions for the team",
      "type": "public",
      "creator_id": 1,
      "created_at": "2024-01-07T08:10:00Z",
      "updated_at": "2024-01-07T10:00:00Z"
    }
    ```

**`DELETE /api/channels/:channelID`**

*   **Description:** Deletes a channel by its ID.
*   **Path Parameters:**
    *   `id`: The ID of the channel to delete.
*   **Response:** `204 No Content` on successful deletion.

**`GET /api/workspaces/:workspaceID/channels`**

*   **Description:** Retrieves all channels within a specific workspace.
*   **Path Parameters:**
    *   `workspaceID`: The ID of the workspace.
*   **Response Body Example (200 OK):**
    ```json
    [
      {
        "id": 1,
        "workspace_id": 1,
        "name": "general",
        "description": "General discussions",
        "type": "public",
        "creator_id": 1,
        "created_at": "2024-01-07T08:10:00Z",
        "updated_at": "2024-01-07T08:10:00Z"
      },
      {
        "id": 2,
        "workspace_id": 1,
        "name": "random",
        "description": "Random thoughts",
      "type": "public",
      "creator_id": 1,
      "created_at": "2024-01-07T08:15:00Z",
      "updated_at": "2024-01-07T08:15:00Z"
      }
    ]
    ```

---

### Channel Member Management

**`POST /api/channels/:channelID/members`**

*   **Description:** Adds a user as a member to a specific channel.
*   **Path Parameters:**
    *   `channelID`: The ID of the channel.
*   **Request Body Example:**
    ```json
    {
      "user_id": 2
    }
    ```
*   **Response Body Example (201 Created):**
    ```json
    {
      "id": 1,
      "channel_id": 1,
      "user_id": 2,
      "joined_at": "2024-01-07T08:20:00Z"
    }
    ```

**`DELETE /api/channels/:channelID/members/:userID`**

*   **Description:** Removes a user from a specific channel.
*   **Path Parameters:**
    *   `channelID`: The ID of the channel.
    *   `userID`: The ID of the user to remove.
*   **Response:** `204 No Content` on successful deletion.

**`GET /api/channels/:channelID/members`**

*   **Description:** Retrieves all members of a specific channel.
*   **Path Parameters:**
    *   `channelID`: The ID of the channel.
*   **Response Body Example (200 OK):**
    ```json
    [
      {
        "id": 1,
        "channel_id": 1,
        "user_id": 1,
        "joined_at": "2024-01-07T08:10:00Z"
      },
      {
        "id": 2,
        "channel_id": 1,
        "user_id": 2,
        "joined_at": "2024-01-07T08:20:00Z"
      }
    ]
    ```

---

### Message Management

**`POST /api/messages`**

*   **Description:** Sends a new message to a channel.
*   **Request Body Example:**
    ```json
    {
      "channel_id": 1,
      "sender_id": 1,
      "content": "Hello team, this is a test message!"
    }
    ```
*   **Response Body Example (201 Created):**
    ```json
    {
      "id": 1,
      "channel_id": 1,
      "sender_id": 1,
      "content": "Hello team, this is a test message!",
      "created_at": "2024-01-07T08:25:00Z",
      "updated_at": "2024-01-07T08:25:00Z"
    }
    ```

**`GET /api/messages/:messageID`**

*   **Description:** Retrieves a message by its ID.
*   **Path Parameters:**
    *   `messageID`: The ID of the message.
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "channel_id": 1,
      "sender_id": 1,
      "content": "Hello team, this is a test message!",
      "created_at": "2024-01-07T08:25:00Z",
      "updated_at": "2024-01-07T08:25:00Z"
    }
    ```

**`PUT /api/messages/:messageID`**

*   **Description:** Updates an existing message's content.
*   **Path Parameters:**
    *   `messageID`: The ID of the message to update.
*   **Request Body Example:**
    ```json
    {
      "content": "Hello team, this is an updated message!"
    }
    ```
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "channel_id": 1,
      "sender_id": 1,
      "content": "Hello team, this is an updated message!",
      "created_at": "2024-01-07T08:25:00Z",
      "updated_at": "2024-01-07T08:30:00Z"
    }
    ```

**`DELETE /api/messages/:messageID`**

*   **Description:** Deletes a message by its ID.
*   **Path Parameters:**
    *   `id`: The ID of the message to delete.
*   **Response:** `204 No Content` on successful deletion.

**`GET /api/channels/:channelID/messages?limit={limit}&offset={offset}`**

*   **Description:** Retrieves messages from a specific channel, with optional pagination.
*   **Path Parameters:**
    *   `channelID`: The ID of the channel.
*   **Query Parameters:**
    *   `limit`: (Optional) Maximum number of messages to return (default: 100).
    *   `offset`: (Optional) Number of messages to skip (default: 0).
*   **Response Body Example (200 OK):**
    ```json
    [
      {
        "id": 1,
        "channel_id": 1,
        "sender_id": 1,
        "content": "Hello team, this is a test message!",
        "created_at": "2024-01-07T08:25:00Z",
        "updated_at": "2024-01-07T08:25:00Z"
      },
      {
        "id": 2,
        "channel_id": 1,
        "sender_id": 2,
        "content": "Hey John, great message!",
        "created_at": "2024-01-07T08:26:00Z",
        "updated_at": "2024-01-07T08:26:00Z"
      }
    ]
    ```

---

### Attachment Management

**`POST /api/attachments`**

*   **Description:** Uploads a new attachment. (Note: Actual file upload mechanisms often involve `multipart/form-data` and may be handled differently than pure JSON. This example assumes metadata submission).
*   **Request Body Example:**
    ```json
    {
      "message_id": 1,
      "file_name": "report.pdf",
      "file_type": "application/pdf",
      "file_size": 102400,
      "url": "https://example.com/attachments/report.pdf",
      "uploaded_by": 1
    }
    ```
*   **Response Body Example (201 Created):**
    ```json
    {
      "id": 1,
      "message_id": 1,
      "file_name": "report.pdf",
      "file_type": "application/pdf",
      "file_size": 102400,
      "url": "https://example.com/attachments/report.pdf",
      "uploaded_by": 1,
      "uploaded_at": "2024-01-07T08:35:00Z"
    }
    ```

**`GET /api/attachments/:attachmentID`**

*   **Description:** Retrieves an attachment by its ID.
*   **Path Parameters:**
    *   `attachmentID`: The ID of the attachment.
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "message_id": 1,
      "file_name": "report.pdf",
      "file_type": "application/pdf",
      "file_size": 102400,
      "url": "https://example.com/attachments/report.pdf",
      "uploaded_by": 1,
      "uploaded_at": "2024-01-07T08:35:00Z"
    }
    ```

**`GET /api/messages/:messageID/attachments`**

*   **Description:** Retrieves all attachments associated with a specific message.
*   **Path Parameters:**
    *   `messageID`: The ID of the message.
*   **Response Body Example (200 OK):**
    ```json
    [
      {
        "id": 1,
        "message_id": 1,
        "file_name": "report.pdf",
        "file_type": "application/pdf",
        "file_size": 102400,
        "url": "https://example.com/attachments/report.pdf",
        "uploaded_by": 1,
        "uploaded_at": "2024-01-07T08:35:00Z"
      },
      {
        "id": 2,
        "message_id": 1,
        "file_name": "image.png",
        "file_type": "image/png",
        "file_size": 51200,
        "url": "https://example.com/attachments/image.png",
        "uploaded_by": 1,
        "uploaded_at": "2024-01-07T08:36:00Z"
      }
    ]
    ```

---

### Reaction Management

**`POST /api/messages/:messageID/reactions`**

*   **Description:** Adds a reaction to a message.
*   **Path Parameters:**
    *   `messageID`: The ID of the message.
*   **Request Body Example:**
    ```json
    {
      "user_id": 1,
      "emoji": "👍"
    }
    ```
*   **Response Body Example (201 Created):**
    ```json
    {
      "id": 1,
      "message_id": 1,
      "user_id": 1,
      "emoji": "👍",
      "created_at": "2024-01-07T08:40:00Z"
    }
    ```

**`DELETE /api/messages/:messageID/reactions/:userID/:emoji`**

*   **Description:** Removes a specific reaction from a message by a user.
*   **Path Parameters:**
    *   `messageID`: The ID of the message.
    *   `userID`: The ID of the user who added the reaction.
    *   `emoji`: The URL-encoded emoji character (e.g., `%F0%9F%91%8D` for 👍).
*   **Response:** `204 No Content` on successful deletion.

**`GET /api/messages/:messageID/reactions`**

*   **Description:** Retrieves all reactions for a specific message.
*   **Path Parameters:**
    *   `messageID`: The ID of the message.
*   **Response Body Example (200 OK):**
    ```json
    [
      {
        "id": 1,
        "message_id": 1,
        "user_id": 1,
        "emoji": "👍",
        "created_at": "2024-01-07T08:40:00Z"
      },
      {
        "id": 2,
        "message_id": 1,
        "user_id": 2,
        "emoji": "❤️",
        "created_at": "2024-01-07T08:41:00Z"
      }
    ]
    ```