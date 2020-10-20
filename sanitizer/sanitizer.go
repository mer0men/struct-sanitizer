package sanitizer

import (
	"encoding/json"
	"reflect"
)

type Sanitizer struct {
	tagsHandlers map[string]TagHandler
	tagsNames []string
}

type TagHandler func(string, interface{}) error

func (s *Sanitizer) UnmarshalAndSanitize(data []byte, v interface{}) error {
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	err = s.Sanitize(v)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sanitizer) Sanitize(v interface{}) error {

	structVal := reflect.ValueOf(v)
	err := s.sanitizeValues(structVal)
	if err != nil {
		return err
	}

	return nil
}

func (s Sanitizer) sanitizeValues(v reflect.Value) error {
	var structType reflect.Type
	var workVal reflect.Value
	valKind := v.Kind().String()

	if valKind != "ptr" {
		structType = v.Type()
		workVal = v
	} else {
		structType = v.Elem().Type()
		workVal = v.Elem()
	}

	for i := 0; i < structType.NumField(); i++ {
		valueField := workVal.Field(i)
		kind := valueField.Kind().String()
		if kind != "struct" && kind != "ptr" {
			structField := structType.Field(i)
			for _, tagName := range s.tagsNames {
				value, ok := structField.Tag.Lookup(tagName)
				if ok {
					err := s.HandleTag(tagName, value, valueField.Interface())
					if err != nil {
						return err
					}
				}
			}
		} else {
			i := v.Elem().Field(i).Interface()
			err := s.sanitizeValues(reflect.ValueOf(i))
			if err != nil {
				return err
			}
		}
	}
	return nil
}


func (s *Sanitizer) AddTag(tagName string, f TagHandler)  error {
	s.tagsHandlers[tagName] = f
	s.tagsNames = append(s.tagsNames, tagName)
	return nil
}

func (s *Sanitizer) HandleTag(tagName, tagValue string, fieldValue interface{}) error {
	f := s.tagsHandlers[tagName]
	err := f(tagValue, fieldValue)
	if err != nil {
		return err
	}

	return nil
}

func NewSanitizer() Sanitizer {
	return Sanitizer{
		tagsHandlers: make(map[string]TagHandler),
		tagsNames: make([]string, 0),
	}
}

