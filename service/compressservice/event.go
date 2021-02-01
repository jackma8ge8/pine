package compressservice

type eventCompress struct {
	eventToCode map[string]byte
	codeToEvent map[byte]string
	events      []string
}

// AddEventRecord 添加需要压缩的客户端监听的事件
func (ec *eventCompress) AddEventRecord(eventName string) {
	if _, exist := ec.eventToCode[eventName]; !exist {
		code := byte(len(ec.eventToCode) + 1)
		ec.eventToCode[eventName] = code
		ec.codeToEvent[code] = eventName
		ec.events = append(ec.events, eventName)
	}
}

// AddEventCompressRecords 添加路由压缩记录
func (ec *eventCompress) AddEventCompressRecords(eventNames ...string) {
	for _, eventName := range eventNames {
		ec.AddEventRecord(eventName)
	}
}

// GetEventByCode 获取真实Event
// func GetEventByCode(code byte) string {
// 	if value, exist := codeToEvent[code]; exist {
// 		return value
// 	}
// 	return ""
// }

// // GetCodeByEvent 获取真实Event对应的Code
func (ec *eventCompress) GetCodeByEvent(eventName string) byte {
	if value, exist := ec.eventToCode[eventName]; exist {
		return value
	}
	return 0
}

// GetEvents 获取Events切片
func (ec *eventCompress) GetEvents() []string {
	return ec.events
}

// Event event map
var Event = eventCompress{
	eventToCode: make(map[string]byte),
	codeToEvent: make(map[byte]string),
	events:      make([]string, 0, 10),
}
