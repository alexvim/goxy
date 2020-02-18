package msg

type MessageType uint8

const (
	AuthReq MessageType = 0x1
	CmdReq              = 0x3
)

type Message interface {
	GetType() MessageType
}

type Serializeable interface {
	Serialize() []byte
}
