package datetime

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"pkg/errors"
)

const JSONNull = "null"
const DateFormat = "2006-01-02"

type Date struct {
	time.Time
}

func NewDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func Today() Date {
	now := time.Now()
	return Date{time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)}
}

func (d Date) AddDate(year int, month int, day int) Date {
	return Date{d.Time.AddDate(year, month, day)}
}

func Parse(str string) (date Date, err error) {

	time, err := time.Parse(DateFormat, str)
	if err != nil {
		return date, errors.BadRequest.Wrap(err)
	}

	return Date{Time: time}, nil
}

func (d *Date) UnmarshalJSON(b []byte) (err error) {

	s := strings.Trim(string(b), "\"") // remove quotes
	if s == JSONNull || s == "" {
		return nil
	}

	d.Time, err = time.Parse(DateFormat, s)
	if err != nil {
		return err
	}

	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {

	if d.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%v"`, d.Format(DateFormat))), nil
}

func (d *Date) UnmarshalText(b []byte) (err error) {

	s := strings.Trim(string(b), "\"") // remove quotes
	if s == JSONNull || s == "" {
		return nil
	}

	d.Time, err = time.Parse(DateFormat, s)
	if err != nil {
		return err
	}

	return nil
}

func (d *Date) Scan(src any) error {
	if t, ok := src.(time.Time); ok {
		d.Time = t
	}
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}
