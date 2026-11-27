package wrap

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/tekkamanendless/restapiclient"
)

type RequestInput struct {
	Method      string
	Path        string
	QueryParams url.Values
	Headers     http.Header
	Body        any
}

func WrapRequest(ctx context.Context, client restapiclient.Client, input any, output any) error {
	requestInput, err := parseRequestInput(input)
	if err != nil {
		return err
	}

	return doRequest(ctx, client, requestInput, output)
}

func doRequest(ctx context.Context, client restapiclient.Client, requestInput RequestInput, output any) error {
	var options []restapiclient.Option
	for name, values := range requestInput.Headers {
		options = append(options, restapiclient.OptionHeader(name, values...))
	}
	for name, values := range requestInput.QueryParams {
		options = append(options, restapiclient.OptionQuery(name, values...))
	}
	err := client.Do(ctx, requestInput.Method, requestInput.Path, requestInput.Body, output, options...)
	if err != nil {
		return err
	}
	return nil
}

func parseRequestInput(input any) (RequestInput, error) {
	requestInput := RequestInput{
		Method:      http.MethodGet,
		Path:        "",
		QueryParams: url.Values{},
		Headers:     http.Header{},
		Body:        nil,
	}

	inputValue := reflect.ValueOf(input)
	for inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}
	if inputValue.Kind() != reflect.Struct {
		return requestInput, fmt.Errorf("input must be a struct")
	}

	return requestInput, nil
}
