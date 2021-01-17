package mycrypto

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/tidusant/c3m/common/log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/tidusant/c3m/common/mystring"
)

func StringRand(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
func NumRand(min, max int) int {
	if max <= min {
		max = min + 1
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
func Base64fix(s string) string {
	n := 4 - len(s)%4
	if n < 4 {
		for i := 0; i < n; i++ {
			s += "="
		}
	}
	return s
}
func Base64Compress(s string) string {
	return CompressToBase64(s)
}
func MD5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func Base64Decompress(s string) string {
	byteDecode, _ := DecompressFromBase64(s)
	return byteDecode
}
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
func Base64Decode(s string) string {
	byteDecode, _ := base64.StdEncoding.DecodeString(Base64fix(s))
	return string(byteDecode)
}
func CampaignDecode(data string) string {
	code := data[:len(data)/2]
	code = strings.Replace(data, code, "", 1) + code
	strbytes, _ := base64.StdEncoding.DecodeString(Base64fix(code))
	return string(strbytes)
}

//Encode: encoding data, if div=2: compress data to base64 before decode, else base64 data.
// - Add a random string with x len (x random from 2 to 9) to data
// - split data by x into 2 string and reverse 2 string and join into 1 string
// - base64 x into x2
// - put x2 into datastring at position of l=len(data)/div
func Encode(data string, div int) string {
	if data == "" {
		return data
	}
	var x = NumRand(2, 9)
	//log.Debugf("random x :%s", x)
	var x2 = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(x)))
	x2 = strings.Replace(x2, "=", "", -1)
	//log.Debugf("x2 base 64 :%s", x2)
	//x2b := []byte(x2)
	if div == 2 {
		data = CompressToBase64(data)
	} else {
		data = base64.StdEncoding.EncodeToString([]byte(data))
		//log.Debugf("data base 64 :%s", data)
	}

	data = strings.Replace(data, "=", "", -1)

	xstr := mystring.RandString(x)
	//log.Debugf("xstr RandString(x):%s", xstr)

	data += xstr

	b := []byte(data)

	var l = int(math.Floor(float64(len(data) / div)))
	//log.Debugf("l :%s", l)
	var result1 []byte
	var result2 []byte

	for i := len(b) - 1; i >= 0; i-- {

		if i%x == 0 {
			result1 = append(result1, b[i]) // string([]rune(data)[i])

		} else {
			result2 = append(result2, b[i])
		}
	}
	//log.Debugf("string(result1):%s", string(result1))
	//log.Debugf("string(result2):%s", string(result2))
	strb64 := string(result1) + string(result2)
	strb64 = strb64[:l] + x2 + strb64[l:]
	//log.Debugf("strb64:%s", strb64)
	return strb64
}

func EncodeLight1(data string, div int) string {
	if data == "" {
		return data
	}
	var x = NumRand(2, 9)
	//log.Debugf("random x :%s", x)
	var x2 = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(x)))
	x2 = strings.Replace(x2, "=", "", -1)
	//log.Debugf("x2 base 64 :%s", x2)
	//x2b := []byte(x2)
	if div == 2 {
		data = CompressToBase64(data)

	} else {
		data = base64.StdEncoding.EncodeToString([]byte(data))
		//log.Debugf("data base 64 :%s", data)
	}

	data = strings.Replace(data, "=", "", -1)

	xstr := mystring.RandString(x)
	//log.Debugf("xstr RandString(x):%s", xstr)

	data += xstr

	var l = int(math.Floor(float64(len(data) / div)))

	strb64 := data[:l] + x2 + data[l:]
	//log.Debugf("strb64:%s", strb64)

	return strb64
}

