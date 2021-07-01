package maintenance

type Healthcheck struct {
	Message string `json:"message" description:"The status of the API." example:"The API is healthy"`
}
