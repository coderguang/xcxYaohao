package httpHandle

import (
	"strconv"
)

const (
	HTTP_ARGS_KEY            string = "op"
	HTTP_ARGS_TIME           string = "time"
	HTTP_ARGS_SEARCH         string = "search"
	HTTP_ARGS_CODE           string = "code"
	HTTP_ARGS_CITY           string = "city"
	HTTP_ARGS_MATCH_KEY      string = "key"
	HTTP_ARGS_BIND_REQUIRE   string = "require"
	HTTP_ARGS_BIND_CONFIRM   string = "confirm"
	HTTP_ARGS_BIND_GET_DATA  string = "getData"
	HTTP_ARGS_BIND_CANCEL    string = "cancel"
	HTTP_ARGS_SHARE          string = "share"
	HTTP_ARGS_BIND_ONE_KEY   string = "reBind"
	HTTP_ARGS_BIND_PHONE     string = "phone"
	HTTP_ARGS_BIND_CARD_TYPE string = "card"
	HTTP_ARGS_BIND_CODE      string = "bindCode"
	HTTP_ARGS_DATA           string = "data"
)

const (
	HTTP_RETURN_ERR_CODE string = "errcode"
	HTTP_RETURN_TIME     string = "time"
	HTTP_RETURN_Data     string = "data"
	HTTP_RETURN_STATUS   string = "status"
	HTTP_RETURN_TIPS     string = "tips"
)

type YaoHaoNoticeError int

const (
	YAOHAO_OK                                    YaoHaoNoticeError = iota //0
	YAOHAO_ERR_DATA_EXISTS                                                //1  数据仍然有效
	YAOHAO_ERR_TITLE                                                      //2  不支持的城市
	YAOHAO_ERR_TOKEN                                                      //3  token错误
	YAOHAO_ERR_PHONE                                                      //4  phone错误
	YAOHAO_ERR_CODE                                                       //5  code错误
	YAOHAO_ERR_LEFT_TIME                                                  //6  剩余时长错误
	YAOHAO_ERR_GM_LIMIT                                                   //7  验证码错误次数过多
	YAOHAO_ERR_CONFIRM_MORE_TIMES                                         //8
	YAOHAO_ERR_TOKEN_STILL_VALID                                          //9  数据依然有效
	YAOHAO_ERR_CODE_STILL_VALID                                           //10 数据依然有效
	YAOHAO_ERR_PHONE_STILL_VALID                                          //11 数据依然有效
	YAOHAO_ERR_REQUIRE_HAD_CONFIRM                                        //12 已应答
	YAOHAO_ERR_REQUIRE_WAIT_ANSWER                                        //13 等待应答
	YAOHAO_ERR_REQUIRE_HAD_LOCK                                           //14 锁定
	YAOHAO_ERR_CONFIRM_NOT_REQUIRE                                        //15 未请求
	YAOHAO_ERR_CONFIRM_RANDOMCODE                                         //16 错误的验证码
	YAOHAO_ERR_HTTP_NO_KEY                                                //17 http参数错误
	YAOHAO_ERR_HTTP_PARAM_NUM                                             //18 http参数错误
	YAOHAO_ERR_HTTP_REQ_TYPE                                              //19 http type 错误
	YAOHAO_ERR_HTTP_REQ_MAX_TIMES                                         //20 请求过于频繁
	YAOHAO_ERR_HTTP_RANDOM_CODE_TIME_OUT                                  //21 验证码已过期
	YAOHAO_ERR_SMS_CLIENT                                                 //22 SMS错误
	YAOHAO_ERR_SMS_PROCESS                                                //23 SMS错误
	YAOHAO_ERR_SMS_OTHER                                                  //24 SMS错误
	YAOHAO_ERR_SMS_RESULT_PARSE_ERROR                                     //25 SMS发送失败
	YAOHAO_ERR_YAOHAO_SERVER_INT_FORMAT                                   //26
	YAOHAO_ERR_YAOHAO_SERVER_DATA_SIZE_NOT_MATCH                          //27
	YAOHAO_ERR_YAOHAO_SERVER_TIME_NOT_MATCH                               //28
	YAOHAO_ERR_HTTP_REQ_CARD_TYPE                                         //29 cardType错误
	YAOHAO_ERR_OPEN_ID_PARAM_NUM                                          //30 token错误
	YAOHAO_ERR_NOT_BIND_DATA                                              //31 没有绑定数据
	YAOHAO_ERR_STATUS_NOT_NORMAL                                          //32 status异常
	YAOHAO_ERR_SMS_SERVER_CLOSE                                           //33 sms服务关闭
	YAOHAO_ERR_PHONE_BIND_TOO_MANY                                        //34 phone绑定过多
	YAOHAO_ERR_WX_ERROR_CODE                                              //35 wx.login code error
	YAOHAO_ERR_DO_NOT_THING                                               //36 default
	YAOHAO_ERR_NO_MATCH_DATA                                              //37 not match code or name data
	YAOHAO_ERR_CODE_HAD_LUCK                                              //38 已中签
)

func errorToStr(index YaoHaoNoticeError) string {
	return strconv.Itoa(int(index))
}
