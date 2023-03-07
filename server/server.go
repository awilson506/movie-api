package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	api "github.com/awilson506/movie-api/pkg"
)

// Server - struct to hold a few services needed in the server
type Server struct {
	client *api.Client
	server *http.Server
	mux    *http.ServeMux
}

// NewServer - get an instance of the server
func NewServer() *Server {
	s := &Server{
		client: api.New(),
		mux:    http.NewServeMux(),
	}

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: s.mux,
	}

	s.mux.HandleFunc("/production-company-details", s.getProductionCompanyDetailsHandler)
	s.mux.HandleFunc("/genre-details", s.genreDetailsHandler)

	return s
}

// Start - start the mux service
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// WriteErrorResponse - handle writing back an error response
func (s *Server) WriteErrorResponse(w http.ResponseWriter, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errors)
}

// getProductionCompanyDetailsHandler - send the company details request off to the service
func (s *Server) getProductionCompanyDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		msg := api.ErrorMessageContainer{}

		productionCompanyId := r.URL.Query().Get("production_id")
		year := r.URL.Query().Get("year")
		page := r.URL.Query().Get("page")

		parsedPageId, _ := api.ValidateOptionalStringParam("page", page, &msg)
		parsedId, _ := api.ValidateOptionalStringParam("production_id", productionCompanyId, &msg)
		parsedYear, _ := api.ValidateOptionalStringParam("year", year, &msg)

		if len(msg.Errors) != 0 {
			s.WriteErrorResponse(w, msg.Errors)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.client.GetProductionCompanyDetails(&parsedId, &parsedYear, getPageId(parsedPageId)))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// genreDetailsHandler - send the genre details request off to the service
func (s *Server) genreDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		msg := api.ErrorMessageContainer{}
		page := r.URL.Query().Get("page")
		year := r.URL.Query().Get("year")

		parsedPageId, _ := api.ValidateOptionalStringParam("page", page, &msg)
		parsedYear, _ := api.ValidateOptionalStringParam("year", year, &msg)

		if len(msg.Errors) != 0 {
			s.WriteErrorResponse(w, msg.Errors)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.client.GetGenreDetails(&parsedYear, getPageId(parsedPageId)))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getPageId - convert the pageId to an int pointer for later use if needed
func getPageId(pageId string) *int {
	var pageIdInt int = 1
	if pageId == "" {
		return &pageIdInt
	} else {
		pageIdInt, _ := strconv.Atoi(pageId)
		return &pageIdInt
	}
}
