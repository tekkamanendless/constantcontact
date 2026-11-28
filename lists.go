package constantcontact

import (
	"context"
	"fmt"

	"github.com/tekkamanendless/constantcontact/wrap"
)

type GetListsInput struct {
	_ string `api:"httppath:/v3/contact_lists"`
	_ string `api:"httpmethod:GET"`

	ChannelType               string `api:"query:channel_type"`
	IncludeCount              bool   `api:"query:include_count"`
	IncludeMembershipCount    bool   `api:"query:include_membership_count"`
	IncludeSMSMembershipCount bool   `api:"query:include_sms_membership_count"`
	Limit                     int    `api:"query:limit"`
	Name                      string `api:"query:name"`
	Status                    string `api:"query:status"`
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
	err := wrap.WrapRequest(ctx, c.restClient, input, &output)
	if err != nil {
		return output, fmt.Errorf("failed to get lists: %w", err)
	}

	// TODO: Handle pagination.

	return output, nil
}
