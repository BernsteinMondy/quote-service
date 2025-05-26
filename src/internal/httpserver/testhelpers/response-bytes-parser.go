package testhelpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func ParseResponseBody[T any](resp *http.Response) (_ T, err error) {
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close response body: %w", closeErr))
		}
	}()

	var ret T
	err = json.NewDecoder(resp.Body).Decode(&ret)

	return ret, err
}