func DecodeLight1(code string, div int) string {

	if code == "" {
		return code
	}
	var rt string = ""

	l := int(math.Floor(float64((len(code) - 2) / div)))
	//log.Debugf("l :%s", l)
	//var l = int(math.Floor(float64(len(data) / div)))

	//log.Debugf("x2 :%s", x2)
	key := code[:l] + code[l+2:]

	x2 := code[l : l+2]
	//log.Debugf("rs1 + rs 2 :%s", key)
	byteDecode, _ := base64.StdEncoding.DecodeString(Base64fix(x2))
	x2 = string(byteDecode)

	floatNum, _ := strconv.ParseFloat(x2, 64)
	intNum := (int)(floatNum)
	key = key[:len(key)-intNum]

	// //log.Debugf("x :%s", intNum)
	// if intNum > 0 {
	// 	//print_r($num);print_r("\r\n");
	// 	//get odd string
	// 	lf := math.Ceil((float64)(len(key)) / floatNum)
	// 	//log.Debugf("lf :%s", lf)
	// 	oddstr := key[:int(lf)]
	// 	//log.Debugf("oddstr key[:int(lf)] :%s", oddstr)
	// 	ukey := strings.Replace(key, oddstr, "", 1)
	// 	//log.Debugf("ukey Replace(key, oddstr :%s", ukey)
	// 	base64str := ""

	// 	for i := len(oddstr) - 1; i >= 0; i-- {
	// 		base64str += string(oddstr[len(oddstr)-1:])
	// 		oddstr = oddstr[:len(oddstr)-1]
	// 		if len(ukey)-intNum+1 > 0 {
	// 			base64str += mystring.Reverse(string(ukey[len(ukey)-intNum+1:]))
	// 		} else {
	// 			base64str += mystring.Reverse(ukey)
	// 		}
	// 		if i > 0 {
	// 			ukey = ukey[:len(ukey)-intNum+1]
	// 		}
	// 	}
	// 	//log.Debugf("base64str :%s", base64str)
	// 	base64str = base64str[:len(base64str)-intNum]

	if div == 2 {
		byteDecode, _ := DecompressFromBase64(key)
		rt = byteDecode
	} else {
		byteDecode, _ := base64.StdEncoding.DecodeString(Base64fix(key))
		rt = string(byteDecode)
	}

	//}
	return rt
}

//for encode
func DecodeOld(code string, div int) string {
	if code == "" {
		return code
	}
	var rt string = ""
	key := code
	//key = "kZXUuYkRWzUgQk92YoNwRdh92Q3SZtFmb9Wa0NW"
	if key == rt {
		return rt
	}

	oddstr := "d"
	l := int(math.Floor(float64((len(key) - 2) / div)))
	//log.Debugf("l :%s", l)
	//var l = int(math.Floor(float64(len(data) / div)))
	x2 := key[l : l+2]
	//log.Debugf("x2 :%s", x2)
	key = key[:l] + key[l+2:]
	//log.Debugf("rs1 + rs 2 :%s", key)
	byteDecode, _ := base64.StdEncoding.DecodeString(Base64fix(x2))
	x2 = string(byteDecode)

	floatNum, _ := strconv.ParseFloat(x2, 64)
	intNum := (int)(floatNum)
	//log.Debugf("x :%s", intNum)
	if intNum > 0 {
		//print_r($num);print_r("\r\n");
		//get odd string
		lf := math.Ceil((float64)(len(key)) / floatNum)
		//log.Debugf("lf :%s", lf)
		oddstr = key[:int(lf)]
		//log.Debugf("oddstr key[:int(lf)] :%s", oddstr)
		ukey := strings.Replace(key, oddstr, "", 1)
		//log.Debugf("ukey Replace(key, oddstr :%s", ukey)
		base64str := ""

		for i := len(oddstr) - 1; i >= 0; i-- {
			base64str += string(oddstr[len(oddstr)-1:])
			oddstr = oddstr[:len(oddstr)-1]
			checknum := len(ukey) - intNum + 1
			if checknum > 0 {
				base64str += mystring.Reverse(string(ukey[checknum]))
			} else {
				base64str += mystring.Reverse(ukey)
			}
			if i > 0 {
				ukey = ukey[:checknum]
			}
		}
		//log.Debugf("base64str :%s", base64str)
		base64str = base64str[:len(base64str)-intNum]
		//log.Debugf("base64str :%s", base64str)
		//log.Debugf("lzjs %s", base64str)
		//log.Debugf("lzjs %s", Base64fix(base64str))
		//byteDecode, _ := base64.StdEncoding.DecodeString(Base64fix(base64str))
		//base64str = Base64fix(base64str)
		//byteDecode, _ := DecompressFromBase64(base64str)
		//data, _ := base64.StdEncoding.DecodeString(base64str)
		if div == 2 {
			byteDecode, _ := DecompressFromBase64(base64str)
			rt = byteDecode
		} else {
			byteDecode, _ := base64.StdEncoding.DecodeString(Base64fix(base64str))
			rt = string(byteDecode)
		}

		//log.Debugf("rt :%s", rt)
		//log.Debugf("data decompress %s", rt)
		//rt, _ = url.QueryUnescape(rt)
		//log.Debugf("data decompress %s", rt)
	}
	return rt
}

