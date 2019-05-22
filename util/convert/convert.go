package convert

import (
	"strconv"
)

func IntToDollar(amount int64) string {
	intStr := strconv.FormatInt(amount, 10)
	amtStr := intStr[:1] + "." + intStr[1:]
	return amtStr
}
