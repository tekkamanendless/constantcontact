package constantcontact

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type GetCustomFieldsInput struct {
	_ string `url:"/v3/contact_custom_fields"`

	Limit int `url:"limit"`
}

type GetCustomFieldsOutput struct {
	HasLinks

	CustomFields []GetCustomFieldsOutputCustomField `json:"custom_fields"`
}

type GetCustomFieldsOutputCustomField struct {
	CustomFieldID string `json:"custom_field_id"`
	Label         string `json:"label"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	// TODO: "metadata"
	Version int `json:"version"`
	// TODO: "choices"
}

func (c *Client) GetCustomFields(ctx context.Context, input GetCustomFieldsInput) (GetCustomFieldsOutput, error) {
	var output GetCustomFieldsOutput

	queryParams := url.Values{}
	if input.Limit > 0 {
		queryParams.Add("limit", strconv.Itoa(input.Limit))
	}

	err := c.restClient.Do(ctx, http.MethodGet, "/v3/contact_custom_fields?"+queryParams.Encode(), nil, &output)
	if err != nil {
		return output, fmt.Errorf("failed to get custom fields: %w", err)
	}

	// TODO: HANDLE PAGINATION

	return output, nil
}
