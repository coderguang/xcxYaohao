package httpHandle

var testPhone string = "15711839048"

func SendTestRandom(cmd []string) {
	SendRandomCode(testPhone, "1311")
}

func SendTestLuck(cmd []string) {
	SendLuck(testPhone, "123141414", "201909")
}

func SendTestUnLuck(cmd []string) {
	SendUnLuck(testPhone, "123141414")
}
