package tools

import (
	"fmt"
	"math"
	"time"
)

type DtTimeRange struct {
	Period       string `json:"period"`
	TimeType     string `json:"time_type"`
	Offset       int    `json:"offset"`
	Year         int    `json:"year"`
	YearTo       int    `json:"year_to"`
	Month        int    `json:"month"`
	MonthTo      int    `json:"month_to"`
	Week         int    `json:"week"`
	Day          int    `json:"day"`
	DayTo        int    `json:"day_to"`
	Hour         int    `json:"hour"`
	HourTo       int    `json:"hour_to"`
	Minute       int    `json:"minute"`
	MinuteTo     int    `json:"minute_to"`
	Second       int    `json:"second"`
	SecondTo     int    `json:"second_to"`
	Shift        int    `json:"shift"`
	Quarter      int    `json:"quarter"`
	ShrinkedFrom int64  `json:"shrinked_from"`
	ShrinkedTo   int64  `json:"shrinked_to"`
	CustomFrom   int64  `json:"custom_from"`
	CustomTo     int64  `json:"custom_to"`
}

// DtTimeRangeRunningStatus defines possible states
type DtTimeRangeRunningStatus int

const (
	IN DtTimeRangeRunningStatus = iota
	OUT
)

// NewDtTimeRange constructor
func NewDtTimeRange() *DtTimeRange {
	now := GetTimeNow()

	return &DtTimeRange{
		Period:       "month",
		TimeType:     "day",
		Offset:       0,
		Year:         now.Year(),
		YearTo:       now.Year(),
		Month:        int(now.Month()),
		MonthTo:      int(now.Month()),
		Week:         GetWeek(now),
		Day:          0,
		DayTo:        0,
		Hour:         0,
		HourTo:       23,
		Minute:       0,
		MinuteTo:     59,
		Second:       0,
		SecondTo:     59,
		Shift:        0,
		Quarter:      1,
		ShrinkedFrom: math.MaxInt64,
		ShrinkedTo:   0,
		CustomFrom:   0,
		CustomTo:     0,
	}
}

// FromStartEnd sets the start and end datetime range
func (dtr *DtTimeRange) FromStartEnd(start, end time.Time) {
	dtr.Period = "custom"
	dtr.TimeType = "custom"

	dtr.Year = start.Year()
	dtr.YearTo = end.Year()
	dtr.Month = int(start.Month())
	dtr.MonthTo = int(end.Month())
	dtr.Day = start.Day()
	dtr.Hour = start.Hour()
	dtr.Minute = start.Minute()
	dtr.DayTo = end.Day()
	dtr.HourTo = end.Hour()
	dtr.MinuteTo = end.Minute()
	dtr.Second = start.Second()
	dtr.SecondTo = end.Second()
}

// FromStartEndMillis sets the start and end datetime range from milliseconds
func (dtr *DtTimeRange) FromStartEndMillis(start, end int64) {
	dtr.CustomFrom = start
	dtr.CustomTo = end

	startDate := time.Unix(start/1000, 0).UTC()
	endDate := time.Unix(end/1000, 0).UTC()
	dtr.FromStartEnd(startDate, endDate)
}

// IsAll checks if the period is set to "all"
func (dtr *DtTimeRange) IsAll() bool {
	return dtr.Period == "all"
}

// Shrink shrinks the time range based on the stored shrinked values
func (dtr *DtTimeRange) Shrink() {
	if dtr.ShrinkedTo == 0 {
		return
	}
	dtr.FromStartEndMillis(dtr.ShrinkedFrom, dtr.ShrinkedTo)
}

func (self *DtTimeRange) AddShrinkTime(time_stamp int64) {
	if time_stamp > self.ShrinkedTo {
		self.ShrinkedTo = time_stamp
	}
	if time_stamp < self.ShrinkedFrom {
		self.ShrinkedFrom = time_stamp
	}
}

