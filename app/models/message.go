package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type Message struct {
	Header         Header   `json:"header"`
	Question       Question `json:"question"`
	Answer         Answer   `json:"answer"`
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

func (h *Header) SetFlagsWithResponseBytes(responseBytes []byte) ([]byte) {
	flags := make([]byte, 2)
	//flags will contain byte1 and 2 of response byte
	flags[0] = responseBytes[2]
	flags[1] = responseBytes[3]
	// fmt.Print("------")
	// fmt.Print(flags)
	// fmt.Print("------")
	// 15 14 13 12 11 10 9 8 7 6 5 4 3 2 1 0
	flagsToBeReturned := uint16(0)
	flagsToBeReturned |= (uint16(1) << 15)
	// fmt.Printf("flag to be returned %v",flagsToBeReturned)
	//mimic next 4 bits
	for _, val := range []int{14, 13, 12, 11} {
		if binary.BigEndian.Uint16(flags)&(uint16(1)<<val) != 0 {
			flagsToBeReturned |= (uint16(1) << val)
		}
	}

	//mimic 7th bit
	opcode := false
	if binary.BigEndian.Uint16(flags)&(uint16(1)<<8) != 0 {
		flagsToBeReturned |= (uint16(1) << 8)
		opcode = true
	}

	if opcode {
		//make last 4 bits as 0010 --> 4
		flagsToBeReturned |= (uint16(1) << 1)
	}

	commonFlagsBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(commonFlagsBytes, flagsToBeReturned)

	return commonFlagsBytes
}

func (h *Header) SetRemainingDataAndReturnBytes(responseBytes []byte) []byte {
	returnResponseBytes := make([]byte, 12)
	returnResponseBytes[0] = responseBytes[0]
	returnResponseBytes[1] = responseBytes[1]

	// fmt.Print("my resp")
	// fmt.Print(returnResponseBytes)

	flagResponse := h.SetFlagsWithResponseBytes(responseBytes)
	returnResponseBytes[2] = flagResponse[0]
	returnResponseBytes[3] = flagResponse[1]

	questionBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(questionBytes, 1)
	// returnResponseBytes = append(returnResponseBytes, questionBytes...)
	returnResponseBytes[4] = questionBytes[0]
	returnResponseBytes[5] = questionBytes[1]

	answerRecordBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(answerRecordBytes, 1)
	// returnResponseBytes = append(returnResponseBytes, answerRecordBytes...)
	returnResponseBytes[6] = answerRecordBytes[0]
	returnResponseBytes[7] = answerRecordBytes[1]

	authorityRecordBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(authorityRecordBytes, 0)
	// returnResponseBytes = append(returnResponseBytes, authorityRecordBytes...)
	returnResponseBytes[8] = authorityRecordBytes[0]
	returnResponseBytes[9] = authorityRecordBytes[1]

	additionalRecordBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(additionalRecordBytes, 0)
	// returnResponseBytes = append(returnResponseBytes, additionalRecordBytes...)
	returnResponseBytes[10] = additionalRecordBytes[0]
	returnResponseBytes[11] = additionalRecordBytes[1]

	return returnResponseBytes
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

func SetName(name string) []byte {
	namee := []byte{} //the name which we gonna assign to the question
	url := strings.Split(name, ".")
	for _, val := range url {
		namee = append(namee, byte(len(val)))
		namee = append(namee, []byte(val)...)
	}
	namee = append(namee, 0x00)
	// fmt.Printf("name is %v", []byte(namee))
	return namee
}

func (q *Question) SetAllDataAndReturnQuestionBytes(name string, typ uint16, clas uint16) []byte {
	nameByte := SetName(name)

	typeBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBuf, typ)

	classBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(classBuf, clas)

	commonBuf := []byte{}
	commonBuf = append(commonBuf, nameByte...)
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
		str = fmt.Sprintf(`\x%02x`, b)
	}
	return str
}

type Answer struct {
	Name   string `json:"name"`
	Type   uint16 `json:"type"`
	Class  uint16 `json:"class"`
	TTL    uint32 `json:"ttl"`
	Length uint16 `json:"length"`
	Data   uint32 `json:"data"`
}

func (a *Answer) FillAnswerAndReturnBytes() []byte {
	nameBytes := SetName("codecrafters.io")

	typeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBytes, uint16(1))

	classBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(classBytes, uint16(1))

	ttlBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ttlBytes, uint32(60))

	lengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthBytes, uint16(4))

	dataBytes := make([]byte, 4)

	ip := "8.8.8.8"
	urlSplit := strings.Split(ip, ".")
	for _, val := range urlSplit {
		// bigEndianForm := make([]byte,1)
		num, _ := strconv.Atoi(val)
		dataBytes = append(dataBytes, byte(num))
	}

	commonBytes := []byte{}
	commonBytes = append(commonBytes, nameBytes...)
	commonBytes = append(commonBytes, typeBytes...)
	commonBytes = append(commonBytes, classBytes...)
	commonBytes = append(commonBytes, ttlBytes...)
	commonBytes = append(commonBytes, lengthBytes...)
	commonBytes = append(commonBytes, dataBytes...)

	return commonBytes

}
