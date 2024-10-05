FROM golang:1.22.7

WORKDIR /app
COPY go.mod ./

COPY *.go ./

RUN CGO_ENABLED=1 GOOS=linux go build -o /padmd_app

EXPOSE 8080
