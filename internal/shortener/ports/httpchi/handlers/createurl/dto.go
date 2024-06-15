package createurl

type aliasResp struct {
	Alias string `json:"alias"`
}

type urlRequest struct {
	Original string `json:"original" validate:"required,url"`
}