func (calc *DtTimeRange) CalculateStartEnd() ([]time.Time, error) {
	var startEnd []time.Time

	// Check if period is "all"
	if calc.Period == "all" {
		startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 23, 59, 59, 0, time.UTC)
		return []time.Time{startDate, endDate}, nil
	}

	// Check timeType and period
	if calc.TimeType == "day" {
		switch calc.Period {
		case "day":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 23, 59, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "week":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 0, 0, 0, 0, time.UTC)
			startDate = startDate.AddDate(0, 0, -(int(startDate.Weekday()) - 1)) // Adjust to the first day of the week
			endDate := startDate.AddDate(0, 0, 6).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
			startEnd = append(startEnd, startDate, endDate)

		case "month":
			daysInMonth := time.Date(calc.Year, time.Month(calc.Month)+1, 0, 23, 59, 59, 0, time.UTC).Day()
			startDate := time.Date(calc.Year, time.Month(calc.Month), 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, time.Month(calc.Month), daysInMonth, 23, 59, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "quarter":
			startMonth := 3*(calc.Quarter-1) + 1
			startDate := time.Date(calc.Year, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(0, 3, 0).Add(-time.Second) // End at the last second of the quarter
			startEnd = append(startEnd, startDate, endDate)

		case "year":
			startDate := time.Date(calc.Year, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, 12, 31, 23, 59, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "custom":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(calc.YearTo, time.Month(calc.MonthTo), calc.DayTo, 23, 59, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)
		}

	} else { // Handle timeType != "day"
		switch calc.Period {
		case "day":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, calc.Hour, calc.Minute, 0, 0, time.UTC)
			endDate := time.Date(calc.YearTo, time.Month(calc.MonthTo), calc.DayTo, calc.HourTo, calc.MinuteTo, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "week":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, calc.Hour, calc.Minute, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, calc.HourTo, calc.MinuteTo, 59, 0, time.UTC)
			endDate = endDate.AddDate(0, 0, 6) // Add 6 days to the end date
			startEnd = append(startEnd, startDate, endDate)

		case "month":
			daysInMonth := time.Date(calc.Year, time.Month(calc.Month)+1, 0, 23, 59, 59, 0, time.UTC).Day()
			startDate := time.Date(calc.Year, time.Month(calc.Month), 1, calc.Hour, calc.Minute, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, time.Month(calc.Month), daysInMonth, calc.HourTo, calc.MinuteTo, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "year":
			startDate := time.Date(calc.Year, 1, 1, calc.Hour, calc.Minute, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, 12, 31, calc.HourTo, calc.MinuteTo, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "custom":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, calc.Hour, calc.Minute, 0, 0, time.UTC)
			endDate := time.Date(calc.YearTo, time.Month(calc.MonthTo), calc.DayTo, calc.HourTo, calc.MinuteTo, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)
		}
	}

	// Adjust the time range based on the offset
	if len(startEnd) > 0 {
		if startEnd[1].Before(startEnd[0]) {
			startEnd[1] = startEnd[1].Add(time.Hour * 24) // Add one day if end is before start
		}

		// Apply offset (in hours) to both start and end times
		startEnd[0] = startEnd[0].Add(time.Duration(calc.Offset) * time.Hour)
		startEnd[1] = startEnd[1].Add(time.Duration(calc.Offset) * time.Hour)

		return startEnd, nil
	}

	return nil, fmt.Errorf("invalid period or time range")
}

func (self *DtTimeRange) GetStartEndMilis() []int64 {
	if self.Period != "all" {
		tmp, _ := self.CalculateStartEnd()
		return []int64{tmp[0].UnixMilli(), tmp[1].UnixMilli()}
	} else {
		return []int64{0, math.MaxInt64}
	}
}

func (self *DtTimeRange) GetDurationMilis() int64 {
	tmp := self.GetStartEndMilis()
	return tmp[1] - tmp[0]
}

// FromJSON populates the fields from a JSON object
func (dtr *DtTimeRange) FromJSON(jObj map[string]interface{}) {

	if v, ok := jObj["period"].(string); ok {
		dtr.Period = v
	}
	if v, ok := jObj["time_type"].(string); ok {
		dtr.TimeType = v
	}
	if v, ok := jObj["offset"].(float64); ok {
		dtr.Offset = int(v)
	}
	if v, ok := jObj["year"].(float64); ok {
		dtr.Year = int(v)
	}
	if v, ok := jObj["year_to"].(float64); ok {
		dtr.YearTo = int(v)
	}
	if v, ok := jObj["month"].(float64); ok {
		dtr.Month = int(v)
	}
	if v, ok := jObj["month_to"].(float64); ok {
		dtr.MonthTo = int(v)
	}
	if v, ok := jObj["week"].(float64); ok {
		dtr.Week = int(v)
	}
	if v, ok := jObj["day"].(float64); ok {
		dtr.Day = int(v)
	}
	if v, ok := jObj["hour"].(float64); ok {
		dtr.Hour = int(v)
	}
	if v, ok := jObj["minute"].(float64); ok {
		dtr.Minute = int(v)
	}
	if v, ok := jObj["day_to"].(float64); ok {
		dtr.DayTo = int(v)
	}
	if v, ok := jObj["hour_to"].(float64); ok {
		dtr.HourTo = int(v)
	}
	if v, ok := jObj["minute_to"].(float64); ok {
		dtr.MinuteTo = int(v)
	}
	if v, ok := jObj["quarter"].(float64); ok {
		dtr.Quarter = int(v)
	}
	if v, ok := jObj["second"].(float64); ok {
		dtr.Second = int(v)
	}
	if v, ok := jObj["second_to"].(float64); ok {
		dtr.SecondTo = int(v)
	}
}

// ToJSON converts the structure to JSON format
func (dtr *DtTimeRange) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"period":    dtr.Period,
		"time_type": dtr.TimeType,
		"offset":    dtr.Offset,
		"year":      dtr.Year,
		"year_to":   dtr.YearTo,
		"month":     dtr.Month,
		"month_to":  dtr.MonthTo,
		"week":      dtr.Week,
		"day":       dtr.Day,
		"day_to":    dtr.DayTo,
		"quarter":   dtr.Quarter,
		"hour":      dtr.Hour,
		"minute":    dtr.Minute,
		"hour_to":   dtr.HourTo,
		"minute_to": dtr.MinuteTo,
		"second":    dtr.Second,
		"second_to": dtr.SecondTo,
	}
}

