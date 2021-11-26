package toold

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

//TimeType TimeType
type TimeType string

//Set Set
func (t *TimeType) Set(val string) {
	*t = TimeType(val)
}

//ConversionToDuration ConversionToDuration
func (t *TimeType) ConversionToDuration() time.Duration {
	dura := fmt.Sprintf("%v", *t)
	if strings.Contains(dura, "s") {
		return time.Duration(ConversionToInt64(strings.Replace(dura, "s", "", -1))) * time.Second
	} else if strings.Contains(dura, "m") {
		return time.Duration(ConversionToInt64(strings.Replace(dura, "m", "", -1))) * time.Minute
	} else if strings.Contains(dura, "h") {
		return time.Duration(ConversionToInt64(strings.Replace(dura, "h", "", -1))) * time.Hour
	} else {
		return time.Duration(ConversionToInt64(dura)) * time.Second
	}
}

//ConversionToSecond ConversionToSecond
func (t *TimeType) ConversionToSecond() int64 {
	return int64(t.ConversionToDuration()) / 1000000000
}

//ConversionToString ConversionToString
func (t *TimeType) ConversionToString() string {
	dura := fmt.Sprintf("%v", *t)
	if strings.Contains(dura, "s") {
		return strings.Replace(dura, "s", "秒", -1)
	} else if strings.Contains(dura, "m") {
		return strings.Replace(dura, "m", "分钟", -1)
	} else if strings.Contains(dura, "h") {
		return strings.Replace(dura, "h", "小时", -1)
	} else {
		return fmt.Sprintf("%v秒", dura)
	}
}

//时间戳转时间格式
func ConversionChinaToTime(dur int64) string {
	if dur < 60 {
		return fmt.Sprintf("%v秒", dur)
	} else if dur < 3600 && dur >= 60 {
		return fmt.Sprintf("%v分%v秒", dur/60, dur%60)
	} else {
		return fmt.Sprintf("%v时%v分%v秒", dur/3600, (dur%3600)/60, (dur%3600)%60)
	}
}

/*
ConversionToString 转换字符串
*/
func ConversionToString(num interface{}) string { //reflect.TypeOf(num)
	var number string
	switch num.(type) {
	case int:
		number = strconv.Itoa(num.(int))
		break
	case float64:
		number = strconv.FormatFloat(num.(float64), 'f', 0, 64)
		break
	case float32:
		number = strconv.FormatFloat(num.(float64), 'f', 0, 32)
		break
	case string:
		number = num.(string)
	case int64:
		number = strconv.FormatInt(num.(int64), 10)
		break
	case int32:
		number = strconv.FormatInt(int64(num.(int32)), 10)
		break
	case int16:
		number = strconv.FormatInt(int64(num.(int16)), 10)
	case map[string]string:
		bt, errors := json.Marshal(num)
		if errors != nil {
			number = errors.Error()
		} else {
			number = string(bt)
		}
		break
	case map[string]interface{}:
		bt, errors := json.Marshal(num)
		if errors != nil {
			number = errors.Error()
		} else {
			number = string(bt)
		}
		break
	case interface{}:
		bt, errors := json.Marshal(num)
		if errors != nil {
			number = errors.Error()
		} else {
			number = string(bt)
		}
		break
	default:
		bt, errors := json.Marshal(num)
		if errors != nil {
			number = errors.Error()
		} else {
			number = string(bt)
		}
		break
	}
	return number
}

/*
ConversionToInt 转换int
*/
func ConversionToInt(num interface{}) int {
	var number int
	switch num.(type) {
	case int:
		number = num.(int)
		break
	case int64:
		number = int(num.(int64))
		break
	case float64:
		number = int(num.(float64))
		break
	case float32:
		number = int(num.(float32))
		break
	case string:
		number, _ = strconv.Atoi(num.(string))
		break
	}
	return number
}

/*
ConversionToInt64 转64
*/
func ConversionToInt64(num interface{}) int64 {
	var number int64
	switch num.(type) {
	case int:
		number = int64(num.(int))
		break
	case float64:
		number = int64(num.(float64))
		break
	case float32:
		number = int64(num.(float32))
		break
	case string:
		number, _ = strconv.ParseInt(num.(string), 10, 64)
		break
	case interface{}:
		number, _ = strconv.ParseInt(fmt.Sprintf("%v", num), 10, 64)
		break
	}
	return number
}

/*
ConversionToInt16 转64
*/
func ConversionToInt16(num interface{}) int16 {
	var number int16
	switch num.(type) {
	case int:
		number = int16(num.(int))
		break
	case float64:
		number = int16(num.(float64))
		break
	case float32:
		number = int16(num.(float32))
		break
	case string:
		number16, _ := strconv.ParseInt(num.(string), 10, 32)
		number = int16(number16)
		break
	}
	return number
}

