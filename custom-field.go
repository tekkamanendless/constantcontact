package constantcontact

import (
	"context"
	"fmt"

	"github.com/tekkamanendless/constantcontact/wrap"
)

type GetCustomFieldsInput struct {
	_ string `api:"httppath:/v3/contact_custom_fields"`
	_ string `api:"httpmethod:GET"`

	Limit int `api:"query:limit"`
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
	err := wrap.WrapRequest(ctx, c.restClient, input, &output)
	if err != nil {
		return output, fmt.Errorf("failed to get custom fields: %w", err)
	}

	// TODO: HANDLE PAGINATION

	return output, nil
}