// CreateMongoTimeQuery creates a MongoDB time query based on the period and time type
func (dtr *DtTimeRange) CreateMongoTimeQuery() map[string]interface{} {
	//TODO sem mato treba dorobit aby sa generovalo mongo query v $IN

	query := make(map[string]interface{})

	if dtr.Period == "custom" && dtr.TimeType == "custom" {
		query["$gte"] = dtr.CustomFrom
		query["$lte"] = dtr.CustomTo
	} else {
		var startDate, endDate time.Time

		switch dtr.Period {
		case "day":
			startDate = time.Date(dtr.Year, time.Month(dtr.Month), dtr.Day, 0, 0, 0, 0, time.UTC)
			endDate = time.Date(dtr.Year, time.Month(dtr.Month), dtr.Day, 23, 59, 59, 0, time.UTC)

		case "week":
			startDate, _ = time.Parse("2006-W01-1", fmt.Sprintf("%d-W%d-1", dtr.Year, dtr.Week))
			endDate = startDate.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

		case "month":
			daysInMonth := getDaysInMonth(dtr.Year, dtr.Month)
			startDate = time.Date(dtr.Year, time.Month(dtr.Month), 1, 0, 0, 0, 0, time.UTC)
			endDate = time.Date(dtr.Year, time.Month(dtr.Month), daysInMonth, 23, 59, 59, 0, time.UTC)

		case "year":
			startDate = time.Date(dtr.Year, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate = time.Date(dtr.YearTo, 12, 31, 23, 59, 59, 0, time.UTC)
		case "all":
			return map[string]interface{}{
				"$gte": 0,
				"$lte": math.MaxInt64,
			}
		}

		query["$gte"] = startDate.UnixNano() / int64(time.Millisecond)
		query["$lte"] = endDate.UnixNano() / int64(time.Millisecond)
	}

	return query
}

// getDaysInMonth returns the number of days in a month
func getDaysInMonth(year, month int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}

// GetWeek returns the ISO week number
func GetWeek(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}
