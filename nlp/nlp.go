package nlp

import "strconv"

var exStrOne = []string{"일", "일본", "일정"}
var exStrTwo = []string{"이", "이정"}
var exStrThr = []string{"삼", "상점", "암점", "삼원", "한번"}
var exStrFor = []string{"사", "사정", "4동", "서점", "서번", "사본", "사전", "화정", "화성", "다정", "아점", "상가점", "카본", "서본"}
var exStrFif = []string{"오", "오정", "호점", "오전"}

var exStrArr = [][]string{exStrOne, exStrTwo, exStrThr, exStrFor, exStrFif}

func ConvertInquiryScore(str string) string {
	for i := 0; i < len(exStrArr); i++ {
		if strIn(str, exStrArr[i]) {
			return strconv.Itoa(i + 1)
		}
	}
	return str
}

func strIn(str string, arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if str == arr[i] {
			return true
		}
	}
	return false
}
