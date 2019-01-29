// question
package question

const PATTERN_NUM = 5                                                           // 변증
var PATTERN_NAME = []string{"칠정", "노권", "담음", "식적", "어혈"}                       // 변증 이름
var PATTERN_INDEX = map[string]int{"칠정": 0, "노권": 1, "담음": 2, "식적": 3, "어혈": 4} // 변증 인덱스 : 이름
const BI_CRITERIA = 3                                                           // 이분화 기준
const SCORE_MAX = 5                                                             // 점수 최댓값
var CATEGORY_NUM = []int{4, 6, 6, 4, 3}

const REP_HALF = 11     // 간단진단 중 질문 수 확인 지점 1
const REP_FINAL = 18    // 간단진단 중 질문 수 확인 지점 2
const DETAIL_GAP = 12   // 정밀진단 중 질문 수 확인 간격
const PROB_PLAYUPTO = 3 // 질문마다 맞장구 쳐주는 확률의 수치 , 3 => rand(3) : 0~4 => 4 => 1/4 (25%)확률로 맞장구 쳐줌.

const FIRST_IDX = 1 // QCWP 실제 데이터 시작 행번호
const QUESTION = 0  // QCWP question 열번호
const CATEGORY = 1  // QCWP category 열번호
const WEIGHT = 2    // QCWP weight 열번호
const PATTERN = 3   // QCWP pattern 열번호

const PTOC_PATTERN = 0 // PtoC pattern 열번호
const PTOC_CUTOFF = 1  // PtoC cutoff 열번호

var RAW_DATA qDataConst = loadData() // raw data

type qDataConst struct {
	QCWP [][]string     // 질문, 카테고리, 가중치, 변증
	PtoC map[string]int // 변증 : 컷오프
}

type QData struct {
	QRepIdx                []int       // 각 변증의 각 카테고리별 대표 질문들에 대한 QCWP 인덱스 슬라이스
	QDetailIdx             [][]int     // [칠정에 대한 QCWP 인덱스 슬라이스, 노권에 대한 QCWP 인덱스 슬라이스, ..., 어혈에 대한 QCWP 인덱스 슬라이스]
	QDetailNum             int         // 정밀 진단 질문 개수
	QDetailCount           int         // 정밀 진단 질문 카운트
	Answer                 map[int]int // QCWP 인덱스 : 응답점수
	SQSProbPatternIdx      []int       // 간단한 문진 이후 컷오프 값을 넘긴 Pattern의 인덱스 슬라이스
	FinalScore             []float64   // 간단한 문진 이후 컷오프 값을 넘긴 Pattern 인덱스 에 대한 표준점수 슬라이스
	RepIdx                 int
	RepMax                 int
	DetPat                 int // GetDQSAnswer()에서 사용. 값 : 0~len(SQSProbPatternIdx)
	DetIdx                 int
	DetMax                 int
	FinalScoreNotification string // 최종 결과
}
