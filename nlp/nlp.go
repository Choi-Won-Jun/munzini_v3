package nlp

import (
	"math/rand"
	"strconv"
	"time"
)

var exStrOne = []string{"일", "일본", "일정"}
var exStrTwo = []string{"이", "이정"}
var exStrThr = []string{"삼", "상점", "암점", "삼원", "한번"}
var exStrFor = []string{"사", "사정", "4동", "서점", "서번", "사본", "사전", "화정", "화성", "다정", "아점", "상가점", "카본", "서본"}
var exStrFif = []string{"오", "오정", "호점", "오전"}

var exStrArr = [][]string{exStrOne, exStrTwo, exStrThr, exStrFor, exStrFif}

var playUptoOnePoint = []string{"다행이에요.", "아주 좋은편에 속하시네요.", "건강왕이시네요."}
var playUptoTwoPoint = []string{"괜찮으시네요.", "좋은편에 속하시네요.", "굿잡"}
var playUptoThreePoint = []string{"평범하시군요.", "건강한 편에 속하시네요.", "좋아요!"}
var playUptoFourPoint = []string{"굿", "와우", "예아"}
var playUptoFivePoint = []string{"오우", "오호", "굳"}

var playUptoPoint = [][]string{playUptoOnePoint, playUptoTwoPoint, playUptoThreePoint, playUptoFourPoint, playUptoFivePoint}

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

func PlayUpto(score int) string { // 사용자가 1~5점 사이로 대답을 하였을 때, 맞장구를 쳐주는 말

	scoreIdx := ConvertScore(score)

	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // Randomly Choose Answer
	answerPicker := r.Intn(len(playUptoPoint[scoreIdx]))

	playUptoMessage := playUptoPoint[scoreIdx][answerPicker] // Choose Answer

	return playUptoMessage
}

func ConvertScore(score int) int { // score을 배열의 인덱스화 시켜주는 함수.
	switch score {
	case 1:
		return 0
	case 2:
		return 1
	case 3:
		return 2
	case 4:
		return 3
	case 5:
		return 4
	default:
		return -1
	}
}
