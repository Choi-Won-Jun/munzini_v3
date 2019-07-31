// package recommendation
package recommendation

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
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

func makeQueries(patterns []string) {
	// TODO : 입력받은 pattern들에 따라서 query list를 만들어 반환한다.
	//  return queries
}

func RequestQueries(queries []string) {
	// TODO : 생산된 query list들을 mongoDB에 요청한다. 그에 따라 나온 결과값들을 JSON string리스트 형식으로 받아온다.
	// return dbResponses
}

func makeRecJson(dbResponses []string) {
	// TODO: DB에서 뽑아온 데이터를 가공하여 추천 JSON을 만든다.
	// return recJson
}

func makeRecScript(recJson []string) {
	// TODO: 가공된 추천 JSON에서 추천의 말씀 제작
	// return recScript (string)
}

func InsertAndGetRecommendation(patterns []string) {
	// TODO: recommendation 외부 패키지에서 이 함수를 호출하여 추천의 말씀을 얻을 수 있음.
	// return recScript
	queries = makeQueries(patterns)
	dbResponses = RequestQueries(queries)
	recJson = makeRecJson(dbResponses)
	recScript = makeRecScript(recJson)
	return recScript
}