//for decode old
func EncodeApp(data string, div int) string {
	if data == "" {
		return data
	}
	var x = NumRand(10, 99)
	//log.Debugf("random x :%s", x)

	data = strings.Replace(data, "=", "", -1)

	xstr := mystring.RandString(x)
	//log.Debugf("xstr RandString(x):%s", xstr)

	data += xstr

	b := []byte(data)

	var l = int(math.Floor(float64(len(data) / div)))
	//log.Debugf("l :%s", l)
	var result1 []byte
	var result2 []byte

	for i := len(b) - 1; i >= 0; i-- {

		if i%div == 0 {
			result1 = append(result1, b[i]) // string([]rune(data)[i])

		} else {
			result2 = append(result2, b[i])
		}
	}
	log.Debugf("string(result1):%s", string(result1))
	log.Debugf("string(result2):%s", string(result2))
	strb64 := string(result1) + string(result2)
	strb64 = strb64[:l] + strconv.Itoa(x) + strb64[l:]
	log.Debugf("strb64:%s", strb64)
	return strb64
}

func DecodeApp(data string, div int) string {
	if data == "" {
		return data
	}
	var rt string = ""
	key := data
	//key = "kZXUuYkRWzUgQk92YoNwRdh92Q3SZtFmb9Wa0NW"
	if data == rt {
		return rt
	}

	oddstr := "d"
	l := int(math.Floor(float64((len(data) - 2) / div)))
	//log.Debugf("l :%s", l)
	//var l = int(math.Floor(float64(len(data) / div)))
	x, _ := strconv.Atoi(data[l : l+2])
	//log.Debugf("x2 :%s", x2)
	key = data[:l] + data[l+2:]

	floatNum, _ := strconv.ParseFloat(strconv.Itoa(div), 64)
	intNum := (int)(floatNum)
	//log.Debugf("x :%s", intNum)
	if intNum > 0 {
		//print_r($num);print_r("\r\n");
		//get odd string
		lf := math.Ceil((float64)(len(key)) / floatNum)
		log.Debugf("lf :%s", lf)
		oddstr = key[:int(lf)]
		log.Debugf("oddstr key[:int(lf)] :%s", oddstr)
		ukey := strings.Replace(key, oddstr, "", 1)
		log.Debugf("ukey Replace(key, oddstr :%s", ukey)
		base64str := ""

		for i := len(oddstr) - 1; i >= 0; i-- {
			base64str += string(oddstr[len(oddstr)-1:])
			oddstr = oddstr[:len(oddstr)-1]
			if len(ukey)-intNum+1 > 0 {
				base64str += mystring.Reverse(string(ukey[len(ukey)-intNum+1:]))
			} else {
				base64str += mystring.Reverse(ukey)
			}
			if i > 0 {
				ukey = ukey[:len(ukey)-intNum+1]
			}
		}
		//log.Debugf("base64str :%s", base64str)
		rt = base64str[:len(base64str)-x]

	}
	return rt
}

