package service

import (
	"logger-app/models"
	"sort"
)

func CalculateRangeSummary(events []models.Event) models.RangeSummary {
	dailyMap := make(map[string][]models.Event)

	// group events by date
	for _, e := range events {
		date := e.Timestamp.Format("2006-01-02")
		dailyMap[date] = append(dailyMap[date], e)
	}

	var result models.RangeSummary

	for date, dayEvents := range dailyMap {

		// sort just in case
		sort.Slice(dayEvents, func(i, j int) bool {
			return dayEvents[i].Timestamp.Before(dayEvents[j].Timestamp)
		})

		// append end of day
		endTime := EndOfDay(date)
		last := dayEvents[len(dayEvents)-1]

		dayEvents = append(dayEvents, models.Event{
			UserID:    last.UserID,
			Place:     last.Place,
			Timestamp: endTime,
		})

		summary := CalculateSummary(dayEvents)

		result.TotalHome += summary.HomeTime
		result.TotalOffice += summary.OfficeTime
		result.TotalOutside += summary.OutsideTime
		result.TotalCommute += summary.CommuteTime

		result.Days = append(result.Days, models.DailySummary{
			Date:        date,
			HomeTime:    FormatDuration(summary.HomeTime),
			OfficeTime:  FormatDuration(summary.OfficeTime),
			OutsideTime: FormatDuration(summary.OutsideTime),
			CommuteTime: FormatDuration(summary.CommuteTime),
			Segments:    summary.Segments,
		})
	}

	return result
}
