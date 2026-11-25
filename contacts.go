package constantcontact

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type GetContactsInput struct {
	_ string `url:"/v3/contacts"`

	CreatedAfter  time.Time          `query:"created_after"`
	CreatedBefore time.Time          `query:"created_before"`
	Email         string             `query:"email"`
	Include       []ContactFieldName `query:"include"` // Comma-separated list of fields to include.
	IncludeCount  bool               `query:"include_count"`
	Limit         int                `query:"limit"`
	Lists         []string           `query:"lists"` // Comma-separated list of list IDs.
	OptoutBefore  time.Time          `query:"optout_before"`
	SegmentID     string             `query:"segment_id"`
	SMSStatus     []string           `query:"sms_status"` // Comma-separated list of SMS statuses.
	Status        ContactStatus      `query:"status"`
	Tags          []string           `query:"tags"` // Comma-separated list of tag IDs.
	UpdatedAfter  time.Time          `query:"updated_after"`
	UpdatedBefore time.Time          `query:"updated_before"`
}

type GetContactInput struct {
	_ string `url:"/v3/contacts/{contact_id}"`

	ContactID string `path:"contact_id"`

	Include []ContactFieldName `query:"include"` // Comma-separated list of fields to include.
}

type PutContactInput Contact

type ContactFieldName string

const (
	ContactFieldNameCustomFields    ContactFieldName = "custom_fields"
	ContactFieldNameListMemberships ContactFieldName = "list_memberships"
	ContactFieldNamePhoneNumbers    ContactFieldName = "phone_numbers"
	ContactFieldNameStreetAddresses ContactFieldName = "street_addresses"
	ContactFieldNameTaggings        ContactFieldName = "taggings"
	ContactFieldNameNotes           ContactFieldName = "notes"
)

type ContactStatus string

const (
	ContactStatusAll                 ContactStatus = "all"
	ContactStatusActive              ContactStatus = "active"
	ContactStatusDeleted             ContactStatus = "deleted"
	ContactStatusNotSet              ContactStatus = "not_set"
	ContactStatusPendingConfirmation ContactStatus = "pending_confirmation"
	ContactStatusTempHold            ContactStatus = "temp_hold"
	ContactStatusUnsubscribed        ContactStatus = "unsubscribed"
)

type GetContactsOutput struct {
	HasLinks

	Contacts      []Contact `json:"contacts"`
	ContactsCount int       `json:"contacts_count"`
	Status        string    `json:"status"`
}

type GetContactOutput Contact

type PutContactOutput Contact

type Contact struct {
	ContactID       string               `json:"contact_id"`
	EmailAddress    ContactEmailAddress  `json:"email_address"`
	FirstName       string               `json:"first_name,omitempty"`
	LastName        string               `json:"last_name,omitempty"`
	JobTitle        string               `json:"job_title,omitempty"`
	CompanyName     string               `json:"company_name,omitempty"`
	BirthdayMonth   int                  `json:"birthday_month,omitempty"`
	BirthdayDay     int                  `json:"birthday_day,omitempty"`
	Anniversary     string               `json:"anniversary,omitempty"` // TODO: time.Time Just a date?
	UpdateSource    string               `json:"update_source,omitempty"`
	CreateSource    string               `json:"create_source,omitempty"`
	CreatedAt       string               `json:"created_at,omitempty"` // TODO: time.Time
	UpdatedAt       string               `json:"updated_at,omitempty"` // TODO: time.Time
	DeletedAt       string               `json:"deleted_at,omitempty"` // TODO: time.Time Just a date?
	CustomFields    []ContactCustomField `json:"custom_fields,omitempty"`
	PhoneNumbers    []ContactPhoneNumber `json:"phone_numbers,omitempty"`
	StreetAddresses []any                `json:"street_addresses,omitempty"` // TODO: THIS
	ListMemberships []string             `json:"list_memberships,omitempty"`
	Taggings        []string             `json:"taggings,omitempty"`
	Notes           []struct {
		NoteID    string `json:"note_id"`
		CreatedAt string `json:"created_at"` // TODO: time.Time
		Content   string `json:"content"`
	} `json:"notes,omitempty"`
	SMSChannel any `json:"sms_channel"` // TODO: THIS
}

type ContactEmailAddress struct {
	Address          string `json:"address"`
	PermissionToSend string `json:"permission_to_send"`
	CreatedAt        string `json:"created_at"` // TODO: time.Time
	UpdatedAt        string `json:"updated_at"` // TODO: time.Time
	OptInSource      string `json:"opt_in_source,omitempty"`
	OptInDate        string `json:"opt_in_date,omitempty"` // TODO: time.Time
	OptOutSource     string `json:"opt_out_source,omitempty"`
	OptOutDate       string `json:"opt_out_date,omitempty"` // TODO: time.Time
	OptOutReason     string `json:"opt_out_reason,omitempty"`
	ConfirmStatus    string `json:"confirm_status,omitempty"`
}

