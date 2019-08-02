// package recommendation
package recommendation

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"munzini/DB"
	"munzini/random"
	"os"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

// 1. Initialize FoodQueryCore
func PrepareQueryCore() FoodQueryCore {
	// open QCWP file	- Use CWP ( Category-Weight-Pattern )
	qcwp_file, _ := os.Open("resources/data/QCWP.csv")

	// create csv Reader
	qcwp_reader := csv.NewReader(bufio.NewReader(qcwp_file))

	// read csv file
	qcwp, _ := qcwp_reader.ReadAll()
	/*
		TODO
		1. QueryCore를 초기화한다.
			2. QueryCore를 초기화하기 위하여 PatternCat 리스트를 만든다.
			3. QueryCore의 Key값에 PatternCat을 넣고, 이에 따른 QueryData를 작성하는 로직을 만든다.
	*/

	// 1. PatternCat 초기화 ( QueryCore의 Key 값 )
	var patcat []PatternCat

	// TODO: PatternCat리스트 초기화 ( 23개 - 카테고리 개수)
	var row int = FIRST_IDX
	var weight []int // 추후에 QueryCore의 Value값 중 Half_Of_Category_Num에 값을 담아놓기 위해 가중치 값들을 미리 저장해놓는 슬라이스

	for row < len(qcwp) {
		temp_patcat := PatternCat{
			Pattern:  qcwp[row][PATTERN_IDX],
			Category: qcwp[row][CATEGORY_IDX],
		}
		patcat = append(patcat, temp_patcat)
		temp_weight, _ := strconv.Atoi(qcwp[row][WEIGHT_IDX])
		weight = append(weight, temp_weight)
		row_gap, _ := strconv.Atoi(qcwp[row][WEIGHT_IDX])
		row = row + row_gap // 가중치만큼 Forwarding
	}

	// 2. QueryCore 초기화 ( PatternCat - QueryData : Pattern / Category / Half_Of_Category_Num / ShouldBeQueried )
	var queryCore map[PatternCat]QueryData = make(map[PatternCat]QueryData)

	// PatternCat의 값을 QueryCore의 Key값에 넣고, 그에 해당하는 QueryData를 작성한다.
	for qd_idx := 0; qd_idx < len(patcat); qd_idx++ {
		queryCore[patcat[qd_idx]] = QueryData{
			Pattern:              patcat[qd_idx].Pattern,
			Category:             patcat[qd_idx].Category,
			Half_Of_Category_Num: weight[qd_idx] / 2,
			ShouldBeQueried:      true,
		}
	}

	// 3. FoodQueryCore 작성
	var foodQueryCore FoodQueryCore = FoodQueryCore{
		QueryCore: queryCore,
	}
	return foodQueryCore
}

// 2. Calculate FoodQueryCore's Half_Of_Category_Num according to user's Response
func CalculateHOCN(fqcore FoodQueryCore, pattern string, category string, score int) FoodQueryCore {

	// Make QueryCore's Key
	QCkey := PatternCat{
		Pattern:  pattern,
		Category: category,
	}

	new_HOCN := fqcore.QueryCore[QCkey].Half_Of_Category_Num
	new_ShouldBeQueried := fqcore.QueryCore[QCkey].ShouldBeQueried

	if score >= HOCN_CRITERIA { // 3점 이상일 시 QCKey에 해당하는 QueryData의 HOCN을 감소시킨다.
		new_HOCN -= 1
	}
	// HOCN이 음수가 되면, 쿼리 대상에서 제외시킨다.
	if new_HOCN < 0 {
		new_ShouldBeQueried = false
	}

	newQueryData := QueryData{
		Pattern:              pattern,
		Category:             category,
		Half_Of_Category_Num: new_HOCN,
		ShouldBeQueried:      new_ShouldBeQueried,
	}

	fqcore.QueryCore[QCkey] = newQueryData

	return fqcore
}

// 3. Determine Which Pattern-Category should be queried to RMD Database
func extractQPC(fqcore FoodQueryCore, ProbPatternList []string) []PatternCat {
	var patcats []PatternCat

	// A. 정밀 문진 결과 의심되는 패턴이 아닌 것은 모두 쿼리 대상에서 제외시킨다.
	for key, value := range fqcore.QueryCore { // Value Type : QueryData
		if !strIn(value.Pattern, ProbPatternList) { // QueryData의 Pattern이 ProbPatternList에 없다면, (문제있는 변증이 아니라면)
			newQueryData := QueryData{
				Pattern:              value.Pattern,
				Category:             value.Category,
				Half_Of_Category_Num: value.Half_Of_Category_Num,
				ShouldBeQueried:      false,
			}
			fqcore.QueryCore[key] = newQueryData
		}
	}

	for _, value := range fqcore.QueryCore {
		if value.ShouldBeQueried { // QueryData를 검색해야한다면
			temp_patcat := PatternCat{
				Pattern:  value.Pattern,
				Category: value.Category,
			}
			patcats = append(patcats, temp_patcat)
		}
	}

	return patcats
}

func strIn(pattern string, ProbPatternList []string) bool {
	var isIn bool = false

	for i := 0; i < len(ProbPatternList); i++ {
		if pattern == ProbPatternList[i] {
			isIn = true
			break
		}
	}

	return isIn
}

