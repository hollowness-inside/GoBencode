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
	String  ItemType = 4
	List    ItemType = 5
	Map     ItemType = 6
)

type BencodeItem struct {
	Type  ItemType
	Value any
}

func (d *Decoder) decode() BencodeItem {
	ch := d.ReadByte()

	if '0' <= ch && ch <= '9' {
		lengthBytes := d.ReadBytes(':')
		lengthBytes = append([]byte{ch}, lengthBytes[:len(lengthBytes)-1]...)
		length, err := strconv.Atoi(string(lengthBytes))
		if err != nil {
			panic(err)
		}

		bytes := make([]byte, length)
		d.Read(bytes)

		return BencodeItem{
			Value: bytes,
			Type:  String,
		}
	}

	switch ch {
	case 'i':
		rawInt := d.ReadBytes('e')
		intVal, err := strconv.Atoi(string(rawInt[:len(rawInt)-1]))
		if err != nil {
			panic(err)
		}

		return BencodeItem{
			Value: intVal,
			Type:  Integer,
		}

	case 'l':
		list := []BencodeItem{}

		for {
			item := d.decode()
			if item.Type == EOL {
				break
			}

			list = append(list, item)
		}

		return BencodeItem{
			Value: list,
			Type:  List,
		}

	case 'd':
		dict := map[string]BencodeItem{}

		for {
			key := d.decode()
			if key.Type == EOL {
				break
			} else if key.Type != String {
				panic(fmt.Sprintf("dict key type is not value, it is %d", key.Type))
			}

			value := d.decode()
			dict[string(key.Value.([]byte))] = value
		}

		return BencodeItem{
			Value: dict,
			Type:  Map,
		}

	case 'e':
		return BencodeItem{
			Value: nil,
			Type:  EOL,
		}
	default:
		panic(fmt.Sprintf("Wrong format symbol (%c) at %d", ch, d.Cursor))
	}
}

func (v BencodeItem) String() string {
	switch v.Type {
	case Double:
		return fmt.Sprintf("%f", v.Value.(float64))
	case Integer:
		return fmt.Sprintf("%d", v.Value.(int))
	case String:
		return fmt.Sprintf(`"%s"`, string(v.Value.([]byte)))
	case List:
		list := v.Value.([]BencodeItem)
		items := make([]string, len(list))

		for i, v := range list {
			items[i] = v.String()
		}

		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	case Map:
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
