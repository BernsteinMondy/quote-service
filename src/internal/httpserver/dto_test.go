package httpserver

import (
	"github.com/MaksKazantsev/quote-service/src/internal/service"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func Test_quoteFromDomainToReadDTO(t *testing.T) {
	type args struct {
		quote *service.Quote
	}
	tests := []struct {
		name string
		args args
		want quoteReadDTO
	}{
		{
			name: "Valid service.Quote argument should result in valid quoteReadDTO",
			args: args{
				quote: &service.Quote{
					ID:     uuid.MustParse("d45cd206-6495-414c-ab1d-f0b6468264be"),
					Author: "author-1",
					Quote:  "quote-1",
				},
			},
			want: quoteReadDTO{
				ID:     "d45cd206-6495-414c-ab1d-f0b6468264be",
				Author: "author-1",
				Quote:  "quote-1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quoteFromDomainToReadDTO(tt.args.quote); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("quoteFromDomainToReadDTO() = %v, want %v", got, tt.want)
			}
		})
	}
}
