package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var quickNotes = []QuickNote{
	QuickNote{"Good", "I had good contact with the customer."},
	QuickNote{"Bad", "I had Bad contact with the customer."},
	QuickNote{"Okay", "I had okay contact with the customer."},
	QuickNote{"Happy", "The customer was happy when I was finished with them."},
	QuickNote{"Sad", "The customer was sad when I was finished with them."},
	QuickNote{"Mad", "The customer was mad when I was finished with them."},
}

func FormToStruct(ptr interface{}, vals url.Values, start string) {
	var strct reflect.Value
	if reflect.TypeOf(ptr) == reflect.TypeOf(reflect.Value{}) {
		strct = ptr.(reflect.Value)
	} else {
		strct = reflect.ValueOf(ptr).Elem()
	}
	strctType := strct.Type()
	for i := 0; i < strct.NumField(); i++ {
		fld := strct.Field(i)
		if ok, v := GetVal(ToLowerFirst(start+strctType.Field(i).Name), vals); ok || fld.Kind() == reflect.Struct {
			switch fld.Kind() {
			case reflect.String:
				strct.Field(i).SetString(v)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				in, _ := strconv.ParseInt(v, 10, 64)
				strct.Field(i).SetInt(in)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				u, _ := strconv.ParseUint(v, 10, 64)
				strct.Field(i).SetUint(u)
			case reflect.Float32, reflect.Float64:
				f, _ := strconv.ParseFloat(v, 64)
				strct.Field(i).SetFloat(f)
			case reflect.Bool:
				b, _ := strconv.ParseBool(v)
				strct.Field(i).SetBool(b)
			case reflect.Map:
				strct.Field(i).Set(reflect.MakeMap(strct.Field(i).Type()))
			case reflect.Slice:
				ss := reflect.MakeSlice(strct.Field(i).Type(), 0, 0)
				strct.Field(i).Set(genSlice(ss, v))
			case reflect.Struct:
				//st := reflect.Indirect(reflect.New(strct.Field(i).Type()))
				st := reflect.Indirect(strct.Field(i))
				FormToStruct(st, vals, start+ToLowerFirst(strctType.Field(i).Name)+".")
				strct.Field(i).Set(st)
			}
		}
	}
}

func genSlice(sl reflect.Value, val string) reflect.Value {
	vs := strings.Split(val, ",")
	for _, v := range vs {
		switch sl.Type().String() {
		case "[]string":
			sl = reflect.Append(sl, reflect.ValueOf(v))
		case "[]int":
			in, _ := strconv.ParseInt(v, 10, 0)
			sl = reflect.Append(sl, reflect.ValueOf(int(in)))
		case "[]int8":
			in, _ := strconv.ParseInt(v, 10, 8)
			sl = reflect.Append(sl, reflect.ValueOf(int8(in)))
		case "[]int16":
			in, _ := strconv.ParseInt(v, 10, 16)
			sl = reflect.Append(sl, reflect.ValueOf(int16(in)))
		case "[]int32":
			in, _ := strconv.ParseInt(v, 10, 32)
			sl = reflect.Append(sl, reflect.ValueOf(int32(in)))
		case "[]int64":
			in, _ := strconv.ParseInt(v, 10, 64)
			sl = reflect.Append(sl, reflect.ValueOf(int64(in)))
		case "[]uint":
			in, _ := strconv.ParseUint(v, 10, 0)
			sl = reflect.Append(sl, reflect.ValueOf(uint(in)))
		case "[]uint8":
			in, _ := strconv.ParseUint(v, 10, 8)
			sl = reflect.Append(sl, reflect.ValueOf(uint8(in)))
		case "[]uint16":
			in, _ := strconv.ParseUint(v, 10, 16)
			sl = reflect.Append(sl, reflect.ValueOf(uint16(in)))
		case "[]uint32":
			in, _ := strconv.ParseUint(v, 10, 32)
			sl = reflect.Append(sl, reflect.ValueOf(uint32(in)))
		case "[]uint64":
			in, _ := strconv.ParseUint(v, 10, 64)
			sl = reflect.Append(sl, reflect.ValueOf(uint64(in)))
		case "[]float32":
			in, _ := strconv.ParseFloat(v, 32)
			sl = reflect.Append(sl, reflect.ValueOf(float32(in)))
		case "[]float64":
			in, _ := strconv.ParseFloat(v, 64)
			sl = reflect.Append(sl, reflect.ValueOf(float64(in)))
		case "[]bool":
			b, _ := strconv.ParseBool(v)
			sl = reflect.Append(sl, reflect.ValueOf(b))
		}
	}
	return sl
}

func GetVal(key string, v url.Values) (bool, string) {
	if v == nil {
		return false, ""
	}
	vs, ok := v[key]
	if !ok || len(vs) == 0 {
		return false, ""
	}
	return true, vs[0]
}

func ToLowerFirst(s string) string {
	return strings.ToLower(string(s[0])) + s[1:len(s)]
}

func PrettySize(size int64) string {
	c := 0
	var sizef float64 = float64(size)
	for sizef > 1024 {
		sizef = sizef / 1024
		c++
	}
	ind := ""
	switch c {
	case 0:
		ind = "B"
	case 1:
		ind = "KB"
	case 2:
		ind = "MB"
	case 3:
		ind = "GB"
	}
	return fmt.Sprintf("%.1f %s", sizef, ind)
}

func ajaxResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, msg)
}

func FormatDate(d string) string {
	ds := strings.Split(d, "-")
	if len(ds) != 3 {
		return ""
	}
	if ds[1][0] == '0' {
		ds[1] = ds[1][1:]
	}
	return fmt.Sprintf("%s/%s/%s", ds[1], ds[2], ds[0])
}

func ToJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
