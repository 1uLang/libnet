package message

type MessageI interface {
	Marshal() []byte
	MsgId() uint64
	HeaderLength() uint32
	GetLength() uint32
	SetData(buf []byte)
}
type ParseI interface {
	CheckHeader([]byte) (MessageI, error)
}