//not have random url
func EncodeBK(data string, keysalt string) string {
	if data == "" {
		return data
	}
	//data = CompressToBase64(data)
	data = base64.StdEncoding.EncodeToString([]byte(data))
	data = strings.Replace(data, "=", "", -1)
	//log.Debugf("keysalt: %s", keysalt)
	keysalt = base64.StdEncoding.EncodeToString([]byte(keysalt))
	keysalt = strings.Replace(keysalt, "=", "", -1)
	x := mystring.RandString(len(keysalt))
	//log.Debugf("keysalt: %s", keysalt)
	l := len(x)
	if len(x) > len(data) {
		l = len(data)
	}
	data = data[:l] + x + data[l:]
	//log.Println("strReturn: %s", data)
	return data
}

func DecodeBK(data string, keysalt string) string {
	if data == "" {
		return data
	}
	keysalt = base64.StdEncoding.EncodeToString([]byte(keysalt))
	keysalt = strings.Replace(keysalt, "=", "", -1)
	l := len(keysalt)
	if l*2 > len(data) {
		l = len(data) - l
	}

	data = data[:l] + data[l+len(keysalt):]

	data = Base64fix(data)
	//byteDecode, _ := DecompressFromBase64(data)
	byteDecode, _ := base64.StdEncoding.DecodeString(data)
	data = string(byteDecode)
	return data
}

//for encDatA
func EncodeA(data string) string {
	if data == "" {
		return data
	}
	var x = NumRand(2, 9)

	//x2b := []byte(x2)

	xstr := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(x)))
	xstr = strings.Replace(xstr, "=", "", -1)

	data = base64.StdEncoding.EncodeToString([]byte(data))
	data = strings.Replace(data, "=", "", -1)

	b := []byte(data)

	var result1 []byte
	var result2 []byte

	for i := len(b) - 1; i >= 0; i-- {

		if i%x == 0 {
			result1 = append(result1, b[i]) // string([]rune(data)[i])

		} else {
			result2 = append(result2, b[i])
		}
	}

	strb64 := string(result1) + string(result2) + xstr
	return strb64
}

//for encDatA
func DecodeA(data string) string {
	if data == "" {
		return data
	}
	oddb64, _ := base64.StdEncoding.DecodeString(Base64fix(data[len(data)-2:]))
	odd, _ := strconv.Atoi(string(oddb64))

	data = data[:len(data)-2]

	rs1len := int(math.Floor(float64(len(data)) / 2))
	rs1 := data[:rs1len]
	rs2 := data[rs1len:]

	rs := ""
	for i := len(rs2); i > 0; i-- {
		rs += rs2[i-1 : i]
		if (len(rs)+1)%odd == 0 {
			rs += rs1[len(rs1)-1:]
			rs1 = rs1[:len(rs1)-1]
		}
	}
	log.Debugf("odd:%d, rs1:%s, rs2:%s, rs:%s", odd, rs1, rs2, rs)
	datab, _ := base64.StdEncoding.DecodeString(Base64fix(rs))
	return string(datab)
}

func EncDat2(data string) string {
	if data == "" {
		return data
	}
	l := NumRand(1, len(data)) //random from 1 to data len
	oddnumber := 10
	x := mystring.RandString(oddnumber)
	y := base64.StdEncoding.EncodeToString([]byte(x))
	y = strings.Replace(y, "=", "", -1)

	data = data[:l] + y + data[l:]
	data = base64.StdEncoding.EncodeToString([]byte(data))
	data = strings.Replace(data, "=", "", -1)
	return x + data
}

