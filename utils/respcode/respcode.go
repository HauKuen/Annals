package respcode

const (
	SUCCESS = 200
	ERROR   = 500

	ErrorUsernameUsed     = 1001
	ErrorPasswordWrong    = 1002
	ErrorUserNotExist     = 1003
	ErrorTokenExist       = 1004
	ErrorTokenRuntime     = 1005
	ErrorTokenWrong       = 1006
	ErrorTokenTypeWrong   = 1007
	ErrorUserNoRight      = 1008
	ErrorEmailUsed        = 1009
	ErrorInvalidEmail     = 1010
	ErrorInvalidRole      = 1011
	ErrorEmptyDisplayName = 1012
	ErrorInvalidAvatarURL = 1013
	ErrorUserInactive     = 1014

	ErrorArtNotExist = 2001

	ErrorCatenameUsed = 3001
	ErrorCateNotExist = 3002
)

var codeMsg = map[int]string{
	SUCCESS:               "OK",
	ERROR:                 "FAIL",
	ErrorUsernameUsed:     "用户名已存在！",
	ErrorPasswordWrong:    "密码错误",
	ErrorUserNotExist:     "用户不存在",
	ErrorTokenExist:       "TOKEN不存在,请重新登陆",
	ErrorTokenRuntime:     "TOKEN已过期,请重新登陆",
	ErrorTokenWrong:       "TOKEN不正确,请重新登陆",
	ErrorTokenTypeWrong:   "TOKEN格式错误,请重新登陆",
	ErrorUserNoRight:      "该用户无权限",
	ErrorEmailUsed:        "邮箱已被使用",
	ErrorInvalidEmail:     "无效的邮箱地址",
	ErrorInvalidRole:      "无效的用户权限",
	ErrorEmptyDisplayName: "显示名称不能为空",
	ErrorInvalidAvatarURL: "无效的头像URL",
	ErrorUserInactive:     "用户账号已停用",

	ErrorArtNotExist: "文章不存在",

	ErrorCatenameUsed: "该分类已存在",
	ErrorCateNotExist: "该分类不存在",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
