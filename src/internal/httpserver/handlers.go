package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	quoteService "github.com/BernsteinMondy/quote-service/src/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func mapHandlers(router *mux.Router, service QuoteService) {
	// The current endpoint structure follows the technical requirements, but in production
	// systems it's strongly recommended to implement API versioning from the start.
	//
	// Good versioning practice example:
	//   /api/v1/quotes
	//   /api/v1/quotes/{id}
	//
	// Key benefits of versioning:
	// 1. Backward compatibility - old clients continue working when API evolves
	// 2. Clear migration path - clients can upgrade to new versions gradually
	//
	// Current implementation processes all requests directly, but adding caching
	// (using Redis, Memcached or in-memory cache) could significantly improve performance.
	quotesGroup := router.PathPrefix("/quotes").Subrouter()
	quotesGroup.Handle("", PostQuoteHandler(service)).Methods("POST")
	quotesGroup.Handle("", GetQuotesHandler(service)).Methods("GET")
	quotesGroup.Handle("/random", GetRandomQuoteHandler(service)).Methods("GET")
	quotesGroup.Handle("/{id}", DeleteQuoteHandler(service)).Methods("DELETE")
}

type QuoteService interface {
	CreateNewQuote(ctx context.Context, author, quote string) error
	GetQuotesWithFilter(ctx context.Context, author string) ([]quoteService.Quote, error)
	GetRandomQuote(ctx context.Context) (*quoteService.Quote, error)
	DeleteQuoteByID(ctx context.Context, id uuid.UUID) error
}

func PostQuoteHandler(service QuoteService) http.HandlerFunc {
	type request = quoteCreateDTO
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			return
		}

		if len(req.Quote) == 0 {
			http.Error(w, "\"quote\" request field can not be empty", http.StatusBadRequest)
		}
		if len(req.Author) == 0 {
			http.Error(w, "\"author\" request field can not be empty", http.StatusBadRequest)
		}

		err = service.CreateNewQuote(r.Context(), req.Author, req.Quote)
		if err != nil {
			if errors.Is(err, quoteService.ErrAlreadyExists) {
				w.WriteHeader(http.StatusConflict)
				return
			}

			http.Error(w, "service: create new quote", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func GetQuotesHandler(service QuoteService) http.HandlerFunc {
	type response struct {
		Quotes []quoteReadDTO `json:"quotes"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		authorFilter := r.URL.Query().Get("author")

		quotes, err := service.GetQuotesWithFilter(r.Context(), authorFilter)
		if err != nil {
			http.Error(w, "service: get quotes: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := response{
			Quotes: make([]quoteReadDTO, len(quotes)),
		}

		for i, quote := range quotes {
			resp.Quotes[i] = quoteFromDomainToReadDTO(&quote)
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetRandomQuoteHandler(service QuoteService) http.HandlerFunc {
	type response = quoteReadDTO
	return func(w http.ResponseWriter, r *http.Request) {
		quote, err := service.GetRandomQuote(r.Context())
		if err != nil {
			http.Error(w, "service: get random quote: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var resp response
		resp = quoteFromDomainToReadDTO(quote)

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func DeleteQuoteHandler(service QuoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, ok := mux.Vars(r)["id"]
		if !ok {
			http.Error(w, "empty \"id\" parameter", http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "invalid \"id\" parameter", http.StatusBadRequest)
			return
		}

		err = service.DeleteQuoteByID(r.Context(), id)
		if err != nil {
			http.Error(w, "service: delete quote by id: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
