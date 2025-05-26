package testhelpers

import (
	"context"
	"github.com/MaksKazantsev/quote-service/src/internal/httpserver"
	"github.com/MaksKazantsev/quote-service/src/internal/service"
	"github.com/google/uuid"
)

type MockQuoteService struct {
	RetError error
}

var _ httpserver.QuoteService = (*MockQuoteService)(nil)

func (m *MockQuoteService) CreateNewQuote(context.Context, string, string) error {
	return m.RetError
}

func (m *MockQuoteService) GetQuotesWithFilter(context.Context, string) ([]service.Quote, error) {
	if m.RetError != nil {
		return nil, m.RetError
	}
	return QuotesArrayFixture, nil
}

func (m *MockQuoteService) GetRandomQuote(context.Context) (*service.Quote, error) {
	if m.RetError != nil {
		return nil, m.RetError
	}

	return &QuotesArrayFixture[0], nil
}

func (m *MockQuoteService) DeleteQuoteByID(ctx context.Context, id uuid.UUID) error {
	return m.RetError
}
