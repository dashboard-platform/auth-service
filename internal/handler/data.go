package handler

type Response struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
}
