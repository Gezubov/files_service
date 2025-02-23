FROM golang:1.23.5

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 8082

CMD ["go", "run", "cmd/app/main.go"]
