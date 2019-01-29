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

var playUptoOnePoint = []string{"다행이에요. ", "다행이네요. ", "정말 다행이에요. "}
var playUptoTwoPoint = []string{"다행이에요. ", "다행이네요."}
var playUptoThreePoint = []string{"아 그렇군요. ", "아 그러시군요. "}
var playUptoFourPoint = []string{"그렇군요. ", "그러시군요. ", "그랬군요. "}
var playUptoFivePoint = []string{"음... 그렇군요. ", "음... 그러시군요. ", "음... 그랬군요. "}

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
	scoreIdx := score - 1

	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // Randomly Choose Answer
	answerPicker := r.Intn(len(playUptoPoint[scoreIdx]))

	playUptoMessage := playUptoPoint[scoreIdx][answerPicker] + "다음 질문입니다. " // Choose Answer

	return playUptoMessage
}
