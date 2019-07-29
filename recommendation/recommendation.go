package recommendation

type QueryData struct { // Query Data : 총 23개
	Pattern              string // 변증 이름
	Category             string // 카테고리 이름
	Half_Of_Category_Num int    // 카테고리별 질문 수의 절반 = 가중치 / 2
	ShouldBeQueried      bool   // 추천 DB에 쿼리를 날려야하는가?	- 1. 정밀 진단 결과 해당하는 변증인가? ( Key = Pattern ), 2. 진단 결과 HOCN 의 값이 양인가?
}

type PatternCat struct { // Queries의 Key 구조체
	Pattern  string
	Category string
}

type Queries struct {
	QueryCore    map[PatternCat]QueryData // Pattern & Category ( = Key )로 QueryData ( = Value ) 접근
	QueryStrings []string                 // Query문들
}
