FROM node AS node-builder

WORKDIR /web
COPY web/package.json web/package-lock.json /web
RUN npm install
COPY web /web
RUN npm run build

FROM ghcr.io/tobiaskohlbau/golang:tip AS go-builder

WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod download
COPY main.go /app/
COPY web/web.go /app/web/
COPY --from=node-builder /web/dist /app/web/dist
RUN CGO_ENABLED=0 go build

FROM alpine

WORKDIR /app/
COPY --from=go-builder /app/unturned-admin /usr/local/bin/unturned-admin

ENTRYPOINT [ "/usr/local/bin/unturned-admin" ]
