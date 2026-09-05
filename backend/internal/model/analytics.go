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

type ClickTimeline struct {
	Date   time.Time `json:"date"`
	Clicks int64     `json:"clicks"`
}

type TimelineResponse struct {
	Timeline []ClickTimeline `json:"timeline"`
}

type PaginationRequest struct {
	Limit  int64 `form:"limit"`
	Offset int64 `form:"offset"`
}

type AnalyticsDimension struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type AnalyticsSummary struct {
	Clicks    int64                `json:"clicks"`
	Unique    int64                `json:"unique_clicks"`
	Countries []AnalyticsDimension `json:"countries"`
	Devices   []AnalyticsDimension `json:"devices"`
	Browsers  []AnalyticsDimension `json:"browsers"`
	Referrers []AnalyticsDimension `json:"referrers"`
}
