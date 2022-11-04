package message

type MessageI interface {
	Marshal() []byte
	MsgId() uint64
	SetData(buf []byte)
	HeaderLength() uint32
	GetLength() uint32
	GetData() []byte
}
type ParseI interface {
	CheckHeader([]byte) (MessageI, error)
}
