package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"strconv"
)

//Header holds the first byte (aka magic byte)
type Header struct {
	Header uint64 //  0:8
}

//Chunk represents a data byte chunk
type Chunk struct {
	Size uint32
	Type uint32
	Data []byte
	CRC  uint32
}

//MetaChunk inherits a Chunk struct
type MetaChunk struct {
	Chk    Chunk
	Offset int64
}

const key string = "stegasaurs"
const offset string = "0x85258"
const chunkType string = "rNDm"

func (mc *MetaChunk) EncodeImage(b *bytes.Reader, data string) error {
	err := mc.validate(b)

	if err != nil {
		return err
	}

	var m MetaChunk
	m.Chk.Data = XorEncode([]byte(data), key)
	m.Chk.Type = m.strToInt(chunkType)
	m.Chk.Size = m.createChunkSize()
	m.Chk.CRC = m.createChunkCRC()
	bm := m.marshalData()
	bmb := bm.Bytes()
	fmt.Printf("Payload Original: % X\n", []byte(data))
	fmt.Printf("Payload Encode: % X\n", m.Chk.Data)
	WriteData(b, bmb)

	return nil
}

func (mc *MetaChunk) DecodeImage(b *bytes.Reader) ([]byte, error) {
	err := mc.validate(b)

	if err != nil {
		return nil, err
	}

	var m MetaChunk
	offset, _ := strconv.ParseInt(offset, 0, 64)
	b.Seek(offset, 0)
	m.readChunk(b)
	origData := m.Chk.Data
	m.Chk.Data = XorDecode(m.Chk.Data, key)
	m.Chk.CRC = m.createChunkCRC()
	fmt.Printf("Payload Original: % X\n", origData)
	fmt.Printf("Payload Decode: % X\n", m.Chk.Data)

	return m.Chk.Data, nil
}

func (mc *MetaChunk) marshalData() *bytes.Buffer {
	bytesMSB := new(bytes.Buffer)
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Size); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.CRC); err != nil {
		log.Fatal(err)
	}

	return bytesMSB
}

func (mc *MetaChunk) readChunk(b *bytes.Reader) {
	mc.readChunkSize(b)
	mc.readChunkType(b)
	mc.readChunkBytes(b, mc.Chk.Size)
	mc.readChunkCRC(b)
}

func (mc *MetaChunk) readChunkSize(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Size); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkType(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkBytes(b *bytes.Reader, cLen uint32) {
	mc.Chk.Data = make([]byte, cLen)
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkCRC(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.CRC); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) validate(b *bytes.Reader) error {
	var header Header

	if err := binary.Read(b, binary.BigEndian, &header.Header); err != nil {
		log.Fatal(err)
	}

	bArr := make([]byte, 8)
	binary.BigEndian.PutUint64(bArr, header.Header)

	if string(bArr[1:4]) != "PNG" {
		errMsg := "provided file is not a valid PNG format"
		fmt.Println(errMsg)
		return errors.New(errMsg)
	} else {
		fmt.Println("Valid PNG so let us continue!")
		return nil
	}
}

func (mc *MetaChunk) createChunkSize() uint32 {
	return uint32(len(mc.Chk.Data))
}

func (mc *MetaChunk) createChunkCRC() uint32 {
	bytesMSB := new(bytes.Buffer)
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
	return crc32.ChecksumIEEE(bytesMSB.Bytes())
}

func (mc *MetaChunk) strToInt(s string) uint32 {
	t := []byte(s)
	return binary.BigEndian.Uint32(t)
}
