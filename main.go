//go:generate cd web && npm run build

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/tobiaskohlbau/unturned-admin/web"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func errorHandler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func run() error {
	dev := flag.Bool("dev", false, "development mode")
	flag.Parse()

	r := chi.NewRouter()

	if *dev {
		target, err := url.Parse("http://localhost:3000")
		if err != nil {
			return fmt.Errorf("failed to create dev proxy: %w", err)
		}
		r.Get("/*", httputil.NewSingleHostReverseProxy(target).ServeHTTP)
	} else {
		r.Get("/*", http.FileServer(http.FS(web.Content)).ServeHTTP)
	}

	r.Get("/api/screenshots", screenshotHandler)

	if *dev {
		r.Get("/ws", serveWs(establishMockRCONConnection))
	} else {
		r.Get("/ws", serveWs(establishRCONConnection))
	}

	return http.ListenAndServe(":8080", r)
}

func establishRCONConnection() (net.Conn, error) {
	conn, err := net.Dial("tcp", "unturned:27115")
	if err != nil {
		return nil, fmt.Errorf("failed to establish RCON connection")
	}
	return conn, nil
}

func establishMockRCONConnection() (net.Conn, error) {
	server, client := net.Pipe()

	go func() {
		defer server.Close()
		for {
			server.Write([]byte("Hello World\n"))
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		defer server.Close()
		for {
			ignore := make([]byte, 1024)
			server.Read(ignore)
			fmt.Printf("%s", ignore)
			time.Sleep(1 * time.Second)
		}
	}()

	return client, nil
}

func screenshotHandler(w http.ResponseWriter, r *http.Request) {
	dir := os.Getenv("SCREENSHOT_DIRECTORY")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		http.Error(w, "failed to list screenshot directory", http.StatusInternalServerError)
		return
	}

	screenshots := []string{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		screenshots = append(screenshots, file.Name())
	}

	data, err := json.Marshal(screenshots)
	if err != nil {
		http.Error(w, "failed to marshal directory entries", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

var upgrader = websocket.Upgrader{}

func serveWs(connector func() (net.Conn, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer ws.Close()

		conn, err := connector()
		if err != nil {
			internalError(ws, "failed connecting to rcon server", err)
			return
		}
		defer conn.Close()

		fmt.Fprintf(conn, "login %s\n", "changeme")

		done := make(chan struct{})
		go ping(ws, done)

		go pumpOutput(ws, conn, done)

		pumpInput(ws, conn)

		<-done
	}
}

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}

func pumpInput(ws *websocket.Conn, w io.Writer) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = append(message, '\n')
		if _, err := w.Write(message); err != nil {
			break
		}
	}
}

func pumpOutput(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}
	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}
