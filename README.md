
# STRIVE Backend

This is the backend service for the STRIVE app, which tracks and visualizes daily effort and productivity.

## Setup

### Prerequisites

- Go 1.18+
- Docker
- Docker Compose

### Running the Services

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/STRIVEBackend.git
   cd STRIVEBackend

2. Start the service using docker compose
"docker-compose up --build"

3. The service will be available at http://localhost:50051

### Project Structure
`cmd/`: Entry points for each microservice
`pkg/`: Core logic for each microservice
`internal/`: Internal packages (e.g database, grpc, auth)
`api/`: Protocal buffer files for gRPC services
`config/`: Configuration files for each microservice
`scripts/`: Scripts for building and running the services
`docker-compose.yml`: Docker compose file for running the services


echo "# STRIVE Backend

This project contains the backend services for the STRIVE application. It is built using Golang for the backend and is structured using a microservice architecture.

## Services

- **User Service**: Manages user authentication and profiles.
- **Activity Service**: Handles logging and retrieving user activities.
- **Score Service**: Calculates and provides the daily effort score for users.

## Structure

\`\`\`
strive-backend/
├── cmd/
│   ├── user/
│   │   └── main.go
│   ├── activity/
│   │   └── main.go
│   └── score/
│       └── main.go
├── pkg/
│   ├── user/
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── activity/
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   └── score/
│       ├── handler.go
│       ├── model.go
│       ├── repository.go
│       └── service.go
├── internal/
│   ├── database/
│   │   └── connection.go
│   └── auth/
│       └── auth.go
├── api/
│   ├── user/
│   │   └── user.proto
│   ├── activity/
│   │   └── activity.proto
│   └── score/
│       └── score.proto
├── configs/
│   ├── user-config.yaml
│   ├── activity-config.yaml
│   └── score-config.yaml
├── Dockerfile
├── docker-compose.yaml
└── README.md
\`\`\`

## Getting Started

To get started with the development, follow these steps:

1. Clone the repository
2. Navigate to the backend folder
3. Run \`go mod init\` to initialize the Go module
4. Start developing your microservices

" > README.md
