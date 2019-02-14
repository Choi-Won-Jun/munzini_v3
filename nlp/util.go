package nlp

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io" // csv load 관련
	"munzini/question"

	// "math"      // 반올림 관련
	"math/rand" // 임의 추출 관련
	"os"
	"time" // 임의 추출 관련
	//	"strconv" // string 관련 형변환
	//"unicode/utf8"
	//"io/ioutil"
	//"golang.org/x/text/encoding/korean"
	//"golang.org/x/text/transform"
)

// load Data to initialize PlayUptoMessage
func loadData() PlayUptoConst {

	/* old code
	playUptoLow_file, _ := os.Open("resources/data/PlayUptoHigh.csv")  // open PlayUptoLow.csv
	playUptoMid_file, _ := os.Open("resources/data/PlayUptoMid.csv")   // open PlayUptoMid.csv
	playUptoHigh_file, _ := os.Open("resources/data/PlayUptoHigh.csv") // open PlayUptoHigh.csv

	playUptoLow_reader := csv.NewReader(bufio.NewReader(playUptoLow_file))   // create csv reader for PlayUptoLow.csv
	playUptoMid_reader := csv.NewReader(bufio.NewReader(playUptoMid_file))   // create csv reader for PlayUptoMid.csv
	playUptoHigh_reader := csv.NewReader(bufio.NewReader(playUptoHigh_file)) // create csv reader for PlayUptoHigh.csv

	playUptoLow, _ := playUptoLow_reader.ReadAll()   // read playUptoLow.csv
	playUptoMid, _ := playUptoMid_reader.ReadAll()   // read playUptoMid.csv
	playUptoHigh, _ := playUptoHigh_reader.ReadAll() // read playUptoHigh.csv

	var playUptolow [][]string
	var playUptomid [][]string
	var playUptohigh [][]string

	fmt.Print(playUptoLow)
	printArray(playUptoLow)

	for i := FIRST_IDX_R; i < len(playUptoLow); i++ { // initialize playUptolow
		for j := FIRST_IDX_C; j < len(playUptoLow[i]); j++ {
			playUptolow[i][j] = playUptoLow[i][j]
		}
	}
	for i := FIRST_IDX_R; i < len(playUptoMid); i++ { // initialize playUptomid
		for j := FIRST_IDX_C; j < len(playUptoMid[i]); j++ {
			playUptomid[i][j] = playUptoMid[i][j]
		}
	}

	for i := FIRST_IDX_R; i < len(playUptoHigh); i++ { // initialize playUptohigh
		for j := FIRST_IDX_C; j < len(playUptoHigh[i]); j++ {
			playUptohigh[i][j] = playUptoHigh[i][j]
		}
	}

	playUptoconst := PlayUptoConst{
		PlayUptoLowPoint:  playUptolow,
		PlayUptoMidPoint:  playUptomid,
		PlayUptoHighPoint: playUptohigh,
	}

	*/
	playUptoLow_file, _ := os.Open("resources/data/PlayUptoLowV1.csv")   // open PlayUptoLow.csv
	playUptoMid_file, _ := os.Open("resources/data/PlayUptoMidV1.csv")   // open PlayUptoMid.csv
	playUptoHigh_file, _ := os.Open("resources/data/PlayUptoHighV3.csv") // open PlayUptoHigh.csv

	playUptoLow_reader := csv.NewReader(bufio.NewReader(playUptoLow_file))   // create csv reader for PlayUptoLow.csv
	playUptoMid_reader := csv.NewReader(bufio.NewReader(playUptoMid_file))   // create csv reader for PlayUptoMid.csv
	playUptoHigh_reader := csv.NewReader(bufio.NewReader(playUptoHigh_file)) // create csv reader for PlayUptoHigh.csv

	/*
		playUptoLow_reader.Read()  // read one line not to load unnecessary data
		playUptoMid_reader.Read()  // read one line not to load unnecessary data
		playUptoHigh_reader.Read() // read one line not to load unnecessary data
	*/
	var playUptolow [][]string  // make 2-dimensional array to save PlayUptoLow.csv
	var playUptomid [][]string  // make 2-dimensional array to save PlayUptoMid.csv
	var playUptohigh [][]string // make 2-dimensional array to save PlayUptoHigh.csv

	// load playUptoLow.csv
	playUptolow = make([][]string, 0) // make slice
	var row int = 0                   // row value
	for {
		playUptoLow, err := playUptoLow_reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		playUptolow = append(playUptolow, playUptoLow)
		row++

	}

	// load PlayUptoMid.csv
	playUptomid = make([][]string, 0) // make slice
	row = 0                           // reset row to zero

	for {
		playUptoMid, err := playUptoMid_reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		playUptomid = append(playUptomid, playUptoMid)
		row++
	}

	// load PlayUptoHigh.csv
	playUptohigh = make([][]string, 0) // make slice
	row = 0
	for {
		playUptoHigh, err := playUptoHigh_reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		playUptohigh = append(playUptohigh, playUptoHigh)
		row++
	}

	/* Test code to verify that the slice is well created.
	// fmt.Println(playUptolow[0][2])
	*/
	/*
		var playUptolow [][]string


		for i := FIRST_IDX_R; i < len(playUptoLow); i++ { // initialize playUptolow
			for j := FIRST_IDX_C; j < len(playUptoLow[i]); j++ {
				playUptolow[i][j] = playUptoLow[i][j]
			}
		}
	*/
	/* Test Code to see how data looks like
	s := playUptoLow[0]
	fmt.Println(len(playUptoLow))
	fmt.Print(s)
	fmt.Println()
	fmt.Println(playUptoLow[1])
	fmt.Print(playUptoLow)
	fmt.Println()

	playUptoLow, _ = playUptoLow_reader.Read()

	s = playUptoLow[0]
	fmt.Print(s)
	fmt.Print(playUptoLow)
	fmt.Println()

	playUptoLow, _ = playUptoLow_reader.Read()
	s = playUptoLow[0]
	fmt.Print(s)
	fmt.Print(playUptoLow)
	fmt.Println()
	//	y := strcon(playUptoLow[0])
	//	fmt.Print(y)
	*/
	playUptoconst := PlayUptoConst{
		PlayUptoLowPoint:  playUptolow,
		PlayUptoMidPoint:  playUptomid,
		PlayUptoHighPoint: playUptohigh,
	}

	return playUptoconst
}

