package respcode

import (
	"fmt"
	"net/http"
)

const (
	SUCCESS          = 200
	ERROR            = 500
	BadRequest       = 400
	Unauthorized     = 401
	Forbidden        = 403
	NotFound         = 404
	MethodNotAllowed = 405

	UserError             = 1000
	ErrorUsernameUsed     = 1001
	ErrorPasswordWrong    = 1002
	ErrorUserNotExist     = 1003
	ErrorUserInactive     = 1004
	ErrorEmailUsed        = 1005
	ErrorInvalidEmail     = 1006
	ErrorInvalidRole      = 1007
	ErrorEmptyDisplayName = 1008
	ErrorInvalidAvatarURL = 1009

	AuthError         = 2000
	ErrorTokenInvalid = 2001
	ErrorNoPermission = 2002

	CategoryError      = 3000
	ErrorCateNameUsed  = 3001
	ErrorCateNotExist  = 3002
	ErrorEmptyCateName = 3003

	ArticleError       = 4000
	ErrorArtNotExist   = 4001
	ErrorArtTitleEmpty = 4002
	ErrorArtContent    = 4003

	ErrorPasswordTooShort = 1010
)

var codeMsg = map[int]string{
	SUCCESS:               "操作成功",
	ERROR:                 "服务器错误",
	BadRequest:            "Bad Request",
	Unauthorized:          "Unauthorized",
	Forbidden:             "Forbidden",
	NotFound:              "Not Found",
	MethodNotAllowed:      "Method Not Allowed",
	UserError:             "User Error",
	ErrorUsernameUsed:     "用户名已存在！",
	ErrorPasswordWrong:    "密码错误",
	ErrorUserNotExist:     "用户不存在",
	ErrorUserInactive:     "用户账号已停用",
	ErrorEmailUsed:        "邮箱已被使用",
	ErrorInvalidEmail:     "无效的邮箱地址",
	ErrorInvalidRole:      "无效的用户权限",
	ErrorEmptyDisplayName: "显示名称不能为空",
	ErrorInvalidAvatarURL: "无效的头像URL",
	AuthError:             "认证错误",
	ErrorTokenInvalid:     "无效的认证令牌",
	ErrorNoPermission:     "没有操作权限",
	CategoryError:         "分类错误",
	ErrorCateNameUsed:     "该分类已存在",
	ErrorCateNotExist:     "该分类不存在",
	ErrorEmptyCateName:    "分类名称不能为空",
	ArticleError:          "文章错误",
	ErrorArtNotExist:      "文章不存在",
	ErrorArtTitleEmpty:    "文章标题不能为空",
	ErrorArtContent:       "文章内容不能为空",
	ErrorPasswordTooShort: "密码长度太短",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func NewError(code int, detail string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: GetErrMsg(code),
		Detail:  detail,
	}
}

func (e ErrorResponse) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}
	return e.Message
}

func IsSuccess(code int) bool {
	return code == SUCCESS
}

func IsClientError(code int) bool {
	return code >= 400 && code < 500
}

func IsServerError(code int) bool {
	return code >= 500
}

func GetHTTPStatus(code int) int {
	switch {
	case code == SUCCESS:
		return http.StatusOK
	case IsClientError(code):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
