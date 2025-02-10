package model

import "time"

type OAuth2Client struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	OwnerId      uint64 `json:"owner_id,string"`
	Label        string `json:"label"`
}

type OAuth2CodeData struct {
	Code      string    `json:"code"`
	ClientId  string    `json:"client_id"`
	UserId    uint64    `json:"user_id,string"`
	CreatedAt time.Time `json:"created_at"`
}
