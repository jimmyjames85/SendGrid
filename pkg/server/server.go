package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi"
)

type Config struct {
	Port        int  `envconfig:"PORT" required:"false" default:"5555"` // service port to run on
	PrettyPrint bool `envconfig:"PRETTY_PRINT" required:"false" default:"true"`
}

func (c *Config) ToJSON() string {
	copy := *c
	// copy.DBPswd = "****"
	b, _ := json.Marshal(copy)
	return string(b)
}

type Server struct {
	cfg    Config
	router *chi.Mux
	server *http.Server
}

func (s *Server) handleEventWebhook(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.handleError(w, r, err, http.StatusInternalServerError)
		return
	}

	output := string(b)
	if s.cfg.PrettyPrint {

		var pp []map[string]interface{}
		err := json.Unmarshal(b, &pp)
		if err != nil {
			log.Printf("unable to marshal: %v", err)
		} else {
			b, err := json.MarshalIndent(pp, " ", "  ")
			if err == nil {
				output = string(b)
			}
		}
	}
	log.Printf("%s", output)

	w.WriteHeader(http.StatusAccepted)
}

// handleError provides a uniform way to emit errors out of our handlers. You should ALWAYS call
// return after calling it.
func (s *Server) handleError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)
	m := make(map[string]string)
	if err != nil {
		m["error"] = err.Error()
	}
	m["status_code"] = fmt.Sprintf("%d", statusCode)
	enc.Encode(m)

	b, _ := json.Marshal(m)
	log.Printf("%s", string(b))
}

func New(cfg Config) (*Server, error) {
	return &Server{cfg: cfg}, nil
	// defer s.Close()
}

func (s *Server) Serve() error {

	s.router = chi.NewRouter()
	s.router.MethodFunc("POST", "/eventwebhook", s.handleEventWebhook)
	s.server = &http.Server{Addr: fmt.Sprintf(":%d", s.cfg.Port), Handler: panicMW(s.router)}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return err
	}
	if err := s.server.Serve(listener); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}

func panicMW(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println(rec, string(debug.Stack()))
			}
		}()
		h.ServeHTTP(w, r)
	})
}
