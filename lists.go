package constantcontact

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type GetListsInput struct {
	_ string `url:"/v3/contact_lists"`

	ChannelType               string `query:"channel_type"`
	IncludeCount              bool   `query:"include_count"`
	IncludeMembershipCount    bool   `query:"include_membership_count"`
	IncludeSMSMembershipCount bool   `query:"include_sms_membership_count"`
	Limit                     int    `query:"limit"`
	Name                      string `query:"name"`
	Status                    string `query:"status"`
}

type GetListsOutput struct {
	HasLinks

	Lists      []GetListsOutputList `json:"lists"`
	ListsCount int                  `json:"lists_count"`
}

type GetListsOutputList struct {
	ListID          string `json:"list_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Favorite        bool   `json:"favorite"`
	CreatedAt       string `json:"created_at"` // TODO: time.Time
	UpdatedAt       string `json:"updated_at"` // TODO: time.Time
	DeletedAt       string `json:"deleted_at"` // TODO: time.Time
	MembershipCount int    `json:"membership_count"`
}

func (c *Client) GetLists(ctx context.Context, input GetListsInput) (GetListsOutput, error) {
	var output GetListsOutput

	queryParams := url.Values{}
	if input.Limit > 0 {
		queryParams.Add("limit", strconv.Itoa(input.Limit))
	}
	if input.IncludeCount {
		queryParams.Add("include_count", "true")
	}
	if input.IncludeMembershipCount {
		queryParams.Add("include_membership_count", "true")
	}
	if input.Name != "" {
		queryParams.Add("name", input.Name)
	}
	if input.Status != "" {
		queryParams.Add("status", input.Status)
	}
	if input.ChannelType != "" {
		queryParams.Add("channel_type", input.ChannelType)
	}
	if input.IncludeSMSMembershipCount {
		queryParams.Add("include_sms_membership_count", "true")
	}

	targetURL := "/v3/contact_lists?" + queryParams.Encode()

	err := c.restClient.Do(ctx, http.MethodGet, targetURL, nil, &output)
	if err != nil {
		return output, fmt.Errorf("failed to get lists: %w", err)
	}

	// TODO: Handle pagination.

	return output, nil
}
