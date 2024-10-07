FROM golang:1.22.7

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=1 GOOS=linux go build -o /padmd_app

EXPOSE 8080
