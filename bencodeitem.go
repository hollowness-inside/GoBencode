package gobencode

import (
	"fmt"
	"strings"
)

type ItemType byte

const (
	EOL     ItemType = 0
	Integer ItemType = 1
	Bytes   ItemType = 2
	List    ItemType = 3
	Dict    ItemType = 4
)

type BencodeItem struct {
	Type  ItemType
	Value any
}

func (v BencodeItem) String() string {
	switch v.Type {
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
