package responses

import "net/http"

type TokenResponse struct {
	BaseResponse
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r TokenResponse) Write(w http.ResponseWriter) {
	r.BaseResponse.Write(w, r)
}

type Message struct {
	BaseResponse
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (r Message) Write(w http.ResponseWriter) {
	r.BaseResponse.Write(w, r)
}

type DataMessage struct {
	BaseResponse
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}

func (r DataMessage) Write(w http.ResponseWriter) {
	r.BaseResponse.Write(w, r)
}

type ClientCreatedResponse struct {
	BaseResponse
	Client       interface{} `json:"client"`
	ClientSecret string      `json:"client_secret"`
}

func (r ClientCreatedResponse) Write(w http.ResponseWriter) {
	r.BaseResponse.Write(w, r)
}
