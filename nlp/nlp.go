package nlp

import (
	"fmt"
	"math/rand"
	"time"
)

func ConvertInquiryScore(str string) string {
	switch str {
	case "일":
		str = "1"
	case "이":
		str = "2"
	case "삼":
		str = "3"
	case "사":
		str = "4"
	case "오":
		str = "5"
	}
	return str
}

var playUptoOnePoint = []string{"다행이에요.", "아주 좋은편에 속하시네요.", "건강왕이시네요."}
var playUptoTwoPoint = []string{"괜찮으시네요.", "좋은편에 속하시네요.", "굿잡"}
var playUptoThreePoint = []string{"평범하시군요.", "건강한 편에 속하시네요.", "좋아요!"}
var playUptoFourPoint = []string{"굿", "와우", "예아"}
var playUptoFivePoint = []string{"오우", "오호", "굳"}

var playUptoPoint = [][]string{playUptoOnePoint, playUptoTwoPoint, playUptoThreePoint, playUptoFourPoint, playUptoFivePoint}

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
