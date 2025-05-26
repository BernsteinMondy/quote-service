package httpserver_test

import (
	"bytes"
	"errors"
	"github.com/BernsteinMondy/quote-service/src/internal/httpserver"
	"github.com/BernsteinMondy/quote-service/src/internal/httpserver/testhelpers"
	"github.com/BernsteinMondy/quote-service/src/internal/service"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDeleteQuoteHandler(t *testing.T) {
	type testCase struct {
		name               string
		service            httpserver.QuoteService
		wantRespStatusCode int
		quoteID            string
	}

	testCases := []testCase{
		{
			name:               "Smoke test",
			service:            &testhelpers.MockQuoteService{},
			wantRespStatusCode: http.StatusOK,
			quoteID:            "4937a248-cb08-46de-8789-493904914cc6",
		},
		{
			name:               "Empty quote ID results in status code 404",
			service:            &testhelpers.MockQuoteService{},
			wantRespStatusCode: http.StatusNotFound,
			quoteID:            "",
		},
		{
			name:               "Non-uuid quote id results in status code 400",
			service:            &testhelpers.MockQuoteService{},
			wantRespStatusCode: http.StatusBadRequest,
			quoteID:            "non-uuid",
		},
		{
			name:               "Service call ended with error results in status code 500",
			service:            &testhelpers.MockQuoteService{RetError: errors.New("some error")},
			wantRespStatusCode: http.StatusInternalServerError,
			quoteID:            "4937a248-cb08-46de-8789-493904914cc6",
		},
	}

	for _, tc := range testCases {
		router := mux.NewRouter()
		router.Handle("/{id}", httpserver.DeleteQuoteHandler(tc.service)).Methods("DELETE")

		server := httptest.NewServer(router)

		url := server.URL + "/" + tc.quoteID
		req, err := http.NewRequest(http.MethodDelete, url, http.NoBody)
		if err != nil {
			t.Fatal("Failed to create request", err)
		}

		resp, err := server.Client().Do(req)
		if err != nil {
			t.Fatal("Failed to make test request", err)
		}

		if resp.StatusCode != tc.wantRespStatusCode {
			t.Errorf("DeleteQuoteHandler returned wrong status code: got %d want %d", resp.StatusCode, tc.wantRespStatusCode)
		}
	}
}

func TestPostQuoteHandler(t *testing.T) {
	type testCase struct {
		name               string
		service            httpserver.QuoteService
		wantRespStatusCode int
		body               []byte
	}

	testCases := []testCase{
		{
			name:               "Smoke test",
			service:            &testhelpers.MockQuoteService{},
			wantRespStatusCode: http.StatusCreated,
			body:               []byte(`{"author":"test author","quote":"test quote"}`),
		},
		{
			name:               "Invalid request body results in status code 400",
			service:            &testhelpers.MockQuoteService{},
			wantRespStatusCode: http.StatusBadRequest,
			body:               []byte(`????WHAT????{"author":"test author","quote":"test quote"}`),
		},
		{
			name:               "Empty \"author\" request field results in status code 400",
			service:            &testhelpers.MockQuoteService{},
			wantRespStatusCode: http.StatusBadRequest,
			body:               []byte(`{"author":"","quote":"test quote"}`),
		},
		{
			name:               "Empty \"quote\" request field results in status code 400",
			service:            &testhelpers.MockQuoteService{},
			wantRespStatusCode: http.StatusBadRequest,
			body:               []byte(`{"author":"test author","quote":""}`),
		},
		{
			name:               "service.ErrAlreadyExists error returned from Service results in status code 409",
			service:            &testhelpers.MockQuoteService{RetError: service.ErrAlreadyExists},
			wantRespStatusCode: http.StatusConflict,
			body:               []byte(`{"author":"test author","quote":"test quote"}`),
		},
		{
			name:               "Service call ended with error results in status code 500",
			service:            &testhelpers.MockQuoteService{RetError: errors.New("some error")},
			wantRespStatusCode: http.StatusInternalServerError,
			body:               []byte(`{"author":"test author","quote":"test quote"}`),
		},
	}

	for _, tc := range testCases {
		router := mux.NewRouter()
		router.Handle("/", httpserver.PostQuoteHandler(tc.service)).Methods("POST")

		server := httptest.NewServer(router)

		req, err := http.NewRequest(http.MethodPost, server.URL+"/", bytes.NewReader(tc.body))
		if err != nil {
			t.Fatal("Failed to create request", err)
		}

		resp, err := server.Client().Do(req)
		if err != nil {
			t.Fatal("Failed to make test request", err)
		}

		if resp.StatusCode != tc.wantRespStatusCode {
			t.Errorf("PostQuoteHandler returned wrong status code: got %d want %d", resp.StatusCode, tc.wantRespStatusCode)
		}
	}
}

func TestGetQuotesHandler(t *testing.T) {
	type (
		quote struct {
			ID     string `json:"id"`
			Author string `json:"author"`
			Quote  string `json:"quote"`
		}
		response struct {
			Quotes []quote `json:"quotes"`
		}
	)

	type testCase struct {
		name               string
		service            httpserver.QuoteService
		wantRespStatusCode int
		queryParams        string
		wantRespBody       *response
	}

	testCases := []testCase{
		{
			name:               "Smoke test",
			service:            &testhelpers.MockQuoteService{},
			wantRespStatusCode: http.StatusOK,
			wantRespBody: &response{
				Quotes: []quote{
					{
						ID:     testhelpers.QuotesArrayFixture[0].ID.String(),
						Author: testhelpers.QuotesArrayFixture[0].Author,
						Quote:  testhelpers.QuotesArrayFixture[0].Quote,
					},
					{
						ID:     testhelpers.QuotesArrayFixture[1].ID.String(),
						Author: testhelpers.QuotesArrayFixture[1].Author,
						Quote:  testhelpers.QuotesArrayFixture[1].Quote,
					},
				},
			},
		},
		{
			name:               "Specified \"author\" query paramater results in status code 200",
			service:            &testhelpers.MockQuoteService{},
			wantRespStatusCode: http.StatusOK,
			queryParams:        "author=some-author",
			wantRespBody: &response{
				Quotes: []quote{
					{
						ID:     testhelpers.QuotesArrayFixture[0].ID.String(),
						Author: testhelpers.QuotesArrayFixture[0].Author,
						Quote:  testhelpers.QuotesArrayFixture[0].Quote,
					},
					{
						ID:     testhelpers.QuotesArrayFixture[1].ID.String(),
						Author: testhelpers.QuotesArrayFixture[1].Author,
						Quote:  testhelpers.QuotesArrayFixture[1].Quote,
					},
				},
			},
		},
		{
			name:               "Service call ended with error results in status code 500",
			service:            &testhelpers.MockQuoteService{RetError: errors.New("some error")},
			wantRespStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		router := mux.NewRouter()
		router.Handle("/", httpserver.GetQuotesHandler(tc.service)).Methods("GET")

		server := httptest.NewServer(router)

		req, err := http.NewRequest(http.MethodGet, server.URL+"/?"+tc.queryParams, http.NoBody)
		if err != nil {
			t.Fatal("Failed to create request", err)
		}

		resp, err := server.Client().Do(req)
		if err != nil {
			t.Fatal("Failed to make test request", err)
		}

		if resp.StatusCode != tc.wantRespStatusCode {
			t.Errorf("PostQuoteHandler returned wrong status code: got %d want %d", resp.StatusCode, tc.wantRespStatusCode)
		}
		if tc.wantRespBody != nil {
			gotResp, err := testhelpers.ParseResponseBody[response](resp)
			if err != nil {
				t.Fatalf("Error parsing response body: %v", err)
			}
			if !reflect.DeepEqual(*tc.wantRespBody, gotResp) {
				t.Fatalf("Did not get desired response body: got %v want %v", gotResp, *tc.wantRespBody)
			}
		}
	}

}

func TestGetRandomQuoteHandler(t *testing.T) {
	type (
		response struct {
			ID     string `json:"id"`
			Author string `json:"author"`
			Quote  string `json:"quote"`
		}
	)

	type testCase struct {
		name               string
		service            httpserver.QuoteService
		wantRespStatusCode int
		wantRespBody       *response
	}

	testCases := []testCase{
		{
			name:               "Smoke test",
			service:            &testhelpers.MockQuoteService{},
			wantRespStatusCode: http.StatusOK,
			wantRespBody: &response{
				ID:     testhelpers.QuotesArrayFixture[0].ID.String(),
				Author: testhelpers.QuotesArrayFixture[0].Author,
				Quote:  testhelpers.QuotesArrayFixture[0].Quote,
			},
		},
		{
			name:               "Service call ended with error results in status code 500",
			service:            &testhelpers.MockQuoteService{RetError: errors.New("some error")},
			wantRespStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		router := mux.NewRouter()
		router.Handle("/", httpserver.GetRandomQuoteHandler(tc.service)).Methods("GET")

		server := httptest.NewServer(router)

		req, err := http.NewRequest(http.MethodGet, server.URL+"/", http.NoBody)
		if err != nil {
			t.Fatal("Failed to create request", err)
		}

		resp, err := server.Client().Do(req)
		if err != nil {
			t.Fatal("Failed to make test request", err)
		}

		if resp.StatusCode != tc.wantRespStatusCode {
			t.Errorf("PostQuoteHandler returned wrong status code: got %d want %d", resp.StatusCode, tc.wantRespStatusCode)
		}
		if tc.wantRespBody != nil {
			gotResp, err := testhelpers.ParseResponseBody[response](resp)
			if err != nil {
				t.Fatalf("Error parsing response body: %v", err)
			}
			if !reflect.DeepEqual(*tc.wantRespBody, gotResp) {
				t.Fatalf("Did not get desired response body: got %v want %v", gotResp, *tc.wantRespBody)
			}
		}
	}
}
