package gobencode

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type BencodeItem[T any] struct {
	Value    T
	itemType string
}

func (v BencodeItem[T]) String() string {
	return fmt.Sprintf("%s(%v)", v.itemType, v.Value)
}

func Decode(path string) (*BencodeItem[any], error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	return decode(reader)
}

func decode(r *bufio.Reader) (*BencodeItem[any], error) {
	ch, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if '0' <= ch && ch <= '9' {
		rest, err := r.ReadBytes(':')
		if err != nil {
			return nil, err
		}

		rest = rest[:len(rest)-1]
		rawLength := []byte{ch}
		rawLength = append(rawLength, rest...)
		fmt.Println("STRING STRING", string(rawLength))
		length, err := strconv.Atoi(string(rawLength))
		if err != nil {
			return nil, err
		}

		bytes := make([]byte, length)
		_, err = r.Read(bytes)
		if err != nil {
			return nil, err
		}

		return &BencodeItem[any]{
			Value:    string(bytes),
			itemType: "bytes",
		}, nil
	}

	switch ch {
	case 'i':
		rawInt, err := r.ReadBytes('e')
		if err != nil {
			return nil, err
		}

		fmt.Println("STRING INT", string(rawInt[:len(rawInt)-1]))
		intVal, err := strconv.Atoi(string(rawInt[:len(rawInt)-1]))
		if err != nil {
			return nil, err
		}

		return &BencodeItem[any]{
			Value:    intVal,
			itemType: "int",
		}, nil

	case 'l':
		list := []BencodeItem[any]{}

		for {
			item, err := decode(r)
			if err != nil {
				return nil, err
			}

			if item.itemType == "EOL" {
				break
			}

			list = append(list, *item)
		}

		return &BencodeItem[any]{
			Value:    list,
			itemType: "list",
		}, nil

	case 'd':
		dict := map[string]BencodeItem[any]{}

		for {
			key, err := decode(r)
			if err != nil {
				return nil, err
			}

			if key.itemType == "EOL" {
				break
			}

			value, err := decode(r)
			if err != nil {
				return nil, err
			}

			dict[key.Value.(string)] = *value
		}

		return &BencodeItem[any]{
			Value:    dict,
			itemType: "dict",
		}, nil

	case 'e':
		return &BencodeItem[any]{
			Value:    nil,
			itemType: "EOL",
		}, nil
	}

	return nil, nil
}