// Choose Random Answer according to current_score
func GetPlayUptoMessage(current_score int, current_idx int) string {

	var playUptoMessage string
	fmt.Println(current_idx)

	if current_score < question.BI_CRITERIA && current_score > 0 { // 1~2점 응답
		playUptoMessage = GetPlayUptoLow(current_idx)
	} else if current_score == question.BI_CRITERIA {
		playUptoMessage = GetPlayUptoMid(current_idx)
	} else if current_score > question.BI_CRITERIA && current_score <= question.SCORE_MAX {
		playUptoMessage = GetPlayUptoHigh(current_idx)
	} else {
		playUptoMessage = " "
	}
	return playUptoMessage
}

func GetPlayUptoLow(current_idx int) string {
	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // Randomly Choose Answer
	random_idx := r.Intn(len(PlayUptoMessage.PlayUptoLowPoint[current_idx])-FIRST_IDX_C) + FIRST_IDX_C
	// index range : FIRST_IDX_C ~ len(PlayUptoMessage.PlayUptoLowPoint[current_idx]) - FIRST_IDX_C

	return PlayUptoMessage.PlayUptoLowPoint[current_idx][random_idx]
}

func GetPlayUptoMid(current_idx int) string {
	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // Randomly Choose Answer
	random_idx := r.Intn(len(PlayUptoMessage.PlayUptoMidPoint[current_idx])-FIRST_IDX_C) + FIRST_IDX_C
	// index range : FIRST_IDX_C ~ len(PlayUptoMessage.PlayUptoMidPoint[current_idx]) - FIRST_IDX_C

	return PlayUptoMessage.PlayUptoMidPoint[current_idx][random_idx]
}

func GetPlayUptoHigh(current_idx int) string {
	rand_seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rand_seed) // Randomly Choose Answer

	random_idx := r.Intn(len(PlayUptoMessage.PlayUptoHighPoint[current_idx])-FIRST_IDX_C) + FIRST_IDX_C
	// index range : FIRST_IDX_C ~ len(PlayUptoMessage.PlayUptoHighPoint[current_idx]) - FIRST_IDX_C
	return PlayUptoMessage.PlayUptoHighPoint[current_idx][random_idx]
}

// debug: print 2-dimensional slice
func printArray(arr [][]string) {
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			fmt.Printf("%s ", arr[i][j])
		}
		fmt.Println()
	}
}

/*
func Decode(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalkorean.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}

	return d, nil
}
*/