type ContactCustomField struct {
	CustomFieldID string   `json:"custom_field_id"`
	Value         string   `json:"value"`
	ChoiceIDs     []string `json:"choice_ids"`
}

type ContactPhoneNumber struct {
	PhoneNumberID string `json:"phone_number_id"`
	PhoneNumber   string `json:"phone_number"`
	Kind          string `json:"kind"`
	CreatedAt     string `json:"created_at"` // TODO: time.Time
	UpdatedAt     string `json:"updated_at"` // TODO: time.Time
	UpdateSource  string `json:"update_source"`
	CreateSource  string `json:"create_source"`
}

func (c *Client) GetContacts(ctx context.Context, input GetContactsInput) (GetContactsOutput, error) {
	var output GetContactsOutput

	queryParams := url.Values{}
	if !input.CreatedAfter.IsZero() {
		queryParams.Add("created_after", input.CreatedAfter.UTC().Format(time.RFC3339))
	}
	if !input.CreatedBefore.IsZero() {
		queryParams.Add("created_before", input.CreatedBefore.UTC().Format(time.RFC3339))
	}
	if input.Email != "" {
		queryParams.Add("email", input.Email)
	}
	if len(input.Include) > 0 {
		var stringValues []string
		for _, value := range input.Include {
			stringValues = append(stringValues, string(value))
		}
		queryParams.Add("include", strings.Join(stringValues, ","))
	}
	if input.IncludeCount {
		queryParams.Add("include_count", "true")
	}
	if input.Limit > 0 {
		queryParams.Add("limit", strconv.Itoa(input.Limit))
	}
	if len(input.Lists) > 0 {
		queryParams.Add("lists", strings.Join(input.Lists, ","))
	}
	if !input.OptoutBefore.IsZero() {
		queryParams.Add("optout_before", input.OptoutBefore.UTC().Format(time.RFC3339))
	}
	if input.SegmentID != "" {
		queryParams.Add("segment_id", input.SegmentID)
	}
	if len(input.SMSStatus) > 0 {
		queryParams.Add("sms_status", strings.Join(input.SMSStatus, ","))
	}
	if input.Status != "" {
		queryParams.Add("status", string(input.Status))
	}
	if len(input.Tags) > 0 {
		queryParams.Add("tags", strings.Join(input.Tags, ","))
	}
	if !input.UpdatedAfter.IsZero() {
		queryParams.Add("updated_after", input.UpdatedAfter.UTC().Format(time.RFC3339))
	}
	if !input.UpdatedBefore.IsZero() {
		queryParams.Add("updated_before", input.UpdatedBefore.UTC().Format(time.RFC3339))
	}

	targetURL := "/v3/contacts?" + queryParams.Encode()

	err := c.restClient.Do(ctx, http.MethodGet, targetURL, nil, &output)
	if err != nil {
		return output, fmt.Errorf("failed to get contacts: %w", err)
	}

	// Handle pagination.
	if output.Links != nil && output.Links["next"].Href != "" {
		link := output.Links["next"].Href
		for link != "" {
			var pageOutput GetContactsOutput
			err := c.restClient.Do(ctx, http.MethodGet, link, nil, &pageOutput)
			if err != nil {
				return pageOutput, fmt.Errorf("failed to get page: %w", err)
			}

			output.Contacts = append(output.Contacts, pageOutput.Contacts...)

			link = ""
			if pageOutput.Links != nil && pageOutput.Links["next"].Href != "" {
				link = pageOutput.Links["next"].Href
			}
		}
	}

	return output, nil
}

func (c *Client) GetContact(ctx context.Context, input GetContactInput) (GetContactOutput, error) {
	var output GetContactOutput

	queryParams := url.Values{}
	if len(input.Include) > 0 {
		var stringValues []string
		for _, value := range input.Include {
			stringValues = append(stringValues, string(value))
		}
		queryParams.Add("include", strings.Join(stringValues, ","))
	}

	targetURL := "/v3/contacts/" + url.PathEscape(input.ContactID) + "?" + queryParams.Encode()

	err := c.restClient.Do(ctx, http.MethodGet, targetURL, nil, &output)
	if err != nil {
		return output, fmt.Errorf("failed to get contact: %w", err)
	}

	return output, nil
}

func (c *Client) PutContact(ctx context.Context, input PutContactInput) (PutContactOutput, error) {
	var output PutContactOutput

	targetURL := "/v3/contacts/" + url.PathEscape(input.ContactID)

	err := c.restClient.Do(ctx, http.MethodPut, targetURL, input, &output)
	if err != nil {
		return output, fmt.Errorf("failed to put contact: %w", err)
	}

	return output, nil
}
