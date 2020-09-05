FROM golang:latest

WORKDIR /app

RUN mkdir files

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]