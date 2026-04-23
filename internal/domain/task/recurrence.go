package task

import (
	"encoding/json"
	"fmt"
	"time"
)

type RecurrenceType string

const (
	RecurrenceDaily         RecurrenceType = "daily"
	RecurrenceMonthly       RecurrenceType = "monthly"
	RecurrenceSpecificDates RecurrenceType = "specific_dates"
	RecurrenceEvenOdd       RecurrenceType = "even_odd"
)

type Parity string

const (
	ParityEven Parity = "even"
	ParityOdd  Parity = "odd"
)

type Date struct {
	time.Time
}

func ParseDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, fmt.Errorf("invalid date %q: expected YYYY-MM-DD", s)
	}
	return Date{t}, nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format("2006-01-02"))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date %q: expected YYYY-MM-DD", s)
	}
	d.Time = t
	return nil
}

type Recurrence struct {
	Type RecurrenceType `json:"type"`

	// daily
	EveryNDays int   `json:"every_n_days,omitempty"`
	StartDate  *Date `json:"start_date,omitempty"`

	// monthly
	DayOfMonth int `json:"day_of_month,omitempty"`

	// specific_dates
	Dates []Date `json:"dates,omitempty"`

	// even_odd
	Parity Parity `json:"parity,omitempty"`
}

// IsScheduledFor reports whether the recurrence rule includes the given date
func (r *Recurrence) IsScheduledFor(date time.Time) bool {
	day := date.UTC().Truncate(24 * time.Hour)

	switch r.Type {
	case RecurrenceDaily:
		if r.StartDate == nil || r.EveryNDays <= 0 {
			return false
		}
		start := r.StartDate.UTC().Truncate(24 * time.Hour)
		if day.Before(start) {
			return false
		}
		daysDiff := int(day.Sub(start).Hours() / 24)
		return daysDiff%r.EveryNDays == 0

	case RecurrenceMonthly:
		return date.Day() == r.DayOfMonth

	case RecurrenceSpecificDates:
		for _, d := range r.Dates {
			if d.UTC().Truncate(24 * time.Hour).Equal(day) {
				return true
			}
		}
		return false

	case RecurrenceEvenOdd:
		dayNum := date.Day()
		if r.Parity == ParityEven {
			return dayNum%2 == 0
		}
		return dayNum%2 != 0
	}

	return false
}
