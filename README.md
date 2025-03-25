# Messaging System

This project is a REST API service designed to automatically send unsent messages from the database every 2 minutes.

## Features

- Automatically sends unsent messages from the database every 2 minutes
- Sends a maximum of 2 messages per cycle
- Stores message content, recipient phone number, and delivery status in the database
- Messages that have been sent once are not sent again
- Redis caches message IDs and send times (bonus feature)
- API to start/stop the message sending service and list sent messages

## Technical Details

- Developed with Go 1.22
- Timing implemented using Go's standard library features without external cron packages
- Uses PostgreSQL database
- Integrated Redis caching system
- API documentation with Swagger
- Easy deployment with Docker and Docker Compose
- Automated mock generation using mockery for testing

## API Endpoints

- `POST /api/service?action=start|stop`: Starts or stops the message sending service
- `GET /api/service/status`: Gets the current status of the message service
- `GET /api/messages?page=1&limit=10`: Lists sent messages (with pagination support)

### API Documentation

API documentation is available via Swagger UI at:
```
http://localhost:8080/swagger/
```

The OpenAPI specification is also available directly at:
```
http://localhost:8080/docs/swagger.yaml
```

#### Using Swagger.io

You can also view and interact with the API documentation on swagger.io:

1. View in Swagger Editor:
   ```
   https://editor.swagger.io/?url=http://localhost:8080/docs/swagger.yaml
   ```

2. Validate the specification:
   ```
   https://validator.swagger.io/validator/debug?url=http://localhost:8080/docs/swagger.yaml
   ```

> Note: The application must be running and accessible for the swagger.io links to work properly.

## Installation

### Installation with Docker

1. Clone the project:
   ```bash
   git clone https://github.com/username/messaging-system.git
   cd messaging-system
   ```

2. Run with Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. Access the API:
   ```
   http://localhost:8080
   ```

4. Access Swagger UI:
   ```
   http://localhost:8080/swagger/
   ```

### Manual Installation

1. Make sure PostgreSQL and Redis are running.

2. Edit the `config.json` file:
   ```json
   {
     "server": {
       "port": 8080
     },
     "db": {
       "host": "localhost",
       "port": "5432",
       "user": "postgres",
       "password": "postgres",
       "name": "message_db"
     },
     "redis": {
       "addr": "localhost:6379",
       "password": "",
       "db": 0
     }
   }
   ```

3. Build and run the application:
   ```bash
   go build -o messaging-system ./cmd/server
   ./messaging-system
   ```

## Development

### Makefile Commands

The project includes several useful Makefile commands to simplify development:

- `make build`: Build the application
- `make run`: Build and run the application
- `make test`: Run all tests
- `make clean`: Clean build files
- `make mockery`: Generate mock interfaces for testing (requires Go installed)
- `make docker-build`: Build Docker image
- `make docker-run`: Run application in Docker
- `make docker-compose-up`: Start application with Docker Compose
- `make docker-compose-down`: Stop application with Docker Compose

### Generating Mocks

To generate mocks for testing:

```bash
make mockery
```

This command will:
1. Check if mockery is installed and install it if necessary
2. Generate mocks for repository interfaces
3. Generate mocks for service interfaces

## Database Schema

```sql
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    is_sent BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP,
    external_msg_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Testing

To add a new message, you can connect to PostgreSQL and run the following SQL query:

```sql
INSERT INTO messages (content, phone_number) VALUES ('Hello, this is a test message.', '+905551234567');
```

To check sent messages:

```sql
SELECT * FROM messages WHERE is_sent = true;
```

To check messages cached in Redis:

```bash
redis-cli keys "message:*"
```

## License

This project is distributed under the MIT license. For more information, please see the `LICENSE` file.