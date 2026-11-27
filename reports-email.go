package constantcontact

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type GetEmailSendsReportInput struct {
	_ string `api:"httppath:/v3/reports/email_reports/{campaign_activity_id}/tracking/sends"`
	_ string `api:"httpmethod:GET"`

	CampaignActivityID string `api:"path:campaign_activity_id"`

	Limit int `api:"query:limit"`
}

type GetEmailSendsReportOutput struct {
	HasLinks

	TrackingActivities []GetEmailSendsReportOutputTrackingActivity `json:"tracking_activities"`
}

type GetEmailSendsReportOutputTrackingActivity struct {
	ContactID            string `json:"contact_id"`
	CampaignActivityID   string `json:"campaign_activity_id"`
	TrackingActivityType string `json:"tracking_activity_type"`
	EmailAddress         string `json:"email_address"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	CreatedTime          string `json:"created_time"` // TODO: time.Time
	DeletedAt            string `json:"deleted_at"`   // TODO: time.Time
}

func (c *Client) GetEmailSendsReport(ctx context.Context, input GetEmailSendsReportInput) (GetEmailSendsReportOutput, error) {
	var output GetEmailSendsReportOutput

	queryParams := url.Values{}
	if input.Limit > 0 {
		queryParams.Add("limit", strconv.Itoa(input.Limit))
	}

	err := c.restClient.Do(ctx, http.MethodGet, "/v3/reports/email_reports/"+url.PathEscape(input.CampaignActivityID)+"/tracking/sends?"+queryParams.Encode(), nil, &output)
	if err != nil {
		return output, err
	}

	// Handle pagination.
	if output.Links != nil && output.Links["next"].Href != "" {
		link := output.Links["next"].Href
		for link != "" {
			var pageOutput GetEmailSendsReportOutput
			err := c.restClient.Do(ctx, http.MethodGet, link, nil, &pageOutput)
			if err != nil {
				return pageOutput, fmt.Errorf("failed to get page: %w", err)
			}

			output.TrackingActivities = append(output.TrackingActivities, pageOutput.TrackingActivities...)

			link = ""
			if pageOutput.Links != nil && pageOutput.Links["next"].Href != "" {
				link = pageOutput.Links["next"].Href
			}
		}
	}

	return output, nil
}

type GetEmailOptoutsReportInput struct {
	_ string `api:"httppath:/v3/reports/email_reports/{campaign_activity_id}/tracking/optouts"`
	_ string `api:"httpmethod:GET"`

	CampaignActivityID string `api:"path:campaign_activity_id"`

	Limit int `api:"query:limit"`
}

type GetEmailOptoutsReportOutput struct {
	HasLinks

	TrackingActivities []GetEmailOptoutsReportOutputTrackingActivity `json:"tracking_activities"`
}

type GetEmailOptoutsReportOutputTrackingActivity struct {
	ContactID            string `json:"contact_id"`
	CampaignActivityID   string `json:"campaign_activity_id"`
	TrackingActivityType string `json:"tracking_activity_type"`
	EmailAddress         string `json:"email_address"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	OptOutReason         string `json:"opt_out_reason"`
	CreatedTime          string `json:"created_time"` // TODO: time.Time
	DeletedAt            string `json:"deleted_at"`   // TODO: time.Time
}

func (c *Client) GetEmailOptoutsReport(ctx context.Context, input GetEmailOptoutsReportInput) (GetEmailOptoutsReportOutput, error) {
	var output GetEmailOptoutsReportOutput

	queryParams := url.Values{}
	if input.Limit > 0 {
		queryParams.Add("limit", strconv.Itoa(input.Limit))
	}

	err := c.restClient.Do(ctx, http.MethodGet, "/v3/reports/email_reports/"+url.PathEscape(input.CampaignActivityID)+"/tracking/optouts?"+queryParams.Encode(), nil, &output)
	if err != nil {
		return output, err
	}

	// Handle pagination.
	if output.Links != nil && output.Links["next"].Href != "" {
		link := output.Links["next"].Href
		for link != "" {
			var pageOutput GetEmailOptoutsReportOutput
			err := c.restClient.Do(ctx, http.MethodGet, link, nil, &pageOutput)
			if err != nil {
				return pageOutput, fmt.Errorf("failed to get page: %w", err)
			}

			output.TrackingActivities = append(output.TrackingActivities, pageOutput.TrackingActivities...)

			link = ""
			if pageOutput.Links != nil && pageOutput.Links["next"].Href != "" {
				link = pageOutput.Links["next"].Href
			}
		}
	}

	return output, nil
}
