package wrap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/tekkamanendless/restapiclient"
)

type RequestInput struct {
	Method      string
	Path        string
	QueryParams url.Values
	Headers     http.Header
	Body        []byte
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
		Method:      "",
		Path:        "",
		QueryParams: url.Values{},
		Headers:     http.Header{},
		Body:        nil,
	}
	supportedTags := []string{"form", "json"}

	inputValue := reflect.ValueOf(input)
	for inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}
	if inputValue.Kind() != reflect.Struct {
		return requestInput, fmt.Errorf("input must be a struct")
	}

	tagFieldsMap := map[string][]reflect.StructField{}
	var bodyInterface any

	// First pass: parse all the "api" tagged fields.
	for _, structField := range reflect.VisibleFields(inputValue.Type()) {
		fieldValue := inputValue.FieldByName(structField.Name)

		for _, tagName := range supportedTags {
			if _, ok := structField.Tag.Lookup(tagName); ok {
				tagFieldsMap[tagName] = append(tagFieldsMap[tagName], structField)
			}
		}

		apiTag := structField.Tag.Get("api")
		if apiTag == "" {
			continue
		}

		apiTagParts := strings.SplitN(apiTag, ":", 2)

		apiTagName := strings.TrimSpace(apiTagParts[0])
		var apiTagValue string
		if len(apiTagParts) > 1 {
			apiTagValue = strings.TrimSpace(apiTagParts[1])
		}

		switch apiTagName {
		case "httpmethod":
			requestInput.Method = strings.ToUpper(apiTagValue)
		case "httppath":
			requestInput.Path = apiTagValue
		case "query":
			queryParameterParts := strings.SplitN(apiTagValue, ",", 2)
			queryParameterName := strings.TrimSpace(queryParameterParts[0])
			if len(queryParameterParts) > 1 {
				_ = queryParameterParts[1] // TODO: any parameters
			}

			stringValue := reflectValueToString(fieldValue)
			requestInput.QueryParams.Set(queryParameterName, stringValue)
		case "header":
			headerParts := strings.SplitN(apiTagValue, ",", 2)
			headerName := strings.TrimSpace(headerParts[0])
			if len(headerParts) > 1 {
				_ = headerParts[1] // TODO: any parameters
			}

			stringValue := reflectValueToString(fieldValue)
			requestInput.Headers.Set(headerName, stringValue)
		case "body":
			bodyInterface = fieldValue.Interface()
		}
	}

	// Second pass: apply any logic.
	if len(tagFieldsMap["form"]) > 0 {
		if requestInput.Method == "" {
			requestInput.Method = http.MethodPost
		}
		if requestInput.Headers.Get("Content-Type") == "" {
			requestInput.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
		}

		form := url.Values{}
		for _, structField := range tagFieldsMap["form"] {
			fieldValue := inputValue.FieldByName(structField.Name)

			tag := structField.Tag.Get("form")
			tagParts := strings.SplitN(tag, ",", 2)

			tagName := strings.TrimSpace(tagParts[0])
			var tagValue string
			tagExtraMap := map[string]bool{}
			if len(tagParts) > 1 {
				tagValue = tagParts[1]

				tagExtras := strings.Split(tagValue, ",")
				for i, tagExtra := range tagExtras {
					tagExtra = strings.TrimSpace(tagExtra)
					if len(tagExtra) == 0 {
						return requestInput, fmt.Errorf("invalid form tag syntax: %s", tagValue)
					}
					tagExtras[i] = tagExtra
				}

				for _, tagExtra := range tagExtras {
					tagExtraMap[tagExtra] = true
				}
			}

			if tagName == "-" {
				continue
			}

			skipField := false
			for tagExtra := range tagExtraMap {
				switch tagExtra {
				case "omitempty":
					if fieldValue.IsZero() {
						skipField = true
					}
				default:
					return requestInput, fmt.Errorf("unsupported form tag extra: %s", tagExtra)
				}
			}
			if skipField {
				continue
			}

			form.Set(tagName, reflectValueToString(fieldValue))
		}
		requestInput.Body = []byte(form.Encode())
	} else if len(tagFieldsMap["json"]) > 0 {
		if requestInput.Method == "" {
			requestInput.Method = http.MethodPost
		}
		if requestInput.Headers.Get("Content-Type") == "" {
			requestInput.Headers.Set("Content-Type", "application/json")
		}

		bodyStruct := map[string]any{}
		for _, structField := range tagFieldsMap["json"] {
			fieldValue := inputValue.FieldByName(structField.Name)

			tag := structField.Tag.Get("json")
			tagParts := strings.SplitN(tag, ",", 2)

			tagName := strings.TrimSpace(tagParts[0])
			var tagValue string
			if len(tagParts) > 1 {
				tagValue = tagParts[1]
			}
			if len(tagValue) > 0 {
				return requestInput, fmt.Errorf("unsupported json tag extras: %s", tagValue)
			}

			if tagName == "-" {
				continue
			}

			bodyStruct[tagName] = fieldValue.Interface()
		}
		jsonBody, err := json.Marshal(bodyStruct)
		if err != nil {
			return requestInput, fmt.Errorf("failed to marshal body: %w", err)
		}
		requestInput.Body = jsonBody
	} else {
		if bodyInterface != nil {
			bodyValue := reflect.ValueOf(bodyInterface)

			if requestInput.Method == "" {
				requestInput.Method = http.MethodPost
			}
			if requestInput.Headers.Get("Content-Type") == "" {
				if v, ok := reflect.TypeAssert[[]byte](bodyValue); ok {
					requestInput.Headers.Set("Content-Type", "application/octet-stream")
					requestInput.Body = v
				} else if v, ok := reflect.TypeAssert[string](bodyValue); ok {
					requestInput.Headers.Set("Content-Type", "text/plain")
					requestInput.Body = []byte(v)
				} else if v, ok := reflect.TypeAssert[fmt.Stringer](bodyValue); ok {
					requestInput.Headers.Set("Content-Type", "string")
					requestInput.Body = []byte(v.String())
				}
			} else {
				switch requestInput.Headers.Get("Content-Type") {
				case "application/json":
					v, err := json.Marshal(bodyInterface)
					if err != nil {
						return requestInput, fmt.Errorf("failed to marshal body: %w", err)
					}
					requestInput.Body = v
				default:
					return requestInput, fmt.Errorf("unsupported content type: %s", requestInput.Headers.Get("Content-Type"))
				}
			}
		}
	}

	if requestInput.Method == "" {
		requestInput.Method = http.MethodGet
	}

	return requestInput, nil
}

func reflectValueToString(value reflect.Value) string {
	if v, ok := reflect.TypeAssert[string](value); ok {
		return v
	} else if v, ok := reflect.TypeAssert[int](value); ok {
		return strconv.Itoa(v)
	} else if v, ok := reflect.TypeAssert[int64](value); ok {
		return strconv.FormatInt(v, 10)
	} else if v, ok := reflect.TypeAssert[float64](value); ok {
		return strconv.FormatFloat(v, 'f', -1, 64)
	} else if v, ok := reflect.TypeAssert[fmt.Stringer](value); ok {
		return v.String()
	}

	return ""
}
