package model

import "time"

type Analytics struct {
	ID        int64     `json:"id"`
	URLID     int64     `json:"url_id"`
	IpAddress string    `json:"ip_address"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Browser   string    `json:"browser"`
	Os        string    `json:"os"`
	Device    string    `json:"device"`
	Referer   string    `json:"referer"`
	ClickedAt time.Time `json:"clicked_at"`
}
