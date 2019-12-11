package convert

import (
	"strconv"
)

func IntToDollar(amount int64) string {
	intStr := strconv.FormatInt(amount, 10)
	amtStr := intStr[:1] + "." + intStr[1:]
	return amtStr
}

func DaysInMonth(monthStr string) int64 {
	switch monthStr {
		case "January":
			return 31
		case "February":
			return 28
		case "March":
			return 31
		case "April":
			return 30
		case "May":
			return 31
		case "June":
			return 30
		case "July":
			return 31
		case "August":
			return 31
		case "September":
			return 30
		case "October":
			return 31
		case "November":
			return 30
		case "December":
			return 31
	}
	return 0
}
