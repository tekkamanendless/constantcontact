package constantcontact

import (
	"context"

	"github.com/tekkamanendless/constantcontact/wrap"
)

type GetEmailIDInput struct {
	_ string `api:"httppath:/v3/emails/{campaign_id}"`
	_ string `api:"httpmethod:GET"`

	CampaignID string `api:"path:campaign_id"`
}

type GetEmailIDOutput struct {
	CampaignActivities []struct {
		CampaignActivityID string `json:"campaign_activity_id"`
		Role               string `json:"role"`
		DocumentID         string `json:"document_id"`
	} `json:"campaign_activities"`
	CampaignID    string `json:"campaign_id"`
	CreatedAt     string `json:"created_at"` // TODO: time.Time
	CurrentStatus string `json:"current_status"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	TypeCode      int    `json:"type_code"`
	UpdatedAt     string `json:"updated_at"` // TODO: time.Time
}

func (c *Client) GetEmailID(ctx context.Context, input GetEmailIDInput) (GetEmailIDOutput, error) {
	var output GetEmailIDOutput
	err := wrap.WrapRequest(ctx, c.restClient, input, &output)
	return output, err
}
