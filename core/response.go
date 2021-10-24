package core

type Response struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
	Count   int    `json:"count"`
	Results []struct {
		Id          string `json:"id_message"`
		Status      string `json:"status_code"`
		Message     string `json:"status_message"`
		Destination string `json:"destination"`
	} `json:"data"`
	HasErrors bool `json:"errors"`
}
