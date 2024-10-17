FROM golang:1.22.7

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /padmd_app ./main.go

EXPOSE 8080

# Run the binary
CMD ["/padmd_app"]

