package wrap

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRequestInput(t *testing.T) {
	t.Run("Trivial", func(t *testing.T) {
		input := struct{}{}
		requestInput, err := parseRequestInput(input)
		require.NoError(t, err)
		require.Equal(t, RequestInput{
			Method:      http.MethodGet,
			QueryParams: url.Values{},
			Headers:     http.Header{},
		}, requestInput)
	})
	t.Run("Method", func(t *testing.T) {
		input := struct {
			_ string `api:"httpmethod:GET"`
		}{}
		requestInput, err := parseRequestInput(input)
		require.NoError(t, err)
		require.Equal(t, RequestInput{
			Method:      http.MethodGet,
			QueryParams: url.Values{},
			Headers:     http.Header{},
		}, requestInput)
	})
	t.Run("Path", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			input := struct {
				_ string `api:"httppath:/v1/contacts"`
			}{}
			requestInput, err := parseRequestInput(input)
			require.NoError(t, err)
			require.Equal(t, RequestInput{
				Method:      http.MethodGet,
				Path:        "/v1/contacts",
				QueryParams: url.Values{},
				Headers:     http.Header{},
			}, requestInput)
		})
		t.Run("Missing param", func(t *testing.T) {
			input := struct {
				_ string `api:"httppath:/v1/contacts/{id}"`
			}{}
			requestInput, err := parseRequestInput(input)
			require.NoError(t, err)
			assert.Equal(t, "/v1/contacts/{id}", requestInput.Path)
		})
		t.Run("Proper params", func(t *testing.T) {
			input := struct {
				_ string `api:"httppath:/v1/contacts/{id}"`

				ID string `api:"path:id"`
			}{
				ID: "123",
			}
			requestInput, err := parseRequestInput(input)
			require.NoError(t, err)
			assert.Equal(t, "/v1/contacts/123", requestInput.Path)
		})
	})
	t.Run("QueryParams", func(t *testing.T) {
		input := struct {
			Limit int    `api:"query:limit"`
			Name  string `api:"query:name"`
		}{
			Limit: 10,
			Name:  "John",
		}
		requestInput, err := parseRequestInput(input)
		require.NoError(t, err)
		require.Equal(t, RequestInput{
			Method: http.MethodGet,
			QueryParams: url.Values{
				"limit": {"10"},
				"name":  {"John"},
			},
			Headers: http.Header{},
		}, requestInput)
	})
	t.Run("Headers", func(t *testing.T) {
		input := struct {
			Authorization string `api:"header:Authorization"`
			ContentType   string `api:"header:Content-Type"`
		}{
			Authorization: "Bearer 1234567890",
			ContentType:   "application/json",
		}
		requestInput, err := parseRequestInput(input)
		require.NoError(t, err)
		require.Equal(t, RequestInput{
			Method:      http.MethodGet,
			QueryParams: url.Values{},
			Headers: http.Header{
				"Authorization": {"Bearer 1234567890"},
				"Content-Type":  {"application/json"},
			},
		}, requestInput)
	})
	t.Run("Body", func(t *testing.T) {
		t.Run("Raw", func(t *testing.T) {
			t.Run("String", func(t *testing.T) {
				input := struct {
					Body string `api:"body"`
				}{
					Body: "Hello, world!",
				}
				requestInput, err := parseRequestInput(input)
				require.NoError(t, err)
				require.Equal(t, RequestInput{
					Method:      http.MethodPost,
					Body:        []byte("Hello, world!"),
					QueryParams: url.Values{},
					Headers: http.Header{
						"Content-Type": {"text/plain"},
					},
				}, requestInput)
			})
			t.Run("Bytes", func(t *testing.T) {
				input := struct {
					Body []byte `api:"body"`
				}{
					Body: []byte("Hello, world!"),
				}
				requestInput, err := parseRequestInput(input)
				require.NoError(t, err)
				require.Equal(t, RequestInput{
					Method:      http.MethodPost,
					Body:        []byte("Hello, world!"),
					QueryParams: url.Values{},
					Headers: http.Header{
						"Content-Type": {"application/octet-stream"},
					},
				}, requestInput)
			})
		})
		t.Run("JSON", func(t *testing.T) {
			input := struct {
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			}{
				FirstName: "John",
				LastName:  "Doe",
			}
			requestInput, err := parseRequestInput(input)
			require.NoError(t, err)
			require.Equal(t, RequestInput{
				Method:      http.MethodPost,
				Body:        []byte(`{"first_name":"John","last_name":"Doe"}`),
				QueryParams: url.Values{},
				Headers: http.Header{
					"Content-Type": {"application/json"},
				},
			}, requestInput)
		})
		t.Run("Form", func(t *testing.T) {
			t.Run("Simple", func(t *testing.T) {
				input := struct {
					FirstName string `form:"first_name"`
					LastName  string `form:"last_name"`
				}{
					FirstName: "John",
					LastName:  "Doe",
				}
				requestInput, err := parseRequestInput(input)
				require.NoError(t, err)
				require.Equal(t, RequestInput{
					Method:      http.MethodPost,
					Body:        []byte("first_name=John&last_name=Doe"),
					QueryParams: url.Values{},
					Headers: http.Header{
						"Content-Type": {"application/x-www-form-urlencoded"},
					},
				}, requestInput)
			})
			t.Run("OmitEmpty", func(t *testing.T) {
				input := struct {
					FirstName  string `form:"first_name"`
					MiddleName string `form:"middle_name,omitempty"`
					LastName   string `form:"last_name"`
				}{
					FirstName:  "John",
					MiddleName: "",
					LastName:   "Doe",
				}
				requestInput, err := parseRequestInput(input)
				require.NoError(t, err)
				require.Equal(t, RequestInput{
					Method:      http.MethodPost,
					Body:        []byte("first_name=John&last_name=Doe"),
					QueryParams: url.Values{},
					Headers: http.Header{
						"Content-Type": {"application/x-www-form-urlencoded"},
					},
				}, requestInput)
			})
		})
	})
	t.Run("Full", func(t *testing.T) {
		input := struct {
			_             string `api:"httpmethod:POST"`
			_             string `api:"httppath:/v1/contacts"`
			Limit         int    `api:"query:limit"`
			Name          string `api:"query:name"`
			Authorization string `api:"header:Authorization"`
			ContentType   string `api:"header:Content-Type"`
			FirstName     string `json:"first_name"`
			LastName      string `json:"last_name"`
		}{
			Limit:         10,
			Name:          "John",
			Authorization: "Bearer 1234567890",
			ContentType:   "application/json",
			FirstName:     "John",
			LastName:      "Doe",
		}
		requestInput, err := parseRequestInput(input)
		require.NoError(t, err)
		require.Equal(t, RequestInput{
			Method: http.MethodPost,
			Path:   "/v1/contacts",
			QueryParams: url.Values{
				"limit": {"10"},
				"name":  {"John"},
			},
			Headers: http.Header{
				"Authorization": {"Bearer 1234567890"},
				"Content-Type":  {"application/json"},
			},
			Body: []byte(`{"first_name":"John","last_name":"Doe"}`),
		}, requestInput)
	})
}
