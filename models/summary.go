package models

import "time"

type TimeSegment struct {
	Place    string `json:"place"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Duration string `json:"duration"`
}

type DailySummary struct {
	Date        string        `json:"date"`
	HomeTime    string        `json:"home_time"`
	OfficeTime  string        `json:"office_time"`
	OutsideTime string        `json:"outside_time"`
	CommuteTime string        `json:"commute_time"`
	Segments    []TimeSegment `json:"segments"`
}

type RangeSummary struct {
	TotalHome    time.Duration
	TotalOffice  time.Duration
	TotalOutside time.Duration
	TotalCommute time.Duration

	Days []DailySummary `json:"days"`
}

type TodayTotals struct {
	Home    string `json:"home"`
	Office  string `json:"office"`
	Outside string `json:"outside"`
	Commute string `json:"commute"`
}

type TodaySummaryResponse struct {
	CurrentTime  string        `json:"current_time"`
	TotalElapsed string        `json:"total_elapsed"`
	Totals       TodayTotals   `json:"totals"`
	Segments     []TimeSegment `json:"segments"`
}
