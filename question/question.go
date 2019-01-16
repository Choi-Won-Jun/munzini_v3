package question

type qDataConst struct {
  QtoC  map[string]string   // 질문 : 카테고리
  CtoP  map[string]string   // 카테고리 : 변증
  CtoV  map[string]int      // 카테고리 : 가중치(질문 문항 수)
  PtoC  map[string]int      // 변증 : 컷오프
}

type qData struct {
  q       [][]string        // [qA[], qB[], qC[], qD[], qE[]]
  qRep    []string          // [q1, q2, q3, ..., q23]
  repIdx  int               // q에서 뽑은 대표 질문들의 인덱스 값
  qOther  [][]string        // [qAother[], qBother[], ... qEother[]
  answer  map[string]int    // 질문 : 응답점수
}
