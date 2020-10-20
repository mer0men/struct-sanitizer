package sanitizer

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

const RuPhoneNumberRegExp = `^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$`

var phoneTagHandler = func(tagValue string, fieldValue interface{}) error {
	var numberRegExp *regexp.Regexp

	if tagValue == "ru" {
		var err error
		numberRegExp, err = regexp.Compile(RuPhoneNumberRegExp)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Unknown country name\n")
	}

	numberString, ok := (fieldValue).(string)
	if !ok {
		return errors.New("Could not validate number\n")
	}

	valid := numberRegExp.MatchString(numberString)

	if !valid {
		return errors.New("Could not validate number\n")
	}

	return nil
}

type testStruct struct {
	Foo int  `json:"foo,string,omitempty"`
	Bar string `json:"bar"`
	Baz string `json:"baz" phone:"ru"`
}

func TestSanitizer(t *testing.T) {
	//var testData = []string{
	//	`{"foo": "123", "bar": "test", "baz": "+79261234567"}`,
	//	`{"foo": "123", "bar": "test", "baz": "89261234567"}`,
	//	`{"foo": "123", "bar": "test", "baz": "8(926)123-45-67"}`,
	//	`{"foo": "1234567890", "bar": "test2", "baz": "8 (950) 288-56-23"}`,
	//	`{"foo": "123.0", "bar": "test3", "baz": "8 (950) 288-56-23"}`,
	//	`{"foo": "", "bar": "test4", "baz": "8 (950) 288-56-23"}`,
	//	`{"foo": "666", "bar": "", "baz": "8 (950) 288-56-23"}`,
	//	`{"foo": "12312412126787657898765498765456789", "bar": "test5", "baz": "8 (950) 288-56-23"}`,
	//	`{"foo": "test string", "bar": "asdvfhjkl", "baz": "8 (950) 288-56-23"}`,
	//	`{"foo": "12547890358y", "bar": "F", "baz": "8 (950) 288-56-23"}`,
	//	`{"foo": "y2y54y6y1", "bar": "E", "baz": "8 (950) 288-56-23"}`,
	//	`{"foo": ",mxcvlkaso1j234oj124, "bar": "I", "baz": "8 (950) 288-56-23"}`,
	//	`{"foo": "2...2", "bar": "P", "baz": "8 (950) 288-56-23"}`}

	s := NewSanitizer()
	err := s.AddTag("phone", phoneTagHandler)
	if err != nil {
		t.Fatal(err)
	}

	var testObject testStruct
	var data []byte

	data = []byte(`{"foo": "123", "bar": "test", "baz": "+79261234567"}`)
	err = s.UnmarshalAndSanitize(data, &testObject)
	if !assert.Equal(t, nil, err) {
		t.Error()
	}

	data = []byte(`{"foo": "123", "bar": "test", "baz": "89261234567"}`)
	err = s.UnmarshalAndSanitize(data, &testObject)
	if !assert.Equal(t, nil, err)  {
		t.Error()
	}

	data = []byte(`{"foo": "123", "bar": "test", "baz": "8(926)123-45-67"}`)
	err = s.UnmarshalAndSanitize(data, &testObject)
	if !assert.Equal(t, nil, err)  {
		t.Error()
	}

	data = []byte(`{"foo": "1234567890", "bar": "test2", "baz": "8 (950) 288-56-23"}`)
	err = s.UnmarshalAndSanitize(data, &testObject)
	if !assert.Equal(t, nil, err)  {
		t.Error()
	}

	data = []byte(`{"foo": "666", "bar": "", "baz": "8 (950) 288-56-23"}`)
	err = s.UnmarshalAndSanitize(data, &testObject)
	if !assert.Equal(t, nil, err)  {
		t.Error()
	}

	data = []byte(`{"foo": "", "bar": "test4", "baz": "8 (950) 288-56-23"}`)
	err = s.UnmarshalAndSanitize(data, &testObject)
	if !assert.Error(t, err)  {
		t.Error()
	}

	data = []byte(`{"foo": "12312412126787657898765498765456789", "bar": "test5", "baz": "8 (950) 288-56-23"}`)
	err = s.UnmarshalAndSanitize(data, &testObject)
	if !assert.Error(t, err)  {
		t.Error()
	}

	data = []byte(`{"foo": "123", "bar": "test7", "baz": "8 (950) "}`)
	err = s.UnmarshalAndSanitize(data, &testObject)
	if !assert.Error(t, err)  {
		t.Error()
	}

	data = []byte(`{"foo": "123", "bar": "test8", "baz": "asdjhsdf"}`)
	err = s.UnmarshalAndSanitize(data, &testObject)
	if !assert.Error(t, err)  {
		t.Error()
	}

	data = []byte(`{"foo": "123", "bar": "test8", "baz": ""}`)
	err = s.UnmarshalAndSanitize(data, &testObject)
	if !assert.Error(t, err)  {
		t.Error()
	}
}
