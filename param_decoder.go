package handy

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type paramDecoder struct {
	handler   Handler
	uriParams map[string]string
}

func newParamDecoder(h Handler, uriParams map[string]string) paramDecoder {
	return paramDecoder{handler: h, uriParams: uriParams}
}

func (c *paramDecoder) Decode(w http.ResponseWriter, r *http.Request) {
	st := reflect.ValueOf(c.handler).Elem()
	c.unmarshalURIParams(st)

	m := strings.ToLower(r.Method)
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("request")
		if value == "all" || strings.Contains(value, m) {
			c.unmarshalURIParams(st.Field(i))
		}
	}
}

func (c *paramDecoder) unmarshalURIParams(st reflect.Value) {
	if st.Kind() == reflect.Ptr {
		return
	}

	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("param")

		if value == "" {
			continue
		}

		param, ok := c.uriParams[value]
		if !ok {
			continue
		}

		s := st.Field(i)
		if s.IsValid() && s.CanSet() {
			switch field.Type.Kind() {
			case reflect.String:
				s.SetString(param)

			case reflect.Bool:
				lower := strings.ToLower(param)
				s.SetBool(lower == "true")

			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				i, err := strconv.ParseInt(param, 10, 64)
				if err != nil {
					if ErrorFunc != nil {
						ErrorFunc(err)
					}
					continue
				}
				s.SetInt(i)

			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				i, err := strconv.ParseUint(param, 10, 64)
				if err != nil {
					if ErrorFunc != nil {
						ErrorFunc(err)
					}
					continue
				}
				s.SetUint(i)
			}
		}
	}
}
