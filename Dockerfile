                   
FROM golang:1.22-alpine AS builder

WORKDIR /byecom

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init

RUN go build -o main .

FROM alpine:latest

WORKDIR /byecom

COPY --from=builder /byecom/main .

COPY --from=builder /byecom/assets ./assets
COPY --from=builder /byecom/docs ./docs
COPY --from=builder /byecom/template ./template

EXPOSE 8080

CMD ["./main"]
