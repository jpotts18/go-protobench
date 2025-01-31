package json

import (
	"encoding/json"
	"net/http"
	"sync"

	"protobench/internal/model"
)

type Server struct {
	server *http.Server
	port   string
	wg     sync.WaitGroup
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/message", s.handleMessage)

	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: mux,
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	if s.server != nil {
		if err := s.server.Close(); err != nil {
			return err
		}
		s.wg.Wait()
	}
	return nil
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg model.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Echo the message back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}
