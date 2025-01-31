package udpfast

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

const (
	MagicBytes = 0x4343  // Different magic bytes from regular UDP
	Version    = 1
	HeaderSize = 12
)

type Header struct {
	Magic      uint16
	Version    uint16
	PayloadLen uint32
	Checksum   uint32
}

type MessageBody struct {
	ID      string
	Content string
	Number  int64
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

func FormatMessage(msg MessageBody) string {
	return fmt.Sprintf("ID:%s|CONTENT:%s|NUMBER:%d", msg.ID, msg.Content, msg.Number)
} 
