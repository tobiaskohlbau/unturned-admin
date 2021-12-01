FROM golang:1.17 AS builder

WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod download
COPY main.go /app/
COPY app /app/app
COPY mock /app/mock
COPY store /app/store
RUN CGO_ENABLED=0 go build

FROM alpine

WORKDIR /app
COPY --from=builder /app/unturned-admin /usr/local/bin/unturned-admin

ENTRYPOINT [ "/usr/local/bin/unturned-admin" ]
