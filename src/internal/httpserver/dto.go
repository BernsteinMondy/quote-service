package httpserver

import "github.com/BernsteinMondy/quote-service/src/internal/service"

type (
	quoteReadDTO struct {
		ID     string `json:"id"`
		Author string `json:"author"`
		Quote  string `json:"quote"`
	}
	quoteCreateDTO struct {
		Author string `json:"author"`
		Quote  string `json:"quote"`
	}
)

func quoteFromDomainToReadDTO(quote *service.Quote) quoteReadDTO {
	return quoteReadDTO{
		ID:     quote.ID.String(),
		Author: quote.Author,
		Quote:  quote.Quote,
	}
}
