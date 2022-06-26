package gobencode

import (
	"bufio"
	"bytes"
	"os"
)

type Decoder struct {
	b      *bufio.Reader
	Cursor int64
}

func DecodeFile(filepath string) BencodeItem {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	d := NewDecoder(reader)
	return d.decode()
}

func DecodeString(bencode string) BencodeItem {
	buff := new(bytes.Buffer)
	buff.WriteString(bencode)

	reader := bufio.NewReader(buff)
	d := NewDecoder(reader)
	return d.decode()
}

func (d *Decoder) ReadByte() byte {
	b, err := d.b.ReadByte()
	if err != nil {
		panic(err)
	}

	d.Cursor += 1
	return b
}

func (d *Decoder) ReadBytes(delim byte) []byte {
	b, err := d.b.ReadBytes(delim)
	if err != nil {
		panic(err)
	}

	d.Cursor += int64(len(b))
	return b
}

func (d *Decoder) Read(p []byte) {
	n, err := d.b.Read(p)
	if err != nil {
		panic(err)
	}

	d.Cursor += int64(n)
}

func NewDecoder(r *bufio.Reader) Decoder {
	return Decoder{
		r,
		0,
	}
}
