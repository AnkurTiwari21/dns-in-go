package models

type Message struct {
	Header         Header `json:"header"`
	Question       string `json:"question"`
	Answer         string `json:"answer"`
	Authority      string `json:"authority"`
	AdditionalData string `json:"additional"`
}

type Header struct {
	PacketIdentifier      uint16 `json:"packet_identifier"`
	QueryIndicator        uint8  `json:"query_indicator"`
	OperationCode         uint8  `json:"opcode"`
	AuthoritativeAnswer   uint8  `json:"aa"`
	Truncation            uint8  `json:"truncation"`
	RecursionDesired      uint8  `json:"recursion_desired"`
	RecursionAvailable    uint8  `json:"recursion_available"`
	Reserved              uint8  `json:"reserved"`
	ResponseCode          uint8  `json:"response_code"`
	QuestionCount         uint16 `json:"question_count"`
	AnswerRecordCount     uint16 `json:"answer_record_count"`
	AuthorityRecordCount  uint16 `json:"authority_record_count"`
	AdditionalRecordCount uint16 `json:"additional_record_count"`
}
