package app

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/websocket"
)

type apiServer struct {
	appServer *appServer

	upgrader websocket.Upgrader

	// websocket settings
	wsMaxMessageSize   int64
	wsWriteWait        time.Duration
	wsPongWait         time.Duration
	wsPingPeriod       time.Duration // Send pings to peer with this period. Must be less than pongWait.
	wsCloseGracePerdio time.Duration

	updateCancelFunc func()
}

func (a *apiServer) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		switch d := data.(type) {
		case error:
			http.Error(w, d.Error(), http.StatusInternalServerError)
		default:
			err := json.NewEncoder(w).Encode(data)
			if err != nil {
				http.Error(w, "failed to create response", http.StatusInternalServerError)
			}
		}
	}
}

func (a *apiServer) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *apiServer) wsError(ws *websocket.Conn, msg string, err error) {
	log.Info().Err(err).Msg("websocket error")
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
						log.Info().Err(err).Msg("ping failure")
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
		log.Info().Err(scanner.Err()).Msg("failed to scan")
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

func (s *apiServer) handleUpdate(connector func() (net.Conn, error)) http.HandlerFunc {
	type response struct {
		Status string `json:"status"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if s.updateCancelFunc != nil {
			s.respond(w, r, response{Status: "update already scheduled"}, http.StatusOK)
			return
		}

		conn, err := connector()
		if err != nil {
			s.respond(w, r, fmt.Errorf("failed to establish RCON connection: %w", err), http.StatusInternalServerError)
			return
		}

		_, err = fmt.Fprintf(conn, "login %s\n", "changeme")
		if err != nil {
			s.respond(w, r, fmt.Errorf("failed to write login command: %w", err), http.StatusInternalServerError)
			return
		}

		go func() {
			defer conn.Close()
			_, err = fmt.Fprintf(conn, "say \"update scheduled in 5 minutes\"\n")
			if err != nil {
				log.Error().Err(err).Msg("could not announce update")
				return
			}

			updateContext, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()
			s.updateCancelFunc = cancel
			updateTimer := time.NewTimer(5 * time.Minute)
			log.Info().Time("update_date", time.Now().Add(5*time.Minute)).Msg("update scheduled")

			select {
			case <-updateTimer.C:
				_, err = fmt.Fprintf(conn, "save\n")
				if err != nil {
					log.Error().Err(err).Msg("could not announce update")
					return
				}
				time.Sleep(5 * time.Second)

				_, err := backupServer(filepath.Join(os.Getenv("BACKUP_DIRECTORY"), "update"))
				if err != nil {
					log.Error().Err(err).Msg("failed to backup server before update")
					return
				}

				_, err = fmt.Fprintf(conn, "shutdown\n")
				if err != nil {
					s.respond(w, r, fmt.Errorf("failed to write shutdown command: %w", err), http.StatusInternalServerError)
					return
				}
			case <-updateContext.Done():
				_, err = fmt.Fprintf(conn, "say \"update canceled\"\n")
				if err != nil {
					log.Error().Err(err).Msg("could not announce update cancelation")
				}
				log.Info().Msg("update canceled")
				updateTimer.Stop()
				s.updateCancelFunc = nil
			}
		}()

		s.respond(w, r, response{Status: "update scheduled"}, http.StatusOK)
	}
}

func (s *apiServer) handleUpdateCancel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.updateCancelFunc()
	}
}

func (s *apiServer) handleBackup() http.HandlerFunc {
	type response struct {
		Status string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := backupServer(filepath.Join(os.Getenv("BACKUP_DIRECTORY"), "manual"))
		if err != nil {
			s.respond(w, r, fmt.Errorf("failed to backup server: %w", err), http.StatusInternalServerError)
			return
		}
		log.Info().Str("output", fmt.Sprintf("%s", output)).Msg("server backup completed")
		s.respond(w, r, response{Status: "OK"}, http.StatusOK)
	}
}

func backupServer(dst string) ([]byte, error) {
	cmd := exec.Command("rsync", "-r", os.Getenv("SERVER_DIRECTORY"), dst)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return data, fmt.Errorf("failed to run rsync: %w", err)
	}
	return data, nil
}

func (s *apiServer) handleFiles() http.HandlerFunc {
	baseDir := os.Getenv("SERVER_DIRECTORY")
	type File struct {
		Path        string `json:"path"`
		Name        string `json:"name"`
		ContentType string `json:"content_type"`
	}
	type response struct {
		Files   []File   `json:"files"`
		Folders []string `json:"folders"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")

		if strings.Contains(path, "..") {
			s.respond(w, r, errors.New(".. not allowed in path"), http.StatusBadRequest)
			return
		}

		fp := baseDir
		if path != "" {
			fp = filepath.Join(baseDir, path)
		}

		log.Info().Str("path", fp).Msg("opening")
		fs, err := os.Stat(fp)
		if err != nil {
			s.respond(w, r, fmt.Errorf("failed to stat path: %w", err), http.StatusBadRequest)
			return
		}

		if fs.IsDir() {
			files, err := ioutil.ReadDir(fp)
			if err != nil {
				s.respond(w, r, errors.New("failed to list directory"), http.StatusInternalServerError)
				return
			}

			res := response{
				Files:   []File{},
				Folders: []string{},
			}

			for _, item := range files {
				if item.IsDir() {
					res.Folders = append(res.Folders, item.Name())
					continue
				}

				contentType := guessContentType(filepath.Join(fp, item.Name()))
				res.Files = append(res.Files, File{Path: path, Name: item.Name(), ContentType: contentType})
			}

			s.respond(w, r, res, http.StatusOK)
			return
		}

		file, err := os.Open(fp)
		if err != nil {
			s.respond(w, r, fmt.Errorf("failed to open screenshot: %w", err), http.StatusBadRequest)
			return
		}
		defer file.Close()

		w.Header().Add("Content-Type", guessContentType(fp))

		_, err = io.Copy(w, file)
		if err != nil {
			s.respond(w, r, fmt.Errorf("failed to copy screenshot to response: %w", err), http.StatusBadRequest)
		}

	}
}

func (s *apiServer) handleFilesSave() http.HandlerFunc {
	baseDir := os.Getenv("SERVER_DIRECTORY")
	type response struct {
		Status string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")

		if strings.Contains(path, "..") {
			s.respond(w, r, errors.New(".. not allowed in path"), http.StatusBadRequest)
			return
		}

		fp := baseDir
		if path != "" {
			fp = filepath.Join(baseDir, path)
		}
		defer r.Body.Close()

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read request body")
			s.respond(w, r, response{Status: fmt.Errorf("failed to read request body: %w", err).Error()}, http.StatusBadRequest)
			return
		}

		if err := ioutil.WriteFile(fp, data, 0); err != nil {
			log.Info().Str("path", fp).Msg("failed to save file")
			s.respond(w, r, response{Status: fmt.Errorf("failed to save file: %w", err).Error()}, http.StatusBadRequest)
			return
		}

		log.Info().Str("path", fp).Msg("saved file")
		s.respond(w, r, response{Status: "OK"}, http.StatusOK)
	}
}

