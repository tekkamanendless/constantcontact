package constantcontact

import (
	"context"
	"net/http"
	"net/url"
)

type GetEmailIDInput struct {
	_ string `url:"/v3/emails/{campaign_id}"`

	CampaignID string `path:"campaign_id"`
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

	err := c.restClient.Do(ctx, http.MethodGet, "/v3/emails/"+url.PathEscape(input.CampaignID), nil, &output)
	if err != nil {
		return output, err
	}

	return output, nil
}
