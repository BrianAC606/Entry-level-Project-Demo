package dns_body

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type Header struct {
	TranscationId uint16
	Flags         uint16
	QuestionNum   uint16
	AnswerNum     uint16
	AuthorityNum  uint16
	AdditionalNum uint16
}

func NewHeader() *Header {
	return &Header{
		TranscationId: 0xFFCA,
		Flags:         0,
		QuestionNum:   1,
		AnswerNum:     0,
		AuthorityNum:  0,
		AdditionalNum: 0,
	}
}

func (h *Header) SetFlags(QR, Opcode, AA, TC, RD, RA, rcode uint16) {
	h.Flags = QR<<15 + Opcode<<11 + AA<<10 + TC<<9 + RD<<8 + RA<<7 + rcode
}

type queries struct {
	domain string
	Type   uint16
	Class  uint16
}

func NewQueries(domainName string) *queries {
	return &queries{
		domain: domainName,
		Type:   1,
		Class:  1,
	}
}

func (q *queries) GetQueriesBtyes() []byte {
	var (
		buffer   bytes.Buffer
		segments = strings.Split(q.domain, ".")
	)
	for _, seg := range segments {
		binary.Write(&buffer, binary.BigEndian, byte(len(seg)))
		binary.Write(&buffer, binary.BigEndian, []byte(seg))
	}
	binary.Write(&buffer, binary.BigEndian, byte(0))
	binary.Write(&buffer, binary.BigEndian, uint16(q.Type))
	binary.Write(&buffer, binary.BigEndian, uint16(q.Class))
	return buffer.Bytes()
}

func ParseResponse(res []byte) {

}
