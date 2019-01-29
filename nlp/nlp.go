package nlp

func ConvertInquiryScore(str string) string {
	switch str {
	case "일":
		str = "1"
	case "이":
		str = "2"
	case "삼":
		str = "3"
	case "사":
		str = "4"
	case "오":
		str = "5"
	}
	return str
}
