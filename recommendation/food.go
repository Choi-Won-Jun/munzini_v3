package recommendation

const (
	FIRST_IDX    = 1 // QCWP.csv를 담아올 때 접근해야하는 첫번째 인덱스
	CATEGORY_IDX = 1 // QCWP.csv에서 Category에 접근하기 위한 인덱스
	PATTERN_IDX  = 3 // QCWP.csv에서 Pattern에 접근하기 위한 인덱스
	WEIGHT_IDX   = 2 // QCWP.csv에서 Weight에 접근하기 위한 인덱스

	HOCN_CRITERIA = 3 // Half_Of_Category_Number을 감소시킬지 말지를 판단하는 점수 기준

	RMD_PER_CAT         = 1                                  // (변증, 카테고리)별 음식 추천 수
	RMD_COLLECTION_NAME = "FOOD_RECOMMEND_COLLECTION_SIMPLE" // 음식 추천 DB 이름 (조회할 DB)
	// RMD_STORE_COLLECTION_NAME = "FOOD_RECOMMEND_STORE_COLLECTION"  // 음식 추천 기록 DB 이름 (저장할 DB)

)

type RecJson struct { // 추천의 말씀을 뽑아내기 위한 구조체
	Pattern  string
	Category string
	FoodNm   string
}

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

// CEKSessionAttributes를 통하여 주고받아야할 구조체
type FoodQueryCore struct {
	QueryCore map[PatternCat]QueryData // Pattern & Category ( = Key )로 QueryData ( = Value ) 접근

	// 확장을 위하여 남겨두었음.
	// QueryStrings []string
	// QueryOutput [][]SimpleDoc
	// QueryStrings map[PatternCat]string                 // Query문들
}
