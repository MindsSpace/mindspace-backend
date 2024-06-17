
FROM golang:1.22

RUN apt-get update && apt-get install -y \
    git \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . ./

ENV GOROOT /usr/local/go

RUN CGO_ENABLED=0 GOOS=linux go build -o mindspace-backend .

EXPOSE 8080

CMD ["./mindspace-backend"]