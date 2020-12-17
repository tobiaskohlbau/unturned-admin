# FROM node AS node-builder

# WORKDIR /web
# COPY web/package.json web/package-lock.json /web
# RUN npm install
# COPY web /web
# RUN npm run build

FROM golang:1.16-rc AS go-builder

WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod download
COPY main.go /app/
COPY web/web.go /app/web/
COPY app /app/app
COPY mock /app/mock
# COPY --from=node-builder /web/dist /app/web/dist
COPY ./web/dist /app/web/dist
RUN CGO_ENABLED=0 go build

FROM alpine

RUN apk add --no-cache rsync

WORKDIR /app/
COPY --from=go-builder /app/unturned-admin /usr/local/bin/unturned-admin

ENTRYPOINT [ "/usr/local/bin/unturned-admin" ]
