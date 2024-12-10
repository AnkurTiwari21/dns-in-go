package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type Message struct {
	Header         Header   `json:"header"`
	Question       Question `json:"question"`
	Answer         string   `json:"answer"`
	Authority      string   `json:"authority"`
	AdditionalData string   `json:"additional"`
}

type Header struct {
	PacketIdentifier      uint16 `json:"packet_identifier"`
	Flags                 uint16 `json:"flags"`
	QuestionCount         uint16 `json:"question_count"`
	AnswerRecordCount     uint16 `json:"answer_record_count"`
	AuthorityRecordCount  uint16 `json:"authority_record_count"`
	AdditionalRecordCount uint16 `json:"additional_record_count"`
}

// 15 14 13 12 11 10 9 8 7 6 5 4 3 2 1 0
func (h *Header) SetFlags(QueryIndicator uint16, OperationCode uint16, AuthoritativeAnswer uint16, Truncation uint16, RecursionDesired uint16, RecursionAvailable uint16, Reserved uint16, ResponseCode uint16) {
	//start appending flags
	h.Flags |= (QueryIndicator << 15)
	h.Flags |= (OperationCode << 11)
	h.Flags |= (AuthoritativeAnswer << 10)
	h.Flags |= (Truncation << 9)
	h.Flags |= (RecursionDesired << 8)
	h.Flags |= (RecursionAvailable << 7)
	h.Flags |= (Reserved << 4)
	h.Flags |= (ResponseCode << 0)
}

func (h *Header) Bytes(PacketIdentifier, Flags, QuestionCount, AnswerRecordCount, AuthorityRecordCount, AdditionalRecordCount uint16) []byte {
	//start appending flags
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf[0:2], PacketIdentifier)
	binary.BigEndian.PutUint16(buf[2:4], Flags)
	binary.BigEndian.PutUint16(buf[4:6], QuestionCount)
	binary.BigEndian.PutUint16(buf[6:8], AnswerRecordCount)
	binary.BigEndian.PutUint16(buf[8:10], AuthorityRecordCount)
	binary.BigEndian.PutUint16(buf[10:12], AdditionalRecordCount)
	return buf
}

func (m *Message) Bytes(headerBytes []byte) []byte {
	b := new(bytes.Buffer)
	b.Write(headerBytes)
	return b.Bytes()
}

type Question struct {
	Name  string `json:"name"`
	Type  uint16 `json:"type"`
	Class uint16 `json:"class"`
}

func (q *Question) SetName(name string) {
	namee := `` //the name which we gonna assign to the question
	url := strings.Split(name, ".")
	for _, val := range url {
		namee += ConvertNumToHexString(uint8(len(val)))
		namee += val
	}
	namee += `\x00`
	fmt.Printf("name is %v", []byte(namee))
	q.Name = namee
}

func (q *Question) SetTypeAndClassAndReturnQuestionBytes(typ uint16, clas uint16) []byte {
	typeBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBuf, typ)

	classBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(classBuf, clas)

	commonBuf := []byte(q.Name)
	commonBuf = append(commonBuf, typeBuf...)
	// commonBuf.Write([]byte(q.Name))
	// commonBuf.Write(typeBuf)
	// commonBuf.Write(classBuf)

	commonBuf = append(commonBuf, classBuf...)
	return commonBuf
}

func ConvertNumToHexString(num uint8) string {
	buf := new(bytes.Buffer)

	// Convert the number to a byte slice (big-endian)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// Get the byte slice
	byteSlice := buf.Bytes()

	str := ``
	// Print the byte slice with \x format
	for _, b := range byteSlice {
		str = fmt.Sprintf(`\\x%02x`, b)
	}
	return str
}
