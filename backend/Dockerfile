FROM golang:1.22.2

WORKDIR /app

RUN apt-get update && apt-get install -y make && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/url-shortener ./cmd/api

EXPOSE 8080

CMD ["make", "serve"]
