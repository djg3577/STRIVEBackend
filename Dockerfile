FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /strive-backend ./cmd/user

EXPOSE 50051

CMD ["/strive-backend"]
