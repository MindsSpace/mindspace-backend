
FROM golang:1.21.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./mindspace-backend

EXPOSE 8080
ENV PORT 8080

CMD ["./mindspace-backend"]