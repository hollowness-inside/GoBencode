package gobencode

import (
	"fmt"
	"io"
	"reflect"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (e *Encoder) Encode(v any) error {
	return encodeValue(e.w, reflect.ValueOf(v))
}

func encodeValue(w io.Writer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if _, err := fmt.Fprintf(w, "i%de", v.Int()); err != nil {
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if _, err := fmt.Fprintf(w, "i%de", v.Uint()); err != nil {
			return err
		}
	case reflect.String:
		str := v.String()
		if _, err := fmt.Fprintf(w, "%d:%s", len(str), str); err != nil {
			return err
		}
	case reflect.Slice, reflect.Array:
		if _, err := w.Write([]byte{'l'}); err != nil {
			return err
		}

		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if err := encodeValue(w, elem); err != nil {
				return err
			}
		}

		if _, err := w.Write([]byte{'e'}); err != nil {
			return err
		}
	case reflect.Map:
		w.Write([]byte{'d'})
		mr := v.MapRange()
		for mr.Next() {
			encodeValue(w, mr.Key())
			encodeValue(w, mr.Value())
		}
		w.Write([]byte{'e'})

	case reflect.Struct:
		if v.Type() == reflect.TypeOf(BencodeItem{}) {
			tp := v.FieldByName("Type").Uint()
			vl := v.FieldByName("Value")

			switch ItemType(tp) {
			case Integer:
				fmt.Fprintf(w, "i%de", vl.Interface().(int))
			case Bytes:
				data := vl.Interface().([]byte)
				fmt.Fprintf(w, "%d:", len(data))
				w.Write(data)
			case List:
				w.Write([]byte{'l'})
				for _, v := range vl.Interface().([]BencodeItem) {
					encodeValue(w, reflect.ValueOf(v))
				}
				w.Write([]byte{'e'})
			case Dict:
				w.Write([]byte{'d'})
				for k, v := range vl.Interface().(map[string]BencodeItem) {
					encodeValue(w, reflect.ValueOf(k))
					encodeValue(w, reflect.ValueOf(v))
				}
				w.Write([]byte{'e'})
			}

			return nil
		}
	}

	return nil
}
