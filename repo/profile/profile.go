package profileservice

type profile struct {
	ProfileID   string `json:"profileId"`
	ProfileName string `json:"profileName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
}
