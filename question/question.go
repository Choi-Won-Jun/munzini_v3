// question
package question

const PATTERN_NUM = 5                                                           // 변증
var PATTERN_NAME = []string{"칠정", "노권", "담음", "식적", "어혈"}                       // 변증 이름
var PATTERN_INDEX = map[string]int{"칠정": 0, "노권": 1, "담음": 2, "식적": 3, "어혈": 4} // 변증 인덱스 : 이름
const BI_CRITERIA = 3                                                           // 이분화 기준
const SCORE_MAX = 5                                                             // 점수 최댓값
var CATEGORY_NUM = []int{4, 6, 6, 4, 3}

const SQ_NUM = 23 // 간단진단 질문 개수
const Q_NUM = 91  // 전체 문진 질문 개수

const REP_HALF = 11         // 간단진단 중 질문 수 확인 지점 1
const REP_FINAL = 18        // 간단진단 중 질문 수 확인 지점 2
const YES_SCORE = 4         // 간단진단의 질문에 대해 해당한다고 답하였을 때의 점수
const NO_SCORE = 2          // 간단진단의 질문에 대해 해당하지 않는다고 답하였을 때의 점수
const SERIOUS_SQS = 3       // 간단진단의 결과 문제가 되는 패턴의 개수가 SERIOUS_SQS개 이상일 시, 동일한 간단문진 결과 출력
const SERIOUS_DQS = 3       // 정밀진단의 결과 문제가 되는 패턴의 개수가 SERIOUS_DQS개 이상일 시, 동일한 간단문진 결과 출력
const PROB_EXIST_SCORE = 60 // 정밀진단을 시행하는 패턴의 FinalScore를 계산했을 때, 해당 패턴이 문제가 있는지에 대한 기준.

const DETAIL_GAP = 12   // 정밀진단 중 질문 수 확인 간격
const PROB_PLAYUPTO = 1 // 질문마다 맞장구 쳐주는 확률의 수치 , 1 => rand(1) == 0 일 때 확률, 100%

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
	QRepIdx           []int       // 각 변증의 각 카테고리별 대표 질문들에 대한 QCWP 인덱스 슬라이스
	QDetailIdx        [][]int     // [칠정에 대한 QCWP 인덱스 슬라이스, 노권에 대한 QCWP 인덱스 슬라이스, ..., 어혈에 대한 QCWP 인덱스 슬라이스
	QDetailNum        int         // 정밀 진단 질문 개수
	QDetailCount      int         // 정밀 진단 질문 카운트
	Answer            map[int]int // QCWP 인덱스 : 응답점수
	SQSProbPatternIdx []int       // 간단한 문진 이후 컷오프 값을 넘긴 Pattern의 인덱스 슬라이스
	FinalScore        []float64   // 간단한 문진 이후 컷오프 값을 넘긴 Pattern 인덱스 에 대한 표준점수 슬라이스
	SQSProb           bool        // 간단한 문진 이후 문제가 없는 지의 여부
	// 간단한 문진 이후 문제가 없지만 정밀 진단을 받겠다고 한 것의 여부
	NoSQSProbPatternIdx    []int // 간단한 문진 이후 문제가 없지만, 정밀 진단을 받겠다고 했을 때, 출력할 질문 패턴의 순서를 섞어놓은 슬라이스
	RepIdx                 int
	RepMax                 int
	DetPat                 int // GetDQSAnswer()에서 사용. 값 : 0~len(SQSProbPatternIdx) / 0~len(NoSQSProbPatternIdx)
	DetIdx                 int
	DetMax                 int
	FinalScoreNotification string // 최종 결과
}
