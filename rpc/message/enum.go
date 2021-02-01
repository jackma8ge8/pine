package message

// TypeEnum 消息类型枚举
var TypeEnum = struct {
	// TextMessage denotes a text data me  ssage. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage int
	// BinaryMessage denotes a binary data message.
	BinaryMessage int
	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage int
	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage int
	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage int
}{
	TextMessage:   1,
	BinaryMessage: 2,
	CloseMessage:  8,
	PingMessage:   9,
	PongMessage:   10,
}

// RemoterTypeEnum 消息类型枚举
var RemoterTypeEnum = struct {
	HANDLER int32
	REMOTER int32
}{
	REMOTER: 1,
	HANDLER: 2,
}

// StatusCode 消息状态码
var StatusCode = struct {
	Successful int
	Fail       int
}{
	Successful: 0,
	Fail:       200000002,
}
