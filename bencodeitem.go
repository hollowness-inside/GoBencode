package gobencode

import (
	"fmt"
	"strconv"
	"strings"
)

type ItemType byte

const (
	EOL     ItemType = 0
	Integer ItemType = 1
	Double  ItemType = 2
	Bytes   ItemType = 4
	List    ItemType = 5
	Dict    ItemType = 6
)

type BencodeItem struct {
	Type  ItemType
	Value any
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

func (v BencodeItem) String() string {
	switch v.Type {
	case Double:
		return fmt.Sprintf("%f", v.Value.(float64))
	case Integer:
		return fmt.Sprintf("%d", v.Value.(int))
	case Bytes:
		return fmt.Sprintf(`"%s"`, string(v.Value.([]byte)))
	case List:
		list := v.Value.([]BencodeItem)
		items := make([]string, len(list))

		for i, v := range list {
			items[i] = v.String()
		}

		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	case Dict:
		mp := v.Value.(map[string]BencodeItem)
		items := make([]string, len(mp))

		i := 0
		for key, value := range mp {
			items[i] = fmt.Sprintf(`"%s": %s`, key, value.String())
			i += 1
		}
		return fmt.Sprintf("{%s}", strings.Join(items, ", "))
	default:
		panic(fmt.Sprintf("Wrong item type %d", v.Type))
	}
}