//Decode: decode for func encDat2 of javascript and EncDat2
func Decode(data string) string {
	if data == "" {
		return data
	}
	if len(data) < 10 {
		log.Errorf("cannot decode %s", data)
		return data
	}
	//fix + char
	data = strings.Replace(data, " ", "+", -1)

	x := 10
	xstr := data[:x]
	data = data[x:]
	//data, _ = DecompressFromBase64(data)
	data = Base64fix(data)
	datab, _ := base64.StdEncoding.DecodeString(data)
	xb64 := base64.StdEncoding.EncodeToString([]byte(xstr))
	xb64 = strings.Replace(xb64, "=", "", -1)
	data = strings.Replace(string(datab), xb64, "", 1)
	return data
}

func Encode2(data string) string {
	if data == "" {
		return data
	}
	oddnumber := NumRand(1, 9)
	x := mystring.RandString(oddnumber)
	//log.Debugf("x: %s", x)
	y := base64.StdEncoding.EncodeToString([]byte(x))
	y = strings.Replace(y, "=", "", -1)
	//log.Debugf("y: %s", y)
	data = CompressToBase64(data)
	//log.Debugf("datacomp: %s", data)
	l := NumRand(2, len(data))
	data = data[:l] + y + data[l:]
	//log.Debugf("data: %s", data)

	data = strings.Replace(data, "=", "", -1)
	oddb64 := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(oddnumber)))
	oddb64 = strings.Replace(oddb64, "=", "", -1)
	return x + data + oddb64
}
func Decode2(data string) string {
	if data == "" {
		return data
	}
	if len(data) < 10 {
		log.Errorf("cannot decode %s", data)
		return data
	}
	oddb64, _ := base64.StdEncoding.DecodeString(Base64fix(data[len(data)-2:]))
	odd, _ := strconv.Atoi(string(oddb64))

	x := data[:odd]
	log.Debugf("x: %s", x)
	y := base64.StdEncoding.EncodeToString([]byte(x))
	y = strings.Replace(y, "=", "", -1)
	log.Debugf("y: %s", y)
	data = data[odd : len(data)-2]
	ycount := strings.Count(data, y)
	i := 0
	ypos := 0
	for {
		if i > ycount {
			log.Errorf("cannot decode %s in loop ypos", data)
			break
		}
		i++
		ypos = strings.Index(data[ypos:], y)
		if ypos == -1 {
			log.Errorf("cannot decode %s y not found", data)
			break
		}
		//test decode
		datatest := data[:ypos] + data[ypos+len(y):]
		log.Debugf("datatest: %s", datatest)
		datatest, _ = DecompressFromBase64(datatest)
		if datatest != "" {
			return datatest

		}
	}
	return ""
}

func Encode3(data string) string {
	if data == "" {
		return data
	}
	data = base64.StdEncoding.EncodeToString([]byte(data))
	data = strings.Replace(data, "=", "", -1)
	var datalen = len(data)
	oddnumber := NumRand(1, 9)
	if datalen < oddnumber {
		oddnumber = datalen
	}

	x := mystring.RandString(oddnumber)

	y := base64.StdEncoding.EncodeToString([]byte(x))
	y = strings.Replace(y, "=", "", -1)

	oddb64 := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(oddnumber)))
	oddb64 = strings.Replace(oddb64, "=", "", -1)

	data = y + data[:datalen-oddnumber] + x + data[datalen-oddnumber:] + oddb64

	return data
}

func Decode3(data string) string {
	if data == "" {
		return data
	}
	oddb64 := data[len(data)-2:]
	data = data[:len(data)-2]

	oddb, _ := base64.StdEncoding.DecodeString(Base64fix(oddb64))
	oddnumber, _ := strconv.Atoi(string(oddb))

	data2 := data[len(data)-(oddnumber*2):]
	data = data[:len(data)-(oddnumber*2)]
	x := data2[:oddnumber]

	data2 = data2[oddnumber:]
	y := base64.StdEncoding.EncodeToString([]byte(x))
	y = strings.Replace(y, "=", "", -1)

	data1 := strings.Replace(data, y, "", 1)
	data = data1 + data2
	datab, _ := base64.StdEncoding.DecodeString(Base64fix(data))
	data = string(datab)
	return data
}

