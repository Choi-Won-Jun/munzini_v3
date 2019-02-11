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

var playUptoOnePoint = []string{"정말 다행이에요. ", "컨디션이 좋으신가봐요! ", "건강관리 잘 하시나봐요! ", "오호~ ", " 좋으시겠어요. ", "걱정 없네요! ", "다행인데요?", "걱정없겠네요. ", "계속 이렇게만 대답하셨으면 좋겠어요. "}
var playUptoTwoPoint = []string{"좋으시겠네요. ", "건강하신 거 같은데요? ", "좋은 거 드셨나봐요. ", "놀러 나가도 괜찮겠어요! ", "몸이 튼튼하신가봐요. ", "좋았어! ", "계속 문제 없었으면 좋겠네요. "}
var playUptoThreePoint = []string{"아 그러시군요. ", "애매하네요. ", "아하. ", "오홍. ", "오~ ", "음~ ", "그렇군요. ", "무난하시네요. ", "오케이. ", "알겠어요. ", "넵. "}
var playUptoFourPoint = []string{"기운내세요. ", "괜찮으세요? ", "운동 좀 하셔야 겠어요. ", "아이구야. ", "나가서 산책 한 번 해보세요. ", "아..그렇군요 ", "무리하셨나봐요. ", "오늘은 그냥 집에 계요. ", "힘내세요. "}
var playUptoFivePoint = []string{"어머....", "에휴..그러시군요. ", "건강관리 하셔야 겠어요. ", "음..걱정되네요. ", "병원 가보셔야 되는거 아니에요? ", "심각하신데요? ", "걱정되네요..괜찮으신거죠? ", "집에서 좀 쉬세요. ", "장난이시죠? ", "에구머니나. "}

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
