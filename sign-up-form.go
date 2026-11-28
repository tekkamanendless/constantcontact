package constantcontact

import (
	"context"

	"github.com/tekkamanendless/constantcontact/wrap"
)

type PostSignUpFormInput struct {
	_ string `api:"httppath:/v3/contacts/sign_up_form"`
	_ string `api:"httpmethod:POST"`

	EmailAddress    *string                      `json:"email_address,omitempty"`
	FirstName       *string                      `json:"first_name,omitempty"`
	LastName        *string                      `json:"last_name,omitempty"`
	JobTitle        *string                      `json:"job_title,omitempty"`
	CompanyName     *string                      `json:"company_name,omitempty"`
	PhoneNumber     *string                      `json:"phone_number,omitempty"`
	ListMemberships []string                     `json:"list_memberships,omitempty"`
	CustomFields    []SignUpFormInputCustomField `json:"custom_fields,omitempty"`
	// TODO: "anniversary": "11-15-2006"`
	BirthdayMonth *int `json:"birthday_month,omitempty"`
	BirthdayDay   *int `json:"birthday_day,omitempty"`
	// TODO: "street_address": {},
	// TODO: "sms_channel": {}
}

type SignUpFormInputCustomField struct {
	CustomFieldID string `json:"custom_field_id"`
	Value         string `json:"value"`
}

type PostSignUpFormOutput struct {
	ContactID string `json:"contact_id"`
	Action    string `json:"action"`
}

func (c *Client) PostSignUpForm(ctx context.Context, input PostSignUpFormInput) (PostSignUpFormOutput, error) {
	var output PostSignUpFormOutput
	err := wrap.WrapRequest(ctx, c.restClient, input, &output)
	return output, err
}
