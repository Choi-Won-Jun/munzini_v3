// question
package question

type qDataConst struct {
	QCWP [][]string     // 질문, 카테고리, 가중치, 변증
	PtoC map[string]int // 변증 : 컷오프
}

type qData struct {
	qIdx       []int       // QCWP 인덱스 슬라이스
	qRepIdx    []int       // 각 변증의 각 카테고리별 대표 질문들에 대한 QCWP 인덱스 슬라이스
	qDetailIdx [][]int     // [칠정에 대한 QCWP 인덱스 슬라이스, 노권에 대한 QCWP 인덱스 슬라이스, ..., 어혈에 대한 QCWP 인덱스 슬라이스]
	answer     map[int]int // QCWP 인덱스 : 응답점수
}
