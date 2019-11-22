package app

var Message = map[int]string{
	Success:       "成功",
	Error:         "失败",
	InvalidParams: "请求参数错误",
	InterError:    "内部请求异常",
	NotFound:      "数据未找到",
}

// GetMsg 根据错误码获取错误信息
func GetMsg(code int) string {
	msg, ok := Message[code]
	if ok {
		return msg
	}

	return Message[Error]
}