func Encode4(data string) string {

	var x = NumRand(2, 9)

	//x2b := []byte(x2)

	xstr := mystring.RandString(x)

	data = xstr + data

	//data = CompressToBase64(data)
	data = base64.StdEncoding.EncodeToString([]byte(data))
	data = strings.Replace(data, "=", "", -1)

	b := []byte(data)

	var result1 []byte
	var result2 []byte

	for i := len(b) - 1; i >= 0; i-- {

		if i%x == 0 {
			result1 = append(result1, b[i]) // string([]rune(data)[i])

		} else {
			result2 = append(result2, b[i])
		}
	}

	strb64 := string(result1) + string(result2)

	var x2 = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(x)))
	x2 = strings.Replace(x2, "=", "", -1)

	strb64 = strb64[:x] + xstr + strb64[x:]
	strb64 = strb64[:len(strb64)/2] + x2 + strb64[len(strb64)/2:]
	return strb64
}

func Decode4(data string) string {
	if data == "" {
		return data
	}
	ld := len(data)/2 - 1
	var x2 = data[ld : ld+2]

	data = data[:ld] + data[ld+2:]
	xb, _ := base64.StdEncoding.DecodeString(Base64fix(x2))
	x, _ := strconv.Atoi(string(xb))
	data = data[:x] + data[2*x:]

	y := int(math.Ceil(float64(len(data)) / float64(x)))
	rs := data[y-1 : y]
	y--
	leny := y

	for i := len(data) - 1; i > leny; i-- {
		rs += data[i : i+1]
		if len(rs)%x == 0 && y > 0 {
			rs += data[y-1 : y]
			y--
		}
	}

	//data, _ = DecompressFromBase64(rs)
	tmp, _ := base64.StdEncoding.DecodeString(Base64fix(rs))
	data = string(tmp)

	data = data[x:]
	return data
}

//encode for wapi - not completed
func EncodeW(data string) string {
	if data == "" {
		return data
	}
	var x2 = base64.StdEncoding.EncodeToString([]byte(data))
	x2 = strings.Replace(x2, "=", "", -1)
	return x2
}

//decode for wapi - not completed
func DecodeW(code string) string {
	if code == "" {
		return code
	}
	var rt string = ""
	key := code
	//key = "kZXUuYkRWzUgQk92YoNwRdh92Q3SZtFmb9Wa0NW"
	if key == rt {
		return rt
	}

	oddstr := "d"
	l := int(math.Floor((float64)(len(key)-2) / 2))
	num := key[l : l+2]

	key = key[:l] + key[l+2:]

	byteDecode, _ := base64.StdEncoding.DecodeString(Base64fix(num))
	num = string(byteDecode)

	floatNum, _ := strconv.ParseFloat(num, 64)
	intNum := (int)(floatNum)
	if intNum > 0 {
		//print_r($num);print_r("\r\n");
		//get odd string
		lf := math.Ceil((float64)(len(key)) / floatNum)
		oddstr = key[:int(lf)]
		ukey := strings.Replace(key, oddstr, "", 1)
		base64str := ""

		for i := len(oddstr) - 1; i >= 0; i-- {
			base64str += string(oddstr[len(oddstr)-1:])
			oddstr = oddstr[:len(oddstr)-1]
			if len(ukey)-intNum+1 > 0 {
				base64str += mystring.Reverse(string(ukey[len(ukey)-intNum+1:]))
			} else {
				base64str += mystring.Reverse(ukey)
			}
			if i > 0 {
				ukey = ukey[:len(ukey)-intNum+1]
			}
		}
		base64str = base64str[:len(base64str)-intNum]

		byteDecode, _ := base64.StdEncoding.DecodeString(base64str)
		//byteDecode, _ := DecompressFromBase64(base64str)

		rt = string(byteDecode)

	}
	return rt
}
