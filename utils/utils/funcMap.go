package utils

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"reflect"
	"strings"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func MapToStruct2(srcMap interface{}, toStruct interface{}, isMustValidate ...bool) (Mp, error) {
	_mp, _ok := srcMap.(map[string]interface{})
	if !_ok {
		return nil, errors.New("parse map error, data not map")
	}
	if _err := MapToStruct(_mp, toStruct, isMustValidate...); _err != nil {
		return nil, _err
	}
	return _mp, nil
}

func MapToStruct(srcMap interface{}, toStruct interface{}, isMustValidate ...bool) error {
	_mp, _ok := srcMap.(map[string]interface{})
	if !_ok {
		_v, _alias := srcMap.(Mp)
		if !_alias {
			return errors.New("parse map error, data not map")
		}
		_mp = _v
	}
	v := reflect.ValueOf(toStruct)
	if !v.IsNil() {
		t := reflect.TypeOf(toStruct)
		fieldNum := t.Elem().NumField()
		for i := 0; i < fieldNum; i++ {
			_type := t.Elem().Field(i).Type
			_k := t.Elem().Field(i).Tag.Get("json")
			if _k != "" {
				_v, _ok := _mp[_k]
				if !_ok {
					if len(isMustValidate) > 0 && isMustValidate[0] == false {
						//非强制性验证字段属性
						_defaultValue := t.Elem().Field(i).Tag.Get("default")
						if _defaultValue != "" {
							_v = _defaultValue
						} else {
							continue
						}
					} else {
						//强制验证字段必填
						_defaultValue := t.Elem().Field(i).Tag.Get("default")
						if _defaultValue != "" {
							_v = _defaultValue
						} else {
							return errors.New("response error:" + _k)
						}
					}
				}
				//字段所属类型
				switch _type.String() {
				case "int64":
					v.Elem().Field(i).SetInt(ToInt64(_v))
				case "string":
					v.Elem().Field(i).SetString(ToString(_v, true))
				case "bool":
					if _bool, _ok := _v.(bool); _ok {
						v.Elem().Field(i).SetBool(_bool)
					} else {
						if _str, _ok := _v.(string); _ok {
							if _str == "true" || _str == "1" {
								v.Elem().Field(i).SetBool(true)
							}
						} else {
							v.Elem().Field(i).SetBool(false)
						}
					}
				case "float64":
					v.Elem().Field(i).SetFloat(ToFloat64(_v))
				}
			}
		}
	}
	return nil
}

func XmlToMap(xmlStr string) Mp {
	params := make(Mp)
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))
	var (
		key   string
		value string
	)
	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement: // 开始标签
			key = token.Name.Local
		case xml.CharData: // 标签内容
			content := string([]byte(token))
			value = content
		}
		if key != "xml" {
			if value != "\n" {
				params.Set(key, value)
			}
		}
	}
	return params
}

/*
 *将大写字段转小写，例如 appId => app_id
 */
func CamelString(str ...string) string {
	var s string
	if len(str) > 0 {
		s = str[0]
	}
	data := make([]byte, 0, len(s))
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if d >= 'A' && d <= 'Z' {
			d = d + 32
			data = append(data, '_')
		}
		data = append(data, d)
	}
	return string(data[:])
}

// EncodeXMLFromMap encodes map[string]string to io.Writer with xml format.
//
//	NOTE: This function requires the rootname argument and the keys of m (type map[string]string) argument
//	are legitimate xml name string that does not contain the required escape character!
func EncodeXMLFromMap(w io.Writer, m map[string]string, rootname string) (err error) {
	switch v := w.(type) {
	case *bytes.Buffer:
		bufw := v
		if err = bufw.WriteByte('<'); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}

		for k, v := range m {
			if err = bufw.WriteByte('<'); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}

			if err = xml.EscapeText(bufw, []byte(v)); err != nil {
				return
			}

			if _, err = bufw.WriteString("</"); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}
		}

		if _, err = bufw.WriteString("</"); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}
		return nil
	case *strings.Builder:
		bufw := v
		if err = bufw.WriteByte('<'); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}

		for k, v := range m {
			if err = bufw.WriteByte('<'); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}

			if err = xml.EscapeText(bufw, []byte(v)); err != nil {
				return
			}

			if _, err = bufw.WriteString("</"); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}
		}

		if _, err = bufw.WriteString("</"); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}
		return nil

	case *bufio.Writer:
		bufw := v
		if err = bufw.WriteByte('<'); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}

		for k, v := range m {
			if err = bufw.WriteByte('<'); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}

			if err = xml.EscapeText(bufw, []byte(v)); err != nil {
				return
			}

			if _, err = bufw.WriteString("</"); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}
		}

		if _, err = bufw.WriteString("</"); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}
		return bufw.Flush()

	default:
		bufw := bufio.NewWriterSize(w, 256)
		if err = bufw.WriteByte('<'); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}

		for k, v := range m {
			if err = bufw.WriteByte('<'); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}

			if err = xml.EscapeText(bufw, []byte(v)); err != nil {
				return
			}

			if _, err = bufw.WriteString("</"); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}
		}

		if _, err = bufw.WriteString("</"); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}
		return bufw.Flush()
	}
}

