package testhelpers

import (
	"github.com/BernsteinMondy/quote-service/src/internal/service"
	"github.com/google/uuid"
)

var QuotesArrayFixture = []service.Quote{
	{
		ID:     uuid.MustParse("d45cd206-6495-414c-ab1d-f0b6468264be"),
		Author: "author-1",
		Quote:  "quote-1",
	},
	{
		ID:     uuid.MustParse("f48a5cda-ed11-4403-acaf-a770c05a9d6f"),
		Author: "author-2",
		Quote:  "quote-2",
	},
}
