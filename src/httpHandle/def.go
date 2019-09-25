package httpHandle

import (
	"strconv"
)

const (
	HTTP_OP_KEY    string = "op"
	HTTP_OP_TIME   string = "time"
	HTTP_OP_SEARCH string = "search"
	HTTP_OP_CODE   string = "code"
)

const (
	HTTP_RETURN_ERR_CODE   string = "errcode"
	HTTP_RETURN_TIME       string = "time"
	HTTP_RETURN_MATCH_Data string = "data"
)

type YaoHaoNoticeError int

const (
	YAOHAO_NOTICE_OK                                    YaoHaoNoticeError = iota //0
	YAOHAO_NOTICE_ERR_DATA_EXISTS                                                //1
	YAOHAO_NOTICE_ERR_TITLE                                                      //2
	YAOHAO_NOTICE_ERR_TOKEN                                                      //3
	YAOHAO_NOTICE_ERR_PHONE                                                      //4
	YAOHAO_NOTICE_ERR_CODE                                                       //5
	YAOHAO_NOTICE_ERR_LEFT_TIME                                                  //6
	YAOHAO_NOTICE_ERR_GM_LIMIT                                                   //7
	YAOHAO_NOTICE_ERR_CONFIRM_MORE_TIMES                                         //8
	YAOHAO_NOTICE_ERR_TOKEN_STILL_VALID                                          //9
	YAOHAO_NOTICE_ERR_CODE_STILL_VALID                                           //10
	YAOHAO_NOTICE_ERR_PHONE_STILL_VALID                                          //11
	YAOHAO_NOTICE_ERR_REQUIRE_HAD_CONFIRM                                        //12 已应答
	YAOHAO_NOTICE_ERR_REQUIRE_WAIT_ANSWER                                        //13 等待应答
	YAOHAO_NOTICE_ERR_REQUIRE_HAD_LOCK                                           //14 锁定
	YAOHAO_NOTICE_ERR_CONFIRM_NOT_REQUIRE                                        //15 未请求
	YAOHAO_NOTICE_ERR_CONFIRM_RANDOMCODE                                         //16 错误的验证码
	YAOHAO_NOTICE_ERR_HTTP_NO_KEY                                                //17
	YAOHAO_NOTICE_ERR_HTTP_PARAM_NUM                                             //18
	YAOHAO_NOTICE_ERR_HTTP_REQ_TYPE                                              //19
	YAOHAO_NOTICE_ERR_HTTP_REQ_MAX_TIMES                                         //20
	YAOHAO_NOTICE_ERR_HTTP_RANDOM_CODE_TIME_OUT                                  //21
	YAOHAO_NOTICE_ERR_SMS_CLIENT                                                 //22
	YAOHAO_NOTICE_ERR_SMS_PROCESS                                                //23
	YAOHAO_NOTICE_ERR_SMS_OTHER                                                  //24
	YAOHAO_NOTICE_ERR_SMS_RESULT_PARSE_ERROR                                     //25
	YAOHAO_NOTICE_ERR_YAOHAO_SERVER_INT_FORMAT                                   //26
	YAOHAO_NOTICE_ERR_YAOHAO_SERVER_DATA_SIZE_NOT_MATCH                          //27
	YAOHAO_NOTICE_ERR_YAOHAO_SERVER_TIME_NOT_MATCH                               //28
	YAOHAO_NOTICE_ERR_HTTP_REQ_CARD_TYPE                                         //29
	YAOHAO_NOTICE_ERR_OPEN_ID_PARAM_NUM                                          //30
	YAOHAO_NOTICE_ERR_NOT_BIND_DATA                                              //31
	YAOHAO_NOTICE_ERR_STATUS_NOT_NORMAL                                          //32
	YAOHAO_NOTICE_ERR_SMS_SERVER_CLOSE                                           //33
	YAOHAO_NOTICE_ERR_PHONE_BIND_TOO_MANY                                        //34
	YAOHAO_NOTICE_ERR_WX_ERROR_CODE                                              //35 wx.login code error
	YAOHAO_NOTICE_ERR_DO_NOT_THING                                               //36 default
)

func errorToStr(index YaoHaoNoticeError) string {
	return strconv.Itoa(int(index))
}
