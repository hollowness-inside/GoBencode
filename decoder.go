package gobencode

import (
	"bufio"
	"fmt"
	"strconv"
)

type BencodeItem struct {
	Value any
	Type  string
}

func (v *BencodeItem) String() string {
	return fmt.Sprintf("%v", v.Value)
}

type Decoder struct {
	r bufio.Reader
}

func NewDecoder(reader *bufio.Reader) *Decoder {
	return &Decoder{
		r: *reader,
	}
}

func (d *Decoder) Decode() (*BencodeItem, error) {
	r := d.r

	ch, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if '0' <= ch && ch <= '9' {
		lengthBytes, err := r.ReadBytes(':')
		if err != nil {
			return nil, err
		}

		lengthBytes = append([]byte{ch}, lengthBytes[:len(lengthBytes)-1]...)
		length, err := strconv.Atoi(string(lengthBytes))
		if err != nil {
			return nil, err
		}

		bytes := make([]byte, length)
		_, err = r.Read(bytes)
		if err != nil {
			return nil, err
		}

		return &BencodeItem{
			Value: string(bytes),
			Type:  "string",
		}, nil
	}

	switch ch {
	case 'i':
		rawInt, err := r.ReadBytes('e')
		if err != nil {
			return nil, err
		}

		intVal, err := strconv.Atoi(string(rawInt[:len(rawInt)-1]))
		if err != nil {
			return nil, err
		}

		return &BencodeItem{
			Value: intVal,
			Type:  "int",
		}, nil

	case 'l':
		list := []BencodeItem{}

		for {
			item, err := d.Decode()
			if err != nil {
				return nil, err
			}

			if item.Type == "EOL" {
				break
			}

			list = append(list, *item)
		}

		return &BencodeItem{
			Value: list,
			Type:  "list",
		}, nil

	case 'd':
		dict := map[string]BencodeItem{}

		for {
			key, err := d.Decode()
			if err != nil {
				return nil, err
			}

			if key.Type == "EOL" {
				break
			}

			value, err := d.Decode()
			if err != nil {
				return nil, err
			}

			dict[key.Value.(string)] = *value
		}

		return &BencodeItem{
			Value: dict,
			Type:  "dict",
		}, nil

	case 'e':
		return &BencodeItem{
			Value: nil,
			Type:  "EOL",
		}, nil
	}

	return nil, nil
}
