package util

import "strconv"

func ShouldPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func MustParseInt(i string) int {
	atoi, e := strconv.Atoi(i)
	ShouldPanic(e)
	return atoi
}
