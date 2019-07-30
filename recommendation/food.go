package recommendation

const FIRST_IDX = 1    // QCWP.csv를 담아올 때 접근해야하는 첫번째 인덱스
const CATEGORY_IDX = 1 // QCWP.csv에서 Category에 접근하기 위한 인덱스
const PATTERN_IDX = 3  // QCWP.csv에서 Pattern에 접근하기 위한 인덱스
const WEIGHT_IDX = 2   // QCWP.csv에서 Weight에 접근하기 위한 인덱스

var queries

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

type SimpleDoc struct {
   Pattern  string `bson:"pattern"`
   Category string `bson:"category"`
   FoodNm   string `bson:"foodNm"`
}

type Queries struct {
	QueryCore    map[PatternCat]QueryData // Pattern & Category ( = Key )로 QueryData ( = Value ) 접근
	// QueryStrings []string
	// QueryOutput [][]SimpleDoc
	// QueryStrings map[PatternCat]string                 // Query문들
}

