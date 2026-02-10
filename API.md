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

**`POST /api/workspaces/:workspaceID/join`**

*   **Description:** Allows an authenticated user to join a specific workspace.
*   **Authentication:** Required.
*   **Path Parameters:**
    *   `workspaceID`: The ID of the workspace to join.
*   **Request Body:** None (User ID is taken from the JWT).
*   **Response Body Example (201 Created):**
    ```json
    {
      "workspace_id": 1,
      "user_id": 2,
      "role": "member",
      "created_at": "2024-02-09T10:00:00Z"
    }
    ```
*   **Error Responses:**
    *   `401 Unauthorized`: If authentication fails.
    *   `404 Not Found`: If the `workspaceID` does not exist.
    *   `409 Conflict`: If the user is already a member of the workspace.
    *   `500 Internal Server Error`: For other server-side errors.

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

*   **Description:** Sends a new message to a meeting.
*   **Request Body Example:**
    ```json
    {
      "meeting_id": 1,
      "content": "Hello team, this is a test message!"
    }
    ```
*   **Response Body Example (201 Created):**
    ```json
    {
      "id": 1,
      "meeting_id": 1,
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
      "meeting_id": 1,
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
      "meeting_id": 1,
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

**`GET /api/meetings/:meetingID/messages?limit={limit}&offset={offset}`**

*   **Description:** Retrieves messages from a specific meeting, with optional pagination.
*   **Path Parameters:**
    *   `meetingID`: The ID of the meeting.
*   **Query Parameters:**
    *   `limit`: (Optional) Maximum number of messages to return (default: 100).
    *   `offset`: (Optional) Number of messages to skip (default: 0).
*   **Response Body Example (200 OK):**
    ```json
    [
      {
        "id": 1,
        "meeting_id": 1,
        "sender_id": 1,
        "content": "Hello team, this is a test message!",
        "created_at": "2024-01-07T08:25:00Z",
        "updated_at": "2024-01-07T08:25:00Z"
      },
      {
        "id": 2,
        "meeting_id": 1,
        "sender_id": 2,
        "content": "Hey John, great message!",
        "created_at": "2024-01-07T08:26:00Z",
        "updated_at": "2024-01-07T08:26:00Z"
      }
    ]
    ```

---

### Meeting Management

**`POST /api/meetings`**

*   **Description:** Creates a new meeting.
*   **Request Body Example:**
    ```json
    {
      "name": "Team Standup",
      "description": "Daily team synchronization meeting",
      "channel_id": 1,
      "start_time": "2024-01-07T09:00:00Z",
      "end_time": "2024-01-07T09:30:00Z",
      "participant_ids": [1, 2]
    }
    ```
*   **Response Body Example (201 Created):**
    ```json
    {
      "id": 1,
      "name": "Team Standup",
      "description": "Daily team synchronization meeting",
      "channel_id": 1,
      "creator_id": 1,
      "start_time": "2024-01-07T09:00:00Z",
      "end_time": "2024-01-07T09:30:00Z",
      "created_at": "2024-01-07T08:50:00Z",
      "updated_at": "2024-01-07T08:50:00Z"
    }
    ```

**`GET /api/meetings/:meetingID`**

*   **Description:** Retrieves a meeting by its ID.
*   **Path Parameters:**
    *   `meetingID`: The ID of the meeting.
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "name": "Team Standup",
      "description": "Daily team synchronization meeting",
      "channel_id": 1,
      "creator_id": 1,
      "start_time": "2024-01-07T09:00:00Z",
      "end_time": "2024-01-07T09:30:00Z",
      "created_at": "2024-01-07T08:50:00Z",
      "updated_at": "2024-01-07T08:50:00Z"
      // Includes participants and potentially messages (if relations are loaded)
    }
    ```

**`PUT /api/meetings/:meetingID`**

*   **Description:** Updates an existing meeting's information.
*   **Path Parameters:**
    *   `meetingID`: The ID of the meeting to update.
*   **Request Body Example:**
    ```json
    {
      "name": "Daily Team Standup",
      "description": "Updated daily synchronization meeting",
      "channel_id": 1,
      "start_time": "2024-01-07T09:00:00Z",
      "end_time": "2024-01-07T09:45:00Z"
    }
    ```
*   **Response Body Example (200 OK):**
    ```json
    {
      "id": 1,
      "name": "Daily Team Standup",
      "description": "Updated daily synchronization meeting",
      "channel_id": 1,
      "creator_id": 1,
      "start_time": "2024-01-07T09:00:00Z",
      "end_time": "2024-01-07T09:45:00Z",
      "created_at": "2024-01-07T08:50:00Z",
      "updated_at": "2024-01-07T09:05:00Z"
    }
    ```

**`DELETE /api/meetings/:meetingID`**

*   **Description:** Deletes a meeting by its ID.
*   **Path Parameters:**
    *   `meetingID`: The ID of the meeting to delete.
*   **Response:** `204 No Content` on successful deletion.

