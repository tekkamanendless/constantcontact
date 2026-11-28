package constantcontact

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tekkamanendless/constantcontact/wrap"
)

type GetEmailsInput struct {
	_ string `api:"httppath:/v3/emails"`
	_ string `api:"httpmethod:GET"`

	Limit      int       `api:"query:limit"`
	BeforeDate time.Time `api:"query:before_date"`
	AfterDate  time.Time `api:"query:after_date"`
}

type GetEmailsOutput struct {
	HasLinks

	Campaigns []GetEmailsOutputCampaign `json:"campaigns"`
}

type GetEmailsOutputCampaign struct {
	CampaignID    string `json:"campaign_id"`
	CreatedAt     string `json:"created_at"` // TODO: time.Time
	CurrentStatus string `json:"current_status"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	TypeCode      int    `json:"type_code"`
	UpdatedAt     string `json:"updated_at"` // TODO: time.Time
}

func (c *Client) GetEmails(ctx context.Context, input GetEmailsInput) (GetEmailsOutput, error) {
	var output GetEmailsOutput
	err := wrap.WrapRequest(ctx, c.restClient, input, &output)
	if err != nil {
		return output, err
	}

	// Handle pagination.
	if output.Links != nil && output.Links["next"].Href != "" {
		link := output.Links["next"].Href
		for link != "" {
			var pageOutput GetEmailsOutput
			err := c.restClient.Do(ctx, http.MethodGet, link, nil, &pageOutput)
			if err != nil {
				return pageOutput, fmt.Errorf("failed to get page: %w", err)
			}

			output.Campaigns = append(output.Campaigns, pageOutput.Campaigns...)

			link = ""
			if pageOutput.Links != nil && pageOutput.Links["next"].Href != "" {
				link = pageOutput.Links["next"].Href
			}
		}
	}

	return output, nil
}
