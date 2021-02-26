package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

/*
Unicode: 定义0x000000-0x10FFFF的数字与字符的关系, 数字叫字符的码点, 分为17个平面, 码点的存储方式由各种编码方式进行定义
plane 0(基本平面): 0x000000 - 0x00FFFF, 0x00D800-0x00D8FF(高位代理)和0x00DC00-0x00DFFF(低位代理)不表示任何字符
plane 1-17(扩展平面): 0x010000 - 0x10FFFF
*/
type Encoder interface {
	Encode(uint64) string
}

type Decoder interface {
	Decode(string) uint64
}

type Codec interface {
	Encoder
	Decoder
}

/*
UTF8: 使用1-4个可变的字节进行进行编码存储
0x000000-0x00007F: 0xxxxxxx
0x000080-0x0007FF: 110xxxxx 10xxxxxx
0x000800-0x00FFFF: 1110xxxx 10xxxxxx 10xxxxxx
0x010000-0x10FFFF: 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
*/
type UTF8 struct{}

func (u *UTF8) Encode(codepoint uint64) string {
	if codepoint <= 0x7F {
		return fmt.Sprintf("%08b", codepoint)
	} else if codepoint <= 0x7FF {
		txt := fmt.Sprintf("%011b", codepoint)
		return fmt.Sprintf("110%5s106%s", txt[:5], txt[5:])
	} else if codepoint <= 0xFFFF {
		txt := fmt.Sprintf("%016b", codepoint)
		return fmt.Sprintf("1110%4s10%6s10%6s", txt[:4], txt[4:10], txt[10:])
	} else if codepoint <= 0x10FFFF {
		txt := fmt.Sprintf("%021b", codepoint)
		return fmt.Sprintf("11110%3s10%6s10%6s10%6s", txt[:3], txt[3:9], txt[9:15], txt[15:])
	}
	return ""
}

func (u *UTF8) Decode(txt string) (codepoint uint64) {
	if strings.HasPrefix(txt, "11110") {
		codepoint, _ = strconv.ParseUint(txt[5:8]+txt[10:16]+txt[18:24]+txt[26:], 2, 64)
	} else if strings.HasPrefix(txt, "1110") {
		codepoint, _ = strconv.ParseUint(txt[4:8]+txt[10:16]+txt[18:], 2, 64)
	} else if strings.HasPrefix(txt, "110") {
		codepoint, _ = strconv.ParseUint(txt[3:8]+txt[10:], 2, 64)
	} else {
		codepoint, _ = strconv.ParseUint(txt, 2, 64)
	}
	return
}

/*
UTF16: 使用1-2个可变的字节进行进行编码存储
0x000000-0x00FFFF: xxxxxxxx xxxxxxxx
0x00FFFF-0x10FFFF: 110110xx xxxxxxxx 110111xx xxxxxxxx
*/
type UTF16 struct{}

func (u *UTF16) Encode(codepoint uint64) string {
	if codepoint <= 0xFFFF {
		if codepoint >= 0xD800 && codepoint <= 0xDFFF {
			return ""
		}
		return fmt.Sprintf("%016b", codepoint)
	} else if codepoint <= 0x10FFFF {
		codepoint -= 0x10000
		txt := fmt.Sprintf("%020b", codepoint)
		// high, _ := strconv.ParseUint(txt[:10], 2, 64)
		// low, _ := strconv.ParseUint(txt[10:], 2, 64)
		// return fmt.Sprintf("%016b", high|0xD800) + fmt.Sprintf("%016b", low|0xDC00)
		return fmt.Sprintf("110110%10s110111%10s", txt[:10], txt[10:])
	}
	return ""
}

func (u *UTF16) Decode(txt string) (codepoint uint64) {
	if strings.HasPrefix(txt, "11011") {
		codepoint, _ = strconv.ParseUint(txt[6:16]+txt[22:], 2, 64)
		codepoint += 0x10000
	} else {
		codepoint, _ = strconv.ParseUint(txt, 2, 64)
	}
	return
}

/*
UTF32: 使用4个字节进行进行编码存储
0x000000-0x10FFFF: xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx
*/
type UTF32 struct{}

func (u *UTF32) Encode(codepoint uint64) string {
	return fmt.Sprintf("%032b", codepoint)
}

func (u *UTF32) Decode(txt string) uint64 {
	codepoint, _ := strconv.ParseUint(txt, 2, 64)
	return codepoint
}

func main() {
	codepoints := []uint64{
		uint64('我'),
		uint64('A'), uint64(0x10000), uint64(0x10FFFF),
		uint64(0x7f), uint64(0x80), uint64(0x7ff), uint64(0x800), uint64(0xffff),
	}
	fmt.Println("codepoint:")
	for _, codepoint := range codepoints {
		fmt.Println(codepoint)
	}

	utf32 := new(UTF32)

	fmt.Println("utf32:")
	for _, codepoint := range codepoints {
		txt := utf32.Encode(codepoint)
		fmt.Println(codepoint, txt, utf32.Decode(txt))
	}

	utf16 := new(UTF16)

	fmt.Println("utf16:")
	for _, codepoint := range codepoints {
		txt := utf16.Encode(codepoint)
		fmt.Println(codepoint, txt, utf16.Decode(txt))
	}

	utf8 := new(UTF8)

	fmt.Println("utf8:")
	for _, codepoint := range codepoints {
		txt := utf8.Encode(codepoint)
		fmt.Println(codepoint, txt, utf8.Decode(txt))
	}

	fmt.Println(ioutil.ReadFile("txts/utf8.txt"))
	fmt.Println(ioutil.ReadFile("txts/utf16be.txt"))
	fmt.Println(ioutil.ReadFile("txts/utf16le.txt"))
}