func (s *apiServer) handleFilesDelete() http.HandlerFunc {
	baseDir := os.Getenv("SERVER_DIRECTORY")
	type response struct {
		Status string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")

		if strings.Contains(path, "..") {
			s.respond(w, r, errors.New(".. not allowed in path"), http.StatusBadRequest)
			return
		}

		fp := baseDir
		if path != "" {
			fp = filepath.Join(baseDir, path)
		}

		if err := os.Remove(fp); err != nil {
			log.Info().Str("path", fp).Msg("failed to delete file")
			s.respond(w, r, response{Status: fmt.Errorf("failed to remove file: %w", err).Error()}, http.StatusBadRequest)
			return
		}

		log.Info().Str("path", fp).Msg("deleted file")
		s.respond(w, r, response{Status: "OK"}, http.StatusOK)
	}
}

func guessContentType(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("could not guess content-type using application/octet-stream")
		return "application/octet-stream"
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("could not guess content-type using application/octet-stream")
		return "application/octet-stream"
	}

	return http.DetectContentType(buffer)
}

func (s *apiServer) handleLogin() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type user struct {
		Username     string
		PasswordHash []byte
		Permissions  []string
	}
	var users []user

	userFile, err := os.Open("users.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to open users list")
		return func(w http.ResponseWriter, r *http.Request) {
		}
	}

	if err := json.NewDecoder(userFile).Decode(&users); err != nil {
		log.Error().Err(err).Msg("failed to read users")
		return func(w http.ResponseWriter, r *http.Request) {
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		log.Info().Msg("got login request")

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("login failed")
			s.respond(w, r, fmt.Errorf("failed to decode login request"), http.StatusBadRequest)
			return
		}

		var usr *user
		for _, u := range users {
			if u.Username == req.Username {
				usr = &user{Username: u.Username, PasswordHash: u.PasswordHash, Permissions: u.Permissions}
				break
			}
		}

		if usr == nil {
			log.Error().Str("username", req.Username).Msg("user does not exist")
			s.respond(w, r, fmt.Errorf("username or password wrong"), http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(req.Password)); err != nil {
			log.Error().Str("username", req.Username).Msg("password did not match")
			s.respond(w, r, fmt.Errorf("username or password wrong"), http.StatusBadRequest)
			return
		}

		log.Info().Strs("permissions", usr.Permissions).Msg("login succeeded")
		if err := s.appServer.setAuthCookies(w, r, Token{Username: usr.Username, Permissions: usr.Permissions}); err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
