package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/codecrafters-io/dns-server-starter-go/app/models"
)

// Ensures gofmt doesn't remove the "net" import in stage 1 (feel free to remove this!)
var _ = net.ListenUDP

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	//setting up a custom resolver
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, os.Args[2])
		},
	}

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
		fmt.Print(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create an empty response
		response := models.Message{
			Question: models.Question{},
		}

		//creating header from request from client

		//setting up flag
		initialPos := 12
		domain := []string{}
		questionBytes := []byte{}
		answerBytes := []byte{}
		for initialPos < size {
			domainNameBytes, pos := DecodeDNSName(buf[:size], uint16(initialPos))
			intermediateQuestionBytes := response.Question.SetAllDataAndReturnQuestionBytes(string(domainNameBytes), 1, 1)
			questionBytes = append(questionBytes, intermediateQuestionBytes...)

			requestToResolver := models.Message{}
			headerForResolver := requestToResolver.Header.SetRemainingDataAndReturnBytes(buf[:size], int(1))

			finalQueryToSendToResolver := headerForResolver
			finalQueryToSendToResolver = append(finalQueryToSendToResolver, intermediateQuestionBytes...)

			//establish conn with resolver
			// resolverUrl := os.Args[2]
			// fmt.Printf("%v,,,,", resolverUrl)
			ips, err := resolver.LookupIP(context.Background(), "ip4", string(domainNameBytes))
			if err != nil {
				fmt.Errorf("error in ip lookup")
				return
			}
			fmt.Printf("domain %v : %s ", string(domainNameBytes), ips)

			//send this data to resolver
			//get the response
			//add that to answerbytes

			domain = append(domain, string(domainNameBytes))

			answerBytes = append(answerBytes, response.Answer.FillAnswerAndReturnBytes(string(domainNameBytes), 1, 1, 60, 4, ips[0].To4().String())...)

			initialPos = int(pos)
			initialPos += 4
		}
		// for _, val := range domain {
		// 	fmt.Print(val)
		// answerBytes = append(answerBytes, response.Answer.FillAnswerAndReturnBytes(val, 1, 1, 60, 4, "")...)
		// }

		headerBytes := response.Header.SetRemainingDataAndReturnBytes(buf[:size], int(len(domain))) //sending remaining data and getting header bytes
		responseBytes := response.Bytes(headerBytes)
		fmt.Print("response byte is --", responseBytes)
		responseBytes = append(responseBytes, questionBytes...)
		responseBytes = append(responseBytes, answerBytes...)
		fmt.Print("final response ?? ", responseBytes)
		_, err = udpConn.WriteToUDP(responseBytes, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

// [23 200 1 0 0 2 0 0 0 0 0 0   3 97 98 99 17 108 111 110 103 97 115 115 100 111 109 97 105 110 110 97 109 101 3 99 111 109 0 0   1 0 1 3 100 101 102 192 16 0  1   0 1]
//  0  1   2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18  19  20 21  22  23  24  25  26  27 28  29  30  31 32  33  34 35 36 37 38 39 40 414243 44 45  46  47  48 49 50 51 52

func DecodeDNSName(data []byte, start uint16) (string, uint16) {
	var name bytes.Buffer
	pos := start

	for {
		// fmt.Print("pos is ", pos)
		// fmt.Print(" ")
		length := uint16(data[pos])
		if length == 0 {
			// End of the name (0x00)
			pos++
			break
		}
		//if msb 2 bits are set then add a dot and make start as the number with msb removed
		if ((length & (uint16(1) << 7)) != 0) && ((length & (uint16(1) << 6)) != 0) {
			//msb 2 bits are simultaneously set
			bufferTransport := []byte{}
			bufferTransport = append(bufferTransport, data[pos])
			bufferTransport = append(bufferTransport, data[pos+1])

			transfer := (binary.BigEndian.Uint16(bufferTransport) & 0x3FFF)
			str, _ := DecodeDNSName(data, transfer)
			pos += 2
			return name.String() + str, pos
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
