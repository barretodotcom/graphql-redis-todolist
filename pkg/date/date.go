package date

import "time"

var dateLayout = "2006-01-02 15:04:05"

func ParseStringToDate(date string) (time.Time, error) {
	return time.Parse(dateLayout, date)
}
