package simulate

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	utilbyte "lucy/utils/byte"
)

var (
	pool = utilbyte.NewBufferPool()
)

type Packet struct {
	Version int8
	CmdId   int16
	BodyLen int32
	Body    []byte
}

func NewPacket(cmdId int32, body proto.Message) Packet {
	pack := Packet{}
	pack.Version = int8(1)
	pack.CmdId = int16(cmdId)
	m, _ := proto.Marshal(body)
	pack.BodyLen = int32(len(m))
	pack.Body = m
	//fmt.Printf("%v\n", pack)
	return pack
}

func (p Packet) encode() ([]byte, error) {
	buf := pool.Get()
	defer pool.Put(buf)

	err := binary.Write(buf, binary.BigEndian, p.Version)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.CmdId)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.BodyLen)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
