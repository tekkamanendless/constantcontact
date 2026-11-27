package constantcontact

import (
	"context"
	"net/http"
)

type ActivitiesContactDeleteInput struct {
	_ string `api:"httppath:/v3/activities/contact_delete"`
	_ string `api:"httpmethod:POST"`

	ContactIDs []string `json:"contact_ids,omitempty"` // Every contact ID in the list will be deleted.  Can only be used if ListIDs is not used.
	ListIDs    []string `json:"list_ids,omitempty"`    // Every contact in these lists will be deleted.  Can only be used if ContactIDs is not used.
}

type ActivitiesContactDeleteOutput struct {
	HasLinks

	ActivityID     string   `json:"activity_id"`
	State          string   `json:"state"`
	StartedAt      string   `json:"started_at"`   // TODO: time.Time
	CompletedAt    string   `json:"completed_at"` // TODO: time.Time
	CreatedAt      string   `json:"created_at"`   // TODO: time.Time
	UpdatedAt      string   `json:"updated_at"`   // TODO: time.Time
	PercentDone    int      `json:"percent_done"`
	ActivityErrors []string `json:"activity_errors"`
	Status         struct {
		ListCount           int `json:"list_count"`
		ItemsTotalCount     int `json:"items_total_count"`
		ItemsCompletedCount int `json:"items_completed"`
	} `json:"status"`
}

func (c *Client) ActivitiesContactDelete(ctx context.Context, input ActivitiesContactDeleteInput) (ActivitiesContactDeleteOutput, error) {
	var output ActivitiesContactDeleteOutput

	err := c.restClient.Do(ctx, http.MethodPost, "/v3/activities/contact_delete", input, &output)
	if err != nil {
		return output, err
	}

	return output, nil
}
