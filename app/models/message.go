package models

import (
	"bytes"
	"encoding/binary"
)

type Message struct {
	Header         Header `json:"header"`
	Question       string `json:"question"`
	Answer         string `json:"answer"`
	Authority      string `json:"authority"`
	AdditionalData string `json:"additional"`
}

type Header struct {
	PacketIdentifier uint16 `json:"packet_identifier"`
	// QueryIndicator        uint8  `json:"query_indicator"`
	// OperationCode         uint8  `json:"opcode"`
	// AuthoritativeAnswer   uint8  `json:"aa"`
	// Truncation            uint8  `json:"truncation"`
	// RecursionDesired      uint8  `json:"recursion_desired"`
	// RecursionAvailable    uint8  `json:"recursion_available"`
	// Reserved              uint8  `json:"reserved"`
	// ResponseCode          uint8  `json:"response_code"`
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

// type Message struct {
// 	Header *Header
// }

func (m *Message) Bytes() []byte {
	b := new(bytes.Buffer)
	b.Write(m.Header.Bytes(1234,m.Header.Flags,0,0,0,0))
	return b.Bytes()
}

// type Header struct {
// 	ID uint16 // Packet Identifier; A random identifier is assigned to query packets. Response packets must reply with the same id. This is needed to differentiate responses due to the stateless nature of UDP.
// 	/*
// 		1 bit; QR (Query Response) - 0 for queries, 1 for responses
// 		4 bits; OPCODE (Operation Code) - Typically always 0. RFC1035 for details.
// 		1 bit; AA (Authoritative Answer) - Set to 1 if the responding server is authoritative - that is, it "owns" - the domain queried
// 		1 bit; TC (Truncated Message) - Set to 1 if the message length exceeds 512 bytes. Traditionally a hint that the query can be reissued using TCP, for which the length limitation doesn't apply.
// 		1 bit; RD (Recursion Desired) - Set by the sender of the request if the server should attempt to resolve the query recursively if it does not have an answer readily available.
// 		1 bit; RA (Recursion Available) - Set by the server to indicate whether or not recursive queries are allowed.
// 		3 bits; Z (Reserved) - Originally reserved for later use, but now used for DNSSEC queries.
// 		4 bits; RCODE (Response Code) - Set by the server to indicate the status of the response, i.e. whether or not it was successful or failed, and in the latter case providing details about the cause of the failure.
// 	*/
// 	Flags   uint16
// 	QDCOUNT uint16 // Question Count; The number of entries in the Question Section
// 	ANCOUNT uint16 // Answer Count; The number of entries in the Answer Section
// 	NSCOUNT uint16 // Authority Count; The number of entries in the Authority Section
// 	ARCOUNT uint16 // Additional Count; The number of entries in the Additional Section
// }

// func (h *Header) SetFlags(qr, opcode, aa, tc, rd, ra, z, rcode uint16) {
// 	h.Flags |= qr << 15     // QR at bit 15
// 	h.Flags |= opcode << 11 // OPCODE at bit 11 - 14
// 	h.Flags |= aa << 10     // AA at bit 10
// 	h.Flags |= tc << 9      // TC at bit 9
// 	h.Flags |= rd << 8      // RD at bit 8
// 	h.Flags |= ra << 7      // RA at bit 7
// 	h.Flags |= z & 0xf << 4 // Z at bit 4 - 6
// 	h.Flags |= rcode & 0xf  // RCODE at bit 0 - 3
// }
// func (h *Header) Bytes() []byte {
// 	buf := make([]byte, 12)
// 	binary.BigEndian.PutUint16(buf[0:2], h.ID)
// 	binary.BigEndian.PutUint16(buf[2:4], h.Flags)
// 	binary.BigEndian.PutUint16(buf[4:6], h.QDCOUNT)
// 	binary.BigEndian.PutUint16(buf[6:8], h.ANCOUNT)
// 	binary.BigEndian.PutUint16(buf[8:10], h.NSCOUNT)
// 	binary.BigEndian.PutUint16(buf[10:12], h.ARCOUNT)
// 	return buf
// }
