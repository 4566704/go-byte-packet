package bytepacket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"

	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Packet struct {
	buffer *bytes.Buffer
}

// 新建封包
func NewPacket(b []byte) *Packet {
	return &Packet{buffer: bytes.NewBuffer(b)}
}

// 设置数据
func (p *Packet) SetData(b []byte) {
	p.buffer = bytes.NewBuffer(b)
}

// 获取数据
func (p *Packet) GetData() []byte {
	return p.buffer.Bytes()
}

// 读取int32
func (p *Packet) ReadInt32() int32 {
	var val int32 = 0
	binary.Read(p.buffer, binary.LittleEndian, &val)
	return val
}

// 写入int32
func (p *Packet) WriteInt32(val int32) {
	binary.Write(p.buffer, binary.LittleEndian, &val)
}

// 读取GBK字符串
func (p *Packet) ReadStringGbk() (string, error) {
	var len int16 = 0
	binary.Read(p.buffer, binary.LittleEndian, &len)
	buf := p.buffer.Next(int(len))
	bytes, err := GbkToUtf8(buf)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 写入GBK字符串
func (p *Packet) WriteStringGbk(str string) error {
	bytes, err := Utf8ToGbk([]byte(str))
	if err != nil {
		return err
	}
	var len int16 = int16(len(bytes))
	err = binary.Write(p.buffer, binary.LittleEndian, &len)
	if err != nil {
		return err
	}
	_, err = p.buffer.Write(bytes)
	if err != nil {
		return err
	}
	return err
}

// 读取字符串
func (p *Packet) ReadString() (string, error) {
	var len int16 = 0
	err := binary.Read(p.buffer, binary.LittleEndian, &len)
	if err != nil {
		return "", err
	}
	buf := p.buffer.Next(int(len))
	return string(buf), nil
}

// 写入字符串
func (p *Packet) WriteString(str string) error {
	bytes := []byte(str)
	var len int16 = int16(len(bytes))
	err := binary.Write(p.buffer, binary.LittleEndian, &len)
	if err != nil {
		return err
	}
	_, err = p.buffer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// 读取字节集
func (p *Packet) ReadBytes() []byte {
	var len int32 = 0
	binary.Read(p.buffer, binary.LittleEndian, &len)
	buf := p.buffer.Next(int(len))
	return buf
}

// 写入字节集
func (p *Packet) WriteBytes(data []byte) error {
	var len int32 = int32(len(data))
	err := binary.Write(p.buffer, binary.LittleEndian, &len)
	if err != nil {
		return err
	}
	_, err = p.buffer.Write(data)
	if err != nil {
		return err
	}
	return err
}

// 读取易语言日期时间
func (p *Packet) ReadElangDateTime() time.Time {
	var f float64
	binary.Read(p.buffer, binary.LittleEndian, &f)
	t := time.Date(1899, 12, 30, 0, 0, 0, 0, time.Local)
	msStr := fmt.Sprint("+", uint64(math.Floor(f*86400*1000+0.5)), "ms")
	ms, _ := time.ParseDuration(msStr)
	return t.Add(ms)
}

// GBK 转 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF-8 转 GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
