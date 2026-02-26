# Application Error Codes

This document lists the internal error codes used in the application. These codes are returned in the JSON response when an error occurs.

## Error Response Format

```json
{
  "code": "ERROR_CODE",
  "message": "Human readable error message"
}
```

## Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `BAD_REQUEST` | 400 | The request was invalid or cannot be served. |
| `INVALID_INPUT` | 400 | The input provided is invalid (e.g., missing fields, wrong format). |
| `UNAUTHORIZED` | 401 | Authentication is required and has failed or has not been provided. |
| `FORBIDDEN` | 403 | The request was valid, but the server is refusing action. |
| `NOT_FOUND` | 404 | The requested resource could not be found. |
| `INTERNAL_ERROR` | 500 | An internal server error occurred. |

## User Domain Errors

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `USER_NOT_FOUND` | 404 | The requested user ID does not exist. |
| `USER_ALREADY_EXISTS` | 409 | A user with the same email or username already exists. |

## File Upload Errors

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `FILE_TOO_LARGE` | 400 | The uploaded file exceeds the maximum allowed size. |
| `UPLOAD_FAILED` | 500 | Failed to process or save the uploaded file. |
