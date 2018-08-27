package benchs

import (
	"strings"
)

func getRandomString() string {
	// Try adding a 1ns sleep here. Don't forget to import "time".
	// time.Sleep(100 * time.Nanosecond)
	return "0123456789012345"
}

func ConcatCopy(n int) string {
	buf := make([]byte, n*len(getRandomString()))
	count := 0

	for i := 0; i < n; i++ {
		count += copy(buf[count:], getRandomString())
	}

	return string(buf)
}

func ConcatAppend(n int) string {
	buf := make([]byte, 0, n*len(getRandomString()))

	for i := 0; i < n; i++ {
		buf = append(buf, getRandomString()...)
	}

	return string(buf)
}

func ConcatBuilderPreGrow(n int) string {
	var b strings.Builder
	b.Grow(n * len(getRandomString()))

	for i := 0; i < n; i++ {
		b.WriteString(getRandomString())
	}

	return b.String()
}