// 4. DB 모듈로부터 받은 쿼리 결과 데이터를 가공
func makeQueries(patterncats []PatternCat) []bson.M {
	var queries []bson.M // Queries를 담는 배열
	for i := 0; i < len(patterncats); i++ {
		// PatternCat의 Pattern과 Category를 사용하여 쿼리를 작성한다.
		queries = append(queries, bson.M{"pattern": patterncats[i].Pattern, "category": patterncats[i].Category})
	}

	return queries
}

// 6. DB 모듈로부터 받은 쿼리 결과 데이터를 가공
func makeRecJsonSet(dbResponses [][]interface{}) [][]RecJson {
	var recJsonSet [][]RecJson
	numRes := len(dbResponses) // (변증, 카테고리) 쌍의 개수
	for i := 0; i < numRes; i++ {
		var recJsonSubSet []RecJson
		numDoc := len(dbResponses[i])
		randIndexes := random.RangeInt(0, numDoc, RMD_PER_CAT)
		for j := 0; j < RMD_PER_CAT; j++ {
			doc := dbResponses[i][randIndexes[j]].(bson.M)
			recJson := RecJson{
				Pattern:  doc["pattern"].(string),
				Category: doc["category"].(string),
				FoodNm:   doc["foodNm"].(string),
			}
			recJsonSubSet = append(recJsonSubSet, recJson)
		}
		recJsonSet = append(recJsonSet, recJsonSubSet)
	}

	return recJsonSet
}

// 7. 가공한 데이터로 추천 스크립트 작성
func makeRecScript(recJsonSet [][]RecJson) string {
	var CRec []RecJson // 칠정 식품
	var NRec []RecJson // 노권 식품
	var DRec []RecJson // 담음 식품
	var SRec []RecJson // 식적 식품
	var URec []RecJson // 어혈 식품
	numSet := len(recJsonSet)
	for i := 0; i < numSet; i++ {
		for j := 0; j < RMD_PER_CAT; j++ {
			recJson := recJsonSet[i][j]
			switch recJson.Pattern {
			case "칠정":
				CRec = append(CRec, recJson)
			case "노권":
				NRec = append(NRec, recJson)
			case "담음":
				DRec = append(DRec, recJson)
			case "식적":
				SRec = append(SRec, recJson)
			case "어혈":
				URec = append(URec, recJson)
			default:
				fmt.Println("Error: document has been crashed.")
			}
		}
	}
	script := "제가 더욱 건강한 삶을 위해 추천드릴 음식들을 정리해봤어요. "
	if len(CRec) != 0 {
		script += "칠정에 좋은 "
		for i := 0; i < len(CRec); i++ {
			script += CRec[i].FoodNm
			script += ", "
		}
	}
	if len(NRec) != 0 {
		script += "노권에 좋은 "
		for i := 0; i < len(NRec); i++ {
			script += NRec[i].FoodNm
			script += ", "
		}
	}
	if len(DRec) != 0 {
		script += "담음에 좋은 "
		for i := 0; i < len(DRec); i++ {
			script += DRec[i].FoodNm
			script += ", "
		}
	}
	if len(SRec) != 0 {
		script += "식적에 좋은 "
		for i := 0; i < len(SRec); i++ {
			script += SRec[i].FoodNm
			script += ", "
		}
	}
	if len(URec) != 0 {
		script += "어혈에 좋은 "
		for i := 0; i < len(URec); i++ {
			script += URec[i].FoodNm
			script += ", "
		}
	}
	script += "이와 같은 음식을 드셔보실 것을 권해드립니다."

	return script
}

// 추천 스크립트 제작 후 저장 및 반환
func GetAndSaveFoodRecommendation(fqCore FoodQueryCore, probPatternList []string) string {
	// 3. 문제 패턴 및 카테고리 정보 생성
	patternCatList := extractQPC(fqCore, probPatternList)

	// 4. 추천 DB에 요청할 쿼리 작성
	fmt.Print("building queries..")
	queries := makeQueries(patternCatList)
	fmt.Println("done")

	// 5. 작성한 쿼리로 DB 모듈을 통해 쿼리 결과 수신 (DB 모듈 호출)
	fmt.Print("loading data..")
	dbResponses := DB.RequestQueries(RMD_COLLECTION_NAME, queries)
	fmt.Println("done")
	numDoc := 0
	for i := 0; i < len(dbResponses); i++ {
		numDoc += len(dbResponses[i])
	}
	fmt.Println(strconv.Itoa(numDoc) + " docs fetched")

	// 6. DB 모듈로부터 받은 쿼리 결과 데이터를 가공
	fmt.Print("building food recommendation documents..")
	recJsonSet := makeRecJsonSet(dbResponses)
	fmt.Println("done")
	fmt.Println("documents content:")
	fmt.Println(recJsonSet)

	// 7. 가공한 데이터로 추천 스크립트 작성
	fmt.Print("building food recommendation script..")
	recScript := makeRecScript(recJsonSet)
	fmt.Println("done")
	fmt.Println("script:")
	fmt.Println(recScript)

	return recScript
}
