package httpHandle

import (
	"net/http"

	"github.com/coderguang/GameEngine_go/sglog"
)

type web_server struct{}

func (h *web_server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chanFlag := make(chan bool)
	go logicHandle(w, r, chanFlag)
	<-chanFlag
}

func NewWebServer(port string) {
	http.Handle("/", &web_server{})
	port = "0.0.0.0:" + port
	sglog.Info("start web server.listen port:", port)
	http.ListenAndServe(port, nil)
}