func ConversionTOFloat64(num string) float64 {
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0.0
	}
	return f
}

/*
ConversionToMap 转map
*/
func ConversionToMap(num string) (map[string]interface{}, error) {
	var con map[string]interface{}
	byt := []byte(num)
	errs := json.Unmarshal(byt, &con)
	if errs != nil {
		return nil, errs
	}
	return con, nil
}

/*
ConversionToMaps 转Maps
*/
func ConversionToMaps(num string) ([]map[string]interface{}, error) {
	var con []map[string]interface{}
	byt := []byte(num)
	errs := json.Unmarshal(byt, &con)
	if errs != nil {
		return nil, errs
	}
	return con, nil
}

/*
ConversionToObj 返回自定义对象
*/
func ConversionToObj(num string, v interface{}) error {
	byt := []byte(num)
	errs := json.Unmarshal(byt, &v)
	if errs != nil {
		return errs
	}
	return nil
}

/*
ConversionByteToObj 返回自定义对象
*/
func ConversionByteToObj(num []byte, v interface{}) error {
	errs := json.Unmarshal(num, &v)
	if errs != nil {
		return errs
	}

	return nil
}

/*
ConversionTopByte 转字节
*/
func ConversionTopByte(vs interface{}) ([]byte, error) {
	return json.Marshal(vs)
}

/*
ConversionIPFromItoa 将整数转为ip
*/
func ConversionIPFromItoa(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

/*
ConversionBufferFromByte Buffer
*/
func ConversionBufferFromByte(body []byte) *bytes.Buffer {
	return bytes.NewBuffer(body)
}

//ConversionReaderFromByte ConversionReaderFromByte
func ConversionReaderFromByte(body []byte) *bytes.Reader {
	return bytes.NewReader(body)
}

//ConversionBodyFromReader ConversionBodyFromReader
func ConversionBodyFromReader(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}

//ConversionBodyFromWrite ConversionBodyFromWrite
func ConversionBodyFromWrite(out io.Writer) ([]byte, error) {
	var in io.Reader
	_, err = io.Copy(out, in)
	if err != nil {
		return []byte{}, err
	}
	body, err := ConversionBodyFromReader(in)
	return body, err
}

/*
ConversionBase64 转bsase54
*/
func ConversionBase64(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

/*
ConversionBodyFromBase64 base64转解吗
*/
func ConversionBodyFromBase64(base64s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64s)
}

/*
ConversionEncoding URLEncoding
*/
func ConversionEncoding(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

/*
ConversionBodyFromEncoding URLEncoding
*/
func ConversionBodyFromEncoding(str string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(str)
}

//ConversionBodyFromInt ConversionBodyFromInt
func ConversionBodyFromInt(num int) []byte {
	x := uint32(num)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//ConversionBodyFromData 任意数据转[]byte
func ConversionBodyFromData(data interface{}) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, data)
	return bytesBuffer.Bytes()
}

//ConversionBodyLittleFromData 任意数据转[]byte
func ConversionBodyLittleFromData(data interface{}) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, data)
	return bytesBuffer.Bytes()
}

//ConversionDataLittleFromBody ConversionDataLittleFromBody
func ConversionDataLittleFromBody(by []byte, data interface{}) {
	bytesBuffer := bytes.NewBuffer(by)
	binary.Read(bytesBuffer, binary.LittleEndian, data)
}

// value := reflect.ValueOf(*userInfo)
// 	for i := 0; i < value.NumField(); i++ {
// 		fmt.Printf("Field %d: %+v\n", i, value.Field(i))
// 	}

//ConvertToString ConvertToString
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

//ConvertCharacterDecoderToString 将字符解码未对应字符 UTF-8 gbk
func ConvertCharacterDecoderToString(charater string, con string) string {
	dec := mahonia.NewDecoder(charater)
	return dec.ConvertString(con)
}

//ConvertCharacterEncoderToString 将字符编码未对应字符
func ConvertCharacterEncoderToString(charater string, con string) string {
	dec := mahonia.NewEncoder(charater)
	return dec.ConvertString(con)
}

func ConversionGbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func ConversionUtf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func ConversionUTF16LEToUTF8(body []byte) ([]byte, error) {
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	bs2, err := decoder.Bytes(body)
	if err != nil {
		return []byte{}, err
	}
	return bs2, nil
}

func ConversionUTF8ToUTF16LE(body []byte) ([]byte, error) {
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM).NewEncoder()
	bs2, err := decoder.Bytes(body)
	if err != nil {
		return []byte{}, err
	}
	return bs2, nil
}
