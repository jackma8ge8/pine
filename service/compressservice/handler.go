package compressservice

type handlerCompress struct {
	handlerToCode map[string]byte
	codeToHandler map[byte]string
	handlers      []string
}

// AddHandlerRecord 添加记录
func (hs *handlerCompress) AddHandlerRecord(handlerName string) {
	if _, exist := hs.handlerToCode[handlerName]; !exist {
		code := byte(len(hs.handlerToCode) + 1)
		hs.handlerToCode[handlerName] = code
		hs.codeToHandler[code] = handlerName
		hs.handlers = append(hs.handlers, handlerName)
	}
}

// GetHandlerByCode 获取真实Handler
func (hs *handlerCompress) GetHandlerByCode(code byte) string {
	if value, exist := hs.codeToHandler[code]; exist {
		return value
	}
	return ""
}

// GetCodeByHandler 获取真实Handler对应的Code
func (hs *handlerCompress) GetCodeByHandler(handlerName string) byte {
	if value, exist := hs.handlerToCode[handlerName]; exist {
		return value
	}
	return 0
}

// GetHandlers 获取Handlers切片
func (hs *handlerCompress) GetHandlers() []string {
	return hs.handlers
}

// AddEventRecord 添加需要压缩的客户端监听的事件
func (hs *handlerCompress) AddEventRecord(eventName string) {
	if _, exist := hs.handlerToCode[eventName]; !exist {
		code := byte(len(hs.handlerToCode) + 1)
		hs.handlerToCode[eventName] = code
		hs.codeToHandler[code] = eventName
		hs.handlers = append(hs.handlers, eventName)
	}
}

// Handler handler map
var Handler = handlerCompress{
	handlerToCode: make(map[string]byte),
	codeToHandler: make(map[byte]string),
	handlers:      make([]string, 0, 10),
}
