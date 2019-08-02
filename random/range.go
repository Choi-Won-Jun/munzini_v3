package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func RangeInt(min int, max int, n int) []int {
	var rangeBar []int
	for i := min; i < max; i++ {
		rangeBar = append(rangeBar, i)
	}

	var randArr []int
	for i := 0; i < n; i++ {
		index := rand.Intn(len(rangeBar))
		randArr = append(randArr, rangeBar[index])
		rangeBar = remove(rangeBar, index)
	}
	return randArr
}