**`GET /api/channels/:channelID/meetings`**

*   **Description:** Retrieves all meetings within a specific channel.
*   **Path Parameters:**
    *   `channelID`: The ID of the channel.
*   **Response Body Example (200 OK):**
    ```json
    [
      {
        "id": 1,
        "name": "Team Standup",
        "description": "Daily team synchronization meeting",
        "channel_id": 1,
        "creator_id": 1,
        "start_time": "2024-01-07T09:00:00Z",
        "end_time": "2024-01-07T09:30:00Z",
        "created_at": "2024-01-07T08:50:00Z",
        "updated_at": "2024-01-07T08:50:00Z"
      }
    ]
    ```

**`POST /api/meetings/:meetingID/participants`**

*   **Description:** Adds a user as a participant to a specific meeting.
*   **Path Parameters:**
    *   `meetingID`: The ID of the meeting.
*   **Request Body Example:**
    ```json
    {
      "participant_id": 3
    }
    ```
*   **Response Body Example (200 OK):**
    ```json
    {
      "message": "Participant added successfully"
    }
    ```

**`DELETE /api/meetings/:meetingID/participants/:participantID`**

*   **Description:** Removes a user from a specific meeting.
*   **Path Parameters:**
    *   `meetingID`: The ID of the meeting.
    *   `participantID`: The ID of the user to remove.
*   **Response:** `200 OK` on successful removal.

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
---

### WebSocket Chat API

This section describes the real-time WebSocket communication for chat functionality. The API provides a simplified interface for easier usage while maintaining all real-time features.

#### Authentication for WebSockets

*   **Method:** Include an `Authorization` header in WebSocket handshake request.
*   **Header Name:** `Authorization`
*   **Header Value:** `Bearer <YOUR_JWT_TOKEN>` (replace `<YOUR_JWT_TOKEN>` with a valid JWT obtained from the `/api/login` endpoint).

#### `GET /ws/meeting/:meeting_id/chat`

*   **Description:** Establishes a WebSocket connection for real-time chat within a specific meeting. Supports sending messages, reactions, typing indicators, and fetching message history.
*   **Path Parameters:**
    *   `meeting_id`: The ID of the meeting for which to join the chat.
*   **Connection URL Example:** `ws://localhost:8080/ws/meeting/123/chat`

---

#### WebSocket Message Format

All WebSocket messages use a simple, consistent format:

```json
{
  "type": "message_type",
  "data": { /* type-specific data */ }
}
```

#### Available Message Types

**Client to Server (Incoming):**
- `message` - Send a chat message
- `reaction` - Add/remove reaction to a message
- `typing` - Send typing indicator
- `history` - Request message history

**Server to Client (Outgoing):**
- `message` - New chat message broadcast
- `reaction` - Reaction update broadcast
- `typing` - Typing indicator broadcast
- `room` - User join/leave notifications
- `history` - Message history response
- `error` - Error message

---

### Incoming WebSocket Messages (Client to Server)

#### 1. Send Message (`type: "message"`)

Sends a new chat message to the meeting.

**Request:**
```json
{
  "type": "message",
  "data": {
    "content": "Hello team!",
    "room_id": 123,
    "type": "text",                    // "text", "file", "system"
    "reply_to": 456,                   // Optional: ID of message to reply to
    "files": [                         // Optional: File attachments
      {
        "id": 1,
        "name": "document.pdf",
        "type": "application/pdf",
        "size": 1024000,
        "url": "https://example.com/files/document.pdf"
      }
    ]
  }
}
```

#### 2. Add/Remove Reaction (`type: "reaction"`)

Adds or removes a reaction from a message.

**Request:**
```json
{
  "type": "reaction",
  "data": {
    "message_id": 789,
    "emoji": "👍",
    "action": "add"                    // "add" or "remove"
  }
}
```

#### 3. Typing Indicator (`type: "typing"`)

Sends a typing indicator to other users.

**Request:**
```json
{
  "type": "typing",
  "data": {
    "room_id": 123,
    "is_typing": true
  }
}
```

#### 4. Request Message History (`type: "history"`)

Requests previous messages in the meeting.

**Request:**
```json
{
  "type": "history",
  "data": {
    "room_id": 123,
    "offset": 0                        // Pagination offset
  }
}
```

---

### Outgoing WebSocket Messages (Server to Client)

#### 1. New Message Broadcast (`type: "message"`)

Broadcasted when a new message is sent to the meeting.

**Response:**
```json
{
  "type": "message",
  "data": {
    "id": 101,
    "content": "Hello team!",
    "room_id": 123,
    "user_id": 1,
    "timestamp": "2024-01-07T10:00:00Z",
    "type": "text",
    "reply_to": null,
    "files": [],
    "user": {
      "id": 1,
      "name": "John Doe",
      "username": "johndoe"
    }
  }
}
```

#### 2. Reaction Update (`type: "reaction"`)

Broadcasted when a reaction is added or removed.

