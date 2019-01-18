// question
package question

const QUESTION = 0
const CATEGORY = 1
const WEIGHT = 2
const PATTERN = 3

const PTOC_PATTERN = 0
const PTOC_CUTOFF = 1

type qDataConst struct {
	QCWP [][]string     // 질문, 카테고리, 가중치, 변증
	PtoC map[string]int // 변증 : 컷오프
}

type QData struct {
	RawData    qDataConst
	QIdx       []int       // QCWP 인덱스 슬라이스
	QRepIdx    []int       // 각 변증의 각 카테고리별 대표 질문들에 대한 QCWP 인덱스 슬라이스
	QDetailIdx [][]int     // [칠정에 대한 QCWP 인덱스 슬라이스, 노권에 대한 QCWP 인덱스 슬라이스, ..., 어혈에 대한 QCWP 인덱스 슬라이스]
	Answer     map[int]int // QCWP 인덱스 : 응답점수
}
