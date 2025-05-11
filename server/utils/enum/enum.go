package enum

const (
	NORMAL   = iota // 正常
	DISABLE         // 禁用
	DISSOLVE        // 解散
)

// contact_type_enum 联系类型
const (
	USER  = iota // 用户
	GROUP        // 群聊
)

// contact_status_enum 联系状态
const (
	NORMAL_ = iota
	BE_BLACK
	BLACK
	BE_DELETE
	DELETE
	SILENCE
	QUIT_GROUP
	KICK_OUT_GROUP
)

// contact_apply_status_enum 申请状态
const (
	PENDING = iota
	AGREE
	REFUSE
	BLACK_
)

// message_status_enum 消息状态
const (
	// 未发送
	Unsent = iota
	// 已发送
	Sent
)

// message_type_enum 消息类型
const (
	Text = iota
	// 语音
	Voice
	// 文件
	File
	// 通话
	AudioOrVideo
)
