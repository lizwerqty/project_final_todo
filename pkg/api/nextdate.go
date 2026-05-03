package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)
const layout = "20060102"

func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(layout, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	startDate, err := time.Parse(layout, dstart)
	if err != nil {
		return "", err
	}

	if repeat == "y" {

		date := startDate

		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format(layout), nil
			}
		}
	}

	if strings.HasPrefix(repeat, "d ") {

		parts := strings.SplitN(repeat, " ", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repeat format")
		}

		N, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}

		if N <= 0 || N > 400 {
			return "", fmt.Errorf("invalid day interval")
		}

		date := startDate

		for {
			date = date.AddDate(0, 0, N)
			if date.After(now) {
				return date.Format(layout), nil
			}
		}
	}

	if strings.HasPrefix(repeat, "w ") {
		parts := strings.SplitN(repeat, " ", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repeat format")
		}
		var n int
		allowed := map[int]bool{}
		secParts := strings.Split(parts[1], ",")
		for _, v := range secParts {
			n, err = strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			if n < 1 || n > 7 {
				return "", fmt.Errorf("invalid weekday")
			}
			allowed[n] = true
		}

		date := startDate
		for {
			date = date.AddDate(0, 0, 1)
			weekDay := int(date.Weekday())
			if weekDay == 0 {
				weekDay = 7
			}
			if date.After(now) && allowed[weekDay] {
				return date.Format(layout), nil
			}
		}
	}

	if strings.HasPrefix(repeat, "m ") {
		parts := strings.Split(repeat, " ")
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("invalid repeat format")
		}
		days := strings.Split(parts[1], ",")
		var months []string
		hasLast := false
		hasPrevLast := false

		if len(parts) == 3 {
			months = strings.Split(parts[2], ",")
		}

		dayMap := map[int]bool{}

		for _, v := range days {
			n, err := strconv.Atoi(v)
			if err != nil {
				return "", fmt.Errorf("invalid day")
			}

			if n == 0 || n < -2 || n > 31 {
				return "", fmt.Errorf("invalid day")
			}

			if n == -1 {
				hasLast = true
				continue
			}
			if n == -2 {
				hasPrevLast = true
				continue
			}

			dayMap[n] = true
		}

		monthMap := map[int]bool{}

		if len(parts) == 3 {
			for _, v := range months {
				n, err := strconv.Atoi(v)
				if err != nil {
					return "", err
				}
				if n < 1 || n > 12 {
					return "", fmt.Errorf("invalid month")
				}
				monthMap[n] = true
			}
		}

		date := startDate
		dateOK := false
		monthOK := false

		for {
			date = date.AddDate(0, 0, 1)
			month := int(date.Month())

			firstOfNextMonth := time.Date(
				date.Year(),
				date.Month()+1,
				1,
				0, 0, 0, 0,
				date.Location(),
			)

			lastDay := firstOfNextMonth.AddDate(0, 0, -1)

			dateOK = dayMap[date.Day()]

			if hasLast && date.Day() == lastDay.Day() {
				dateOK = true
			}

			if hasPrevLast && date.Day() == lastDay.Day()-1 {
				dateOK = true
			}

			if len(monthMap) == 0 {
				monthOK = true
			} else {
				monthOK = monthMap[month]
			}
			if date.After(now) && dateOK && monthOK {
				return date.Format(layout), nil
			}
		}
	}
	return "", fmt.Errorf("unsupported repeat rule")
}
