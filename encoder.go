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
	case reflect.Struct:
		if _, err := w.Write([]byte{'d'}); err != nil {
			return err
		}

		for i := 0; i < v.NumField(); i++ {
			fmt.Fprintf(w, "%s:", v.Type().Field(i).Name)

			if err := encodeValue(w, v.Field(i)); err != nil {
				return err
			}
		}

		if _, err := w.Write([]byte{'e'}); err != nil {
			return err
		}
	}

	return nil
}
