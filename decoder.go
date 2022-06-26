package gobencode

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
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

func (d *Decoder) decode() BencodeItem {
	ch := d.ReadByte()
	bi := BencodeItem{}

	if '0' <= ch && ch <= '9' {
		lengthBytes := d.ReadBytes(':')
		lengthBytes = append([]byte{ch}, lengthBytes[:len(lengthBytes)-1]...)
		length, err := strconv.Atoi(string(lengthBytes))
		if err != nil {
			panic(err)
		}

		bytes := make([]byte, length)
		d.Read(bytes)

		bi.Type = Bytes
		bi.Value = bytes
		return bi
	}

	switch ch {
	case 'i':
		rawInt := d.ReadBytes('e')
		intVal, err := strconv.Atoi(string(rawInt[:len(rawInt)-1]))
		if err != nil {
			panic(err)
		}

		bi.Type = Integer
		bi.Value = intVal

	case 'l':
		list := []BencodeItem{}

		for {
			item := d.decode()
			if item.Type == EOL {
				break
			}

			list = append(list, item)
		}

		bi.Type = List
		bi.Value = list

	case 'd':
		dict := map[string]BencodeItem{}

		for {
			key := d.decode()
			if key.Type == EOL {
				break
			} else if key.Type != Bytes {
				panic(fmt.Sprintf("dict key type is not value, it is %d", key.Type))
			}

			value := d.decode()
			dict[string(key.Value.([]byte))] = value
		}

		bi.Type = Dict
		bi.Value = dict

	case 'e':
		bi.Type = EOL
	default:
		panic(fmt.Sprintf("Wrong format symbol (%c) at %d", ch, d.Cursor))
	}

	return bi
}