**Response:**
```json
{
  "type": "reaction",
  "data": {
    "message_id": 789,
    "user_id": 1,
    "emoji": "👍",
    "action": "add",
    "timestamp": "2024-01-07T10:01:00Z",
    "user": {
      "id": 1,
      "name": "John Doe",
      "username": "johndoe"
    }
  }
}
```

#### 3. Typing Indicator (`type: "typing"`)

Broadcasted when a user starts or stops typing.

**Response:**
```json
{
  "type": "typing",
  "data": {
    "user_id": 1,
    "room_id": 123,
    "is_typing": true,
    "user": {
      "id": 1,
      "name": "John Doe",
      "username": "johndoe"
    }
  }
}
```

#### 4. Room Events (`type: "room"`)

Broadcasted when a user joins or leaves the meeting chat.

**Response:**
```json
{
  "type": "room",
  "data": {
    "room_id": 123,
    "user_id": 1,
    "action": "join",                  // "join" or "leave"
    "timestamp": "2024-01-07T10:02:00Z",
    "user": {
      "id": 1,
      "name": "John Doe",
      "username": "johndoe"
    }
  }
}
```

#### 5. Message History Response (`type: "history"`)

Response to a history request.

**Response:**
```json
{
  "type": "history",
  "data": {
    "room_id": 123,
    "messages": [
      {
        "id": 100,
        "content": "Previous message",
        "room_id": 123,
        "user_id": 2,
        "timestamp": "2024-01-07T09:30:00Z",
        "type": "text",
        "reply_to": null,
        "files": [],
        "user": {
          "id": 2,
          "name": "Jane Smith",
          "username": "janesmith"
        }
      }
    ],
    "has_more": true,
    "offset": 50
  }
}
```

#### 6. Error Messages (`type: "error"`)

Sent when there's an error with a client request.

**Response:**
```json
{
  "type": "error",
  "data": {
    "code": "ROOM_MISMATCH",
    "message": "Room ID mismatch",
    "details": "Expected: 123, Got: 456"
  }
}
```

---

### Error Codes

| Code | Description |
|------|-------------|
| `INVALID_FORMAT` | Message format is invalid or malformed |
| `INVALID_MESSAGE_DATA` | Message data is invalid or missing required fields |
| `INVALID_REACTION_DATA` | Reaction data is invalid |
| `INVALID_TYPING_DATA` | Typing data is invalid |
| `INVALID_HISTORY_DATA` | History request data is invalid |
| `ROOM_MISMATCH` | Room ID in message doesn't match connection room |
| `UNKNOWN_TYPE` | Unknown message type |
| `SEND_FAILED` | Failed to send message |
| `REACTION_FAILED` | Failed to process reaction |
| `HISTORY_FAILED` | Failed to retrieve message history |
| `INVALID_REACTION_ACTION` | Invalid reaction action (must be "add" or "remove") |

---

### Usage Examples

#### JavaScript Client Example

```javascript
// Connect to WebSocket
const token = 'your-jwt-token';
const meetingId = 123;
const ws = new WebSocket(`ws://localhost:8080/ws/meeting/${meetingId}/chat`, [], {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

// Handle incoming messages
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  switch (message.type) {
    case 'message':
      console.log('New message:', message.data);
      displayMessage(message.data);
      break;
    case 'reaction':
      console.log('Reaction update:', message.data);
      updateReaction(message.data);
      break;
    case 'typing':
      console.log('Typing indicator:', message.data);
      showTypingIndicator(message.data);
      break;
    case 'room':
      console.log('Room event:', message.data);
      updateParticipantList(message.data);
      break;
    case 'error':
      console.error('WebSocket error:', message.data);
      showError(message.data);
      break;
  }
};

// Send a message
function sendMessage(content, replyTo = null) {
  ws.send(JSON.stringify({
    type: 'message',
    data: {
      content: content,
      room_id: meetingId,
      type: 'text',
      reply_to: replyTo
    }
  }));
}

// Add a reaction
function addReaction(messageId, emoji) {
  ws.send(JSON.stringify({
    type: 'reaction',
    data: {
      message_id: messageId,
      emoji: emoji,
      action: 'add'
    }
  }));
}

// Send typing indicator
function sendTyping(isTyping) {
  ws.send(JSON.stringify({
    type: 'typing',
    data: {
      room_id: meetingId,
      is_typing: isTyping
    }
  }));
}

// Request message history
function requestHistory(offset = 0) {
  ws.send(JSON.stringify({
    type: 'history',
    data: {
      room_id: meetingId,
      offset: offset
    }
  }));
}
```

---

### WebSocket Features

- **Real-time messaging**: Instant delivery of chat messages
- **Reactions**: Add/remove emoji reactions to messages
- **Typing indicators**: See when other users are typing
- **Message history**: Paginated access to previous messages
- **User presence**: Join/leave notifications
- **Error handling**: Structured error messages with codes
- **File attachments**: Support for file sharing in messages
- **Message replies**: Threaded conversations with reply functionality
- **Robust error handling**: Automatic reconnection support recommended