package constantcontact

type HasLinks struct {
	Links map[string]Link `json:"_links"`
}

type Link struct {
	Href string `json:"href"`
}
