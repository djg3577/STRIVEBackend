
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
