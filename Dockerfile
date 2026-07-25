ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN go build -v -o /run-app .


FROM debian:bookworm

# Копируем бинарник
COPY --from=builder /run-app /usr/local/bin/

# Добавь вот эту строку - копируем папку web
COPY --from=builder /usr/src/app/web /web

CMD ["run-app"]