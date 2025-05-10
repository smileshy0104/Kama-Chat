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