// DecodeXMLToMap decodes xml reading from io.Reader and returns the first-level sub-node key-value set,
// if the first-level sub-node contains child nodes, skip it.
func DecodeXMLToMap(r io.Reader) (m map[string]string, err error) {
	m = make(map[string]string)
	var (
		decoder = xml.NewDecoder(r)
		depth   = 0
		token   xml.Token
		key     string
		value   strings.Builder
	)
	for {
		token, err = decoder.Token()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		switch v := token.(type) {
		case xml.StartElement:
			depth++
			switch depth {
			case 2:
				key = v.Name.Local
				value.Reset()
			case 3:
				if err = decoder.Skip(); err != nil {
					return
				}
				depth--
				key = "" // key == "" indicates that the node with depth==2 has children
			}
		case xml.CharData:
			if depth == 2 && key != "" {
				value.Write(v)
			}
		case xml.EndElement:
			if depth == 2 && key != "" {
				m[key] = value.String()
			}
			depth--
		}
	}
}

// 复制结构体的相同属性值
func StructCopy(DstStructPtr interface{}, SrcStructPtr interface{}) {
	srcv := reflect.ValueOf(SrcStructPtr)
	dstv := reflect.ValueOf(DstStructPtr)
	srct := reflect.TypeOf(SrcStructPtr)
	dstt := reflect.TypeOf(DstStructPtr)
	if srct.Kind() != reflect.Ptr || dstt.Kind() != reflect.Ptr ||
		srct.Elem().Kind() == reflect.Ptr || dstt.Elem().Kind() == reflect.Ptr {
		panic("Fatal error:type of parameters must be Ptr of value")
	}
	if srcv.IsNil() || dstv.IsNil() {
		panic("Fatal error:value of parameters should not be nil")
	}
	srcV := srcv.Elem()
	dstV := dstv.Elem()
	srcfields := DeepFields(reflect.ValueOf(SrcStructPtr).Elem().Type())
	for _, v := range srcfields {
		if v.Anonymous {
			continue
		}
		dst := dstV.FieldByName(v.Name)
		src := srcV.FieldByName(v.Name)
		if !dst.IsValid() {
			continue
		}
		if src.Type() == dst.Type() && dst.CanSet() {
			dst.Set(src)
			continue
		}
		if src.Kind() == reflect.Ptr && !src.IsNil() && src.Type().Elem() == dst.Type() {
			dst.Set(src.Elem())
			continue
		}
		if dst.Kind() == reflect.Ptr && dst.Type().Elem() == src.Type() {
			dst.Set(reflect.New(src.Type()))
			dst.Elem().Set(src)
			continue
		}
	}
	return
}

// DeepFields ...
func DeepFields(ifaceType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for i := 0; i < ifaceType.NumField(); i++ {
		v := ifaceType.Field(i)
		if v.Anonymous && v.Type.Kind() == reflect.Struct {
			fields = append(fields, DeepFields(v.Type)...)
		} else {
			fields = append(fields, v)
		}
	}
	return fields
}

// Delete m k v elements indexed by d.
// eg:DeleteMpKey(m, "k0", []int{0, 3}) =>  删除m对象中k0索引的0, 3号元素
func DeleteMpKey(m map[string][]interface{}, k string, d []int) {
	v, ok := m[k]
	if !ok {
		return
	}
	for _, i := range d {
		if 0 <= i && i < len(v) {
			v[i] = nil
		}
	}
	lw := 0
	for i := range v {
		if v[i] != nil {
			lw++
		}
	}
	if lw == 0 {
		delete(m, k)
		return
	}
	w := make([]interface{}, 0, lw)
	for i := range v {
		if v[i] != nil {
			w = append(w, v[i])
		}
	}
	m[k] = w
}
