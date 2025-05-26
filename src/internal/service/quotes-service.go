package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrAlreadyExists = errors.New("already exists")

	ErrRepoAlreadyExists = errors.New("repository: already exists")
)

type QuoteRepository interface {
	// CreateNewQuote must return ErrRepoAlreadyExists if the quote already exists.
	CreateNewQuote(ctx context.Context, quote *Quote) error
	DeleteQuoteByID(ctx context.Context, id uuid.UUID) error
	GetQuotesWithFilter(ctx context.Context, authorFilter string) ([]Quote, error)
	GetRandomQuote(ctx context.Context) (*Quote, error)
}

type Service struct {
	QuoteRepository QuoteRepository
}

type Quote struct {
	ID     uuid.UUID
	Author string
	Quote  string
}

func (s *Service) CreateNewQuote(ctx context.Context, author, quoteText string) error {
	quote := &Quote{
		ID:     uuid.New(),
		Author: author,
		Quote:  quoteText,
	}

	err := s.QuoteRepository.CreateNewQuote(ctx, quote)
	if err != nil {
		if errors.Is(err, ErrRepoAlreadyExists) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("quote repository: create new quote: %w", err)
	}

	return nil
}

func (s *Service) DeleteQuoteByID(ctx context.Context, id uuid.UUID) error {
	err := s.QuoteRepository.DeleteQuoteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("quote repository: delete quote by id: %w", err)
	}

	return nil
}

func (s *Service) GetQuotesWithFilter(ctx context.Context, authorFilter string) ([]Quote, error) {
	quotes, err := s.QuoteRepository.GetQuotesWithFilter(ctx, authorFilter)
	if err != nil {
		return nil, fmt.Errorf("quote repository: get quotes with filter: %w", err)
	}

	return quotes, nil
}

func (s *Service) GetRandomQuote(ctx context.Context) (*Quote, error) {
	quote, err := s.QuoteRepository.GetRandomQuote(ctx)
	if err != nil {
		return nil, fmt.Errorf("quote repository: get random quote: %w", err)
	}

	return quote, nil
}

func New(quoteRepo QuoteRepository) *Service {
	return &Service{
		QuoteRepository: quoteRepo,
	}
}
