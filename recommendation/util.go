// package recommendation
package recommendation

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
)

// Initialize PatternCat
func loadData() {
	// open QCWP file	- Use CWP ( Category-Weight-Pattern )
	qcwp_file, _ := os.Open("resources/data/QCWP.csv")

	// create csv Reader
	qcwp_reader := csv.NewReader(bufio.NewReader(file))

	// read csv file
	qcwp, _ := qcwp_reader.ReadAll()

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
