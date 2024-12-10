package main

import (
	"bytes"
	"fmt"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/app/models"
)

// Ensures gofmt doesn't remove the "net" import in stage 1 (feel free to remove this!)
var _ = net.ListenUDP

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		// fmt.Print(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create an empty response
		response := models.Message{
			Question: models.Question{},
		}
		bytesExceptHeader := buf[12:size]
		domainNameBytes,_ := DecodeDNSName(bytesExceptHeader,0)
		                             //setting up flag
		headerBytes := response.Header.SetRemainingDataAndReturnBytes(buf[:size]) //sending remaining data and getting header bytes
		responseBytes := response.Bytes(headerBytes)

		questionBytes := response.Question.SetAllDataAndReturnQuestionBytes(string(domainNameBytes), 1, 1)
		responseBytes = append(responseBytes, questionBytes...) //appending question bytes

		answerBytes := response.Answer.FillAnswerAndReturnBytes(string(domainNameBytes))
		responseBytes = append(responseBytes, answerBytes...)

		// fmt.Print("<------->")
		// fmt.Print(responseBytes)
		// fmt.Print(len(responseBytes))
		// fmt.Print("<>")
		_, err = udpConn.WriteToUDP(responseBytes, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func DecodeDNSName(data []byte, start int) (string, int) {
	var name bytes.Buffer
	pos := start

	for {
		length := int(data[pos])
		if length == 0 {
			// End of the name (0x00)
			pos++
			break
		}
		pos++
		name.Write(data[pos : pos+length])
		pos += length
		if data[pos] != 0 {
			name.WriteByte('.')
		}
	}
	return name.String(), pos
}
