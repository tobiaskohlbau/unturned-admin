package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type apiServer struct {
	upgrader websocket.Upgrader

	// websocket settings
	wsMaxMessageSize   int64
	wsWriteWait        time.Duration
	wsPongWait         time.Duration
	wsPingPeriod       time.Duration // Send pings to peer with this period. Must be less than pongWait.
	wsCloseGracePerdio time.Duration
}

func (a *apiServer) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, "failed to create response", http.StatusInternalServerError)
		}
	}
}

func (a *apiServer) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *apiServer) wsError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

func (s *apiServer) handleRCON(connector func() (net.Conn, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.respond(w, r, err, http.StatusInternalServerError)
			return
		}
		defer ws.Close()

		conn, err := connector()
		if err != nil {
			s.wsError(ws, "failed connecting to rcon server", err)
			return
		}
		defer conn.Close()

		fmt.Fprintf(conn, "login %s\n", "changeme")

		done := make(chan struct{})
		// ping regulary to keep websocket open
		go func() {
			ticker := time.NewTicker(s.wsPingPeriod)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(s.wsWriteWait)); err != nil {
						log.Println("ping:", err)
					}
				case <-done:
					return
				}
			}
		}()

		go s.pumpOutput(ws, conn, done)
		s.pumpInput(ws, conn)

		<-done
	}
}

func (s *apiServer) pumpInput(ws *websocket.Conn, w io.Writer) {
	defer ws.Close()
	ws.SetReadLimit(s.wsMaxMessageSize)
	ws.SetReadDeadline(time.Now().Add(s.wsPongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(s.wsPongWait)); return nil })
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

func (s *apiServer) pumpOutput(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ws.SetWriteDeadline(time.Now().Add(s.wsWriteWait))
		if err := ws.WriteMessage(websocket.TextMessage, scanner.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if scanner.Err() != nil {
		log.Println("scan:", scanner.Err())
	}
	close(done)

	ws.SetWriteDeadline(time.Now().Add(s.wsWriteWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(s.wsCloseGracePerdio)
	ws.Close()
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

func (s *apiServer) handleScreenshots() http.HandlerFunc {
	dir := os.Getenv("SCREENSHOT_DIRECTORY")
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			s.respond(w, r, errors.New("failed to list screenshot directory"), http.StatusInternalServerError)
			return
		}

		screenshots := []string{}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			screenshots = append(screenshots, file.Name())
		}

		s.respond(w, r, screenshots, http.StatusOK)
	}
}
