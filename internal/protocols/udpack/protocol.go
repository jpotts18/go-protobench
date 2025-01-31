package udp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

const (
	MagicBytes = 0x4242  // Protocol identifier
	Version    = 1       // Protocol version
	HeaderSize = 12      // 2 magic + 2 version + 4 payload size + 4 checksum
)

type Header struct {
	Magic      uint16
	Version    uint16
	PayloadLen uint32
	Checksum   uint32
}

type MessageBody struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Number  int64  `json:"number"`
}

func calculateChecksum(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

func EncodeMessage(payload []byte) ([]byte, error) {
	checksum := calculateChecksum(payload)
	
	header := Header{
		Magic:      MagicBytes,
		Version:    Version,
		PayloadLen: uint32(len(payload)),
		Checksum:   checksum,
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, header); err != nil {
		return nil, fmt.Errorf("failed to write header: %w", err)
	}

	buf.Write(payload)
	return buf.Bytes(), nil
}

func DecodeMessage(data []byte) ([]byte, error) {
	if len(data) < HeaderSize {
		return nil, fmt.Errorf("message too short")
	}

	var header Header
	reader := bytes.NewReader(data[:HeaderSize])
	if err := binary.Read(reader, binary.BigEndian, &header); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	if header.Magic != MagicBytes {
		return nil, fmt.Errorf("invalid magic bytes")
	}

	if header.Version != Version {
		return nil, fmt.Errorf("unsupported version: %d", header.Version)
	}

	expectedLen := int(header.PayloadLen) + HeaderSize
	if len(data) != expectedLen {
		return nil, fmt.Errorf("invalid message length: got %d, want %d", len(data), expectedLen)
	}

	payload := data[HeaderSize:]
	checksum := calculateChecksum(payload)
	if checksum != header.Checksum {
		return nil, fmt.Errorf("checksum mismatch: got %x, want %x", checksum, header.Checksum)
	}

	return payload, nil
}

// FormatMessage creates a simple ASCII representation of the message
func FormatMessage(msg MessageBody) string {
	return fmt.Sprintf("ID:%s|CONTENT:%s|NUMBER:%d", msg.ID, msg.Content, msg.Number)
}

// ParseMessage parses the ASCII representation back into a MessageBody
func ParseMessage(data string) (MessageBody, error) {
	var msg MessageBody
	var id, content string
	var number int64
	
	_, err := fmt.Sscanf(data, "ID:%s|CONTENT:%s|NUMBER:%d", &id, &content, &number)
	if err != nil {
		return msg, fmt.Errorf("failed to parse message: %w", err)
	}

	msg.ID = id
	msg.Content = content
	msg.Number = number
	return msg, nil
} 
