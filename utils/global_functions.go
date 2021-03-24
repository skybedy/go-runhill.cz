package utils

import (
	"strconv"
)

func SecToTime(secinput int) string {
	var hh string
	h := secinput / 3600

	if h > 0 {
		hh = strconv.Itoa(h) + ":"
	}

	hz := secinput % 3600
	m := strconv.Itoa(hz / 60)
	s := strconv.Itoa(hz % 60)
	if len(m) == 1 {
		m = "0" + m
	}
	if len(s) == 1 {
		s = "0" + s
	}

	return hh + m + ":" + s
}
