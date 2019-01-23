// question
package question

const PATTERN_NUM = 5 // 패턴 수
const PATTERN_NAME = {"칠정", "노권", "담음", "식적", "어혈"}
const BI_CRITERIA = 3 // 이분화 기준

const FIRST_IDX = 1 // QCWP 실제 데이터 시 행번호
const QUESTION = 0  // QCWP question 열번호
const CATEGORY = 1  // QCWP category 열번호
const WEIGHT = 2    // QCWP weight 열번호
const PATTERN = 3   // QCWP pattern 열번호

const PTOC_PATTERN = 0 // PtoC pattern 열번호
const PTOC_CUTOFF = 1  // PtoC cutoff 열번호

type qDataConst struct {
	QCWP [][]string     // 질문, 카테고리, 가중치, 변증
	PtoC map[string]int // 변증 : 컷오프
}

type QData struct {
	RawData           qDataConst
	QRepIdx           []int       // 각 변증의 각 카테고리별 대표 질문들에 대한 QCWP 인덱스 슬라이스
	QDetailIdx        [][]int     // [칠정에 대한 QCWP 인덱스 슬라이스, 노권에 대한 QCWP 인덱스 슬라이스, ..., 어혈에 대한 QCWP 인덱스 슬라이스]
	Answer            map[int]int // QCWP 인덱스 : 응답점수
	SQSProbPatternIdx []int       // 간단한 문진 이후 컷오프 값을 넘긴 Pattern의 인덱스 슬라이스
	FinalScore	[]int	// 간단한 문진 이후 컷오프 값을 넘긴 Pattern에 대한 표준점수
}
