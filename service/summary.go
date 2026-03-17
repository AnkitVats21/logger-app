package service

import (
	"logger-app/models"
	"time"
)

type Summary struct {
	HomeTime    time.Duration
	OfficeTime  time.Duration
	OutsideTime time.Duration
	CommuteTime time.Duration
	Segments    []models.TimeSegment
}

func CalculateSummary(events []models.Event) Summary {
	var summary Summary

	if len(events) < 2 {
		return summary
	}

	const commuteThreshold = 2 * time.Hour

	for i := 0; i < len(events)-1; i++ {
		curr := events[i]
		next := events[i+1]

		duration := next.Timestamp.Sub(curr.Timestamp)
		if duration <= 0 {
			continue
		}

		segment := models.TimeSegment{
			Start:    curr.Timestamp.Format("15:04"),
			End:      next.Timestamp.Format("15:04"),
			Duration: FormatDuration(duration),
			Place:    curr.Place,
		}

		// Commute detection
		isCommute := false
		if curr.Place == "outside" && i > 0 && i < len(events)-1 {
			prev := events[i-1]
			nextPlace := events[i+1].Place

			if (prev.Place == "home" && nextPlace == "office") ||
				(prev.Place == "office" && nextPlace == "home") {

				if duration <= commuteThreshold {
					summary.CommuteTime += duration
					segment.Place = "commute"
					isCommute = true
				}
			}
		}

		if !isCommute {
			// Normal allocation
			switch curr.Place {
			case "home":
				summary.HomeTime += duration
			case "office":
				summary.OfficeTime += duration
			case "outside":
				summary.OutsideTime += duration
			}
		}

		summary.Segments = append(summary.Segments, segment)
	}

	return summary
}

func CalculateTodaySummary(events []models.Event, now time.Time) models.TodaySummaryResponse {
	// Filter out future events that might exist due to mock data seeding
	// TODO: remove in production
	var pastEvents []models.Event
	for _, e := range events {
		if e.Timestamp.After(now) {
			continue
		}
		pastEvents = append(pastEvents, e)
	}
	events = pastEvents

	// Start of current day based on 'now'
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	totalElapsed := now.Sub(startOfDay)
	commuteThreshold := 90 * time.Minute

	var totals Summary
	var segments []models.TimeSegment

	if len(events) == 0 {
		return models.TodaySummaryResponse{
			CurrentTime:  now.Format("2006-01-02T15:04:05"),
			TotalElapsed: FormatDuration(totalElapsed),
			Totals: models.TodayTotals{
				Home: "0h 0m", Office: "0h 0m", Outside: "0h 0m", Commute: "0h 0m",
			},
			Segments: segments, // empty list
		}
	}

	for i := 0; i < len(events); i++ {
		curr := events[i]

		var endTime time.Time
		var isLast bool

		if i < len(events)-1 {
			endTime = events[i+1].Timestamp
		} else {
			endTime = now
			isLast = true
		}

		duration := endTime.Sub(curr.Timestamp)
		if duration <= 0 {
			continue
		}

		segment := models.TimeSegment{
			Start:    curr.Timestamp.Format("15:04"),
			End:      endTime.Format("15:04"),
			Duration: FormatDuration(duration),
			Place:    curr.Place,
		}

		// Commute detection
		isCommute := false
		if curr.Place == "outside" && i > 0 && (!isLast || i < len(events)-1) {
			prev := events[i-1]
			nextPlace := ""

			if !isLast {
				nextPlace = events[i+1].Place
			} else {
				// If we are currently "outside" and it's the last event, we can't definitively call it a commute yet
				// since they haven't arrived at the destination. We'll leave it as "outside".
			}

			if (prev.Place == "home" && nextPlace == "office") ||
				(prev.Place == "office" && nextPlace == "home") {

				if duration <= commuteThreshold {
					totals.CommuteTime += duration
					segment.Place = "commute"
					isCommute = true
				}
			}
		}

		if !isCommute {
			switch curr.Place {
			case "home":
				totals.HomeTime += duration
			case "office":
				totals.OfficeTime += duration
			case "outside":
				totals.OutsideTime += duration
			}
		}

		segments = append(segments, segment)
	}

	return models.TodaySummaryResponse{
		CurrentTime:  now.Format("2006-01-02T15:04:05"),
		TotalElapsed: FormatDuration(totalElapsed),
		Totals: models.TodayTotals{
			Home:    FormatDuration(totals.HomeTime),
			Office:  FormatDuration(totals.OfficeTime),
			Outside: FormatDuration(totals.OutsideTime),
			Commute: FormatDuration(totals.CommuteTime),
		},
		Segments: segments,
	}
}
