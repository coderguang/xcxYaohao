package httpHandle

import (
	"encoding/json"
	"net/http"

	"github.com/coderguang/GameEngine_go/sglog"
)

func logicHandle(w http.ResponseWriter, r *http.Request, flag chan bool) {
	defer func() {
		flag <- true
	}()

	r.ParseForm()

	sglog.Debug("require data is ", r.Form)

	op := r.FormValue(HTTP_OP_KEY)
	returnData := make(map[string]interface{})
	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_NOTICE_ERR_DO_NOT_THING
	switch op {
	case HTTP_OP_TIME:
		requireLastestTime(r, returnData)
	}

	str, err := json.Marshal(returnData)
	if err != nil {
		sglog.Error("parse returnData to string error", err)
		return
	}
	sglog.Info("return str is", string(str))
	w.Write([]byte(string(str)))
}
