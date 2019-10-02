package httpHandle

var testPhone string = "15711839048"

func SendTestRandom(cmd []string) {
	//SendRandomCode(testPhone, "guanzhou", "1234567890123", "1311")
	SendRandomCode(testPhone, "guangzhou", "1234", "1311")
}

func SendTestLuck(cmd []string) {
	SendLuck(testPhone, "shenzhen", "201909")
}

func SendTestUnLuck(cmd []string) {
	SendUnLuck(testPhone, "hangzhou", "201912")
}
