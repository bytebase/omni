package catalog

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"io"
	"strings"
	"sync"
	"unicode/utf8"
)

// decodeBytesInColumnCharset validates and decodes raw binary literal bytes
// (0x... / b'...') according to the column's declared MySQL 8.0 character set
// (measured on mysql:8.0.46) into a UTF-8 Go string.
func decodeBytesInColumnCharset(raw []byte, charset string) (string, bool) {
	cs := strings.ToLower(strings.TrimSpace(charset))
	switch cs {
	case "", "utf8mb4", "utf8", "utf8mb3":
		if !utf8.Valid(raw) {
			return "", false
		}
		s := string(raw)
		if cs == "utf8" || cs == "utf8mb3" {
			for _, r := range s {
				if r > 0xFFFF {
					return "", false
				}
			}
		}
		return s, true
	case "binary":
		return string(raw), true
	case "ascii", "swe7":
		for _, b := range raw {
			if b >= 0x80 {
				return "", false
			}
		}
		return string(raw), true
	case "ucs2", "utf16", "utf16le":
		if len(raw)%2 != 0 {
			raw = append([]byte{0}, raw...)
		}
		u16 := make([]uint16, len(raw)/2)
		for i := 0; i < len(u16); i++ {
			if cs == "utf16le" {
				u16[i] = uint16(raw[2*i]) | uint16(raw[2*i+1])<<8
			} else {
				u16[i] = uint16(raw[2*i])<<8 | uint16(raw[2*i+1])
			}
		}
		var b strings.Builder
		for i := 0; i < len(u16); i++ {
			w := u16[i]
			if cs == "ucs2" {
				if w >= 0xD800 && w <= 0xDFFF {
					return "", false
				}
				b.WriteRune(rune(w))
				continue
			}
			if w >= 0xD800 && w <= 0xDBFF {
				if i+1 >= len(u16) || u16[i+1] < 0xDC00 || u16[i+1] > 0xDFFF {
					return "", false
				}
				r := 0x10000 + ((rune(w)-0xD800)<<10 | (rune(u16[i+1]) - 0xDC00))
				b.WriteRune(r)
				i++
			} else if w >= 0xDC00 && w <= 0xDFFF {
				return "", false
			} else {
				b.WriteRune(rune(w))
			}
		}
		return b.String(), true
	case "utf32":
		if rem := len(raw) % 4; rem != 0 {
			raw = append(make([]byte, 4-rem), raw...)
		}
		var b strings.Builder
		for i := 0; i < len(raw); i += 4 {
			cp := uint32(raw[i])<<24 | uint32(raw[i+1])<<16 | uint32(raw[i+2])<<8 | uint32(raw[i+3])
			if cp > 0x10FFFF || (cp >= 0xD800 && cp <= 0xDFFF) {
				return "", false
			}
			b.WriteRune(rune(cp))
		}
		return b.String(), true
	case "sjis", "cp932", "ujis", "eucjpms", "gbk", "gb2312", "gb18030", "big5", "euckr":
		return decodeLegacyMultibyteBytes(raw, cs)
	default:
		if sbIdx, ok := singleByteCharsetIndex[cs]; ok {
			ensureCharsetTablesLoaded()
			tbl := charsetTables.sb[sbIdx*128 : (sbIdx+1)*128]
			var b strings.Builder
			for _, c := range raw {
				if c < 0x80 {
					b.WriteByte(c)
				} else {
					b.WriteRune(tbl[c-0x80])
				}
			}
			return b.String(), true
		}
		if !utf8.Valid(raw) {
			return "", false
		}
		return string(raw), true
	}
}

func decodeLegacyMultibyteBytes(raw []byte, cs string) (string, bool) {
	ensureCharsetTablesLoaded()
	var b strings.Builder
	for i := 0; i < len(raw); {
		b0 := raw[i]
		if b0 <= 0x7F {
			b.WriteByte(b0)
			i++
			continue
		}
		switch cs {
		case "sjis", "cp932":
			tbl := charsetTables.sjis
			if cs == "cp932" {
				tbl = charsetTables.cp932
			}
			if b0 >= 0xA1 && b0 <= 0xDF {
				b.WriteRune(tbl[int(b0-0xA1)])
				i++
				continue
			}
			if ((b0 >= 0x81 && b0 <= 0x9F) || (b0 >= 0xE0 && b0 <= 0xFC)) && i+1 < len(raw) {
				b1 := raw[i+1]
				var col int
				if b1 >= 0x40 && b1 <= 0x7E {
					col = int(b1 - 0x40)
				} else if b1 >= 0x80 && b1 <= 0xFC {
					col = 63 + int(b1-0x80)
				} else {
					return "", false
				}
				var row int
				if b0 <= 0x9F {
					row = int(b0 - 0x81)
				} else {
					row = 31 + int(b0-0xE0)
				}
				b.WriteRune(tbl[63+row*188+col])
				i += 2
				continue
			}
			return "", false
		case "ujis", "eucjpms":
			tbl := charsetTables.ujis
			if cs == "eucjpms" {
				tbl = charsetTables.eucjpms
			}
			if b0 == 0x8E && i+1 < len(raw) && raw[i+1] >= 0xA1 && raw[i+1] <= 0xDF {
				b.WriteRune(tbl[int(raw[i+1]-0xA1)])
				i += 2
				continue
			}
			if b0 >= 0xA1 && b0 <= 0xFE && i+1 < len(raw) && raw[i+1] >= 0xA1 && raw[i+1] <= 0xFE {
				idx := 63 + int(b0-0xA1)*94 + int(raw[i+1]-0xA1)
				b.WriteRune(tbl[idx])
				i += 2
				continue
			}
			if b0 == 0x8F && i+2 < len(raw) && raw[i+1] >= 0xA1 && raw[i+1] <= 0xFE && raw[i+2] >= 0xA1 && raw[i+2] <= 0xFE {
				idx := 63 + 8836 + int(raw[i+1]-0xA1)*94 + int(raw[i+2]-0xA1)
				b.WriteRune(tbl[idx])
				i += 3
				continue
			}
			return "", false
		case "gb2312":
			if b0 >= 0xA1 && b0 <= 0xF7 && i+1 < len(raw) && raw[i+1] >= 0xA1 && raw[i+1] <= 0xFE {
				idx := int(b0-0xA1)*94 + int(raw[i+1]-0xA1)
				b.WriteRune(charsetTables.gb2312[idx])
				i += 2
				continue
			}
			return "", false
		case "gbk":
			if b0 >= 0x81 && b0 <= 0xFE && i+1 < len(raw) {
				b1 := raw[i+1]
				var col int
				if b1 >= 0x40 && b1 <= 0x7E {
					col = int(b1 - 0x40)
				} else if b1 >= 0x80 && b1 <= 0xFE {
					col = 63 + int(b1-0x80)
				} else {
					return "", false
				}
				b.WriteRune(charsetTables.gbk[int(b0-0x81)*190+col])
				i += 2
				continue
			}
			return "", false
		case "gb18030":
			if b0 >= 0x81 && b0 <= 0xFE && i+1 < len(raw) {
				b1 := raw[i+1]
				if (b1 >= 0x40 && b1 <= 0x7E) || (b1 >= 0x80 && b1 <= 0xFE) {
					col := int(b1 - 0x40)
					if b1 >= 0x80 {
						col = 63 + int(b1-0x80)
					}
					b.WriteRune(charsetTables.gb18030_2b[int(b0-0x81)*190+col])
					i += 2
					continue
				}
				if b1 >= 0x30 && b1 <= 0x39 && i+3 < len(raw) && raw[i+2] >= 0x81 && raw[i+2] <= 0xFE && raw[i+3] >= 0x30 && raw[i+3] <= 0x39 {
					lin := ((int(b0-0x81)*10+int(b1-0x30))*126+int(raw[i+2]-0x81))*10 + int(raw[i+3]-0x30)
					if lin < len(charsetTables.gb18030_4b) {
						b.WriteRune(charsetTables.gb18030_4b[lin])
					} else if lin >= 189000 && lin <= 1237575 {
						b.WriteRune(rune(0x10000 + (lin - 189000)))
					} else {
						b.WriteRune('?')
					}
					i += 4
					continue
				}
			}
			return "", false
		case "big5":
			if b0 >= 0xA1 && b0 <= 0xF9 && i+1 < len(raw) {
				b1 := raw[i+1]
				var col int
				if b1 >= 0x40 && b1 <= 0x7E {
					col = int(b1 - 0x40)
				} else if b1 >= 0xA1 && b1 <= 0xFE {
					col = 63 + int(b1-0xA1)
				} else {
					return "", false
				}
				b.WriteRune(charsetTables.big5[int(b0-0xA1)*157+col])
				i += 2
				continue
			}
			return "", false
		case "euckr":
			if b0 >= 0x81 && b0 <= 0xFE && i+1 < len(raw) {
				b1 := raw[i+1]
				var col int
				if b1 >= 0x41 && b1 <= 0x5A {
					col = int(b1 - 0x41)
				} else if b1 >= 0x61 && b1 <= 0x7A {
					col = 26 + int(b1-0x61)
				} else if b1 >= 0x81 && b1 <= 0xFE {
					col = 52 + int(b1-0x81)
				} else {
					return "", false
				}
				b.WriteRune(charsetTables.euckr[int(b0-0x81)*178+col])
				i += 2
				continue
			}
			return "", false
		}
		return "", false
	}
	return b.String(), true
}

var singleByteCharsetIndex = map[string]int{
	"armscii8": 0,
	"cp1250":   1,
	"cp1251":   2,
	"cp1256":   3,
	"cp1257":   4,
	"cp850":    5,
	"cp852":    6,
	"cp866":    7,
	"dec8":     8,
	"geostd8":  9,
	"greek":    10,
	"hebrew":   11,
	"hp8":      12,
	"keybcs2":  13,
	"koi8r":    14,
	"koi8u":    15,
	"latin1":   16,
	"latin2":   17,
	"latin5":   18,
	"latin7":   19,
	"macce":    20,
	"macroman": 21,
	"tis620":   22,
}

var charsetTablesOnce sync.Once

var charsetTables struct {
	sb         []rune
	sjis       []rune
	cp932      []rune
	ujis       []rune
	eucjpms    []rune
	gbk        []rune
	gb18030_2b []rune
	gb2312     []rune
	gb18030_4b []rune
	big5       []rune
	euckr      []rune
}

func unpackDeltaRunes(b64 string) []rune {
	comp, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	zr, err := zlib.NewReader(bytes.NewReader(comp))
	if err != nil {
		panic(err)
	}
	raw, err := io.ReadAll(zr)
	_ = zr.Close()
	if err != nil {
		panic(err)
	}
	n := len(raw) / 4
	out := make([]rune, n)
	var prev uint32
	for i := 0; i < n; i++ {
		prev += binary.LittleEndian.Uint32(raw[4*i : 4*i+4])
		out[i] = rune(prev)
	}
	return out
}

func applySparseRuneDiff(base []rune, b64 string) []rune {
	out := make([]rune, len(base))
	copy(out, base)
	comp, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	zr, err := zlib.NewReader(bytes.NewReader(comp))
	if err != nil {
		panic(err)
	}
	raw, err := io.ReadAll(zr)
	_ = zr.Close()
	if err != nil {
		panic(err)
	}
	for i := 0; i+6 <= len(raw); i += 6 {
		idx := int(binary.LittleEndian.Uint16(raw[i : i+2]))
		r := rune(binary.LittleEndian.Uint32(raw[i+2 : i+6]))
		out[idx] = r
	}
	return out
}

func ensureCharsetTablesLoaded() {
	charsetTablesOnce.Do(func() {
		charsetTables.sb = unpackDeltaRunes(csTableSB)
		charsetTables.sjis = unpackDeltaRunes(csTableSJIS)
		charsetTables.cp932 = applySparseRuneDiff(charsetTables.sjis, csDiffCP932)
		charsetTables.ujis = unpackDeltaRunes(csTableUJIS)
		charsetTables.eucjpms = applySparseRuneDiff(charsetTables.ujis, csDiffEUCJPMS)
		charsetTables.gbk = unpackDeltaRunes(csTableGBK)
		charsetTables.gb18030_2b = applySparseRuneDiff(charsetTables.gbk, csDiffGB18030_2B)
		charsetTables.gb2312 = unpackDeltaRunes(csTableGB2312)
		charsetTables.gb18030_4b = unpackDeltaRunes(csTableGB18030_4B)
		charsetTables.big5 = unpackDeltaRunes(csTableBIG5)
		charsetTables.euckr = unpackDeltaRunes(csTableEUCKR)
	})
}

const csTableSB = "" +
	"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" +
	"+P6/S2Pe7c6NGGRzXez336tbYN7PuMbwPwR1o+A="

const csTableSJIS = "" +
	"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" +
	"FdBbGO9sJLlSK/nXM9r5+8QZyUdvF7aQfXun1vI9+QLbpG8qZQ9vCvimVnjqFeg376vXy2/weuNzxcn7BvnClOBpJFiKZUf4/GWs" +
	"3w14OpCVCyU72qmflZqX580q4X+z+m0g+igTjBuC3zj2fuBGwextRx8TaaLnqiAGlxfQYhinC2VDnWSMz5nYOp0MQ4xZtDO/kQx0" +
	"3nbZmN6H2E9tlmiuJRq7THbpEvnWnVSeVZzKx16qZbP113q0EX1tjnfapJ5P7b19A3me9EdP+W0rNF5XzSlP82ktfK7X+BnZRH69" +
	"22rMVaKNCvno7WSz1WhOhUEeKV/jVkuP5AcxQq8/wv/BYe+H2m/imB1n+67W7aTJVMAfqUAXbNX8EvptiGSQg9r1t9TTgT2TF8hw" +
	"LyPLAhsmzGN7P6UwkOEFQfzEx9HyNdYEABi2eWd+uyjQ417HePpZAlF3W+FkW576yu4Sq/Y2z/bAd/XyPBn0mQ3G8nJkV/vC2w7l" +
	"wp337ZdgqLfq5Xjer2Fa/pTXq/2kC/w+DJ8XLJb9Z/tt31nlaLipcGp+zRXC5xatfan418NTrT7bCI9eTr4rWi8PeLogyEWmAhpo" +
	"EOhn7ytXiicygdxuHOA2Fdwng7JUkEtMBrK2RDzhaaIuyHGGub+E8jCHBGuZDebcXPhtJR7YM7B5NgXxJYPX9h3/SnLBx066ylfr" +
	"Kf/qNNVfrD7j4HfbK4TTVdIVPnd3kXwiH1/z8el77DfC1DYlGbqv9EZSvoWPTx4tPrZ46RLFPrOSB/1kO++lMb7VWrTUuDaXUdID" +
	"9p17Vb3z/yaqnv02YTfhY7NiKJMlC40mzb9cLhwX6rtK+6bmcfnfWemUAZrXOsFcLXnmc1UWN7lB/lWVbBO/v2WT/PdYdv4g2fBV" +
	"oi3D0WTbVyh90F0ycankwzbJ6CmyW/L0Lfso8cd3WuNGipl5X97nU8J9LFsko32uwseftwf8t1W82kA8lm9BqL2dTPDyq0Ng6yWC" +
	"PEg3wdtUtLhKfpXtS7yj3uWKfJwnJdz4PGpH4djzYiqAZavkyBrJ127qw5ev0tXmfphot7niiUnFQJaobVZ+TE2wv+Z3wmedaLCn" +
	"cLFRcOyYKwj5qtL1u14wx1rv/prPumBf1AbxW0Jy9yPFKjK7+GMNA9m7Tv0W7LL/x/scBZJjdfKB5kpf7COabyi+qBUPrxYNrRYt" +
	"FAf7ZRoFcfnQPmwVyMet4ptC0Vd+kHPsHdBeJvBvUkHuPikemC0eqw/8zYJdZF8Ye88EsYJMsI+kLtBHqWA/m5eXLYL7wiC/7O18" +
	"7+s3CHJiZXrfQnB7Guwf2BbV4uk6veskOq8M8l6bNaf2gV3YTP37nEOd1r2j9pIktV4r9K6P1q6xYv97BvGKAs3hKvttkU0Ob40D" +
	"HisKYqS1gZ0Qi+43qyzcs5YK9nel/wud9nSw/2zX/KCP2RQGfXi7eXugW73e+lE0eSuC/f5uO2PVVUG+Nz+I73QM9lp2DPzQgmAv" +
	"RV1ATzvyjjSesjaO31YeMKl9mNkA9lLFB3YP9m00DfZ0VIs+KiXrE4Ed0zyIFVQHMXbvI3SUTPMxjPKAZrwNtVXr4cebLt24UfHi" +
	"7eK3TGCrej0eB7kQv45hHsqPuVprHeZ3vV/TPLBrGgdxvIaC1dvynj62S5a10HrnB7GAOuGvVPd1u+TKw9jJPgEOCgN66RjE/f2e" +
	"BYvNXhTQrZeh3ifyemCr5G+rwC+2OpOkv8P/ifOufCfPl8kgluFt5r72bfs2h4va4P/9FQbxyOZB/Cwt38Tuz0XIPbSH+32GhsFe" +
	"1nrhrjjwK5IBHryuqpct0C7IxdcG/1MwzIfkBzyYDX43p0Y05GNdXaWbm4rmV2q9jlGddZpjS9kbVv6y1jYj26wmWNdJssN9XKav" +
	"dOtq5fxWKWfgbYuCIM6alWxMB7r1R+3X2i54Vyiu1TrYU+H3Kvm84ybppVaKjfucaKXmUiCY/F7aNcKvj69vDOyiZOAHTWKgwcmd" +
	"5VnlHzoEdv0Ovkv8tf8V7pfIC/yGMFZTH+Rse8yL488671z7giAumw10fZ3shdCXaxzElbYrJ7xMbcPYbnFAC0WBDGqs+8LAX911" +
	"z0fTYP+WvbsBhN1X4nIin8ivLwraFAWx4PogT18U5PAKAv0S7smuCXzH9C765RfvxPETg3buzSqWnMsEMbVM0G9xsG+6LohDN9jF" +
	"NyvTeTYDPwshrcXhqWjt8NdCc+8huvGxGj8HH6Ntt0sM1+vdvQO6ykjX+H2rHYM9BSHPpoN+E8G+uuJgz1wYu/Mx69ogxufl9cBg" +
	"X1x+sIeno9r5vZDbgjix19EFQe61VPptvubk8z2lAW34PR8V4vXSaOf/C0sFOAnjBYcENoc+xf0r2ivcBaZUgJswHp0N7LpGu8Qv" +
	"Erv0mQjwlwhin3Gw7zUR2J+t5I/6fF1mlxhPfkBrhfIFSoL1DPel+v0L9t7yqOODPZSdA/+hUDRdozY+l1GjfQ1rAx3k9+RWBPkY" +
	"Lw+zwbco1s8lAY94vb9E+i8lO9Zwa99VXy7d0Fn2el3gy8RBXKppwMdr9a5A/Veofif5O0WK7RUGNmcjraHXK/7/yBVItleInjvL" +
	"H8oL7LbWwRqeGMiMMK6dCPIXno/yA3rzPt0YkHZh3l/TSTrQrXFgY+26L+i/ojFPu6NxykZ0dXMqCngnP+Dn7UEuKhNc60Ur3vZo" +
	"G+wRvQ/ise8xH0S5FaGsr4I5J1LhT/fm/m9e7sgduSN35I7ckTtyR+7IHbkjd+SO3JE7ckfuyB25I3fkjtyRO3JH7sgduSN35I7c" +
	"kTtyR+7IHbkjd+SO3JE7ckfuyB25I3fkjtyRO3JH7sgduSN3/P94/CebMekp"

const csDiffCP932 = "" +
	"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" +
	"7H/RaR77X7TFTyI9n/ZZJJ/P/hfn8bvoxG9BP0Rn/Cm64C/RFX+LbvhHdMe/ogcGwj0xQrgXRgz3xkhh8ztqMPwII4ejYpRwRDwU" +
	"DtzjexGOgIfDEfFIOBIeDQfxWDgyHg9HwRPhqHgyHA1PhaPj6XAMPBOOiWfDsfBcODaeD8fBC+G4eDEcDy+F4+PlcAK8Ek6IV8OJ" +
	"8Fo4MV4PJ8Eb4aR4M5wMb4WT4+1wCrwTTol3w6nwXjg13g+nwQfhtPgwHMJHYYGPw2F8ErbwaVjis7CNz8MOvggrfBl28VVY4+uw" +
	"h2/CPr4Np8N34fT4PpwBP4Qz4sdwJvwUzoyfw1nwSzgrfg1nw2/h7Pg9nAN/hHPiz3Au/BXOjb/DefBPOC/+DefDgJUfI1gFMKJV" +
	"ECNZhTBoFcbIVhGMYhXFqFYxjGYVx+hWCYxhlcSYVimMZZXG2FYZjGOVxbhWOYxnlcf4VgVMYFXEhFYlTGRVxsRWFUxiVcWkVjVM" +
	"ZlXH5FYNTGHVxJRWLUxl1cbUVh1MY9XFtFY9DFn1UVgNMGw1RMtqhNJqjLbVBB2rKSqrGbpWc9RWC/SsluhbrTCd1RrTW20wg9UW" +
	"M1rtMJPVHjNbHTCL1RGzWp0wm9UZs1tdMIfVFXNa3TCX1R1zWz0wj9UT81q9MJ/VG/NbfbCA1RcLWv2wkNUfC1sDsIg1EItag7CY" +
	"9Q8WtwZjCWsIlrSGYilrGJa2hmMZawSWtUZiOWsUlrdGYwVrDFa0xmIlaxxWtsZjFWsCVrUmYjVrEla3JmMNawrWtKZiLWsa1ram" +
	"Yx1rBta1ZmI9axbWt2ZjA2sONrTmYiNrHja25mMTawE2tRZiM2sRNrf+xRbWYmxpLcFW1lJsbS3DNtZybGutwHbWSmxvrcIO1mrs" +
	"aK3BTtZa7Gytwy7WeuxqbcBu1kbsbm3CHtZm7GltwV7WVuxtbcM+1nbsa+3AftZO7G/twgHWbhxo/YeDrD34j7UXB1v7cIi1H4da" +
	"B3CYdRCHW4dwhPU/HGkdxlHWERxtHcUx1jEcax3HcdYJHG+dxAnWKZxoncZJ1hmcbJ3FKdY5nGqdx2nWBZxuXcQZ1iWcaV3GWdYV" +
	"nG1dxTnWNZxrXcd51g2cb93EBdYtXGjdxkXWHfzXuouLrXu4xLqPS60HuMx6iMutR7jCeowrrSe4ynqKq61nuMZ6jmutF7jOeonr" +
	"rVe4wXqNG603uMl6i5utd7jFeo9brQ+4zfqI261PuMP6jDutL7jL+oq7rW/4n/Ud91g/cK/1E/dZv3C/9RsPWH/woPUXD1mB++x/" +
	"KwIetiLiESsSHrWCeMyKjMetKHjCioonrWh4yoqOp60YeMaKiWetWHjOio3nrTh4wYqLF614eMmKj5etBHjFSohXrUR4zUqM160k" +
	"eMNKijetZHjLSo63rRR4x0qJd61UeM9KjfetNPjASosPrRA+sgQ+tsL4xLLwqSXxmWXjc8vBF5bCl5aLryyNry0P31g+vrXS4Tsr" +
	"Pb63MuAHKyN+tDLhJyszfray4BcrK361suE3Kzt+t3LgDysn/rRy4S8rN/628uAfKy/+tfJhQObHCLIARpQFMZIshEFZGCPLIhhF" +
	"FsWoshhGk8UxuiyBMWRJjClLYSxZGmPLMhhHlsW4shzGk+UxvqyACWRFTCgrYSJZGRPLKphEVsWkshomk9UxuayBKWRNTClrYSpZ" +
	"G1PLOphG1sW0sh6GZH0UsgGGZUO0ZCOUsjHasgk6sikq2Qxd2Ry1bIGebIm+bIXpZGtML9tgBtkWM8p2mEm2x8yyA2aRHTGr7ITZ" +
	"ZGfMLrtgDtkVc8pumEt2x9yyB+aRPTGv7IX5ZG/ML/tgAdkXC8p+WEj2x8JyABaRA7GoHITF5D9YXA7GEnIIlpRDsZQchqXlcCwj" +
	"R2BZORLLyVFYXo7GCnIMVpRjsZIch5XleKwiJ2BVORGryUlYXU7GGnIK1pRTsZachrXldKwjZ2BdORPryVlYX87GBnIONpRzsZGc" +
	"h43lfGwiF2BTuRCbyUXYXP6LLeRibCmXYCu5FFvLZdhGLse2cgW2kyuxvVyFHeRq7CjXYCe5FjvLddhFrseucgN2kxuxu9yEPeRm" +
	"7Cm3YC+5FXvLbdhHbse+cgf2kzuxv9yFA+RuHCj/w0FyD/4j9+JguQ+HyP04VB7AYfIgDpeHcIT8H46Uh3GUPIKj5VEcI4/hWHkc" +
	"x8kTOF6exAnyFE6Up3GSPIOT5VmcIs/hVHkep8kLOF1exBnyEs6Ul3GWvIKz5VWcI6/hXHkd58kbOF/exAXyFi6Ut3GRvIP/yru4" +
	"WN7DJfI+LpUPcJl8iMvlI1whH+NK+QRXyae4Wj7DNfI5rpUvcJ18ievlK9wgX+NG+QY3ybe4Wb7DLfI9bpUfcJv8iNvlJ9whP+NO" +
	"+QV3ya+4W37D/+R33CN/4F75E/fJX7hf/sYD8g8elH/xkAw8YP/LCHhYRsQjMhIelUE8JiPjcRkFT8ioeFJGw1MyOp6WMfCMjIln" +
	"ZSw8J2PjeRkHL8i4eFHGw0syPl6WCfCKTIhXZSK8JhPjdZkEb8ikeFMmw1syOd6WKfCOTIl3ZSq8J1PjfZkGH8i0+FCG8JEU+FiG" +
	"8Ym08KmU+Eza+Fw6+EIqfCldfCU1vpYevpE+vpXp8J1Mj+9lBvwgM+JHmQk/ycz4WWbBLzIrfpXZ8JvMjt9lDvwhc+JPmQt/ydz4" +
	"W+bBPzIv/pX5MGDnxwh2AYxoF8RIdiEM2oUxsl0Eo9hFMapdDKPZxTG6XQJj2CUxpl0KY9mlMbZdBuPYZTGuXQ7j2eUxvl0BE9gV" +
	"MaFdCRPZlTGxXQWT2FUxqV0Nk9nVMbldA1PYNTGlXQtT2bUxtV0H09h1Ma1dD0N2fRR2AwzbDdGyG6G0G6NtN0HHborKboau3Ry1" +
	"3QI9uyX6ditMZ7fG9HYbzGC3xYx2O8xkt8fMdgfMYnfErHYnzGZ3xux2F8xhd8WcdjfMZXfH3HYPzGP3xLx2L8xn98b8dh8sYPfF" +
	"gnY/LGT3x8L2ACxiD8Si9iAsZv+Dxe3BWMIegiXtoVjKHoal7eFYxh6BZe2RWM4eheXt0VjBHoMV7bFYyR6Hle3xWMWegFXtiVjN" +
	"noTV7clYw56CNe2pWMuehrXt6VjHnoF17ZlYz56F9e3Z2MCegw3tudjInoeN7fnYxF6ATe2F2MxehM3tf7GFvRhb2kuwlb0UW9vL" +
	"sI29HNvaK7CdvRLb26uwg70aO9prsJO9Fjvb67CLvR672huwm70Ru9ubsIe9GXvaW7CXvRV729uwj70d+9o7sJ+9E/vbu3CAvRsH" +
	"2v/hIHsP/mPvxcH2Phxi78eh9gEcZh/E4fYhHGH/D0fah3GUfQRH20dxjH0Mx9rHcZx9AsfbJ3GCfQon2qdxkn0GJ9tncYp9Dqfa" +
	"53GafQGn2xdxhn0JZ9qXcZZ9BWfbV3GOfQ3n2tdxnn0D59s3cYF9Cxfat3GRfQf/te/iYvseLrHv41L7AS6zH+Jy+xGusB/jSvsJ" +
	"rrKf4mr7Ga6xn+Na+wWus1/ievsVbrBf40b7DW6y3+Jm+x1usd/jVvsDbrM/4nb7E+6wP+NO+wvusr/ibvsb/md/xz32D9xr/8R9" +
	"9i/cb//GA/YfPGj/xUN24CH7346Ah+2IeMSOhEftIB6zI+NxOwqesKPiSTsanrKj42k7Bp6xY+JZOxaes2PjeTsOXrDj4kU7Hl6y" +
	"4+NlOwFesRPiVTsRXrMT43U7Cd6wk+JNOxnespPjbTsF3rFT4l07Fd6zU+N9Ow0+sNPiQzuEj2yBj+0wPrEtfGpLfGbb+Nx28IWt" +
	"8KXt4itb42vbwze2j2/tdPjOTo/v7Qz4wc6IH+1M+MnOjJ/tLPjFzopf7Wz4zc6O3+0c+MPOiT/tXPjLzo2/7Tz4x86Lf+18GHDy" +
	"YwSnAEZ0CmIkpxAGncIY2SmCUZyiGNUphtGc4hjdKYExnJIY0ymFsZzSGNspg3GcshjXKYfxnPIY36mACZyKmNCphImcypjYqYJJ" +
	"nKqY1KmGyZzqmNypgSmcmpjSqYWpnNqY2qmDaZy6mNaphyGnPgqnAYadhmg5jVA6jdF2mqDjNEXlNEPXaY7aaYGe0xJ9pxWmc1pj" +
	"eqcNZnDaYkanHWZy2mNmpwNmcTpiVqcTZnM6Y3anC+ZwumJOpxvmcrpjbqcH5nF6Yl6nF+ZzemN+pw8WcPpiQacfFnL6Y2FnABZx" +
	"BmJRZxAWc/7B4s5gLOEMwZLOUCzlDMPSznAs44zAss5ILOeMwvLOaKzgjMGKzlis5IzDys54rOJMwKrORKzmTMLqzmSs4UzBms5U" +
	"rOVMw9rOdKzjzMC6zkys58zC+s5sbODMwYbOXGzkzMPGznxs4izAps5CbOYswubOv9jCWYwtnSXYylmKrZ1l2MZZjm2dFdjOWYnt" +
	"nVXYwVmNHZ012MlZi52dddjFWY9dnQ3YzdmI3Z1N2MPZjD2dLdjL2Yq9nW3Yx9mOfZ0d2M/Zif2dXTjA2Y0Dnf9wkLMH/3H24mBn" +
	"Hw5x9uNQ5wAOcw7icOcQjnD+hyOdwzjKOYKjnaM4xjmGY53jOM45geOdkzjBOYUTndM4yTmDk52zOMU5h1Od8zjNuYDTnYs4w7mE" +
	"M53LOMu5grOdqzjHuYZznes4z7mB852buMC5hQud27jIuYP/OndxsXMPlzj3canzAJc5D3G58whXOI9xpfMEVzlPcbXzDNc4z3Gt" +
	"8wLXOS9xvfMKNzivcaPzBjc5b3Gz8w63OO9xq/MBtzkfcbvzCXc4n3Gn8wV3OV9xt/MN/3O+4x7nB+51fuI+5xfud37jAecPHnT+" +
	"4iEn8Ij970TAw05EPOJEwqNOEI85kfG4EwVPOFHxpBMNTznR8bQTA884MfGsEwvPObHxvBMHLzhx8aITDy858fGykwCvOAnxqpMI" +
	"rzmJ8bqTBG84SfGmkwxvOcnxtpMC7zgp8a6TCu85qfG+kwYfOGnxoRPCR47Ax04YnzgWPnUkPnNsfO44+MJR+NJx8ZWj8bXj4RvH" +
	"x7dOOnznpMf3Tgb84GTEj04m/ORkxs9OFvziZMWvTjb85mTH704O/OHkxJ9OLvzl5MbfTh784+TFv04+DKj8GEEVwIiqIEZShTCo" +
	"CmNkVQSjqKIYVRXDaKo4RlclMIYqiTFVKYylSmNsVQbjqLIYV5XDeKo8xlcVMIGqiAlVJUykKmNiVQWTqKqYVFXDZKo6Jlc1MIWq" +
	"iSlVLUylamNqVQfTqLqYVtXDkKqPQjXAsGqIlmqEUjVGWzVBRzVFpZqhq5qjVi3QUy3RV60wnWqN6VUbzKDaYkbVDjOp9phZdcAs" +
	"qiNmVZ0wm+qM2VUXzKG6Yk7VDXOp7phb9cA8qifmVb0wn+qN+VUfLKD6YkHVDwup/lhYDcAiaiAWVYOwmPoHi6vBWEINwZJqKJZS" +
	"w7C0Go5l1Agsq0ZiOTUKy6vRWEGNwYpqLFZS47CyGo9V1ASsqiZiNTUJq6vJWENNwZpqKtZS07C2mo511Aysq2ZiPTUL66vZ2EDN" +
	"wYZqLjZS87Cxmo9N1AJsqhZiM7UIm6t/sYVajC3VEmyllmJrtQzbqOXYVq3AdmoltlersINajR3VGuyk1mJntQ67qPXYVW3Abmoj" +
	"dlebsIfajD3VFuyltmJvtQ37qO3YV+3Afmon9le7cIDajQPVfzhI7cF/1F4crPbhELUfh6oDOEwdxOHqEI5Q/8OR6jCOUkdwtDqK" +
	"Y9QxHKuO4zh1AserkzhBncKJ6jROUmdwsjqLU9Q5nKrO4zR1AaerizhDXcKZ6jLOUldwtrqKc9Q1nKuu4zx1A+erm7hA3cKF6jYu" +
	"UnfwX3UXF6t7uETdx6XqAS5TD3G5eoQr1GNcqZ7gKvUUV6tnuEY9x7XqBa5TL3G9eoUb1GvcqN7gJvUWN6t3uEW9x63qA25TH3G7" +
	"+oQ71Gfcqb7gLvUVd6tv+J/6jnvUD9yrfuI+9Qv3q994QP3Bg+ovHlKBx+x/FQEPq4h4REXCoyqIx1RkPK6i4AkVFU+qaHhKRcfT" +
	"KgaeUTHxrIqF51RsPK/i4AUVFy+qeHhJxcfLKgFeUQnxqkqE11RivK6S4A2VFG+qZHhLJcfbKgXeUSnxrkqF91RqvK/S4AOVFh+q" +
	"ED5SAh+rMD5RFj5VEp8pG58rB18ohS+Vi6+UxtfKwzfKx7cqHb5T6fG9yoAfVEb8qDLhJ5UZP6ss+EVlxa8qG35T2fG7yoE/VE78" +
	"qXLhL5Ubf6s8+Eflxb8qHwbc/BjBLYAR3YIYyS2EQbcwRnaLYBS3KEZ1i2E0tzhGd0tgDLckxnRLYSy3NMZ2y2ActyzGdcthPLc8" +
	"xncrYAK3IiZ0K2EitzImdqtgErcqJnWrYTK3OiZ3a2AKtyamdGthKrc2pnbrYBq3LqZ162HIrY/CbYBhtyFabiOUbmO03SbouE1R" +
	"uc3QdZujdlug57ZE322F6dzWmN5tgxnctpjRbYeZ3PaY2e2AWdyOmNXthNnczpjd7YI53K6Y0+2GudzumNvtgXncnpjX7YX53N6Y" +
	"3+2DBdy+WNDth4Xc/ljYHYBF3IFY1B2Exdx/sLg7GEu4Q7CkOxRLucOwtDscy7gjsKw7Esu5o7C8OxoruGOwojsWK7njsLI7Hqu4" +
	"E7CqOxGruZOwujsZa7hTsKY7FWu507C2Ox3ruDOwrjsT67mzsL47Gxu4c7ChOxcbufOwsTsfm7gLsKm7EJu5i7C5+y+2cBdjS3cJ" +
	"tnKXYmt3GbZxl2NbdwW2c1die3cVdnBXY0d3DXZy12Jndx12cddjV3cDdnM3Ynd3E/ZwN2NPdwv2crdib3cb9nG3Y193B/Zzd2J/" +
	"dxcOcHfjQPc/HOTuwX/cvTjY3YdD3P041D2Aw9yDONw9hCPc/+FI9zCOco/gaPcojnGP4Vj3OI5zT+B49yROcE/hRPc0TnLP4GT3" +
	"LE5xz+FU9zxOcy/gdPciznAv4Uz3Ms5yr+Bs9yrOca/hXPc6znNv4Hz3Ji5wb+FC9zYucu/gv+5dXOzewyXufVzqPsBl7kNc7j7C" +
	"Fe5jXOk+wVXuU1ztPsM17nNc677Ade5LXO++wg3ua9zovsFN7lvc7L7DLe573Op+wG3uR9zufsId7mfc6X7BXe5X3O1+w//c77jH" +
	"/YF73Z+4z/2F+93feMD9gwfdv3jIDTxh/7sR8LAbEY+4kfCoG8RjbmQ87kbBE25UPOlGw1NudDztxsAzbkw868bCc25sPO/GwQtu" +
	"XLzoxsNLbny87CbAK25CvOomwmtuYrzuJsEbblK86SbDW25yvO2mwDtuSrzrpsJ7bmq876bBB25afOiG8JEr8LEbxieuhU9dic9c" +
	"G5+7Dr5wFb50XXzlanztevjG9fGtmw7fuenxvZsBP7gZ8aObCT+5mfGzmwW/uFnxq5sNv7nZ8bubA3+4OfGnmwt/ubnxt5sH/7h5" +
	"8a+bDwM6P0bQBTCiLoiRdCEM6sIYWRfBKLooRtXFMJoujtF1CYyhS2JMXQpj6dIYW5fBOLosxtXlMJ4uj/F1BUygK2JCXQkT6cqY" +
	"WFfBJLoqJtXVMJmujsl1DUyha2JKXQtT6dqYWtfBNLouptX1MKTro9ANMKwboqUbodSN0dZN0NFNUelm6OrmqHUL9HRL9HUrTKdb" +
	"Y3rdBjPotphRt8NMuj1m1h0wi+6IWXUnzKY7Y3bdBXPorphTd8Ncujvm1j0wj+6JeXUvzKd7Y37dBwvovlhQ98NCuj8W1gOwiB6I" +
	"RfUgLKb/weJ6MJbQQ7CkHoql9DAsrYdjGT0Cy+qRWE6PwvJ6NFbQY7CiHouV9DisrMdjFT0Bq+qJWE1Pwup6MtbQU7Cmnoq19DSs" +
	"radjHT0D6+qZWE/Pwvp6NjbQc7ChnouN9DxsrOdjE70Am+qF2Ewvwub6X2yhF2NLvQRb6aXYWi/DNno5ttUrsJ1eie31KuygV2NH" +
	"vQY76bXYWa/DLno9dtUbsJveiN31JuyhN2NPvQV76a3YW2/DPno79tU7sJ/eif31Lhygd+NA/R8O0nvwH70XB+t9OETvx6H6AA7T" +
	"B3G4PoQj9P9wpD6Mo/QRHK2P4hh9DMfq4zhOn8Dx+iRO0Kdwoj6Nk/QZnKzP4hR9Dqfq8zhNX8Dp+iLO0Jdwpr6Ms/QVnK2v4hx9" +
	"Defq6zhP38D5+iYu0Ldwob6Ni/Qd/FffxcX6Hi7R93GpfoDL9ENcrh/hCv0YV+onuEo/xdX6Ga7Rz3GtfoHr9Etcr1/hBv0aN+o3" +
	"uEm/xc36HW7R73Gr/oDb9Efcrj/hDv0Zd+ovuEt/xd36G/6nv+Me/QP36p+4T//C/fo3HtB/8KD+i4d04Cn7X0fAwzoiHtGR8KgO" +
	"4jEdGY/rKHhCR8WTOhqe0tHxtI6BZ3RMPKtj4TkdG8/rOHhBx8WLOh5e0vHxsk6AV3RCvKoT4TWdGK/rJHhDJ8WbOhne0snxtk6B" +
	"d3RKvKtT4T2dGu/rNPhAp8WHOoSPtMDHOoxPtIVPtcRn2sbn2sEXWuFL7eIrrfG19vCN9vGtTofvdHp8rzPgB50RP+pM+Elnxs86" +
	"C37RWfGrzobfdHb8rnPgD50Tf+pc+Evnxt86D/7RefGvzocBLz9G8ApgRK8gRvIKYdArjJG9IhjFK4pRvWIYzSuO0b0SGMMriTG9" +
	"UhjLK42xvTIYxyuLcb1yGM8rj/G9CpjAq4gJvUqYyKuMib0qmMSrikm9apjMq47JvRqYwquJKT1zxziVZ+4Yp/bMHeM0nrljnNYz" +
	"d4xDnrljLDxzxzjsmTvGlmfuGEvP3DG2PVOiOZ4p0ZRnSjTXMyWa9kyJ5nmmRPM9U6Kl80yJlt4zJVoGz5RoGT1zPzmTZ+4nZ/bM" +
	"/eQsnrmfnNUzVVc2z7Rg2T1TgeXwTIeV0zPtVS7PtFe5PdNe5fFMe5XXM+1VPs+0V/k9014V8Ex7VdAz7VUhz7RXhT3TXhXxTHtV" +
	"1DPtVTHPtFfFPdNelfBMe1XSM+1VKc+0V6U9016V8Ux7VdYz7VU5z7RX5T3TXlXwTHtV0TPtVSXPtFeVPdNeVfFMe1XVM+1VNc+0" +
	"V9U9017V8Ex7VdMz7VUtz7RXtT3TXtXxTHtV1zPtVT3PtFf1PdNeNfBMe9XQM+1VI8+0V40901418Ux71dQz7VUzz7RXzT3TXrXw" +
	"THvV0jPtVSvPtFetPdNetfFMe9XWM+1VO8+0V+0901518Ex71dEz7VUnz7RXnT3TXnXxTHvV1TPtVTfPtFfdPdNe9fBMe9XTM+1V" +
	"L8+0V70901718Ux71dcz7VU/z7RX/T3TXg3wTHs10DPt1SDPtFf/eKa9GuyZ9mqIZ9qroZ5pr4Z5pr0a7pn2aoRn2quRnmmvRnmm" +
	"vRrtmfZqjGfaq7Geaa/Geaa9Gu+Z9mqCZ9qriZ5pryZ5pr2a7Jn2aopn2qupnmmvpnmmvZrumfZqhmfaq5meaa9meaa9mu2Z9mqO" +
	"Z9qruZ5pr+Z5pr2a75n2aoFn2quFnmmvFnmmvfrXM+3VYs+0V0s8014t9Ux7tcwz7dVyz7RXKzzTXq30THu1yjPt1WrPtFdrPNNe" +
	"rfVMe7XOM+3Ves+0Vxs8015t9Ex7tckz7dVmz7RXWzzTXm31THu1zTPt1XbPtFc7PNNe7fRMe7XLM+3Vbs+0V/95pr3a45n2aq9n" +
	"2qt9nmmv9numvTrgmfbqoGfaq0Oeaa/+55n26rBn2qsjnmmvjnqmvTrmmfbquGfaqxOeaa9Oeqa9OuWZ9uq0Z9qrM55pr856pr06" +
	"55n26rxn2qsLnmmvLnqmvbrkmfbqsmfaqyueaa+ueqa9uuaZ9uq6Z9qrG55pr256pr265Zn26rZn2qs7nmmv7nqmvbrnmfbqvmfa" +
	"qweeaa8eeqa9euSZ9uqxZ9qrJ55pr556pr165pn26rln2qsXnmmvXnqmvXrlmfbqtWfaqzeeaa/eeqa9eueZ9uq9Z9qrD55prz56" +
	"pr365Jn26rNn2qsvnmmvvnqmvfrmmfbqu2faqx+eaa9+eqa9+uWZ9uq3Z9qrP55pr/56pr3i11+M4Jv2KqJv2qtIvmmvgr5pryL7" +
	"pr2K4pv2Kqpv2qtovmmvovumvYrhm/Yqpm/aq1i+aa9i+6a9iuOb9iqub9qreL5pr+L7pr1K4Jv2KqFv2qtEvmmvEvumvUrim/Yq" +
	"qW/aq2S+aa+S+6a9SuGb9iqlb9qrVL5pr1L7pr1K45v2Kq1v2quQb9or4Zv2Kuyb9sryTXslfdNe2b5prxzftFfKN+2V65v2Svum" +
	"vfJ80175vmmv0vmmvUrvm/Yqg2/aq4y+aa8y+aa9yuyb9iqLb9qrrL5pr7L5pr3K7pv2Kodv2qucvmmvcvmmvcrtm/Yqj2/aq7y+" +
	"aa/y+aa9yu+b9qqAb9qrgr5prwr5pr0q7Jv2qohv2quivmmvivmmvSrum/aqhG/aq5K+aa9K+aa9Ku2b9qqMb9qrsr5pr8r5pr0q" +
	"75v2qoJv2quKvmmvKvmmvarsm/aqim/aq6q+aa+q+aa9qu6b9qqGb9qrmr5pr2r5pr2q7Zv2qo5v2qu6vmmv6vmmvarvm/aqgW/a" +
	"q4a+aa8a+aa9auyb9qqJb9qrpr5pr5r5pr1q7pv2qoVv2quWvmmvWvmmvWrtm/aqjW/aq7a+aa/a+aa9au+b9qqDb9qrjr5przr5" +
	"pr3q7Jv2qotv2quuvmmvuvmmverum/aqh2/aq56+aa96+aa96u2b9qqPb9qrvr5pr/r5pr3q75v2aoBv2quBvmmvBvmmvfrHN+3V" +
	"YN+0V0N8014N9U17Ncw37dVw37RXI3zTXo30TXs1yjft1WjftFdjfNNejfVNezXON+3VeN+0VxN8015N9E17Nck37dVk37RXU3zT" +
	"Xk31TXs1zTft1XTftFczfNNezfRNezXLN+3VbN+0V3N8017N9U17Nc837dV837RX/wcqqCIN"

const csTableUJIS = "" +
	"eNrtfQeYVtW19ulfmT4DM/QOIoKAYC+JimJFg9gVuzExaGKLGq/YFXvBmhB77MECiAUVEVFQESwoSi9KH9oA087/Lvbazyy2+3yY" +
	"P3/u733uN89z5rR9dll79bX2/s6IHcd1/nceZ/SJY32dne44Hs6r+sZx+8/UswDHwjiOH+mryrX4PI5vb3Sc89w4fgDnp/HsFLx/" +
	"h98/GqvzOTjPbuc43Rc01b8RR8s1cVyHd/FcxzkSdbzX23Fq8ezgeY4ztUO89W891+HzdyvmqnMBt9GN+4nqt55fEGPIddRiTKnO" +
	"cfzo3DjuhvsOVBfGcCCevf29quuQzqqu19GHSpyLAYBGXP9xdhxf117V8znBBH1qj+d/d1X5MxtVOZ9hRt/ejHendXUcGtOHOOg5" +
	"XW9Rw4y/43MFPy/g8/DaOJ5V4jjXrI3jRoYF9e2aYsfpCfg7P/Pv444Khn/h8ddxXc9wu5Vzc9f1Y0f13Q0Er46qDz15XvaZ+/P7" +
	"sQDf3ojzMBzVGFCI80k4MLy4BfdtNK474zyfn9VwXytytLMSExii3K8iBfsGHOO2qO/uR0Uj4m2/fRlzcfTsn99v829G48+jqW5u" +
	"chsLG/976DpM6IP338Trljjbtu/1+d/LY/9fHy9O//k4fGYe7v+R49DtzMFg/99vw8vx7oH63O1f9B9u/7L6/3s++p/8m5Ew7pWQ" +
	"Bc3/m3Cj9N+AzdqE/rchnhr/PD3n3z2q6/71/k/sovAly/pPGgfAENO5FMdmXJfy88m4prKFRtlCoywdX+EaqpCzhc9Udi7rTBv5" +
	"TG02zP5l4mP+L/+X/8v/5f/yf/m//F/+L//3S/w7dZDjXBk5zioo5Pvs4zjXvhPHJ0MhL4Ki/llNHPfczXGO2RDHzzTGcf9NcXzS" +
	"0jjeHbbKACjfH+H+inaO82u8O6ql4/z9kzgeuTqOryx2nM+h2PeA8r4r3r3Q03GuKHecCxbH8SnkR23rOLd8GMcjWjlOC7Q1cyba" +
	"QZ1D0Ic7cX4FxsjwXo7zIc7vw95+4Ic47orvDsRxNPlBcDyHPh+EPo7CdS+yD3D9Nb7vQ75tHE/ieh2ejU05zv247ohnt+H8Kuo/" +
	"ANdf4t3huD4Wzx7F9TGh41AbE9/HGOuUHbJXxnG+6O44JbjZ+cc4LoNRshbjfgUG5QOosHZZHNfAMDkL9shTGH/52Dju7ZF/23G6" +
	"kE8X59ZooxXqnrMnbJkv43gL4HAo2ty1zHEApq39PRnlSt+OYzQRj0R/D0H9q/EyRLvdcH8z3q061HGGzonj47+I4+t6OM7v0Lce" +
	"MG7XATYrAeOWlbCF0IeOVY7z/IQ4/gCwn4D7y/D9ax1wvIc5goH15v6O8wxgdyq+nYB2vnozjr9C+1eij/ujX6Px7gv0qQzP+uP+" +
	"LVwDBZxdUM8wjO0mjGcp5m/C8jj+Pb4J8ewOwGAy7LAtKNMe39V1cpwFmN+T0c8ReN8d779Hv9vjm28b4vh2fNcO5Zrh+CPeH4A6" +
	"d8H1tTgfjramofyxeH4lnt2E+4NwXIXj94QPgP+l6NfJuP8QZT7GfR2ev4fzoLfi+Jgi4GGk6vjye+Af+rQ7yrbGM8K9+fhmP7TT" +
	"B3U04P4OHLNQZjfg4f2A6dcVKu4xDnXOARw3wvhcAtiO+DSOBxY4zvnoeyng+TLoYAXKb/4ujh9H2atR51rU8R7a74bzU4BRi3Fx" +
	"XIDy3wCffkt+ddR7Ns5A97gXyh+F71qhPxGuX3VVnGQnssPx7TjMT3N6Tz4fvP8trsfjfDyOaaizL55fjeNwPO+B/gPUzo44LsZx" +
	"OQ6QWnwfyh6B48/T4nj4bmrc5xH9ARZ/WBPHe30cx13Q3x1R3yHAtYXoy/Po+zeOwoX3aX4Jrnj+PvnmxwCHKL6wF2AOWv8D2r4G" +
	"931RXzXOv6I5wvEZ2r7AVXT2Do7f4fs02u+KfhyG6ymoewvO3+O8F8qtwfd90NYQHPeiPz6elYD3DJkcx5/i+xPwfhT5JjDOs/Bu" +
	"ML5rhfnthfNc9G8/9Pl2jK016twVc1KGcmMA5H3Ad/quiuOpn4PvYHzTiVdhDi8FLQT49r+Ai9OzKtbzDeq93lf+CDq/TnQYKl7U" +
	"GziwB/rxMJ4/hPvlOMbhm7XkpyAfxKQ4Jpqfin58vrPypRyJsd4IWiwAnIfgfnfQ6PH90HfQQcWiON4JdX2Btq9Cn4/C+ROUmYrj" +
	"VtDkWNTdH2P+HLy4K/jzhNdxbq5o6zAwu39gPLuA9/wWE76G+ABg9jLeTQfv+vRg8KWpcbwIiNMTzy4BzZ8FvjUBjKbjujh+6Ns4" +
	"DvG8LdrdBe3uATh/hnbJx9IMsHkUPP1uvHsbZf7hKl8N8d0H8QwgcXp4ihYHOIoGqZ6bUe5CnK+jGBO9wzX53AI8n4brS3H9T8CV" +
	"aIfwbxqeg606GwCDJxyFQx+5im++ieuXCMdwvQeO5oBDS5zn4fgW5T8l/u2quA9YqDMJB+HasTiqQWAdgJfrQ8WTiL/MQUee9BR+" +
	"DlsSxwd0dpwOgNc4wO1oV42NfFHPgVbx2DkLNzt5aj7u5T697yo+SM+uwvVFOA/CuVOg6GtUo4rh7YtnL7rKZ/VXnGesj+P7Uelg" +
	"vDuX/Kpo82PCUzCYUe1UQG0CrmdiHkfg9m+eGg+N5TfkL8uoOdwZx4DX4riBcBTP78PYMih7X6TK/QHXg3B9LiZoLa7Hod4zcT8Y" +
	"5SbifgzF7b6J43YdFH/7HeRJLfr4LOp7DvIoQJkWePc6aGUQ8PUEMMxnQOs74/vLIRR3CxW8yoGPJGtpng/D8QbGf4GjZMw5aLMc" +
	"9fVE2ZU4h6GSbX+imBze3Ybj4kjREcUAyz3FowgfrsW8XQJEO31BHN+KNp8kGguUfITYdFIA7l4twKdRx4sEP9SdAW9/Gvg7BnJ4" +
	"0Gsq/kp0t4ZiYMSTcHMZxRcB2798DVrE9VA8AxuPhxLMSNbjfhb6MRRtHYdzxdb4ldIVwIKdV/FsLxzr8P4ynN1I8WeP2gc8CF+I" +
	"Vx2Mxp9DH5ejTDuUPcFVzx8ADFfh2dHo5xAcK3H9BWjwVVxfj+tmOE+BbtW3GeYXukF7T+Hf62in+67g7+Bhw0MlAzYBLkfiGiqa" +
	"cyfpRyi7BmUrcP0azvPxfDauHwqUA38QeMEZ+wH/IRO/xfPRpOtg7KPaO867KHIMyT+0/zH6XVXQ5MN9GrzmSTCDAZj3wdC9pkKR" +
	"qQf/OBVlLukPXou+vYQyz7UBb9gB/B68BiwjPho4kAK/27A2jhdAGK3Cw0NXxPF88JdG1FEL2rsf9ydCli4A72oNQN70bhzviL7d" +
	"DWF2Hvpzz0bco+1m6MyN6H8nwOCN6ji+FvNXT31G2X7kDw7UWCegD4+Dl5KeWIZ+dXIUT3ZdFWMm+UN0couj9Auaky+ILgHHY9Df" +
	"WajnHU/R0JE035QbgPdP+grOJAf29hWfo28fxHVbzPul+H5HjKcSCNAKePsoeMdD6PfZqGtfwOnbbo7z7AzABHh9FybvJHw7G+86" +
	"4N33xIsA6EV4tgQ4vBr624G4XkXxbLRTheu3WQ8YQPIa9zW4noLvdkZ/x+N8MMX0cQzH8Tz6UoQyQ4AfvZk2l7J+lMW7+8CMRqHP" +
	"C9MKb4ivVuO4B0c56g9xfIfrd3Ccg++mvQpZSbhJPBn9fIH87cWKh2/A+ytJH8uo648Bw1rcH4dxNQInJjCPfoh0E0/R1n14doen" +
	"+DrJjTK8L8Lxe0/poVeinlrUt7xA4cyvSRbh/VC83xPvziA9Cee7cTxG/BXfXDta5UdQfsPegZpvigmsxXhPDZXsaSSbAGP+Ae8+" +
	"Qf07oFxb8Dp8vlW/bYmLz8AHz8fEf4cxk87Tj/UXyifoh+8oL6Ef62Ckg/6N4n84duYyf3GV3CL8IRjOwpEhXHMUP/E4lnEt5bDQ" +
	"3DJfIfuE9J5a1t+gkjvTcYzE9RRH3b+C4wpH8W7iY/3xbiQqnMZ61oeO4r33R0peLGHcftxTdgzxMsKTYxyFlzDL4hc5tvckwQ3H" +
	"X3F0p34ECmZvsiwn/vyZr+ae6O4z4PoEAOFZvH8O9+Nd1ddf4TiL8Bp2D9llRxAd4vwp6usE3FkG2rgVNLoRZU4CfdwLOJMcLsH9" +
	"ax7l30AvAlCWATgXkVwFD5xMujLm6XHIzpWwF+ZQ+xAC/SC/FqCuZ6FvPwq9bgporhvm9xJH0dffcBQDFqfhfAjaJ57Uz1O68MHA" +
	"ibkko3D0x/WXkCUnQmaMJRgHSm7/DnX9meNMRej734nHoPxdgMMg9OdgTOyhQLI/Qbk9Hs9vwPtbVynet4D0bgDwY9yTrDuV8BHH" +
	"G7hejD4Mof6h/G64v5/0IRyHQci8MR/lQwWzgSjTETr/E1UKx5eg/CE434u2PgYP/wD33YFMrdD33cF3KqAb34lvfwSPLwWvPBqV" +
	"9MC3jdAvT8EgVkPAPwyZfg6Q4cxXwGfQ5v4E9xKlv66k/Bwcd2FsM0m3wbMi8OZe4PW/ASHcPg+8exbGC0IjfexZjL8qpeaeZPhA" +
	"svvQpzWwMbuj7XdR/mC0ddkU2JhgRE+T/oW5ug/6xOvgcXehfA+MfzTBGd89j+sZhF+Ec6h7Meo9kfRVjHkFyWIaA8vSUTg/hLb/" +
	"9olK7JmP+8MWqVymc5k23/IUDyO9Zi3Om4n2PCUXiVeQDfsbV9l8Y11lK5GuRjRCcrs7jl6QVcdkld6xE5ebjH62JNokOxh9CCKF" +
	"U+QDIT0B4HP+jOvTWfY8iTIfMb0Tj6B8szNxkAwmGqvGWPugrYsh29tgDETnxD8m4X0qq2hzEOb1DeILqHcew6oalV0fKrlPMoL0" +
	"CuIPM9EfiGHndTx7Hs+yaPhVyJXeYEJ3w845FLTyBeZn+hzwTDCF+wHDoSD+FuDXq1FfxDYY2TXfOQrvaoBfNxOdo+6ewKmxwPkz" +
	"ADCS5feCRq/Ds+fADz6ALbEMivyXsDGqSceDvnA+xnA2YNCI8ruBZmeiP5+nlFwgHGyB+5dw/RTNN+pbRPIUesB3ePdrXLfE+x0z" +
	"avxV7ZXMIJvvTeDRHNR3oa947gGgtb+i/6+i7kmRknkn4jgN95scpQ/vi+s9CpXuMwD4dBXGPJlkEvFNvBuB4xWUuxhwGYJJno1+" +
	"bIAB+Aje98Y7Or9EvhnobBWAIenf/TG+OwvV/NK8F+DZja6yPYnv/4n1b9Lb23mK/sl28EiOuApnTkK/vwaidpkAuiC7Gs/rya5F" +
	"vX9B+Q/xfhjGMhXPTsD1JLT5DPo4jnBlYhwfhGf/QLlyIN/t6MgMlB0GJF8EeuuIcZwMvnIiEH5USo2pL+azB2B6c5nC2z3xfBHm" +
	"cArgsRjvDoiUrkO+CJLLhINkEz6CcpV49ynqG7mL0u+H437+G3FMPpXO4Ee/IT5FMovoE7pam/ZKTt8Emn+9nXrnkh0NvtAHzx/B" +
	"dZtQyezJxH9Q7zVo+5g3lS3wDttkhNNPs/whfR2qY7yUxo+yx7Hv4wqWvRmWa3Sm/pOMpzJ3sy5B9sDxLFOpHrJvQ9bzSZ6nyV+H" +
	"479czt9k+4/kNvk7aA4himKSy6czvZyC4ymSreRXdJR9ST6ZB/D8XvKrot1f4xjLPpBmuH7AUXhBOgDlOgJUzmpql+R5oHwsLVyF" +
	"dzd4auwXjVO6TiXbQlTntwSTUNmUD7GuUekrPYt0JpLJpA8RfV+EYwLpnHh+LPDoTQDn4EjplWWT4pj4ThYDW/MxbOdI+euOx/m3" +
	"wMXOnvKtXYK53QKmeNlqJVtuBQ9pAE2QLgB2EA9GGy2pHfYXzCZegOsF5NPzlJ5BNjjJkTm+4lsHpJRv71zSI1hHPxLvXg4U7Os8" +
	"ZZjTHP4pUPB9EMcL5O9CW68Bd8kfQb7DiyIlu9vjugT43xm4VwZZ9Db6+lao/IFZ8p+B3z4P+bcncP4k6MX3QIaOBb1fCGW/0wfg" +
	"lfj+SaIFfHM7vj0O7c6KVHvH4t3deDYzpeiA/CNp5gHEW88Foz8Qcph4fQovdiE8xLN+QMr7AK/vUO5z0l/BNwdDuX4Ok7ofYHjT" +
	"0jiugQ14NRSJdrivhBHTF+3NrFb20kXs716LfvyT6B7Xy8gnhfp2dZSOdizr210wFwTDJzzFP15iP84ifPcQnj8C/lnUR9kIHfGc" +
	"8q1JltP8XUo+1FD5ePqA39+IPm5Gm2MhtC4C7Y4mv6mr7PgL2Q9yFuj/JE8RH8ncT3B+AvxovKNwmmj0MBDOUMiJmaVqXslffxzG" +
	"fj8I/DiMkXSvlcXswwMsHiB9DP24G21eyPhH9HIgYHqao+ySKa7yoXdk3+JZLIvJz38V+1ffBn8l/r7OVXbRIMB4TDPlS+mFb/7I" +
	"/oCZpDM7Kqd8GMkW1HE63u/gKhuddIzd2FdB8YSdImXvER8dA1xrRLldwB/bQg5+jm+Hkn8COuzV5GcjWsIzD/0Y4qu6TsWz4cCJ" +
	"cVw38T7ikeT/IduOzhQ7mMT+zRTg0g1jeRcwmgI9jmy445lP7QDcfTxUdiLZFKvBVAbBtib+8AZk64dFyv4gGBE/hbkTE3wIf9uQ" +
	"XCM/C8qR7UW+mhao63OUfQr6w+n4oATv/0o+ZPafke+E5OePHvsLXRWPITtgKp69QDa8p3wrpOvdiKMG9xMAr49SyodxEHTSM8n3" +
	"ifsWRG/vKh2IbBfijfWFylauRF2fAtY3t1f+rfEodzPe/QEdDQHbaZBf46GrDgARzP8IeA9ZeP8mlRf2KNnraQWXAughV76DclkF" +
	"i4NRYAEQatXUOCb79G60eTj7lS/A+C8BfXYA/1gBvekGX/H3nWDfvIb6NpOPC5PzOOBzNcof5yrethPOy8ivST7hUPH06jIlb8g2" +
	"ez1Uc3mIr2zHU8hvl1YKYhp2SBHe706+dtD+GZjYd8nHgvvFn8fxYQSXRUqn3Z9iM5HSWYn2p2cUz+yHujeCTjtijCOgA+zbGjIK" +
	"Otlo0HEZCq77ROnKpwIvi/dX6zbeBL+fgjn+BGNZBbvqK9BGa8D1VNhkE0GHRwNuGzHhF2O+RkBOtEd//gl4nY5+nAZ8mcN+3vPI" +
	"x+Apv8Bg0uEAS89XOk8Z+YQZp0P0eyqef4hx34L6SN//M55dGCnYLSd+hmfXYUx3YyyngjdPRP8Oh12xAnUexr7gpWhvJMp0C1S8" +
	"4CbiewRfPLsAR+VC2Bzk88CzE8guw7MjyCezTMmrp9D+C7BrTsOYV30NfCafNfBqIOqbQ7TpK71/E+6vB30X4tkHuK9AH54mfzL6" +
	"MBDzUoJzQP5bEMIfSU9DvYNJN/PVmof3yMeBMh09xb8GOUrv7MjxNNI79vOUT/8eX63VIJxYzPRFPJXij5PZ3m/EeEg+UV3Pecr2" +
	"f4b7fBuOReQDBJ7eCyQfBRzvCN56Mb45CP2uAV9Ygb4Q7b6KMf+Z5Eqp0kv6QkaNQr37A26/Au7/GrT0LMb0KnBlIL45BLb/MtDg" +
	"i8CnRZiP86FY3InrxzBHK0B/FPupBb3cgXHsSv0GDjeC185fAD60L3RTwH7Ly+BTeFXsKlvqM5xP3gv2Cnj892j7Ml/x2stZLyLf" +
	"MsUHhpGd5il4ku5EthDBaiDmK8V601cUG6M4Jca52FH8d7qrdMdLWT7uBj44i3AZ/OFWngfiwXfjmwfx7SCM5Tx8t2w25CFwfwro" +
	"aAWeP9cDuizpjuhDf/Cf9ZGCGcWqD8CzswGnkYDjHeChj4EfdSObDPUOhZw9ey5seeDYOpQlH+60jkrHfZ/8M4HyVRJvGEN2g6/o" +
	"YzLpxDifBDw4G30uJr+Vr+jjaJT7ZrHKZaY1SfWA6xWYB4qTkn7jYL4PcVWsoSP7skp9xQO3rj9ifDsDZe6nefeVTUp2H5WZyf4m" +
	"8nltRtsuYEX6Vk/wCoqVk21+G2yuG8GfHvWV3Cxg2/oG5oPvQhnsB30+chU+9wiU7/IA9jtuYP0wi+ckB8lm9XF+Fri2EbChGP7u" +
	"5O9E/YUo8zDgtgnPrgE87s8qOBHd0biOKlQ61sng4zv8GnYeZMwJhSpfehn6PydQMQTyZXYl/xfsx1aY74mYp+Ogr5F/8UHg7Qgg" +
	"2mvQayuhDL2IejYC32dj/m4A3/xqrsp32ES8hPIEAhWP/wz08IOvdMH9IFQvI//QTo6zN/D5ZBhGC2F3fwTZTHQ1up/yFx+Gejuh" +
	"jRlQRqpAg4Mhx45erHByPvqTAVyHQ2cNMNYrQVe9QY8vF6k8hW9RtiV45nQYAONRRzXxPfRxKWBx549KL2mbVnr+aNj2RdCN7wK8" +
	"DsT794DDQ9A2yYpXMVc3ot5bMeY3MeabMHmTVsbxta6K6xMv2hFtfkV+3xlx/AZw+yjA7R58+yHKkn+iFXD+MJQhe5/42O041uN+" +
	"CXhuEeTZ0bOUDke2w2zIyuo5Ku5HsXfSJxtJhqIPvTHWPUgnwrgHYu4+h37TKlD2EtFBC/Kboa0RwA1a/0b+ywtx7oLn70GerwTs" +
	"H/ZU26tJP6Z78hOAwewNebcacL8V/R6Nfp0LGJ8D3drHGA4AjkTToNP5Kt5+cYmK59Sh7dfIx8W2KelVg4i3YMzrMF+kfy6luCa+" +
	"vxvvKiuVr+R4yO8xkJUTMd8UB5qOZxOhrz0GuF2BudoJhD55g8phL4Cs3QQc60F+f/S3LeT0MSTLxkDf6qPk/GtggMeBDw/A+3tg" +
	"Z1Pc8zYoDd+ggsvHY+49Wmem4oBb/VToXzfA9RxX6R3vBkrPqUUfKeeAdHbKHViF/q8geYXjPZRvjf55GMOTkK/PdVPyfNBUlYfj" +
	"lqjcC4rjXVyobK59ANMdON5NPvYgVPDbjeIirtIfyf+xBP0bECn/29K0smEp1kH+zfW+8rmRDh7hXTH5yEm+kp8rUvZ3H+BdL9DT" +
	"GeAlT4PB/JH4VEr5GU5NK3q4EeeXQIN3sd7XCXrLYMD5TODUFbDjdsf1EIovYE6eAg2T/jEKjHMx+lwDPWcnou+1oC+SESh7AXDk" +
	"n3j+Mb69FGMfnVW+G4pdzydaC1SAfJWr8O9cGBhkT9dAxyvCOcQc/gZAB/ji51Hvb0gPQnsjMV8voC6yq8kH9gae39tGraU8D2Nr" +
	"Bpn0DOqtIn2W/BAZFQ/Y21V28fcbVS7Ww46aD7JLXF/BmuIEZKyRvVINePSIlNyM2XYYwL7NO8l3m1by9PfApwaM8y2UHY+5PxD1" +
	"ZSh3BX1pCZoYDFo6Fvxp5grldx8PHfF0tLcXnv8NtPMR6rkO8F6Odu9MqfjTheBhnTBmyi05C2OtAm0/AiGyI2jiEtA16QR7esrG" +
	"X5tWOQ4vAnGHgG+17qxs130Jr0C3D0JfnQXcfwtwOYUXIk1DXztHCqdJHr2JOs4ET64Cnz2iVMmtFri/E7RUC3h/SPleUJJoDsgG" +
	"uwY48Qf0/1LI6I+BU/T9fmjrmozSzWaWqfgOyb6LIOOz5MtFHd8Ari/COO6DOvcCYn+BtloCP8hH1wx1vI/vX4GetQf6tjPgNQo4" +
	"NQa48n5a8bCn8D3pzKRnPY/69gL90JoWin3t7ypaJF8P4SnFNUk/34t1ntaB8vXNYjk6mOyrUMnOFoHy320JVQzlK/abk5+70lM+" +
	"b6q3A86noGxPtllr2O+yCLoY6QgrtUwOVKyNbGfKWSCwn+upHCviH2PZd/kexTIh40b4KrZ2AsZ/u6v0jHfRlzHk350Qx4UDMDbI" +
	"8ncgC8iPdgTeZQuVHdAPeEaLoSkHi/SlOZCTczEXS4BbV0IOHYF+psHvv8F8TQcfXApesAQy5EnoHcOBt38Arf4W7ZQCj4nXXYB6" +
	"L+Y8tdbQ76tR5irw2oHAqy2Qec37A76A+xCcC3A+dKnyNyyhXC3g4EX47hryXaH8Msz9e+Czg2HvHDIadIe5bVek9McfQcunT4/j" +
	"14H7rVHXP9+Hbgccq4KO2wrzdhzgsjN0iZuJpkCzGzCOe4APG9DfG8BnSzDuN2CXHYrncZXiyVeh3SPAw/bCs+dBe5eAhh5GXXdA" +
	"Vp3H/ol7gJszgD8tWqq43imgrxroGxcDNueBt90CWruqSvE0ysu4DOVOAG68A9vpOdS/AOWug+yi2ORVpco2m8X8/xyyzUjHAlxq" +
	"ML5XQE/N0edLoEeOAj5HYB7toXNUoeyBIJBd8P2j4Le/xViOBQ+7FHVfi/keivminDTKPWmH9usBs4fQ7+NBwxNQ5gmMeyO+ua2t" +
	"4kM3o82TKHaM61Gwo97C86cxF+PYRn0FZR8Gbv0VYzqc4ngYS6dalcMwylU5RRSX+8ZTOsZCR9FQP44n3YH+diU+HPLaf/xrj3pu" +
	"LFExllfxfCD7bafj+mOM9xbg2cWYN/p+LHDxL/uo/K4Z1EfQ/3SK82Oc1B4tYOsUqj0EiAeT3/8CT/FZ6gPJyrm4H4nrm1Huu0DF" +
	"b8lfRXlO4zmOvYx0V6LVSPltLwFsKZfj8bSKU1Es6YRyxes/AP5PBkx3TSvcGQaZdgTm4Dbi5ehvF8zxwZif71FXL/Q1hXkh+fwd" +
	"4DgG9z+gzYXNVYzqSvTrOrS1lPwXqG9ZVsVPhxYoe4rip6vKVVxrBd5Xop4HaL0dZEF9hdJLrke9g+epPnWBfje+lcphHRaqnCrK" +
	"K6B4+wbgw3GQD6+AhncE3xlfpvxhDwJ4PdAHkknEnwaCV7zgq5jxCxll/+6Nus7C0RnPz+f8MpLHt9DztMp7IL75j0DlC10PnrMb" +
	"6PAkyGby45EjiHLoQvaHj6U4B44rKDeQ8h3xbCIY8gLoXCdQPi/wozdkyE3jFH88lvtGMfoMYHMmrh9Anbfi2wfRxzNQ/vpA2ebk" +
	"k7wjpWLzA2ADfEZ9Aq8h3bhrO9QBGk1hjq5DPw/CcV+RshU3YryUZ0oxn+GUE4x2Muh75UToQEDSs8FzxqG9PTHXvSgugnavews8" +
	"gmJXFIdwVJyUaJns5q995aM4xlcxBJK3zQI1r+RLo3l7FDy325F4BxsyhXqewUQ/BZ5zPOjsNJRdArtny3zYPTAiDwQvnYp5nk65" +
	"ZFD2n3lZ0ejl+P5kKOTHgSZoT4mHAMOriY9iPO+gzeVX/89Zk/npUPCUNWpviuZKpdoaG4pUet7WoxnHl3rw/VIoLO0hA26B8KuD" +
	"Ln/tGpUvLtdPF4vvaW4Ibpt4/Ws9r0lOcd7lOt77opxjwUCVeCeuh+LZYFFxb36f4jgXhb/JJkqxvU9rY6s434Tw7G3OSUTX4iz3" +
	"g/bWyPC5Gfd1V+7rSo4x0V4htBdHFe3jwnHwNH9PfnXKq+rEa3D35dzJt/Hh4GVqfxdaw1vOY6zktb70fSuOvfVmn0EVt7+M1/c2" +
	"43IpPtfwvigFXFczXh/sclneyiSO+JuWHG8P+V3A9ZczDNM8zgegl/7pa1VPwO8JzseRDlOr9hKJ+Jkn1sIX83U9w1zPdZrH4/I5" +
	"5GdpUT89uxR6w8hpql8+P8swnuk5pX625PeBwMNKhndKjKmG4Zvld1mGdSOvoU7x9TyGNdVfta+CTThJ9Vevry7kb7MC/6kPKegA" +
	"699QPnXPadqXp4jv+wARvpqjci/1XOhxFQr4aZhEXG4mBrpLqMY4Ap28hX1JIeeznYX39Tz/zbhf+lsd6424H0dBiEwBckWQeW3R" +
	"0XVUeKHKbw25DI1pH+7/yrhp7OsYLygOBlYX9+Xcq6U8DzMYlyOGUyfG15DhX8zjrOHY8Hc8J6t4rlrwPM7n/YGqVNfiVjw3bbh/" +
	"Lufuh1z/WsbtIl4zcDyPoTO/pzp25usy/h6q6Fb668Fj/YFojfYaYr72PfehFY+nlPP0qd5uzD86sg2/gvu4hPlSWrkXt+bmEd9p" +
	"z/2iv3b87bE8L3dBhh67uInPReyLJPhNY56hx17LeEsw28D9r2T80XRVzOVCrr+B50zzZw4Tb4VtAcMp5vKaF/gCZ0Ku9xX8I/8r" +
	"ta9hqOcuZPxoxXkDUFXib3msWT6KOMduLY8x4nufeYXmCZQn8QnPZXd+F7B+MZrxhuaRcvaORgX3AkEJ/vswbvVnfCzkPlFscA7v" +
	"eVTDtK3b03xL484mhlUd8wS9z5Tc30PSJr2/DY0Mr27iNZoH0TESCHL9XDVPjdynRp7DRsalkOeohOeR5nkiGBTl127OcNxQhQm2" +
	"zvkrOI7gtg7nuV7HtL+K8ZvwsyePW9NvT+53H67L4znZxHjUjesv4TmhurpwO2VO075b1L+/4/okIas0rtVx+5oXu2Kfi5TAPzqX" +
	"kp9wo5oDn2m1nGHKKdnbyDJNH1qeRALWDcx3PGMOigQelwr5E4pzA9NlJ4IDl6ec7GbFql8Vkr+jMPnSS3msEZ+LmV4zPIZSnrMG" +
	"3qON4qtDmV9kuQ3KwZvE46lhmitmma3/+nPf6xl/Owo4lnH9hEcduB8lXDedv4iV/1Hz9EDIBD0/paK+Sv4uZp2E+tlayK+Y+1DM" +
	"eNrA+KDndBHzWK0nHADmN+ld1eeO/DxiuMc8h0WMvw2MoynxzhVyX/OxQIxDz/EQ/ma0kKcVXK6A+5oSOLWA5VIodM8043LE5QrF" +
	"eHXMRuuBVQzrmHmFxiPN/7oyf2rGOEH09zXl+nK7ZP+8wbptS5YZhPuUL1UrdMplrNN2YNm0I8OL5MgKQQeVLC+1TtuD+066bnce" +
	"C9Hq37jvy8UY9F53zXn9G+WAzeLvCG8JVvvEak1ZhnXjzhwn/IjlbjXXk2H+EfE18f7ZjNttGNZ7Mr9ZzvDrzLpYA9eT5m/1nPfm" +
	"+Vgg9vMp4zJ1XDeNbzXTQQ3DuAvjy1qWEz9yv6oZZq0Yt9rx3MaMuzOYF/o8Fmrrn22BYwvVXDXwPjoB422B+N4TuKRpQONIStg2" +
	"GYGLUlcuYnoMmQdoHtVM6LwbhR1Cde3HMrGM15zVMmybMZ+fQvSKRjY0qGfVDCdd/5c87gNZz/I51v4dw+sHPp/MMC7hvHtaUzKQ" +
	"+7mO8ZR0lV6sL2o8pmMct3ks20azuM62LKczPF+kWw3mOZnJ8riCYfIDja1RjZP6dIKwfb7B/fc4n8/tlDN85zH/bMf9XM9zRPPa" +
	"lumgN9N0JzyYtVDh2DTAandP0Q3pXlO53kahRxEetOZ2lrF+Fwv8dYVeQecKDKJhjeKdmr8WCH7cnK+z/G0KwrhstuJ3heLwGe99" +
	"oYNkhC7h6vUTFKduVPtnUluX05q8xiabi+ihH+OTxre2gvdqu2WTkEsVgufWMZ0042ddnKa9qDYKWG1gHVDz7NNw/4yWpYLOWwg7" +
	"VtNPGRC4Zrqa800Mr01MayX8jcc8VsOtK+PHniz/VjM/iHm+qM/E50hnpzy+yZzPN4Xz3dFc3JPpSOvEhYzbzXmMxGfHki+4QfHJ" +
	"RcxPi5iXE687h9YUM21lha3egWmwHfczI+aD1lAW4OYt1nG1TbE1V5F5GcUwP8aZ6P5znA9j2L7FNDiR56+C6fUMtmOIRr9gXnkA" +
	"0/0PDMc9GY7Voj9P43iNr/WY1rNM6sFlCSZ7c3/aMF33FfpUF6aNQ1h/6cTyw2MdaQXjSyt0bmGR6vNcxhuaI8qnP4v7WM97wbZn" +
	"n1KWeXhLLks55WsY3zcyfF/EeRjP6Q7cv/Y8D6NpjQHLnL4syyg3+ANeAzed8WUT86GFXJZ4xgVCDizmtlzG9Wvw0VOL1Rhilkd1" +
	"3PddGIbFjE99mFbLeDxl+HbxYmU/0py2RoHV05XcobootnsqGhkLJnviTCX/pS+kjudey5YqoXNJeZQVay0ywlbRZXzBt7SemBU6" +
	"c6HQB/V+dNMxuMM+VTTqij7FjCMajmmm2wKu4wfG8UUM2wLO8W5k3t2LYb+cy61m+3oG86Uibr8j09hyloXko9sZ8Fy2WOFDOcNn" +
	"AeOq9vPUsq+gD9dZwzhC+tSuPI6jAIjHAaxL2Pah+l7gOtuzrCpmm2sd85w23Jf1XM88Htta5iMB64FE5++xPnci21zzmOYO4XeH" +
	"M+7OZJraxHz3a/ZVtGNc1TlgqxmnZT/WcbuaH9Icvckyui3biU8wT6FY+mONildE3MfXeK3lOziTrvBIrPKsv2J+ehfjv9aTad7X" +
	"sO5L+fZXcvsE7++Z1/qc//9UrGxdrZv47H8pYvraVejdP7INoHWxVcxjPeYtVVxPMd9rmXZSrHIZNY2UsUzbwPNdznrPFv7mAOap" +
	"hULW9WXdRvtq+4jvG9gmr2a5Opf18+XCX9OZ11CPYBzVeuROjOMh47CWdbXcVqGwMVPcViT85ZqWpB/XE3qIL/iAtgXbCLtc+yEL" +
	"hO0TMky2cF994bcImG9pf00Lvm4BJEp3VjxO+5BTzIfqmNbq2ebRtvgIDH74l4qXNgp7So+R+jaIZUVr4WvqyuMjeH/AcuL3rGt+" +
	"wXkBabZfl7N9XMtytrXwca7nPmo/htY5Mzxf2ifWlvBZ+Jm1jzYQcEgJP5LmnymhB2r/XIHQ3fQ+nxkhN7VsbMY2ZZ3YK137RLTP" +
	"2mWYVDpN+6WXC1+fbrM9855S4cOtY1zZyHyiQdgqWX7fgdvQcQIdG9ks5t4TsYs044yM65Ty+wKGs9ZBK/l5xpBTsr1A+Cv1WcYD" +
	"fDFOjTPaL6nLyriD9gfoPVQlXOn6IAD8k+qmOWxg+1Hjr9axNwqbtoRxRtNjNcOoFdeh7WvdvxIev5wnbVNo2BQKe0D7V1Lc/2KG" +
	"acTrCcguorU6rfHdoWxT6XhHO57vKqZnre93EtcZphHNu/X7OpZfnZknVQh/aJrHrG2XDobvVPs7tf+pgMvTeJ4CQE+b3eSDigQc" +
	"C0X8ZyPTUFdRh6SzM2K1DkLDSMdlmon4oC7fH5V8MVuNT/uF9Dxo3lbB8vxHMU96/j7AIAcsV2Wlb7xR4Jwv8E36O8uZr9Qx/+vl" +
	"NO1PHDJ/1HphK9Zvb+W5fIdtgS0sb6pZH14v6KaO4dRbyDdtqzey3aDlxmfsJ99FxAt97k8596WQ9Qgdv4x534DHBK+cz7yjA9ez" +
	"jMelaXMjw38D11vC+FHGfVwr9jjWduA67rPWHbVPpBnfbxQ+uJTARemnk7yhgXmMjpnQPJ0IghizoCm+4QmZV886prYh9uZvFvFY" +
	"i7jtMtZLF7Ft3YWfR+yn0r6XGtbZ+vN8tGP82hA36aSaTuk9rYE8h+VHL7bbfuD2uvOYUjyetgzPNdx+yDqRnu/23OZyxo0qttE7" +
	"sM5Wy2PKijhSmtutYTmSFj5CLT/kb3DQ+8G0Jw7pcZR3tboJJ31BH76QBZt5fC7vDeGJGJS5l3og9JmU4OGaR5YJHUbGsbWdkhU8" +
	"PCP8J9qPlua23kAHhmxoim8XCDmuZYzGn0VA6h4/KN6W4roiw1etdZ56Ybtqfu6JOiPRluYjpn6hdYdKhp227RdBUW/TW9G8nsOA" +
	"7SktV/uxLNB5GDouWMj6H+XbvrVc4XBzhinZNZcyPDfx3Bcz/er+1HCd7RiOmk++zbheKWg6I2KRvsCBEiGfta1czTQRCr5dIWDr" +
	"i2tPPPNFLNETvLaIaULjRIOIccrYn8txmAPFXEZizC0Zvm2YBnYVOs964V+i/lLe8R+ZL2jfSXe21XqxfXUKl1/IdcZi3/Yqhuly" +
	"lhU6djeMbSLtX9P+6TtojzD+1mceuifLDY9tC+2fPILpmPyli9j3GTE/6Me6827cxrc8F625XRrLKJYDtM59Y6Oy/yZwOdqbsAfD" +
	"YwP7UCYyLyScJPtyKcM4y+sqaU3NQ2x/RyxT9uJxreY+1zA/07Eq8ptcw/bVRtZNdH7LerbfY9bzB7AOv5Fxi2A0kfIKWR7sxDxx" +
	"MfOHLcyjJ7PekuK17KOYPr7jOS5nn5m25XU8ReaxbGIerWMV2v9cL+hvM9NqCdNYmpxQuyueoPlXJ6HruSIO0oP725xxcTnbVZSX" +
	"eEujihVpP4/PsNFx1M4MY02LvujLZuYjK5m/9uA69PPlfKaxH8y425L9iR77QBbxtxHbMbUiv+avDM8GxsFeDIt13I+tYwVAvqpW" +
	"9a7hPsc83/15PKtFXtRapjeX+e4H7KsIDXusVPDe1Vxvxsj/0TZHhvlYA9tAs1he7ME4X8p0Ucc0vIJxaAXjQqHIlykXfnmpH7YR" +
	"/HEz002W8SstYo59BO6Fwr7xRezeYxqYwTTWKOzNjMH7pO89FL6CUOSRNAh55It8Ns0vW4nrrIgvaz1f2/olIiZWxu9bcb81DvYX" +
	"ukUN03QDv+vCeF4t4l4beEwdhV7YguvXMYcGnvfOnEvi8Xz9wO/68txVsO9/V+GvyPAYLqe9RdYruFUIGisQPtI6oSfEjPcb+JnM" +
	"WfNFfldgkWmPifwzMz6ofTZZUYfWm+uFbNVy60fGyRvB2O/u0eSr3ijivWnh3+ksci07Czs0I3IpGgQ+bY074uPJq+L4TY4DepyH" +
	"GYm+F7N/YEeRt9Fc5HTUMH5UM693hR7TUvgKaoSPXdsInZmnaR9GpcAZrUNt5vnQ7U1l2biO/cX1TG+h0FW1HI9FLETPo4xD6TZX" +
	"8FzL+K62a1oKvaZC+PFKua9al9f4Uc+8rBXPd1r4AhoYfsV83WDEyqXvZA8Bg6zAl87C769zFsg3O0zgreah2ibScmAz8982wi6m" +
	"Mu+y/Ja/ifM2206aLj3hy9A68y60tn2LgkWd+L2/rPBHthT+s4BtE7o+G0zuvp3V/gylIpe1kWFXKOwKT8BBy6pG1gU6iFh8nfhN" +
	"QRkPSQsajMS+ObWMQ9rX1Z1lc3PG+WU8X0dymdU8xtasb9Dz53luQ9bNasW8vst6uPbL7MKydQXH/JZzzEDrFhnhZ42YNwZCtv7I" +
	"+Vr13N8f2K/VVuRU6FwlHXdcz3KpDfvGdUy0mseS4T7pXNqVDF/tX18n9CJP2EHvoqGBXtPziOMPnYRev5Xu3G3tL5kvkRJ2g/TV" +
	"NIqYbc/Zcfxp16a5zwi/bCRkfQPrC9KWqxB+pXqOCS/hb6Vvt1DgQoHgQRV8nRX2qpnz0Vzkb9G7awCwu4pUTOQjtusLxDcFwhfc" +
	"KOL0BSKGlxHyReZk1wrbMTDky+/fiuOHBzTlZhUynwuFTy0U9RaKvOkG4YcuMWyzMj7ORMNPApFWweCpaqvg14rH3pPxRvtq9Bi0" +
	"j7aD4cPVcnd3gVchyxqdt9pZ5BRImg1Eva7IqysUOXPSd6d91nXCx6f59QEiLy4tcng683c6F3KL8BNrGZ0Rsddilm9zeEw63lMs" +
	"cEPnfFQxrRc7Tb8X5guYSH/BgULn4KW42+Be1uiTL2Aj/dGR0OvKDf+Fa9TpCvi5wvcZi7xXV+ifbdge1fG60PDxpAWuZdkWKBLz" +
	"KfNSdf4Cvac46niRQ9lV2A9Zxula/kbHMmo5r2GVkEE6J7dKxGM0P4zEWhSq50JBI1ruL2L557MeS7ClddWXsGzoyvp6g7BlYuGX" +
	"ai7oeBW/y3D9VVy+C9s7Bezbywqds5znUMsV/TtyGebtVYzPXdkeSgm9ra2Yw+MFz5B+bVfELzQdpQW+aZtuDIB2fmpbPAmEbI2F" +
	"jmXmBdlwTOPuaBhlQ7urMRUI2kkLeq4XsahQnBsZV7Tu0V7kiN4F5KH1mPdCuBVAWF8O4pyAAv+4M/+7ef/Jv4nz8r+Vmj/yR/7I" +
	"H/kjf+SP/JE/8kf+yB/545d3jGyT9wf8T/qjvS/Jh7i2IY5pHzvy9czEdU/eV/wpXNPvSm5dG1Gfe27PFLlB18S/TDz4C/tEtW/4" +
	"kg6Oc928OL4I5x3n5XE3/5f/y//l/35pfwP9beMdFNe5v179TslrOB/Fe9qPxPXR9Psvhqwa5m+rp1DOwYec918j9ixYwnk1f6r/" +
	"5ciCSf7P1780jPr8gvpfEfzr/Z9fl5fF/z/+aH055dLQnl+0v0RLyn/BuYhipfx7yQ6vLb8TZ/q9jlV4SL9zczmvP/p32qd9ECgH" +
	"eR7H/nWOi15v3Sj2bnsJ17QHx+040+8+PUjrguL/mXhDv6Wk88caeY802ge0Hz+byHkPBzpNeRtv8b2+1jFjOtPv6ej8QopH7xCr" +
	"vbj1tzq/5jnjO9qz/gyxFqZW5IL44tu0uL5d5ADpZ/T7iMO5fVrbO9xoV/ZD5xlebvTlHkftHT2K+6PX1dP1heJ7evYojt/R72U7" +
	"TfmJ+v0Hjck48YTR744G3Ocz3LuIPs/me30t+9zXgDvtwdKd4dOB1jfwdx8Z39HzC9HQSqdpH53twf0FC9znCrg/i+P0nwH3B42+" +
	"XPQvwv0W7sflRl8nNeZ5eP4v/5f/y//9b/ubdOS2ecl630m51keu9dA5aHo9u35HOp/W9+RanwJRj973Q+6ZUO5su8dq1tk2302X" +
	"0/l95j51pj1irifxjTLN+RnlLH4h9r3zhAzX15FRt7lPtM65SxnfR0bffOMs1wi7Rp2e8Z1ZVu9jZstdtdXrGuU8MV9yrYsnYG/2" +
	"wTfgInUd14Cxbj804KLrLbDAJTDqdY3n5jz4lr7R/VRe1yfbM/fRMGERWp7LPTf0+0xCm74xH3JMgQFzc47lmupAwM23wEgfes+u" +
	"VAIOhAbtujlwS48rLb4z1ynIPe7k2maNQ/I7vWbTbM83aNjExUjwmYzz05x5me9dZMFtvb5dzlFK1FXo/DQf1+Rncu9DV/RVzm8q" +
	"geZs/Cjt/HRfGM8Yt4R34Gy7/kDiZijeuwZum3wiI8Yvx2fLfZc46CbQvmepT65JKhY0KvmiK3hAaOBlZBlbYPDOwBiXCTtf+Bp8" +
	"g57lva5P78dR5jSt44qM7yReFBu4q3OyAwtPM2m6wMIrJU8IDPlk40O+pX0J/5SBE+XCB+OL8dnWN/gW3pJymnL2ZY58aHwbGX30" +
	"DXox8VFftzTkk42vS7wrFPis15B6xvybOGFbw+ELeOn1F5GB2zb/plxv5OeQ0V4O+pblQzE/rkW+BhaeWWKBcyZB7qYSdAEJ70Dg" +
	"fmCBl15vUOj8dI2MDU56/bfcry1r6GFSvoUW+NjW3ngWXSZj6IPyrPdX8BNw0M0xXzZ+rr/VstFP0KX0/QaxTl7ibyjgIflb2sJ7" +
	"pQwILfqJa5Hr3/C+qZpf+QkwCgU/yFpkQGjwDv2+1CJ3ogR4mrSdstCkL8ZO8Nos1quGCbpwYOiaen+3coM/yX716uU4389s8qNJ" +
	"nmmTl4HBc6IEOtdy3rPAStZZYplTT+hZSTZCYNCMCVPfsL1MXpdK4EN6jGnLdwUChr5FxjazzK+m5Yyh+5r4pfW10NCb/QQ9PDTw" +
	"3bPYQKGhZ8h10lLPSFlkmV7zRvHHJcL2k3v56D32sgIX9LusmLeiHHxDrlE156PcAmvPwjPlelXPwt/lPj9ynkucbfekChLkpcnX" +
	"O1t0psBCv4FFL3AN21GPo8ACA9/Zdv2crqvEmNNaXuMn4VRk6JQtnaa9HOXaUN+wS0wa1fvuUh/0PjjmXMr9QzIJ8ilM0P0DYw5T" +
	"Aj+l/pc16KlRrLmW35SIZ/K3QzwL7et9srKi7xpuW8TeCZ4B22wOeal/S0b+Nppr6Dih4UOQ/iQdszDlTkrMOZWrMnAuEu2a+Gjq" +
	"spKWPYuPJLDMrW0/Ttex7xlj6osFhqyTuoie5y1inbOESaHFlgstsM8aviw3h77sGnpkgbPt/mKuxU6Sv0lm28vON3iS5IeR8Nm5" +
	"27FnPUOvCQy4po26fQsNmWMNLb4jiRNFhm4t58ZPqDtJh5eyM7TgouYBUhbrva303nlRgt1l0kKYoOvZZIRv+GFlnwoN32mSr9T8" +
	"jTstozxDl5X+Ts8C+2JDx5e0Flj0qEDALWPYJ3ofGtnnSgOvXcvcZg25VyrkTDrBdvAsekmZBTd9A4eT/DA2mZl1frqXvW0MnvPT" +
	"/QEjCy5IO1zObbHFl2LqWilDV/IssHAF/CNjzmy/lWjqBKGAe6GBM9uzI0yfnqSNtEXephLsBpN+U4ZsChPm0sTtogSfQmjR9cz5" +
	"D3PMranzh0Jem7zWM/y9BcZclhjyLBTyy7fo32mjnwWGnAgNGtXtl1lgZNqngeXaN+bKT9Cdk+ZC7lecsfB+vcdXJkcshH4f90Kx" +
	"n4LeT1ruo2HKY70HnOkPTSXIpMDin4kSeHhg3KfE91JOFwo92vyt2kDwqwLDBydlVtqQ8zLuEBk8yrPYXXI/D9div1Ju5hTxO4Oe" +
	"4a8KDR1A+hs0bCg/7H1j33S951GYwCPSBg7K/QCl/CpN8F/5Fh3A5G/FFj+eZ/Q9EHOXNejaTZAzrkXHtOGTaWtlLDqhhHcri38i" +
	"5Wz7G7M2f4DNRjNjGYGwfc15Sop1pQ28ygg+Fxr8J2Ph5RkDDwsNe983eF/aEqM04eVb9FBTt8wY8JO6h2exOyKLDRVZ9B/fEmfy" +
	"RW5YYLGBZLnQeG/Gy22xD1O3jiw6SGTwppQhmwLDtk9Z/FVpQ78pt8RBbX11nZ/+9pypT8iYdcqYg5QRJ4ksfiPTdk9ZdOMgwfcm" +
	"99a28cjQ8P15CTaoaR97Rp0SdmUW3TpX/NH08XsW37IneKGfAJNCUfevgJNTjPmRv3GestSfFvZHZPjATZ9quVGvab/b/Cg2mPpC" +
	"NzXjS3ruWwmeUWCJL6Yt9BMk8GabT17XmxZwjgSOmXutexY9zLPwUf2ujehLgSUHxTN82p5R3rfQnhzLWrGHoZtgp5u/xxMIWylM" +
	"iEvZfL2mzE6yC6Rv35QDpSIOLXEjsORmyLqzzrZ72oYJce7AEuuSeFLNe/6Zfiwp1wMjluxa5EDKMicVOWhdy7hKQ47Y4gElzrZ7" +
	"iZp2oC3XwpYvYOZdlRm0GBpxrdCwncx8Etfit0wZupHHOG+zc7MWv3Zo8bO6ho7jGTDzE3wZnsVuMmWqiWOhwQvSQk8y8+SkrDL1" +
	"npTIkzN1liCBp4cWeRlYcg9CC12YcRtbva5F37TFNiKLHewZPNbUD216mGfRgW2xN9uc2HITTB3R/L4gQT55OXLhUmLeMgmx5KQj" +
	"bcSlvIS4Qehsu9ermzBOyWMCi/8tib96CTZKkj3iWXAxMPKuAsM/YvKVwOI/iCw8odAS+zfjWKFh2yTZdpGhk0m9pTghhmPWU5AQ" +
	"q48MeR9a5ikyaMCW42fCTO6T7iXMhWfhT0m5D6Fl/sIceKD3dk3KJck62/6WuY3X+iKXWMrGKEf+QNJaQ1/IfFt81cxzMHM/TVsx" +
	"ylGXWa85r2ECngUG3/ITyoUJfMa38GDXyBXzDds9Zcm7lbIrMHA2Y5HXfoIvz7PwJN/i1zZzgQsTbFRpExUm2OJy33cvB38sdLb9" +
	"vUf52xfmvqlBjlzldEIOp2/RxzWf3cS2utlGWvhHasR+wTL/JEoYo2vxLZm+C287OUE2n4dpK6QMuBYl4KTJF8IEe0+2l7LI1tDw" +
	"J9v4V2jUm0qw51zDNowscRJbXMJPwEXXIr/8hBiZ+VsiUUJ+p2ljZQ09L2Wh76xB+74RWzVjfynD1isSNm2QIz/RtD08yxqLlBHb" +
	"9hPi34HQb0Ojj6FFL/dz6HauhUebOQmyr0WGPmjzJXk54gxmnVkDV0071ebXLhR6geQ3UQ49wsZPPUu+XMbQFbKGnh4l2NlhwvqJ" +
	"IEGn9YUtrOFdbPE/egk2W2jQf2DkpnmWPM3Agv82+yAy5jfYTq6ta7FzIotvx4yvZBLyqlxDt6ff61gofg8wSFhjIvXxtMWv6Cb4" +
	"GTxjDURo2ITSz5yyrK3wLTq1b7E/A8taA3r/d6dpL3qbLyFtxC3chHmRcclC43lBgswz47GBIVPSFns7MugwsKyrsPXdNWISfo6Y" +
	"YyohryuyxICS1gTo3zEJDH+dho+0900feJHwKYc5+FfakBU2G1jGhQssY/Is/v7QElPzDX3b9vvH0hec5f5vMHQy36B339DrNKxK" +
	"OTbrWdaQeJa5TvLhJdnaNp01TLAhJR/OJMgu1/CJBYbNkUpYu+Ul+KT8hLVjSetF3Bz+LN+yvisr5jmyxBPMmHeUwNv1/aFGLNiU" +
	"s67Fb+9bxuVactfkfLYU8tFGd6kEH4Rn8el6Il5lrv+KLHrTZfz7j56xtsz0BZo5IiPx3XXOtr+F5v4LPjHbmsIoIQdS8tEiQ96l" +
	"E+yOpPpCSw5wRsQVQgst0G9ofST8B2FCDndhQi5errWUGcP/5ib46HL5E1JG/MDGt22+tsCi83pGbrFNhibpaXSex+s4PEN3MXM0" +
	"IovunLbgtS23LTL651l4TmCxy90Ef4ZNZzdppcIy354l9yrcjt0SWeao0GKX2PJOH+d9UjxLvliQsNY2tKxHCS28LpVjfbUnZK0Z" +
	"l8oa9lva2f4+W4GFR5i0GyT4pjIJOk1gycnyLPQTJcQcPMvaZ3Ndrf5tWnc7enzGAnM/RztJa+jcHGsv5W/c6jpbJPhX00a83Ba/" +
	"t63RTgs5kU7wm4bGPEQJeBQYuJ/O4Xf0LPZ6lGAXhpb+2tYdu4aObepDJs0UJNj0LSwwK7Lonp5lPWRoif+4CbFuz2JnRQk5rjZe" +
	"HyasI5X2RXGOtdI2f36QENuPDDxM8h+HCTm8niXO4wm4BhbfcWTElWRMYwX/fqpv5Bna+JGZR7uBf2czdLbdB8T2e+dJvkw5pocN" +
	"+Kcs8+Ul5CbofPZsQk5vZPGnpS26iWuJ2zR3tt3TxMwvNnO4Uhafp5+gO0u7QuPNYQbfLzPaDC25p65ljX3SHh5hjnx236IHFlh8" +
	"sJ/w7w+7Ri6gyW8iw89k45/ZBP6eJDtdCy91LXlhJq3a8scj56d7PrgJaw8DS55Ekk0mbdy0yIuy4dh6/q1dky5cSy5f0v4DQY4c" +
	"e98iy9M5fMSyD3skrF8psfCgXLnvrmUdR2DkNJu+Rc/wmch9QOi7OmFvRgbspZ5cYuCuGR8KRf1+wlq/pDVb5npHqW/p3xMMLbKj" +
	"KEd8V/92YhfOVZfyKJtgY3gW2LkJayR03rybIG+T9t9wjbzWVELswYzPpCz5WJ4lT9O3wCSyrLvItR4uTLAlpQz3LbLBtubK3U6O" +
	"uJ9jXYOb4KfyLLpUUt5G0v4DXoIt51p0cc+i6yflQSfFupPW6mQSctfDhPXEniWWHDi597Ryf8b6RzeHPyXJB+Ql+F08Y12amWcQ" +
	"WNa1JOGbZ/E3ucYat1SO+Q4Scq/8HDGRICEPxCY7/ITDVi5j8J014nditzdvkSU3W/L8cqfpd3a9BL052I6cSxsyyRaX0HuqmniS" +
	"SfC5hjlgrX8jeKOh65g6eEkCjzZjX4Hhc/VyrBGOtuP73T2BvqMEPmrqO2UWGtLx5six7zdh2n5pg15scjM0/BF63eb8D+O4254/" +
	"1QmzAo8zhg9M5iR7lvV4GSPvK2WJKZixlcDwr9jymyWNF+ZYh+Zb9NxcfCkXP0mKH9hsWP2b3pGhq5i6TpHIHd4gfmvcs+Q82HwL" +
	"sm+UI93S0ANSOfhxOiFG86sE3uYmxEr0XDVPkAFm/rebIxfa5l8KEmzjIOE+TMhpz7VP495GLCOb4PvxDRmRtdginhH3zrVWM7DI" +
	"pMCyrtdLsCtMvzjhUa34XXbbejg9xiO3I9+S1l75CfiZ5Mcxv42c3HvquTl8frnW49j0y9Ci35ljLrDkgiTFA90cuce2fUKjBFz1" +
	"cqyDjiyx/6T1zjbYeMa6e3OPQDN33LZXXpJulElYt+UmxHLcBJ3Atay5Ttp3y82Bg66FR6zjdcQNQu8wYxbnC36k92DYxN/Z5qws" +
	"IU6iYbJZ2JuBobtpP20s9t83fTWm7AgT+G+SrZQxfLNm7olvsRmz21mfGyas75e+oUIx9qTcaSnzTt1OLotnrFGx7c0Q5PAFlxk5" +
	"fVkjRzfJl+wb/qsowU60yY7AgK8thzLJtk/iJYXbyQ/NJvAV16JHl+cYR8qiTzxyev73CfJ/+b/8X/4v/5f/y//l//J/v+y/YYvy" +
	"vzGeP/JH/sgf+SN/5I/8kT/yR/7IH/njl3f8H3IuZa8="

const csDiffEUCJPMS = "" +
	"eNoV03t0znUABvBnenb2l0ouSWlXl1l1mJAm9aQb2ZipEeZ9t727GrvYmBlFK7ltREK9p85JFBUqlEoXp3RRqJRyOuWELutCU4j0" +
	"/P75nOc8v9973u/5Pu+bhyEXgMnIsyH0SAIa0N65Gd/bFvxg9+CIXc9QCrCBYfsc8+3zLLAbWWg3MWJfYJF9kcX2JZbYzSy1W1hm" +
	"t7Lcvswp9hVW2Fc51W7jNLudlXYHq+xrrLavs8buZCgReINh+ybz7VsssLtYaN9mxL7DIvsui+17LLG7WWrf5/D+wAfsZPcwyX7I" +
	"TPsRu9qP2ct+wovsXg60nzLbfsZcu4/t7X72tAeYbD9nmv2CI+yXzLAH+YT9ik/arxm1h7jUfsNm+y1328N82h5lgz3G7n2B44y3" +
	"P7GLT/gz97r/hYnOv/KZdKCV6+xvfNb+zvX2D26wf7KfPcF0e5KD7V+c48+2sdGe4mz7N8d4wX8YtqeZZs/wOnuWHey/7GbPcZ09" +
	"zwT7H+PtBe6yiB1gY2JTbbvYPnZQcvCrGJb8o+2nSp8wXVW2v6rt9aqxAzTdDlStHaQ6e4Nm2MGaaW9Uvc1QsOAQBQvepGDBoQoW" +
	"vFnBgrcoWFAKFrxVwYLDFCx4m4IFb1ecv/0OtbN3KriBuxTc23AFNzZCOyuBu3UuCxipttFAplqzgSzF5QCjlGJH6+KzQLYWjgXG" +
	"aHwukKNL3IzVpfYeRccB96p1ApCrUxOBccqZBIxXBz+9T9snAxPUFAImKiEMTFJWBMhTRpH/NYrYkC7zm2GlngHyFVMCFKijm0J1" +
	"shGVlwFFOmKLdboCKNFJW6rePnmZOvudcnWxUxRTC1RoXh0wVVE7TZe7r1TfeqBKXZ2rdYWtUTc7XVfaWh2fBdTpKucZ6m5n6qlG" +
	"oF4H7Sxd7aZB+fOB2Yp3btTmJmCOEpznKtHer7gFwANKcp6nuYuA+WpbDDyouiVAk5LdP6QU+7B62AX6rgV4RD2dF6rzMmCRejkv" +
	"Vm+7RMtXAkt12LfRrFQ3LcpYBSxTpl2ubPuo+rhfoTS7UtfYxxRdDazSoTXA46peC6zWte7XaH8U+B+1giaK"

const csTableGBK = "" +
	"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" +
	"vVJTeXiW2ta7FX/L1UaWevVRH1x87CbqW0zRPQ9mKvZqqs9ZV/3QH1U3yf/ik3+tcQiu/+HOb21SP32V6hqx995SLBdTHVdL0z6k" +
	"NqG0nejG/urfr1GMV9C4NbHl5btDz+kcgzsn0bC+0xcZnrF/d7zTtbuKxNkTtUA/JH9cMt+2f/bEIsq8Un+1T3qqT1pbsVBebZa1" +
	"6vOWVXtI5rne0LkrSfN0/b1MObS0+sWF1A7f6plbWqY2S7rKo/iGZ6uverv6vGMVswH1S5/VuH/xFa+UPaDUjlmi4zIbNO5ig/J3" +
	"DU3bVj/aHYfqn+bIeqO9Trldn8KNUdvvmT9L0TJOqB75mxrecUozyn6l3hiUU1KObj+HlKPSPHMUfr4fGm2eN9WI/o2K1CixA36+" +
	"5RBJ9/v5/onfmNL/b6vuy1nX786zNJNxlPkOn3ltHNfucOeYRkGIaTscHs5Uv9Tta9n7+AdPnq4PfAn/vHjIsf8Fl49BcLcjcLVU" +
	"V7jx9IfXo5TTtagaY1nEEwu/Wzm5sPalK7dVPfNhbhzombX+/i0y7xygt+0PqX/j6sibD+bY+pE49d11TvzDWNVX6VpHt89Sw7Df" +
	"d6sjvw9InGqT6D5CMJf7TcaLQTcj5BdpX+RQHqwviRdrFG1P21h794TjPphgrI/8rjzT+XatfDvqzw1O/2XouE1QucH0jKWkhfWV" +
	"e+2ABzMiS1V0jHBVoZzvyAkfN4Fov5G9wRs6GN+j/nBn8P11Ocd23q+c6foPjcPKv03HJd1xm506TtOCNGaWy7EPLB07qeKZdzq8" +
	"PoTnsnnuDARo6U5nbriIJ46stJY3Ve3BNarHyuj8RLpyZ0B1l5T3z4W2fS0Fml8rZ87jVIR1VnEnHdF3l65xvoci5b+6VOJrJuJ9" +
	"m/5I2Dyx9q6MtSYmEe5NRoZjxexYPspT4y/HJn1gvTOXZHp4Uv6/u0HktT+L59n2s3D9I+D6OIz/dzF+J7QwjHaTHRy7Y9W2Z77e" +
	"1R2fZjj2736VpUoqv3t03jNF57dWaUyyq4/LaTrt1G8upHZSA7WlG+r9Ssr/pvqoIsNtVT5KaKxEDdUbp1JIsb8ldli+WfLY1pw4" +
	"C5dTtpfOsbktfd8dM3Xl3o17LqqcUcITD5GleknKW+lPJ/1ZJyS2jsNMYK/Wf8thRVlrFc/uyouj4VTHNnDHypfoWLfgVuZ5LlPe" +
	"kxjJjoqdhlpe8QtlfwLpx5KQ5VVyXfbZSXfKnuaJ63LHrd+CAx8p7vBmwDOm7fJxRcV0quIqRf3dyno9oD5tXb0fUpwfp3hLVf+0" +
	"lEdnF9U6bNR5tZDydHWtRwnl8crqv+5Xn+QauPsKj72+a5GD5/C1KsW1rnvUb9jnWZ9eVp/ZrWNDZbQ9Pq5S4Fvm57FzeU4cShHF" +
	"gzc+y+3bLE+8Vj9dF1hE5SXbE3ecqjKzIN3Bh3D+ei5U2eXwc7b6Bmn6+5w9jpx1AnCyp7jXZgv3Eb226CHFfCDMz0r1rGO1PWPx" +
	"sb4XHD6G73d/+HDb3M0v3vpuOT9mRU6/TFiMUtAzxvNPHqMMRonJTHafkiO+b9QBB1tWFB8h3J8IX08eaf19rPU04bhJ8/iUVoT4" +
	"m0CEdrU8Y+zR1sxbcWK/LU++wVzYpTKG+vkqRwaXoXgKrXbkeLfKt9hbp3BNvjPhzqvJvZ2V8nf/43/DOKiVZD1Mz3h+pLVcwtlf" +
	"rnTGooXvH5CYLZT8DK6VPZhj47pr0FoecmTksL1T57/X7ni/BWB79jnw8nbQiP1dOzNCWSONe4TitEtIYykO76GzKGfONEP1TFG1" +
	"S7I85SqkNpLt2f+gpD5X2jOfe0jtKXct7WYdR7XUTtqqdXfHnnboGPAO1Y2S1x06HlooF3sDFBzJH9WWOvb9n9qfVbTv3Zhid13F" +
	"AZ37OUb71uWz/WpLS19ehmFTgQdfUjs81TO/X9hjOxfy7L0g94cXzoktyPTgvIr6n9vsnHjhcmpblzJycHOs+sAbFZ/bdX3VXs9e" +
	"Ols8PoG7J0hNHQ/K4sVLuTBTv8Neh8I3T/v7Xh3R5M/125do+kWMnPWDIbVbJH/5vuGuIjnjYXcFCrCXr/GYh/5uA4eMv++fE8mW" +
	"+XiTw3dlPXs3HVBOd9f3rVE5qqq4zVCsp6p/WlpxXETvl1PsZig2XR2U4on/dG0id8x0ma6dnaHpbdB1YPt0LjeoaWd5MOvG05XV" +
	"mKFUj59aQ/XBRuXlvSrzaRoD+j7/X6ZlyNa52SzdM6KazhtM1WvldM56l8p5Q/Xh96l8PFk8b/ovVds/JUyPuzqscII2dtDn+tV4" +
	"3xb0u7e9n29d+R2/SmQN4ng6eaInBu0v5U353VLt13RPTKO7hrykcmiGtrtXN9RRXGzQPnf7pay20/7dznW3HDP48dfenDUdB1Rf" +
	"FNPxwao61nOi4lj4fItyvqXl8O7zl6nzyC01D3fedpHGXLjxSIV17jlL362sY/du3js15jhTZaym6hJ37jZdsb5Vnymm9c704KWa" +
	"vnNIbadUrdPBSv6+dZDs92vNI2A/WRHs2aDxz/h+SaLHhyudfjqoMZRu/NXY0N/5Xrjr4hRnTXnAMx9qKodZHm7Zq7j1xp64PnGm" +
	"J16ydErO3FQXxcSDHG12O7ojzcPxh9edU8iZFOSMMjn77bnpZnnihIKefAOe9a3FFZuPc6MlglSmXM4YZ5bH9rK86/B58UFPrNOb" +
	"ur9IxSrO8/Kti2KefRyqaKyF5L9E/aFAyCnPIW2XwzEeGQW2Rn4eT+xy/Liiit8Uz9jlTsWVG9uxT7nM9c2KhcVlhvu1cn4bA/tY" +
	"ANdyixNrdoziRnDSUfVFtvJhhmdOyPLsxeDu57pB48JcPnXXBJpqO9XUMpfxxFy46wpdeyPTyNnj9aDa/jtVl7nyWtSzhriWPp+i" +
	"fmlI653h8Z1d3bUk26mDO3/t7jOXuti2T6rhXD+obWoUL8Befh6fbnLwKRg7R3hd4hPl26wem7yM9mOqZ3y8gmfsUOSkXVrsvZ3l" +
	"+u06Nya/h+7J8SuKaf6b0D1VMTa+OODwpshIX42h7KF71wzT+MmSaZH3Ynkbg+TaatHnqkM+9wUIxBmD92sLpEa5f1s+f+fnnzLu" +
	"b4bFBoWP5UXyPwqHrdcPGdG/L5oWdg55ONPvPovhcYVv7Xds/0Ke6+VVVvaqDjm8d6uZE98zkBevOJgTVx/yjHlmqy/hlq+PZ3zd" +
	"9Dwn/7+e7fB7VY3rlHvPBHLSTNdyuOtTCnnmbd0xSzm/pfNH/VbYdocqf5cLefZUFMWUDCfdVmFt4c5TpXrGN137qKkn3iPDE/Pm" +
	"9smA7P/buM/vI/vK//4mQcHfkf8LvluAvfw83i3Afb789SrAfb4eJxfgPl/+fjrKuE/x2G6mJ846fJw/0jqQSGs93f22/HyzJ1Zc" +
	"cyL7dOZm/D78+qs9C3CfH383vpfceLzpY91SIEZskxklpsyKEOecTFy/maAcRHvPTADPiXyL1z0XKcB9vvzteC92nL0fbOd1TNc/" +
	"JQ4tM8H34sV/PnGJsz5VxiDP1XGGzikFuM+Pv5n7ovNW0Ii8P7i3T0sX9FvB37/w7/8BQ1QVGA=="

const csDiffGB18030_2B = "" +
	"eNoV2mOQXGkYhuHYVh/1UXbjjW3b7KPY2WRj27Zt27Zt27a59/fnqqmZyaTqnb6f7q6avaHDD6NE2Rc6gvtDR/FA6BgeDB3HQ6ET" +
	"eDh0Eo+ETuHR0Gk8FjqDx0Nn8UToHJ4MncdToQt4OnQRz4Qu4dnQZTwXuoLnQ1fxQugaXgxdx0uhG3g5dBOvhG7h1dBtvBa6g9dD" +
	"d/FG6B7eDN3HW6EHeDv0EO+EHuHd0GO8F3qC90NP8UHoGT4MPcdHoRf4OPQSn4Re4dPQa3wWeoPPQ2/xRegdvgy9x1ehD/g69BHf" +
	"hD7h29BnfBf6gu9DX/FD6Bt+DH3HT6Ef+Dn0E7+EfuHX0G/8FvqD30NRHkWJ8iMUFX+GouGvUHT8HYqBf0IxMYoUC6NKsTGaFAej" +
	"S3ExhhQPY0rxMZaUAGNLCTGOlAjjSokxnpQE40tJMYGUDBNKyTGRlAITSykxiZQKk0ohTCZJmFySMYWkYEpJxVSShiEpjJKkoywZ" +
	"qEgmqpKFmmRjWEqNuvQXGtLf2FtKg32ktNhXSof9pPTYX8qAA6SMOFDKhIOkf3CwlBmHSFlwqJQVh0nZcLiUHUdIOXCklBNHSblw" +
	"tJQbx0h5cKyUF8dJ+XC8lB8nSAVwolQQJ0mFcLJUGKdIRXCqVBSnScVwulQcZ0glcKZUEmdJpXC2VBrnSGVwrlQW50nlcL5UHhdI" +
	"FXChVBEXSZVwsVQZl0hVcKlUFZdJ1XC5VB1XSDVwpVQTV0m1cLVUG9dIdXCtFMF1koPrJRc3SB5ulHzcJAW4WaqLW6R6uFWqj9uk" +
	"Brhdaog7pEa4U2qMu6QmuFtqinukZrhXao77pBa4X2qJB6RWeFD6Fw9JrfGw1AaPSP/hUaktHpPa4XGpPZ6QOuBJqSOekjrhaakz" +
	"npG64FmpK56TuuF5qTtekHrgRaknXpJ64WWpN16R+uBVqS9ek/rhdak/3pAG4E1pIN6SBuFTqekTepSa4XOpOb6QWuBLqSW+klqh" +
	"Ka80o0Sx5NZ8nElug//I/2F+uS0WkNthQXkwP62QPAQLy0OxiDwMi8rDsZg8AovLI7GEPApLyqOxlDwGS8tjsYw8DsvK47GcPAHL" +
	"yxOxgjwJK8qTsZI8BSvLU7GKPA2rytOxmjwDq8szsYY8C2vKs7GWPAdry3OxjjwPI/J8dOQF6MoL0ZMXoS8vxkBegnXlpVhPXob1" +
	"5eXYQF6BDeWV2EhehY3l1dhEXoNN5bXYTF6HzeX12ELegC3ljdhK3oT/ypuxtbwF28hb8T95G7aVt2M7eQe2l3diB3kXdpR3Yyd5" +
	"D3aW92IXeR92lfdjN/kAdpcPYg/5EPaUD2Mv+YjYBPmo2AT5mNgE+bjYBPmE2AT5pNgE+ZTYBPm02AT5jNgE+azYBPmc2AT5vNgE" +
	"+YLYBPmi2AT5ktgE+bLYBPmK2AT5qtgE+ZrYBPm62AT5htgE+abYBPmW2AT5ttgE+Y7YBPmu2AT5ntgE+b7YBPmB2AT5odgE+RF+" +
	"lx+LrZafiK2Wn4qtlp+JrZafi62WX4itVl6KrVZeia1WXoutVt6IrVbeiq1W3omtVt6LrVY+iK1WPoqtVj6JrVY+i61WvoitVr6K" +
	"rVa+ia1WvoutVn6IrVZ+iq1WfomtVn6LrVb+iK1Wojxmq5WomEqJhiElOkpKDJSVmKgosVBVYqOmxMGwEhd1JR4aSnw0lQRoKQnR" +
	"VhJhaiUx/qUkwb+VpJhGSYZpleSYTkmB6ZWUmEFJhRmVEGZSJPxHkTGzomAWRcWsiobZlDBmV3TMoRiYUzExl2JhbsXGPEpqzKv8" +
	"hfmUvzG/kgYLKGmxoJIOCynpsbCSAYsoGbGokgmLKf9gcSUzllCyYEklK5ZSsmFpJTuWUXJgWSUnllNyYXklN1ZQ8mBFJS9WUvJh" +
	"ZSU/VlEKYFWlIFZTCmF1pTDWUIpgTaUo1lKKYW2lONZRSmBEKYmOUgpXKe1ZodVKB1yjdMS1Sidcp3TG9UoX3KB0xY1KN9ykdMfN" +
	"Sg/covTErUppfs42pQxuV8riDqUc7lTK4y6lAu5WKuIepRLuVSrjPqUK7leq4gGlGh5UquMhpQYeVmriEaUWHlVq4zGlDh5XInhC" +
	"cfCk4uIpxcPTio9nlADPKnXxnFIPzyv18YLSAC8qDfGS0ggvK43xitIErypN8ZrSDK8rzfGG0gJvKi3xltIKbyv/4h2lNd5V2uA9" +
	"5T+8r7TFB0o7fKi0x0dKB3ysdMQnSid8qnTGZ0oXfK50xRdKN3ypdMdXSg98rfTEN0ovfKv0xndKH3yv9MUPSj/8qPTHT8oA/KwM" +
	"xC/KIPyqDMZvyhD8rgzFH8ow/KkMx1/KCPytjMQ/yiiMoo7GqOoYjKaOxejqOIyhjseY6gSMpU7E2OokjKNOxrjqFIynTsX46jRM" +
	"oE7HhOoMTKTOxMTqLEyizsak6hxMps7F5Oo8TKHOF2ugLhBroC4Ua6AuwjZqLx5F/6m9sa3aB9upfbG92g87qP2xozoAO6kDsbO6" +
	"mH/VRV2CXdWl2E1dht3V5dhDXYE91ZXYS12FvdXV2Eddg33VtdhPXYf91fU4QN2AA9WNOEjdhIPVzThE3YJD1a04TN2Gw9XtOELd" +
	"gSPVnThK3YWj1d04Rt2DY9W9OE7dh+PV/ThBPYAT1YM4ST2Ek9XDOEU9glPVozhNPYbT1eM4Qz2BM9WTOEs9hbPV0zhHPYNz1bM4" +
	"Tz2H89XzuEC9gAvVi7hIvYSL1cu4RL2CS9WruEy9hsvV67hCvYEr1ZuiffUWrlZv4xr1Dq5V7+I69R6uV+/jBvUBblQf4ib1EW5W" +
	"H+MW9QluVZ+K9tVnon31uWhffSHaV1+K9tVXon31tWhffSPaV9+K9tV3on31vWhf/SDaVz+K9tVPon31s2hf/SLaV7+K9tVvon31" +
	"u2hf/SHaV3+K9tVfon31t2hf/SPaV6PwCDmnRsXzajS8oEbHi2oMvKTGxNfqIHyjDsa36hB8pw7F9+ow/KAOx4/qCPykjsRE2ihM" +
	"rI3GJNoYTKqNxWTaOEyujccU2gS0tIloa5MwnTYZM2pTxGs/bap47adNw8zadMyizcCs2kzMps3C7NpszKHNwZxaLMylxcbcWhzM" +
	"o8XFvFo8zKfFF68htQTiNaSWULyG1BJhIS0xFtaSYBEtKRbVkmExLTkW11JgCS0lltRSYSkthKU1CctoMpbVFCynqVhe07CCFsaK" +
	"mo6VNAMrayZW0SysqtlYTUuN1bW/sIb2N9bU0mAtLS3W1tJhHS09RrQM6GgZ0dUyoaf9g76WGQMtC9bVsmI9LRvW17JjAy0HNtRy" +
	"YiMtFzbWcmMTLQ821fJiMy0fNtfyYwutALbUCmIrrRD+qxXG1loRbKMVFWujFRNroxUXa6OVEGujlRRro5USa6OVFmujlRFro5XF" +
	"Llo57KqVx25aBeyuVcQeWiXsqVXGXloV7K1VxT5aNeyrVcd+Wg3sr9XEAVotHKjVxkFaHRysRXCI5uBQzcVhmofDNR9HaAGO1Ori" +
	"KK0ejtbq4xitAY7VGuI4rRGO1xrjBK2JeI7W5ornaG2eeI7W5ovnaG0BbtMW4nZtEe7QFuNObQnu0pbibm0Z7tGW415tBe7TVuJ+" +
	"bRUe0FbjI20NPtbW4hNtHT7V1ov3QdoG8T5I2yjeB2mbxPsgbbN4H6RtEcVpW0Vx2jZRnLZdFKftEI+68E7xqAvvEo+68G7xqAvv" +
	"EY+68F7xqAvvE4+68H7xqAsfEI+68EHxqAsfEo+68GHxOw0X0fk9hr9F5TcYPiZ+g+Hj4jcYPiF+g+GTOCN8CmeGT+Os8BmcHT6L" +
	"c8LncG74PM4LX8D54Yu4IHwJF4Yv46LwFVwcvopLwtdwafg6LgvfwOXhm7gifAtXhm/jqvAd8WoqfFe8mgrfw6Ph+3gi/ABPhR/i" +
	"6fAjPBN+jF/DhbPy3B1+m4Xn7vA7/BF+jz/DH/BX+CP+Dn/CP+HPGEX/glH1rxhN/4bR9e8YQ/+BsfSP/MzY+ieMo3/GuPoXjKd/" +
	"xfj6N0ygf8eE+g+xbPpPsWz6L7Fs+m+xbPofsWx6lKc0qEfF+no0bKBHx4Z6DGykx8TGeixsosfGpnocbKbHxeZ6PGyhx8eWegJs" +
	"pSfEf/VE2FpPjKf1KHe5gB4Vz+rR8JweHc/rMfCCHhMv6rHwkh4bL+tx8IoeF6/q8fCaHh+v6wnwhp4Qb+qJ8JaeGG/rSfCOnhTv" +
	"6snwnp4c7+sp8IGeEh/qqfCRHsLHuoRPdBmf6go+01V8rmv4Qg/jS13HV7qBr3UT3+gWvtVtfKenxvf6X/hB/xs/6mnwk54WP+vp" +
	"8IueHr/qGfCbnhG/65nwh/4P/tQz4y89C/7Ws+IfPRtGMbJjVCMHRjNyYnQjF8YwcmNMIw/GMvJibCMfxjHyY1yjAMYzCmJ8oxAm" +
	"MApjQqMIJjKKYmKjGCYximNSowQmM0picqMUpjBKY0qjDKYyymLIKIeSUR5lowIqRkVUjUqoGZUxbFRB3aiKhlENTaM6WkYNtI2a" +
	"mNqohX8ZtfFvow6mMSKY1nAwneFiesPDDIaPGY0ARxp1cZRRD0cb9XGM0QDHGg1xnNEIxxuNcYLRBCcaTXGS0QwnG81xitECpxot" +
	"cZrRCqcb/+IMozXONNrgLOM/nG20xTlGO5xrtMd5Rgecb3TEBUYnXGh0xkVGF1xsdMUlRjdcanTHZUYPXG70xBVGL1xp9MZVRh9c" +
	"bfTFNUY/XGv0x3XGAFxvDMQNxiDcaAzGTcYQ3GwMxS3GMNxqDMdtxgjcbozEHcYo3GmMxl3GGNxtjMU9xjjca4zHfcYE3G9MxAPG" +
	"JDxoTMZDxhQ8bEzFI8Y0PGpMx2PGDDxuzMQTxiw8aczGU8YcPG3MFf0a80S/xnzRr7FA9GssFP0ai0S/xmLRr7FE9GssFf0ay0S/" +
	"xnLRr7FC9GusFP0aq0S/xmrRr7FG9GusFf0a60S/xnrRr7FB9GtsFP0am0S/xmbRr7FF9GtsFf0a20S/xnbRr7EDy5s7sYK5Cyua" +
	"u7GSuQcrm3uxirkPq5r7sZp5AKubB7GGeQhrmoexlnkEa5tHsY55DCPmcXTME+iaJ9EzT6FvnsbAPIN1zbNYzzyH9c3z2MC8gA3N" +
	"i9jIvISNzcvYxLyCTc2r2My8hs3N69jCvIEtzZvYyryF/5q3sbV5B9uYd/E/8x62Ne9jO/MBtjcfYgfzEXY0H2Mn8wl2Np9iF/MZ" +
	"djWfYzfzBXY3X2IP8xX2NF9jL/MN9jbfYh/zHfY132M/8wP2Nz/iAPMTDjQ/4yDzCw42v+IQ8xsONb/jMPMHDjd/4gjzl+jX/C36" +
	"Nf+Ifs0o9+jXjIpjzWg4zoyO480YOMGMiRPNWDjJjI2TzTg4xYyLU814OM2Mj9PNBDjDTIgzzUQ4y0yMs80kOMdMinPNZDjPTI7z" +
	"zRS4wEyJC81UuMgM4WJTwiWmjHEsBeNaKsazNIxvhTGBpWNCy8BElomJLQuTWDYmtVJjMusvTG79jSmsNJjSSouprHQYstKjZGVA" +
	"2cqIipUJVesf1KzMGLayoG5lRcPKhqaVHS0rB9pWTkxt5cK/rNz4t5UH01h5Ma2VD9NZ+TG9VQAzWAUxo1UIM1mF8R+rCGa2imIW" +
	"qxhmtYpjNqsEZrdKYg6rFOa0SmMuqwzmtspiHqsc5rXKYz6rAua3KmIBqxIWtCpjIasKFraqYhGrGha1qmMxqwYWt2piCasWlrRq" +
	"YymrDpa2IljGcrCs5WI5y8Pylo8VrAArWnWxklUPK1v1sYrVAKtaDbGa1QirW42xhtUEa1pNsZbVDGtbzbGO1QIjVkt0rFboWv+i" +
	"Z7VG32qDgfUf1rXaYj2rHda32mMDqwM2tDpiI6sTNrY6YxOrCx62uuIRqxsetbrjMasHHrd64gmrF560euMpqw+etvriGasfnrX6" +
	"4zlrAJ63BuIFaxBetAbjJWsIXraG4hVrGF61huM1awRet0biDWsU3rRG4y1rDN62xuIdaxzetcbjPWsC3rcm4gNrEj60JuMjawo+" +
	"tqbiE2saPrWm4zNrBj63ZuILaxa+tGbjK2sOvrbm4htrHr615uM7awG+txbiB2sRfrQW4ydrCX62luIXaxl+tZbjN2sFfrdW4g9r" +
	"Ff60VuMvaw3+ttbiH2sdRrHXY1R7A0azN2J0exPGsDdjTHsLxrK3Ymx7m+jX3i76tXeIfu2dol97l+jX3i36tfeIfu29ol97n+jX" +
	"3i/6tQ+Ifu2Dol/7kOjXPiz6tY+Ifu2jol/7mOjXPi76tU+Ifu2Tol/7lOjXPi36tc+Ifu2zol/7nOjXPi/6tS+Ifu2Lol/7Eg60" +
	"L+Mg+woOtq/iEPsaDrWv4zD7Bg63b+II+xaOtG/jKPsOjrbviv2074n9tO+L/bQfiP20H4r9tB+J/bQfi/20n4j9tJ+K/bSfif20" +
	"n4v9tF+I/bRfiv20X4n9tF+L/bTfiP2034r9tN+J/bTfi/20P4j9tD+K/bQ/if20P4v9tL+I/bS/iv20v+FS+zsus3/gcvsnrrB/" +
	"4Ur7N66y/+BqO8p9Xv/YUXGtHQ3X2dFxvR0DN9gxcaMdCzfZsXGzHQe32HFxqx0Pt9nxcbudAHfYCXGnnQh32Ylxt50E99hJca+d" +
	"DPfZyXG/nQIP2CnxoJ0KD9khPGxLeMSW8ait4DFbxeO2hifsMJ60dTxlG3jaNvGMbeFZ28Zzdmo8b/+FF+y/8aKdBi/ZafGynQ6v" +
	"2Onxqp0Br9kZ8bqdCW/Y/+BNOzPesrPgbTsr3rGz4V07O96zc+Cookl49zS6aFIcUzQZji2aHMcVTYHrIznF9SK5xPUiucX1InnE" +
	"9SJ5xfUi+cT1IvnF9SIFxPUiBcX1IoXE9SKFxfUiRcT1IkXF9SLFxPUixcX1IiXE9SIlxfUipcT1IqXF9SJlxPUiZcX1IuXE9SLl" +
	"xfUiFcT1IhXF9SKVxPUilcX1IlXE9SJVxfUi1cT1ItXF9SI1xPUiNcX1IrXE9SK1xfUidcT1IhFxvYgjrhdxxfUinrhexBfXiwTi" +
	"epG64nqReuJ6kfriepEG4nqRhuJ6kUZ4P9IYH0Sa4MNIU3wUaYaPI83xSaQFPo20xGeRVvg88i++iLTGl5E2+CryH76OtMU3kXb4" +
	"NtIe30U64PtIR/wQ6YQfI53xU6QLfo50xS+Rbvg10h2/RXrg90hP/BHphT8jvfFXpA/+jvTFP5F+GMXpj1GdARjNGYjRnUEYwxmM" +
	"MZ0hGMsZirGdYRjHGY5xnREYzxmJ8Z1RmMAZjQmdMZjIGYuJnXHY1hmP7ZwJ2N6ZiB2cSdjRmYydnCnY2ZmKXZxp2NWZjt2cGdjd" +
	"mYk9nFnY05mNvZw52NuZi32cedjXmY/9nAXY31mIA5xFONBZjIOcJTjYWYpDnGU41FmOw5wVONxZiSOcVTjSWS0e/84aHO2sxTHO" +
	"OhzrrMdxzgYc72zECc4mnOhsxknOFpzsbMUpzjac6mzHac4OnO7sxBnOLpzp7MZZzh6c7ezFOc4+nOvsx3nOAZzvHMQFziFc6BzG" +
	"Rc4RXOwcxSXOMVzqHMdlzglc7pzEFc4pXOmcxlXOGVztnBXr55wT6+ecF+vnXBDr51wU/TqXRL/OZdGvc0X061wV/TrXRL/OddGv" +
	"c0P069wU/Tq3RL/ObdGvc0f069wV/Tr3RL/OfdGv80D06zwU/TqPRL/OY9Gv80T06zwV/TrPRL/Oc9Gv80L067wU/TqvRL/Oa9Gv" +
	"8wYzu28xi/sOs7rvMZv7AbO7HzGH+wlzup8xl/sFc7tfMY/7DfO63zGf+wPzuz+xgPsLC7q/sZD7Bwu7UR7w6tGNikXdaFjMjY7F" +
	"3RhYwo2JJd1YWMqNjaXdOFjGjYtl3XhYzo2P5d0EWMFNiBXdRFjJTYyV3SRYxU2KVd1kWM1NjtXdFFjDTYk13VRYyw1hbVfCOq6M" +
	"EVdBx1XRdTX03DD6ro6Ba2Bd18R6roX1XRsbuKmxofsXNnL/xsZuGmzipsWmbjps5qbH5m4GbOFmxJZuJmzl/oP/upmxtZsF27hZ" +
	"8T83G7Z1s2M7Nwe2d3NiBzcXdnRzYyc3D3Z282IXNx92dfNjN7cAdncLYg+3EPZ0C2Mvtwj2dotiH7cY9nWLYz+3BPZ3S+IAtxQO" +
	"dEvjILcMDnbL4hC3HA51y+MwtwIOdyviCLcSvnQr4yu3Cr52q+Ibtxq+davjO7cGvndr4ge3Fn50a+Mntw5+diP4xXXwq+viN9fD" +
	"766PP9wAf7p18ZdbD3+79fGP2wCjeA0xqtcIo3mNMbrXBGN4TTGm1wxjec0xttcC43gtMa7XCuN5/2J8rzUm8NpgQu8/TOS1xcRe" +
	"O0zitcekXgdM5nXE5F4nTOF1xpReF0zldcWQ1w0lrzvKXg9UvJ6oer1Q83pj2OuDutcXDa8fml5/tLwBaHsDMbU3CP/yBuPf3hBM" +
	"4w3FtN4wTOcNx/TeCMzgjcSM3ijM5I3Gf7wxmNkbi1m8cZjVG4/ZvAmY3ZuIObxJmNObjLm8KZjbm4p5vGmY15uO+bwZmN+biQW8" +
	"WVjQm42FvDlY2Jsr+vXmiX69+aJfb4Ho11so+vUWiX69xaJfb4no11sq+vWWiX695aJfbwUu91biCm8VrvRW4ypvDa721uIabx2u" +
	"9dbjOm8Drvc24gZvE270NuMmbwtu9rbiFm8bbvW24zZvB273duIObxfu9HbjLm8P7vb24h5vH+719uM+7wDu9w7iAe8QHvQO4yHv" +
	"CB72juIR7xge9Y7jMe8EHvdO4gnvFJ70TuMp7wye9s7iGe8cnvXO4znvAp73LuIF7xJe9C7jJe8KXvau4hXvGl71ruM17wZe927i" +
	"De8W3vRu4y3vDt727uId7x7e9e7jPe8B3vce4gPvET70HuMj7wk+9p7iE+8ZPvWe4zPvBT73XuIL75Xo13st+vXeiH69t6Jf753o" +
	"13sv+vU+iH69j6Jf75Po1/ss+vW+iH69r6Jf75vo1/su+vV+iH69n6Jf75fo1/st+vX+iH69KA/p14uKUfxoGNWPjtH8GBjdj4kx" +
	"/FgY04+Nsfw4GNuPi839eNjCj48t/QTYyk+I//qJsLWfGNv4SfA/Pym29ZNhOz85tvdTYAc/JXb0U2EnP4SdfQm7+DJ29RXs5qvY" +
	"3dewhx/Gnr6OvXwDe/sm9vEt7Ovb2M9Pjf39v3CA/zcO9NPgID8tDvbT4RA/PQ71M+AwPyMO9zPhCP8fHOlnxlF+FhztZ8UxfjYc" +
	"62fHcX4OHO/nxAl+Lpzo58ZJfh6c7OfFKX4+nOrnx2l+AZzuF8QZfiGc6RfGWX4RnO0XxTl+MZzrF8d5fgmc75fEBX4pXOiXxkV+" +
	"GVzsl8Ulfjlc6pfHZX4FXO5XxBV+JVzpV8ZVfhVc7VfFNX41XOtXx3V+DVzv18QNfi3c6NfGTX4d3OxHcIvv4FbfxW2+h9t9H3f4" +
	"Ae706+Iuvx7u9uvjHr8B7vUbir8S9xuJvxL3G4u/EvebiL8S95uKvxL3m4m/4vb7ZaYvPxXvoa76IbzmS3jdH8jnb/gdctKXXwZv" +
	"+UP5zG1/BN7xdb7nrj87F335Sm768hPhA38kX33oT8FHfpt89OVL+MRPw/c/9c/w8TP/Dj73O+SnL/8sH7/0M/LVV34mfO1XLkBf" +
	"fpuC9OXfxXf+En7aez8bX/3gZ8eP/nI+88mvXZS+/HrF6MtfzWe++nlK0Je/gY+/+5vxh78Ff/oF+Fe//HV89be/Ev/4O/h8lOAW" +
	"H0cNLpekr6BBafoKKmKMoATfHzNIXYa+gvQYO+iJcYJRGDc4zr+NF5Qpx/Nj0B0TBL0wYdAfEwUDMHEwCJMEgzFpMAeTBTMwebAF" +
	"UwSbMWVQk/8lVVALQ8GCijw/BnNQDuaiEsxDNeiCWjAfw0HySjw/BinQCFKiGaRCKwihHUiYOpDxr2A1/h005ienCZrzW04btMB0" +
	"QUtMH7TCDMG/mDFojZmCNvhP8B9mDtpilqAdZg3aY7agA2YPOmKOoBPmDDpjrqAL5g66Yp6gG+YNumO+oAfmD3pigaAXFgx6Y6Gg" +
	"DxYO+mKRoB8WDfpjsWAAFg8GYolgEJYMBmOpYAiWDoZimWAYlg2GY7lgBJYPRmKFYBRWDEZjpWAMVg7GYpVgHFYNxmO1YAJWDyZi" +
	"jWAS1gwmY61gCtYOpmKdYBpGgunoBDPQDWaiF8xCP5iNQTAH6wZzsV4wD+sH87FBsAAbBguxUbAIGweLsUmwBJsGS7FZsEzsc7Bc" +
	"7HOwQuxzsFLsc7BK7HOwWuxzsEbsc7BW7HOwTuxzsF7sc7BB7HOwUexzsEnsc7BZ7HOwRexzsFXsc7BN7HOwXexzsEPsc7BT7HOw" +
	"S+xzsFvsc7BH7HOwV+xzsE/sc7Bf7HNwQOxzcFDsc3AI/wcbC2uZ"

const csTableGB2312 = "" +
	"eNrtnXd0FtXW/ye9QUhCKCGQQpPeRZp0QaQIShVpUkRRmoiIIsULioCIjaZSBUEEESnSi3QpoYP0XgQChPRkfp/N2bMyNxe89/7u" +
	"e9/1/iFrPWuezJzZZ+/v/u599jlnnsGqYFkelvmk8NkbZ9uZtm2fy7DtjuUty5NzIfG23eSMZa2Kse0Q/r7LdWnvpfcl5TLHe3r+" +
	"UR8/px9tN3O/bRfnGMOnTqxl3dHz12yjQyTfOdhn+KCKdYHjOI4jOKbxucv3Dzh+ykd0qcznBt/T+dznc5ZPA84Fa7/3Vf4vHPOp" +
	"7GBvy5K++p+w7VHRpt1+aX8a/Txsu5Ztzj2L7W8hsB/Kxnube7fx8VY5YpP8+12PufX82FTbbudh+hA5gucIFCqD7dZ/8K9J0Ufj" +
	"/P/zOf5PfPffknfkd4PDuYr/mty8e0372Nh/T59Dp819mzOzzp0X7v0P2/3X5x8/Y/h0rvAXzv9Tn0X7/vXc8VKFv/D6b3ya/BMf" +
	"POf1n/fh+SfXvkz/8/4H/pf7H5z+n41f/61/cY+w+w/yfJ7/JW6E/AfY3HmE/gX5XPpfGqvi0/59/TfzkXrrJscwjlKLHeFYieMe" +
	"rcM2c7wgtnDsxnGe1HY6Br/NcSLHExyLcTyl/HM+0uYWnzb2v6bbtv+DeS/vvv9uzGwq+lde/nc+CSf+b+awv/799e+vf3/9++vf" +
	"X//++vfXv/+Nf0FtLOt+Xcv62ybb7hlEbX7HtmtyfiX1dlysZa3hGB9sWT/dte1PAi0rV6Jtz02x7ZepOe/Tbvohavrbtj0gwbZz" +
	"5basj5nHxNKmPtdbXLHtXvktqzUTmN9CLKv077bdiTZtc1jWU8j8zN+yyu2w7WUce/K5ltOyUsMsawYTCl/6nY4uNejjTHGuM6mo" +
	"jsyWV237rLdlJZWyrO8DLCt8t21PLGRZDY/a9u91LKsU14+g61O0eY77l9DPjpu2XeqGbS9/ErsOY18+yxrK3y9yrU68bTfOtO1f" +
	"oi2rTV7Lytxr2zuxszjnO9N3byaQS4pZ1sunbbs1OpUpY1knQy3rQyY37WmXl4LyIvJaZTBHu27b93JZ1nbs2IrsWoA3GTBXg4cs" +
	"BD92wLZzVrWs53fadjQTpTeRc3GVbTfE9pb1LKsME6VwP+pT+smF/lvpNxzc3kbnV+7bdmHax2LfjLtmDXpALezGxlbg+Tfue43+" +
	"1hTGpsu2fZG/y3Hvp+i2YD9YlrSsJy/adjr4e5y37b34tip67UhCDt8/xY6GsnaN3lN/Az/6Trtm2725p3SEZZWnj5e5PuWyWfvG" +
	"XLs0esga30D6SAC//sxzfoAbazh3orxlFcZOD/pdh8yTPviDPnIkwwH4kpdzebDrPXwTgcAxyHumVtb6w2w+q/Fpd/Qrh037fS2r" +
	"b7hlLf/VtmfiG7uIZdVHvxjONxI/cOoIvvoILMveNBwOx59t0PfH0pa1DP2eQF4YeiwFw6fyMN8Enz1wrSa+2VOR+47Z9hDav4Ks" +
	"ZmC9jLnkAf6eftC2/0DfSPh3lraDkD2NvubxfUFZyzoPh/3Q6RuU74Nf3wGHsmCbCBdeggsn6GcBupfj/H5srYqsb+l/IpzKSLXt" +
	"w3D5GbizDn2/IX5mw6036Osp5N3GXz+ttu2DTKgT8cVL980eyGRktyBWevB9Nv46RAyUwE+ZyMxBzBUFg/R0225Z0LKa4f9Nj8FR" +
	"YvUSMjdEEj9Mqksjsyc+SMIf3fEXMFpL6GMf+o4lPmaD/XLOvcA9Kdh5hhgsB35j0G0MWKKWPVH8gB7j4ULAPds+gF3LfwHDxpb1" +
	"OY7sjS2BYFAUuRewPZJzjZ+wHiSWWPA9T0xIPF3Bnj/OkGvAowdxaqNrAfQbzvU86FcL2wah8/Pg2gmbDkVZVgB+eRx5sml0DFBy" +
	"ovdrcPk2vthDn73BpSVGjYWTU+DvBrhYHZtu0Eexk7bdBV1Pn7Ltcei/kvPVsCGEa8vgy2nxO35qS1/fYvcM7H8dWR3QMZr54x36" +
	"yQnG+SV/oNtwcKiFf9tAvN3ku4bwLeBH2/6JfkvBrW7IPsv5/OTYD+Hc67IfVcWyvkL2VPQvRLwuKWdZ/uA0Bf+Fw69z6DCkAH7G" +
	"zopgNYHYGU/7k3CsAj5Fdfst2j5BHJwhUQ4jFnbjq7v4Zg7+TgS7y9jQG16X5+/qyDxCXEeh5xlsSQWX2di3HlwvoFsEcrvQ1zT6" +
	"eZf21/DX7gu2XQJde+KzMmttOwq8W3HPH9izGLs6k7Pm0EcBePs0DUdsIVbQaS2+eIP8HApW62i/Fb+9gU+mwt8gfHeAPhbGkVvh" +
	"9xF4Pxlcj2BvXvCLw7Am+MaP/keI/+HxDvzbHHtX0M824iWBeHqXmFmLvmtrE2f4Yg047gaLH5A9BS6+AGZpJNNV4FKZPJkAzwYw" +
	"JvlxbiY2HUSvC/C8p8QTn5Gy1oRN4eAzCQ58ybXD8KEUuB1C/2h0qEyuKwVWk5B1l7+7YlMnMGqNDq35OxX/LKDfu8goFIOvybWD" +
	"kFOQuJkDrzdiTzCxGoOMp9ErP/d+h6x78CSD8x2x51Ow/xCf7cXW6vjGjzjwpP3X2Nadtpewfxfj9VjuOYkNo+mjCvoFgf9IMJxM" +
	"LP+N49mrZi/0JvheRW8v2uWgz73g2g1Mm+Cbmeg9aTs2wv95cDYRnV8iJ+wD/wWcO4v9NYkPT8buQ3D0LDZOB9/niIeF9HkHmZ5w" +
	"7h42exKn44mF87vIr/A8Dv/H0nYDSlx4ijFuI/kZG4Pp9z487II9z8KnobTpij9tBpVd2PQGfR2vjy4kGIZL60NsHkmcjaVO6Y9O" +
	"O+nvdeKkF7YchXtvkfNjsOcgdsSi7y1wsInv6bS/CYZTaBeNrq8zaPqj5wg4SzP7E66tBNsu8O0y8rtj83j4cQjutoQv3cmTtemn" +
	"Bv170HYX916Cv7LWvAl/FON6L3Dupxx7CV1zYX+N30z8H0ReCLZeBKsX4EAOcvM75LkBfO9E/shE5xbgMRTOJBOL3bg3Gb835fte" +
	"BuMTWxm7aDMczJ8hvhOOmBjtuZ56BS51A/sG3NMIztRB39beZk+9G3GQj3h9H1lbqZ1qol+O/aYWKgZ+18gz1/FxJ2qO4UvB+mnL" +
	"Kki8vEbt0wo7gvmeTNuD6L0MecPA9zrYzcGmysRsAvefw0fd0EHir+M5akJ4NIs2Y4mb++i7SXQC46WMA4fROw25EfCoEDr1gIuT" +
	"pJYkfhOIn8/h8BK4MJN7fOjjCnb5UQuUR4/itBsD1uuxoy/X+3HPK/ThS9zcITd+Qd/3wGQweacdvtnIfWnkJuH/z+SZ18gvXuA1" +
	"rAF1Gb4dTfwdrGZZRRizp1HLXcL3Edg6EuDe4P7FfN8HP+6i4w36XUqb79B5DTj3wK4guLEZbv1EW094UQwujoEHY8AlDlsWwvt2" +
	"2HALfCPpZwtcOUf7FtxbjPMTiNu+4DoTno/fY2LpLXi/G9/5kCOacb4s8uLgjQfx1kyezUCH5fjmBPL94M81cNjJ+cfRsRu41pI9" +
	"AvotzLVeYC11Tz/y4Q3ue5WYu4mfd3PPh8illLKaYvPPfClKX2m/Gq5sRc/N1NJD4Mfz8KchfiuA3g2QOQ/ZP9D/UfLKLPzRAKw/" +
	"I34HY+Nl7J+Cn48+Tv47Tk7Bjnz450l8uEy4Rk6Yhz7l7ppnOMrhj2+w8zh9L0Cn+YwVXxEPHsRsMXT5VOKd2P8SP3dHdj50ER9N" +
	"AeuKxMXb5ImfCN7m1CAr8bEP1+uBXxN8W4ncfZt4DQK76ch5mX6Xg68fbR4HpBXoNI9YOI5dScTXXRljJf/Ck0LY+Rn3tKdtbeqX" +
	"E+godVlXeDkPW6ui/xn6m4dN0cRuAJwYhs7tybEt4EJbfNUDzPfio7LMQzwY+0M5boRLdcCsN/geQfabMu8hF1VAdkvsrQYv28DP" +
	"WfCrLvomg9dubDiPPgvo9y4+GUYSLI3ee4nhZOJgDrhNIG/MZ2x4HZn1ya1xP+MLctlp8J1ILniV9nORU53YeBf5e9DlZWquxfh9" +
	"EHhOktwsD+YgrzV5Zii4vgCO+4n1t7GpGDr7o2dhOHwGeZfobxfxMIN464RuF4jZAuhfFXu/Avc29NWafPcKnM4vtRH+aE7+7YiM" +
	"vvjtEBz6lhw7Gb/uQt40eQaHvpty9Kf/5uSSltQA92hbDd5eJV5CyCcziDfJlXvhuQft4rF1KZwOJpcswfbh+DQYHQ6Rf07S11Ng" +
	"5AdPa2BbfeqVj8C1kNS6XJuAD3ZTAwo3fkPfSdw/l/FnD33t4/sw+mrOd3nO6TLcqI2v1pFjviemfgGbk3B8JDrehXszJX+g/1l8" +
	"MZA+cmBvCfh6gthuj/xXkD8Gzu6FizulpiQudoH9SjB6G+yXItPHw8xLvk0z2B1BTx/mTJ+Sj72p4ZbBvfb480N49Tn+XAQ2o4jh" +
	"J+DrHa5HlTe+C8cXs7h/m+Qs4rQncsZi60eSS+mnOP2Ppj71x3edsW04Y8p9rvdFzwPIXIKenfD/V+j8Zmmzx/sZeL9IXKai63j8" +
	"FYfvHqMuOY38IeB9i2v50WMC3B9K7ixOXn0LDOPp8zPseVXGZuLMGxu3w6F7cG41bQfT5y3u74ZtB7FtAnybi22jkV0T7j970cTz" +
	"SM5tQ+4QYmIZ5z5GJy/sOoy/joJ/KvKGwJmq4NsJP+2nfVfsfIVxyRd+tsWXSZwLIP49uGcTdr2EnOU4twf8qgJPvpW5Kz4M5PsF" +
	"7D9I+4PyfBD2TkbHMZKnuL4eLhSSmggedsa217CRcLIuI2s8fv0eDjbl3qvwPRc4zOPaM2CTSt1zGsyf475m2NrltKlRaxHPBeBj" +
	"K+TKGN2AftZyT130Fs5tReZP+Gke2CXx/Tg58nPab4d7l/D/u3B/N7GcUdSM9xa6DUG+PG/2C/quRe7b8NZ7HWMzeXYq3PiZ+z2u" +
	"mDx1Gi4MA4u14FkKDl3BT7e4bxaYWvjgAj7/nrG3CznqMPoNRqcvkTP4lnk+7wz3oKpVCeyvS07GnhaMCa3JVxzsLnBzMZwaJes3" +
	"9NuAe69zXE39VEF0AasEdBkPjmOx4z4xuYb53gbs3EEOa0B8r8F/bcgZY8FwLRgUpe1qbHoB/5bjXFl02I69/vQ1nn6L0GYq+akM" +
	"/W/lXlkfeY1c8x56RiP3deQPgwdS648kX40jl0+VejlNnvWjFoRfr9PfcM43kToD2btpFws3ohgLZ3H/m8gqTNvnsSuNzxW4WwqZ" +
	"dbE1gvOPgZsPwJxiTLlFvLyILX2pSVtJbkK/dPDbybXbtCkGL/dzXEy+LM54NwwfPEnueH8DHJXaHP7dkPkL3JJ1kDfRZTt93iEe" +
	"vuO+1di6kD5noesJ+FAQTgdmmrWZ3uTl7/FnO2JslTzTGWaeC+mCH88gN4ExReppGc8PIEf25f3oMwEOjaJdOLw+gE0h4HGTeUx9" +
	"7FwJZ07SZj7XDuL3RfR1gj4aoU8XMPoZvJozRnQA87LokkLu7UgM1CA37bhnHrycwbjanZhdJjU5PDxGLqlKny2REUmcjEf+VjB6" +
	"E658ju0Fub4Iu5vDl0h4WR+eJVCvRmFPQqR5rrMg+G9DxzjG9wC4EQHmcehWijrrGryLx8YSMh+hjwnE4dgL5jnN8+TF2+h3DZ89" +
	"B07juV4K8BLpqx1x1X4Fc73Cpj68j3/abTbztC/BNw82cZs1llifgx2vkP99aVsNXD3g3kV0+hAeVcHvrZA9mBiqi79fAcM29Nuf" +
	"fhdh9xB8+xVyhpK3JoJRL2ygVLNCkbmZuU8COHnKWiL+XgM+15B9j3nASWLGD4xCyZl14U0fctucSybvdCZ35sIX75IDYtHpc3C/" +
	"Dgdzcl8h+ngVXUbjn9+RlYCPRqDrWfy7D79GIOctjh+g52IMvIMOsq62GzxbILNVoFkT+Jh8NYQa4W30mh1v1uAWY8Ngrm0gIWXS" +
	"fhHxvRHZr+Cnd9BNatQgdLqHD/pyXwK238S2d+FQPuK9B32H086Xfs5yfTPxXgObC5PzNmqMPkX858N/8VLbod+LcGQp/vqc2Owu" +
	"4zh2nAYX4Wlz5B5Bnjf3Ps/5FHRqQPzIM7g9yFlX4XZtxtWfyB2NuT6J8wWJxVbkvNbweTDtu9LmBH/7gVtv/BYDLyrRV32wO4Dc" +
	"kVIHYWd58shA6sK/QchfOPcEum6Cx+sZI0+gy1MlzFziMvGxH72PI+Mejq6znNqnCXGBf14Ei3NgOhs96uLDhfSZTv/DkPUqMupR" +
	"o/eD0404tx4fTSMuSoJ5B9pvoQZoypzsZ5kTYc8H6DKBMeFbZB6ljT/8LAJ2XeBNmSgz39pDTtlN20D4UAW91oJNEWyOwbcTsaM3" +
	"eu5Ff2/iqQm2P4+PEzjfhjhcgm9mILsPcfc2nBgtuuGnX9Hzcex8D2zmU8MIhqHwoxm29IQbd8lVA5BTgXt3RprnjfrFm9wj62Ed" +
	"sP8c/g9hjDhO/83haRHizw97tnBfe/q/h5ynkfcefjlG31Pgf3OOhJ69Br/5M44XJ/aaYUcJzk3AxqfJC3PQpQsfL+GRPKtU3HBZ" +
	"xpDc6BCKvEvgcxV+9mJ8sdCtD/zqTAyWIY98irwO1GbP0/599H+bexoQp6ewYxux3BwFBsALP/xxA/9Fgv8qckoF9FnPvcu5/hI+" +
	"bwJPl+D/Ctjajlw4F/xnIeMTdF8Hv54Ai3zk0GPU7QvANIZPOnicQ5c3dhoc3iOOvicX+iLjd2KkT4yZh/YAl67kqGh8HLWSOCZv" +
	"FUH3m8gOkfWSmuQOBsal6HgRXevw8ZE9G/hfB2zfIN570E8l8tqCKDM3qEb749zXFWwy4dIqfP869nxP+5XIKQpeu8B3A37Igf6F" +
	"8JXMz4aS98q0oI4ByxbYOwdZ74P9dmzNJN47w+Oe2Pi4zEvBMhN/x5GramFDLfJaR2xci95FcFYQOqRgVymuLYV/N+BkJbANA8sU" +
	"5pzz0HUdcf888p+k3Zn9Zt0gDi7eAaerYN+IGnQW488OarRt2FmY+2cmmt8MyFy4ltRFtH/xntkj6YT/I8HnCvx7HMyflTVRmVvJ" +
	"Hgr2lccnW+mrFjlkDzxdxbj9PbIOE0shjE/ecKQ793wNf74GE8H7quxzwM1k9G1HfniKPk7jwz6M+22Ru4hxtTL2ngWvXrQdCb7X" +
	"6RNVrangnw6vfkbet+hdtaVl/cicb6CMGcTuWNljwo4X0eeezIW56Rrc6IWuO9Elk5h/hTZTkbUdG2Px3zaOF8n1fcH0BHkrD3pv" +
	"ItYtOHOD9guYN37G+PkJ4/D3nHuWHLNInnMkD3Unrm5z/w15dpH89hmc6Ij+T2JT70Azzi5hzPkYeRuwYSp2L6DPa9hwjnta1DQ1" +
	"5H0KlBjmVRP4nps2lemnFfjckHV5yafkuv1wpg/5a7uPGUsikVcMx8n6v6yf9UTmTYT50WdL2i0Ev3fhxgfo0Q4MD+PbmQdMLfQu" +
	"cdKVv0/C1U7UhW+S45Lh6HJkNqOtJ+NJbz5d4k0OGFLazIPm8n2xYEw8f4LvuvJ3C/QOhWuTiIG/gfcursfAg/6yj0JObHbMPLc5" +
	"hesSf+/A72ZwYSTybsDLnuT9gbLGDdYZkhMaMQdhDHrT08w5+oLfBjBoDS4DOCbhy2nY+TSxVoj48cNf22gXid7zkftRgqnVD4LT" +
	"NtoNkefd0PMM9v4Cz27Ak15g0JIY8AK/t8knk8GhO5wpy3hcGZkyp61GTOWir7Po/QGfHMSr7NneAes+5KNh2H1K6m70ioav3gDb" +
	"Q/ZiuXcpNvckr+Uj1tZgVzu4NYf4GYmeZ+BuLTg0Dh1/Zi40AMy+xs+LZH0TvyWByzZy9FB0K0L/Lz9h1sli4GEINtxDXk1yzR3G" +
	"75N1zO+QqtFfEXhYjeuZ+NtG7lx06cXfSRyHgUUYvqxXyaw3tpV6HV2a4I/exEV/+qmaYr5341iE42py3LPoPQ5c1qJPAPwLkv0p" +
	"5vz1ZI8YWRPxUyw6jyA2YtG3HLptA+NnyDNDkH0NW2QjMxa+TIYvW2U/E7ubIfN7jmnI3Zxk9pBHoGsyRCnDOVl/2o9/V/L3EHRr" +
	"St4oIjVmARN7W88ZP86G82OIkzNwqiE2huCXycT9S+SiaHyTRh7vCrfygsEF5BQEw0r4ryB5ORc+fVXkg0crdPmCPs9g73Bi+hYx" +
	"NBtb3wHcpcjtLntZ+CH/VrMnswF85zDePAZGdXeb53xPce8PyJT9wgnIP0HMDIKL6Vwbhy3fkpcK/W5ql9Lkk47440P4nsw9Lckd" +
	"9aReJ0+clb1FdI4GuyiZz3L8Aq5NQ9aXxMpk5KbK+hI2TZR6ERt6kXOG4v8A8thefDOOeAyTPQp0HUebDeSqW5KTuT4fOQvx/SbG" +
	"5jRZX0P53+DpxirUw7LuaJt11cLkla5c34ddfTgmy76a1Oz0VYW+g7Ez0tvMh7zBdhZ/P43sTWBWn75+I++v4/oVeNAG/7wrNTKx" +
	"anPvd/jIk5iri7wPGNcG5Da/i2sZa9YTgqipTuLfObJWCh6F+ayh7YvgvgibfsSOOrQZCk970WYH1y7h77PosodYqUb7O+hdAttk" +
	"m/U9/P4eY8Q6sOyGXY2JzfbMZdoyfpySNWzZy14HJ9DbDx3OnzTzolP4tiayDtFPScvsFy0FuxT4VQ1/5SRHfMN3LzCRtc1v4ehd" +
	"aqjqiv85WX8Dux/wa2v8VAeeLa2Oz/BfGeIpD3F/HnvXc89Q9M+AryH0GQlOrQuYejWXrHHQVxX82AKcfwSDw/Anglx1g7iK4TiX" +
	"+IhDRkkwn4/uBWX9CH3GcP/ORLOGshY7vsA/PwHICLgykbFgBVjWxK/yTEdzuDiVWuk5uF9H9g6JiwHk3qG/mzWoJXB/PLVKGcaJ" +
	"ddQmU7GhJHO7JeC5gvjoipw25MH8xHl+fNqc/mdQ631FbmiPbuWQM0P27vF7yZtmjjCeeO0Jrl+D4Qv4ey461gKTI89QX+02tdVz" +
	"+CECjHZjx6e0kTFyG/c9A89D8YuFzCr4fAbX7sJh2fNOJUf1JBcVBZsE+rHh2+H75neKualXCpCfwqkNF8CJLuj6AX6KJe4XYPvb" +
	"+Dc/fC2Jr0vBhUlgN01iEy77gUlhco6XPK8BB7eCYSV8FcS15YyjP3qaNY6MHWb89+e+8oxxm+HfdbjlD04bwf8ZsNrApyAcDefz" +
	"E/iexZQVcKAOPmmGD7uBxQhq0rzo0og+28CX6cIFbFzMWLTKx9RnMqe/Rb+j0Kk1frvJ/WPoM45c0U5qL+Y9ucHxiRum1koBi/mV" +
	"TG6qTd8p2NQVnuTHjkDsWYj+KfSVQZwmXTfrnLFg2BmcS4DhEj6dZG+HfofL/Ae9++Pjr8hhzyEzDj7Ece0u7a7ggx7EyKeyho69" +
	"AeCewVi1Hh1C0Ok4fl+CLxfB6w3cc0K4in9eJ37SG1pW5K+2fRQb9tQ0vPqCOLsKB4elGFukruhAPj1NfG0k/zQHk2pwYTnnRmCP" +
	"+FKe/3iaPtbRxxuMl434Xp7z38LDqbKvTSxEgMNhWQcHh7nc+wkc34mc49hfHhtQ3WqODw5zfhFxGIxOqGJ/DKffYCxvTP+XaFQS" +
	"ndZFmHWEzfiuPT7vhi3V4f/+3/AzfQZJrYW/djFORsPzWfgtilxcgtgvA2/i4dXX8GAf/ZcgvtpQS8nYVhrdlqJ7aakf6UPqjpeI" +
	"/3nE4kdcn00+608e/QFbPpHfSJFnDsK7VXx/krieB8dn4N97xEk9+u+GroXRPZd8x0d5sK8p3LzLWPGWrOv4mefgFoDLs4xd7bB1" +
	"iewxElOhAH+U+/zB3IO6vR5+qo2sDdi2nvuO4uOr6F4YXpTHryHoeR7snpYanjj5mDZ74PXz/F2bvBtHrGyCP7Hc1wJcZR1lJ+2W" +
	"kAsqgutqPpvgY1H0kXXpI/TZGD4WBM9f8OVJsAxGxjnuK0F8/LiW3EXe+BWcKss6qezdoW8MTntPxkjibTZj9DFkduP+DXDrB77f" +
	"gr+XsXeWrGvy91zuf4J7iiG3M74MIrYGgcM1+JeM7l2SzdhRFT+OBINx3LMavXfhpxFwfbjUo7vNcyu4w3oM+2qh1xgwjKaPA/D/" +
	"S+JgIH5qgk02nB2BPq3JxZ3BqLNgBz8CuX81ut0nzxeFr2nEyWXsDeZvWetsi/23iNk14HWPHHYEXW+ga09wW0VfqdRLvvBxEZhv" +
	"Q8eTtOlAP0Xg5RXkhgWYGqoOOT0Hsj/i8wP8S4QTv4NHBLJlzlkW7mOmvQF++2L4MeyUZ+8scOqLL34j5xylz37ce4rYyeB4Fr96" +
	"gPV6sBooe2pwtH+QieGNxJb8ZrwzGI2FI3P4O5Y64Tb3tuP4M3+vhKcJ8OoA9xzD9tPYvEbmLOSZbXC/PsfZcC0MXQ6Sl5ph21z0" +
	"ywE+1eFpGO2jiI868GUj7Q+D0TfoVElyGfYWRvc+6DSNfmrL/iJ6FGOM/Rad95J3bkuuhsc74NUi7qkmNqK7YPkJmMhzdxuxczHy" +
	"ffkuNcgBZPTHhq0FzW/vWqL/d4wllWQdEV+0lHmt7JuSB3/Bx/VkjQZ/t2Eu0o8+5uPTydw7DJ1WYktPeb4RP1jYMoF7c8pzFeTe" +
	"qvApisLkKXJBB8aM/cxpPtzF3Ap/TYI3++BtTrCSPb90fPrVCbNuPIPvvzNGXCPWd4PJesusvX8OR0bBtb5gdJqY2MD9d7CnITxI" +
	"QZ8v8N82mReDxR98/xJc88CDI2AUBTFacX8rxtG82DgKm47DdW+JG3xfztu8i2Cp7J3h23Jwrg+2HcOOpdgXxpjUHTu2A9gFct0K" +
	"fLIX7i7CjqrknQvkrmcZh0pLPQbX1ixHP/LSF+AcRB/zOBaSuRUYPc73/ei+GMzWw5kW3J8HH8heb14Id1X2WPBfeXLDUHSPhSfr" +
	"ks3e4jDZF5XnFTi3hXvfZBxKIo7l+bElt8xzkj/CxY7olRufzkeng7RvTPtt5PGm8myXPHeJXX3R7z1wipc9AalpyTuDuLiB+N5D" +
	"rZS5jJoVjGRNNJHrL2Ob7KnIWkwg+nRDzwy4OAifrEfOKOQnFzO59yLX7jE2nMX/ZemzPLqlkts+AsN0fDToinm2Zxt6DgWT2fim" +
	"N7IH4Ptm8GQuNX8e+h4XbfYuvsAfa/FxkDw7BVaJ8jwYGE3nWI6YWgynffHbEfL8bcayWdj4I3gugvMJ8KEX+vvQ75dgOlveYYH9" +
	"LckXxeDcaGxOQsfb4HXwllkjeoPjbuInA25dQ58OcG4M8vrR/9cZZp8qnPZDGAu3c/47bOgj65tcm8enAn7zBqOS2DACXCqj2/P0" +
	"fw2/N2acr4K/LuKbb4jvwVJLyh5MsnmO+1v0T2uAbP6+w9j9QqLJac3BLwp+TMcPr/CRdfPXkfFdRcYOxpRJ9LulOeMevAznvtfI" +
	"8RVoc0Ge94WvX6HrFDCex7mBMs/Db42x6z30ngQua8jxQfitIDYGEG97aph3RfTDt9XlWUDOHapn9vHy08dnOc1+QAD3voEOk+hn" +
	"JPGbE3tPSC6UZ6SawDFq81exvS9y6+Pzlui1Dc5WJZ/NB59Ga8we/3I42qm2ea6uLPiHhpnnKdfeN+Pur8h4E/3Oy1wAnj1BHOVA" +
	"l8bosAj8/Pg8iR9+XWP27SKR158x7O4h8xzYs7JXiS754M0O7h0JL/ZhYHqgmYfJfLk8XOsAZ33AshP95wPHIcRzHo5N8R0lkB2F" +
	"L3tgwwmwa4TcX2XvXXjM+a9lrwjfHIUTdcFkIXY2gruX8elROH5M43iArMHKc9TMlULIHQPgQVPGu23w8YNT5je2P3JtBnHyIVg0" +
	"gDOryOFH4Edz+LAWvV8PMNx+AZ6uh9/dqaPykP9LYHMO+vgMvEbJ8+rYMVvqVHK5P7ETSDyH0kefomZ95l0wq808Ix29A/hE8/kF" +
	"344Gzxfw2TbJi9gzDewfk7jAv8FgNZjrbbF9HDn9V/wYybX68D4R30SA4yB8OXALYwuyDsjv7qPNc+Mdvc1vFrZxrSJ+erWmef9M" +
	"7hXMGdH7U3nmFpsahJvn7hoQp23xcV78koNaZCdtCjJHfp5zsp/SDE4eAbcCGQbjNvD8Nm0zwWoyujan00BkfCW1DG1G0LYi48Zo" +
	"cnlPfP8ZsRFKXb8dW7vwfZ3sy+CPlYwTH4ORcLTuetN2F+3+Bnd6ETsDkZeEnd/jew/y5Spic+RNs58lc7BKxc3epZwrgT4JyN1J" +
	"3Jzn/ALarsUffjIHJ4/Ew53i8P4sY0sctp+FH7XBoQx9fUqfK6kFDsC/1sSo7HuHgfNweUaJOiAK+VeR70GOmAzGV8m7C7AvL3Vf" +
	"LLlF9JVn+zPgQXuRDz5fgm0L4jYFm+ojdzociISvi1bDK3zTgYKwLOfmM5eLITjm0s9t4siXcXAWcXuHvgugZzN5XgrZ32Hnu/Bq" +
	"6k4zZr6MPR+hkyd9fcH4lo++ZN9a9oQerPWj62p5lqewwbGZPLsiz82DQTnyTk2OK9DNg77eh1dSjxwDk3qFzD7Si7JGzrEi9VUq" +
	"4xZusD35vCzP7lum/jqKjbXKmLWr7XCtKHYeJZffR/YJMFsAh3YTHyfhxzvYeJYcOpKYq0vMNUCfdHJwUez/HPmTvM1cUuZ09eX9" +
	"RdjyAVz7HD71ApsFUjuC2RHqvJ9vmprdN9jsA3UFMx/siKdNHvLGk9hxm3j8Ep/2IKZkvV/26hvT5mlw6QYmbZA5Xdbg8H1j2TPB" +
	"nznh5DDyTZjMgeU5EOwZjkKhjAsb4dbmP8xz4PL8yxByQQmxF5B+BsdG+LYwY/VC/FCFRr9KbUP8riUW/sC+59FlNPd2BZ/a5I0J" +
	"1EtSIyTJb1ho77UKG8mfK2tZVvGVnMO2Gvhi7zkTC/vJ97/xebBHUsr8fqIumDei/SBybnl4vhKcX0bfUuizgb83yd4nvLlBTHQB" +
	"L0+4vYR8cAU5nty/DYyPktemEeMzbLOGc5q+oiBAG7gzi5huj3198MFiZFfHnlbY/qGMAXzuY2NfbL4GtjUgxceQZAW5/2Kwme+c" +
	"l+doSeh3iJkexMg35JUuj5OH+P42dWBbMIyAH73hwXHiYiL5bK7Mb8EhmViYCaZdt5jf50TKPi3xkIbc8fhhOX6qRK7dA/dDNpsx" +
	"eyx6dkOv/mBwsrh5Z4Rg1kneKwHmA7kvBj+mc8934NoIvVeBjxc6nwKva9i6AsVfIk8lpJjnX2pHmrquAn1fJLbr4MsDjCefSC0B" +
	"P/vT70v0dQt9e9JXf/gWyL01SphxKh8yKtDXKbg2h+8XsH0z1/vAwRDZ8yRWKnM+VR7kIA+uoK/F2H+Qa8Xh6xTORYOz/BYrFL+e" +
	"wI+3iHvhtuztXuD7fuS3B4tV4L6J/kehx2xwKMy5TRw/gTNb8PtNbHgNbhdEzmJ0+Ak7/bmel5hcKs+KgPcg2ffGxjqMp/Iel1bE" +
	"RW3a/8K5uxcMb07FkSe4pyc5YBT8KwUeL6BjJvcvAJPX0DVMXjjCcbOs7RNLg+DCAXn2QJ7dx09FkH0I7Itg66tgH8WAdRs/FEV+" +
	"KHn7Sc6XQIf8xPYtREmO6gPXosg5dYnzfrKfRo4fh7yusieGzEH0fVKea8A3M5EVKPvrxKg8Eylk+JQ28nz+FWTuQbd1tO8q71+T" +
	"MZsYeRfuj8YX39Fe1hzHwvk4vh8m3zwN1mvRc7vU/tz3Iv4rw5ieSf/7wTwaHzSQfU7kyRxf8thQiRt0XgfPOoFrQ3QXvSvKe9k4" +
	"Drhg1peFL6dkbCVe5XdWU+BOJPlpQxUz522FDj7yTKbEkzxvvY5ciAx/eNwE7AeDY2v8O41+j2PXvjCzVjJKng+U32yg61U4MEee" +
	"wUVmLWqbZ+DWUtq3lTEKm/bmN++Yexx+XIOj59FXngWW+uMSen+IHTInKUicS75YRZt8yOqLTrvApQNcykT3qeSsN4m1KMbc/DvM" +
	"b7QOolNeMPYib7WVvRo4vhAbX8Onc2TvDh3rcW3TfvPboJbotoNc3J7z+RhfjuHDjuDjTx00k/mqLzhtJOe9Dxaf2KYuTEHHI9x7" +
	"Bv405ZNM3khHh1pSI8Ph8thfHv98xfiwBl79Fm6eh5ly0zzTuCDO/I6tLX4f9od5rrUIx2e4voja8ADx1gbQX+fYU/Yf8PdQmd/K" +
	"76aQf4ac+R74hsHdQ/h2BvoN49oHfO/HcbDsB8PbJOJiO7796aj5ndRbYPEq2M+UZ5jAcxl/ryffvQruH3P+3atm3E0A43Z8psz6" +
	"x/d2tGhBrUKbKti3H11uE9ONEswc+H3L/K4wnHlLjY3mt3YbkP8puX00OesI/q6UputJsl6GXbPpdx/+9cEny/heC459D7ZTwT56" +
	"iW2/gswBLakp+C5L8KVkDxsuHiSXhaaafRt55nwnsuWdib76e0Z5XmwluT4Zbq0DH9lvqijPPlumXRR1SSAxJb+rnIFvfPXeaD3m" +
	"0PpAjmH6Dkb5HeDzxPNe/OKj7x3MpdduI+MAPn15n+nDV9+/WFyPgfreRsoU65i+Z5EhyZrE9wqqM2nBLqzvVxT94/Bf4evmnYty" +
	"PVnfeyhzfhnz5NnQQnpefl8jupJyHtgk8g+obNxqP+d6x8/v+o6mANVNhgNZL4mwTL6ActZ8zlXXd0FGq61xZtiQVGtHqN1tZY+U" +
	"v8PUPrFb2uRSXzl9huq7lUSfK/ruxlT1h8iStcUDKN4VgyIvGF+L7w88a1nziCH5bW0qHJqIkb9R18h60zVql8/g6SWMIMSsVPJF" +
	"wE6zjliO87UApikyV6BAJ45rqpt6MhS/d8GPsn68WNb+4OE++c0HXL0XZ2zIoTYL1vsU2wccqG5867yXM9D1Dq8g/e6v9tj63kr3" +
	"e748Xfc67/X01o9899E2Xg/5uGV4aVsvl0z3e8Sc88IFn2zXvF0yPB5ydNpWpkYRv82pb/oSe0KUM857N331faNB+vFXvZw4DNe/" +
	"A9XPV1IMJnItXeXcd8nzUv2E28nKvUR976dwLEHbiD9yq4xUfU+pt8a2t/LP4WK6YhCscgNVH5GznTz1OAZ0TTHnPZSrsm8Rl2Hs" +
	"KyBxJTJl/x9Bj0FgD8gudlduQt6lnm6K8R3ymnXLTsjqpHEpc//aGms59PcuXqpHgOrovC+1rPz23jZyH9N4uazvZ03UmLmptgSo" +
	"b+X3+lLXX9VcU0jzS06VJcfHXXxw3qEqcmtpTgvQ97bJ+XLaTmI8Xv0Spe1I49ZJfZ5PdMyj/g5VPUswdx+Acpc3m98oOH4O0fuj" +
	"9H4He4cj4tec+rf8FiS/nhe+LFcOJGmuc94ve1Vj86bqL+eKqV+PaJ6V+dMtxay0+uOwyj+t3K2iOVMwc/z0h+ZD6Wu35sfGiqPD" +
	"kUWc+4jjQuWp+OJVyzyPL3kyRvsT7okfz8lvtvg00fe05dQ4quzi/TKOLRSbTD3noeOErFXVdL1jL1Df3xugxyTlguTgQ9p2i2Lm" +
	"q3qHKa7Cn7pmO9uW9uW5IbSMufeA6i5rhY2Ug8cVg72KeQWNQ4e3d3RcEh/LMxKUTrace0f9EqA2FtB84KF+zqXY/6r2FFEe5FU9" +
	"5yhHW2vseOk1wTGvcuK+9ivc36S+lb8vKJaCozxP4QGZS0PEVeoDJ66c9/39qj7spDrFuHKLtCmBwrNQrrxypKTmmTM6hl2Xelb5" +
	"Jv6oon2nmFcnPDgfqTaWUP4VVf6K7puPGQyqpJr4TuRk1HHzXeSsRFDuAIOBn8Zdfl8Tm4+rHbk01jL0vcuiwyF9J2CY8lE4Kva/" +
	"o34JN4+dPnhXcx3V64zGuMRAVeXPVY3DrepHyQ1/4zhI+f+B4rVLxwgf5WkB9ceT2s819cdFPeehsRGr/YSofTnVbicPiq6nwGA7" +
	"x8YknmQEtUgwv1/MoX0574POp76Jd+Uq9xgsmIudVfXvXcrFlorTBfX7PW3jvCv7rubEJ7RGO6V11F3NbYdVh7J6zzeypiW/21G+" +
	"eatdOV267NEYStV6y4kLweY3PXdK/em8n9vJAZ4aD17avoRi/LlyS3wfrtwM1XvS1I/RLkzWKvdJAQ/4LbluqY4F/lr/pKiNokN7" +
	"8zoJW+T2UM54aN4Vv7dSGyWXSH34k/LAQ3kUprF7RfNyouZD9/usI1WGxOKPgPDRNtOft+aUIppTflccHlP+lFY+Of6I1XEqSfNp" +
	"vMZrJdU7p2Ifov0HKO+Fc8Vc45S3xkqijg2+rjrH4d09X5Nn0xTPGNd4kaZ9Ofx2arrSmeZ9m7ldfd116ZOsceqMhzc1j9xWrvm6" +
	"+OfgGKw8vqK4O+NmUdXVX9um6jiaR/kjHE3U/sprH+Gau4L1vaO28idI+3FiTfbqGsaZdos9zLmgbPVsgLZ16jcf5ZfgMU/eIQCB" +
	"5Hd7Xq52HionU22zlWN59b5UFz5OPe7gl6RcKqi5KVjjOFzlJKvMfBobfqpXitqYoeNSgGv8kdj8GkGBDEKbSZhPbjS5pFvdrBo+" +
	"+8ff5Vu3bQ+r7Z17bvxqYjZVx9xcyqEItcWpeWaC+/A445dAjRmRs1/2dU8brJxaQ8aDGzp2ntJ4KODKjxf1nPQXq7jd0TE2Qf+v" +
	"gIIq/4yOdzJnbefSX/SQWqWE1ufO3C63xt5lzW+pemyl+SVR80+avn9WfCTPPb2muN1W/8XpXCxYeR2tMf+E5pbcGmOF1IarOPEI" +
	"Sbyq5gnfbPMwX1dN5JHtXbjZ39Hs9pXbX2EUFbkJriMoGRNodPLXuOmgGAUpDsK5k8rLIFd9XERzmLfmwGesrP/D4bLmuXIat0XV" +
	"L15a4xV2jf81td9LOk6JXRs0H+VRzgg2hTVP+mn9454jiG87a34sr+fWaxsPrUW9tU/JMSs5/qD8iND6cKnm3oqah5O0FpY884W+" +
	"l1j27uUZiAaaWwrp2PODjvvO2Hhax2LBdbqOK9KvPOPVXOcIW1zjYab2FaM8l3y5UWtt+XuLct1D/xYsL+pY78xrpEZZobWQyOnq" +
	"mtP/obkzUPtqp349rWNbhnI0WmMvXnXLpVwN0PbO3LSAXpM1dG+dT2WqLOHzW4Daj8Hoy52MhUmm3/plzbtt2jGQbUTIEm/DeYdP" +
	"AcoVH9Unp8v3tvo90sr6fzwClI9OfZ6kOudx1SuZWl/Hq+3Czc7k7DUlTX/OeJimY8zWvaYfW2VdoKP+GC7PLkv7/Pp/izxYR5Lf" +
	"xZwxeUwwaMzfO06aiUxlCvvcGUZGoPLOiY0E9es55Yj0X075lKF56IrO3zx07lBa18DuK0dzK0YlFb9Tin281nKPaU0YpLKKO36R" +
	"36ofN/w5qfHYUOfx25U3nsq1EmpnsMab8/+7eGv7YMWpoOIo+BxWPtbWvoM1Fhwu5a7z92s/Ptnmjd6umPbNtvbk+Yhc5pktDyZf" +
	"MP15ZrvH22WPO1f6uf72drXzduVYv0eMe5E6Pns/ZI3Mw7VG5pXNjuzrZg6H/uy9+l7/wjv45dxTcNKHAuFFBpfuHNuhZMp5g3Eh" +
	"nQdn6DrrPR3TLmtM5VQ/2hp30nYOyVue5T2j+TJJuVNN5UVqztyttZaXjiWZ2t6p+x7Ttml63lvXaRK11o9Sfjn/h5Cvaxx2atIr" +
	"Oq8orWNnhSJZtnfdZ2I3SOcPEie1dCy7oPmimPJwr8Zgotp/WnOk+GSQ1gD79B31Yvd23fPfo7YX0zGvlI6DWzQ25e8yui5SS+e5" +
	"zjqZh9bC4apTjObbaYplqMZ4cZWdqTnZ1v+vaLzWF5cU44K6bia5XH7387HOMZw5SWmdjzlrcTKvGa/5oZuOiY+pLpGaM67omJtf" +
	"85HMs6bo3EnkNtLvp5RDeXVczqF5ON41tzmlOSNAc7KMTU11rByoY+56xcxDx8UJuu4vY5W8K+tHzSMntC65oeseN5S/RVW2reO4" +
	"U4eN9zdcK59s9HbGFGeNLNU1f/NVHTcVMfHj/Ygc44wJHg/JLZ6unOGOS6fOryJ1xEHjT3eOsV3/T5TEzBoI4X/P8DBIx0WnL3n2" +
	"d7VLtjMGd+APeU9NmOIzmmAZiMOLa6w46+kP9kMizF7kfV1jDHathScqJ3Nq/Dq8iXXNx5x10CbFs3KST7b1fTdumTq+OfsLfdOz" +
	"cv3DctjXV836w3qN1wC10fGhXzbsB8cb/sg7PCZa5lkW+e3Q5Rvm/kCtW7yVo56uWsI/W051zqW5dPZKMZgIz8/nMDp4Kx8rQbSf" +
	"5dnscgbjJB2PW4DvsgiTO1OVM874USGbvXe0LnTqlgStU6ojY0+EGaucOViK+seZ9zzYn6BdBu2exoEnE8zcONi1jpVX9fXTuYXz" +
	"/4Tk0/lJgOZXD41d0fcy9VFvFDpYPGvO01DeWxVm5Ei8v3DJ/B7kjmIlct+7ZuYyni6eyN9Dyv69v534Ok4tOAGujwLXiiT/r0l+" +
	"m6pTJ24zODq1qu2aLzuxszDQ5L9U9WWU8iNJ6ytfnV+d1zVhZ1yOUDn1ddzOoXmqrOZSZy8hSvnvqWOkcKiu+ie3rlUU1bhpKO+y" +
	"EL9b5jcbo+Oz1jmceLibN6vW9NL7nZrV4Z2z7pxLOZvbtR4RonEp+kZdNvJfU586tfIJrXUFN5lndFS/yxphM9W9nOog45Lsz4sc" +
	"eY9odzkvz9kEGJv9XetKTt06DQ6MCjO88XDVtA4fCymmfmqXr4630XreQ8fUUnrdR3GuqPb66fiYx5VTcqkNf+i8zkd5WkTtyK08" +
	"jtbxM1XHpJfhbmdXvr5/1OCZfa8kTG1N0nEjxbU/nl/bJGptlE/xSHWtQwSrPu71IUd2iGu9aJjuywWrvzJc675+6rPDAUY/4fw1" +
	"TsTcN/zM0LHBX78/m2T83ByD5Zlud62ZfYxy172Zinn2/3PRz7WfY7tq8UetK2QfF32sR68jPWyP1yPbmPqwNQrPPxl/vbLV1o8a" +
	"j7PvZ3s+RAevR+jr9RC9/V1jmtdD1l+y71U79zk1tke2OYLnIzB5mC5Ov97WP///tzweMtfw0PndD+cNB04ReDkuGB4lKr8k39Xj" +
	"nPzOwZnXJWqtm6njyXuyZkQA7KY4zJ+eleN8NVfVyDQ+epBvSmTN47Lr5utaF81eW3k/5NmBR3HP6yHPBrh56fUn+PjoWsaDZ2iO" +
	"Zs0ZAzXOcmleCHHplUNzlPv/8QzXdnld89lMzWfOXuotrWO9NE/Fq+1O7XVPa/B7mhscm2J1bpapY/9N5V4u17p+ms49Cqhsx5+p" +
	"mktFVkcSS8E08/+hOjHvzK9zunJnDtezB3J9ac6suX2QC+cYHf/u2FnrtRFW1v+l5qlYldQx+A/F567uLyW7nqW57RoTnGcyimk9" +
	"FMKNL3Bij/4OXd6fWNU/ixe+f+J/p244ofKDray9SB/NvdK//L7vfvA/rnm7eeT7iOdNnL/n3zT+zu96dihNOe3sr11SP8YqboGK" +
	"tZ+Oj3kVx2C9HqHYBSo2Tgz6utY//V17Nb5ayzTQuXhenQ+k6KeSYnVHOe1g5qzn5dc1Iz/XOFlU4+EP5WWycs5f10Bn8ndH1SFD" +
	"56Yh+sxEYZ037NBzETpnv688K6c1RIr6R+qTza41sOvKG/leQ/NXgGtNz9nDDlcOBSru7tgooXrd0D5DFb/8antqojnv+HM3X64n" +
	"Z+0ppGm8hGp9Fqu1TjXF8a6ub+Rxnqmz//45uyCdR9fQPpx561Fd83DWo3Lq3DtE743W2t3pO0HXfIPUx8U0lpy5a4BiHa9tQtXu" +
	"IFfuKqz3ZGru8lObwlXPdF1DdNZ/5J1Ubr6L79r7mj1tD9d80FN96KV2eGl8O9d9XfWce49N+snrmzU3aqk6DZdnjhNN7Pi7OP5g" +
	"31veYYEiT+fLet7NkRviWifydvXr4drfDFNs5L2dNXBkvoisGjPElfuc8f3BhxuH6/6iv84x5fkOeXe7tJffWoS6nmOI0bUO6f+E" +
	"jocePkYf5/96zuOa2+RS/HxdtWOC2uWjNW6K+tIZG0KzrUtmH1flOJ0EWxKDa9w2a10FVG/Rs5nGS4byIdA1J/FyPYvgPE96Q9el" +
	"HD45e3KemruKqc75XGsezr6eE+9BrmdM0zX3J2gsO3zJ5dpDLq7tfXVc9FG7A11jtxO7JzKMDc783XnOzu+4bdcsas6nK6YBrnHj" +
	"WeG1rM/JbzNdOTmfyvFz1ccFXbWb+Eneeen1J7WYnB+oczP5vigpq7YO1f5vkhxiCfbFaYY34qPBuobYVZ9dWaLrh+H+/7i+LJ/p" +
	"JITehf9RF/c6s8ef1EKPqu89/8VaM7tcv39y3eNP5gMeDxln5fu0VJP7c7jOR6qvkjWGHjw76pm1vjOLGzunZ9VSPq6aM0PHEgef" +
	"Qa762tPVTv6enGH4HavrmnJtnEeWzADVw9kfyeGatzo1oxyn6f/fOuyMbTeO+Xu/SNuGBMr2QCO3VjYsnHmSn6u+dPJTZdd6S6Br" +
	"zc3JB/8PibShzA=="

const csTableGB18030_4B = "" +
	"eNrt3U1u00AUAODpG6etqFARsEeqVKReoHsOwAIJiQ0SJ+AA7NhyEs7SU5VESqGyjBPHP5nY3+KT69geu+95Zux0XP9MKZ0VJLbT" +
	"3LJObtnumMfc5Xg3qoZluVZubig/Rj7mGDEWQx53FHAORGH1Z6liB+fEfPIsL/PKZyWX8Nfrtcsey5faxn/QLgAAAAAAAAAAAAAA" +
	"AAAAAAAAAAAAAAAAAACwIL8eUvpx//goFgDQzdN771bp3zuku75Du6nMpvdPnyXv2ZvL+cLwcY0ecQ75a/WyNn810n5WtSnzaN92" +
	"1c+o5f1p/fOWclcD1eOh+tm849zNhecrCuvH9JXiCDC1fa9v78SqKB9r87cHlPHq2c9Vx2u/Q+UBl8fA15vHvBfLM7+GCddnAAAA" +
	"AAAAAMAMxqVGj7ErMdL+244neu4PoMTnujZj/d62jMN7MVIf8HmP8YdfGpb9b3zme7mFRY3Ld33lWTVguucy9lUVVr+jwzpZ26R/" +
	"cawA2qmF5k6eAQDA/QFiOue/Xcgb2gBQzwC0PXIIAACAezpA+wAA6Fvnnk/PVwPqNwAlt+n6CNQp9aLU9koeANBHnsbvEQ3Pm+cR" +
	"YzZW/LJzEwAA9yRHdy4/J+tSDPDdSVF+fxIDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
	"AAAAAAAAADgV369T+rqevjtg25u1quHzaNmmvizvmO9Tdlf5wO3CeQQAAAAAAAAAAAAAAAAAAABAofxvjHLjJTfTxHCqOK+2+8p7" +
	"7Lt6tl4ljwAAALBYm+8H3tQ+u9p+z3C9nYoTQ/omBot3MdF+/gDQNaRd"

const csTableBIG5 = "" +
	"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" +
	"z+MhW8V/00BemOufbQPoZKvsdQ9r/qHg1zGSa7eG+xvRNTV+5vbfetGF3efnvqDR9MesRl8/XeqSqiVv0s+r6ZZ+yjdKrorDfb3E" +
	"TLsGyOWtQv9dcg5FFl9Xu9SSew+bJvb2LOm37v1y/yb92yOpny8RfUL7YrjtflIUxnXfwQreJ0h8o1b6rGwHZ99FdqyWtQxn3QP7" +
	"qpdKjG81uOD25qtx6s912eD2Vy/2y26SK/VJ6eHbSQ5JQXSwUolnt0td2q3sp9cjaO5wL2ehNYh+G0lswHs1Ncn5A+vFDtKzTk8D" +
	"/ydlgu3TImcLuH3lPP1OOf/C61TcR+v6dqXo2ddSI+O9SSueCHVGfjb4m/R0/i2+kh1SZ9IkOnlezqh4hPn9NgrnvbYK3rqu6L0G" +
	"YvGtFeTMkUESg7B1LeAcEa8d6xV7r1dqZcvEb5wf0D+/VPJfCmI3Xsg5IavoZ7INv9J2OUfHv3f90Pjmsdn+NejeV2QeORbDeX1E" +
	"/PA4nrcT+Zqt0gfK6WqbxB4HQXO12RR3/Qy7N+D5i/H918DrvT5tCT6lVWJHTQX2W8QvPioTaqPeR6a2Cz/wM4Oqpf/wx8UX4bmP" +
	"+wvNWc+CMciSSeC46/r3AusO7OYZ6E4Ww/sp+1MhdQaTJca0SXKBJvDuPrK6OPgPOoDRJvEnlsp5Nu73HYU9XSa93V/k/Q3JucmL" +
	"blbgzJFR4GzzgP4hl5BL/Lj4pSoG5OO4vmZ8ZDq90LPi1ypIfrn7q3rJYagQWdEsdZRFkpPlvibj4w+L3vElsadd15oShVr6Vvbc" +
	"z9+pQq9ajm7bLvJrK3NsknoCr306GLnXCr8aLP0cusT/7DV9GyXG67Wrfsa0n0/l9DuHWPO+/G5+lsOjUFfvNlyr5GSUSMxiBz6O" +
	"z4heYEi9XeRPUSac77TBc4mRZV4H5jlWo4FPPXrQCmK8xdgSg1jrDHw83vfP4zxtIufnREHOdotvv0POWV8F/EYzD++VUYGs9Lis" +
	"+8nvIS+rTPDZfZ4N4jtxG9v26B3qN7wfm8u2DcB/vsQx3gM39Mxx9xl7P/VecuA6pNbZz+t1+zSmj025+JPyUhMQDcg5u4b72wf4" +
	"ptxPdIjU9JST77mPrOXsTP9zze37+6lV3SRnTpbSE+JZcPQDiedsx7e+Hl5m/+y6kwfkdLheuz+24QRwZhrw7YZONkF/BfahPBPO" +
	"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" +
	"LW8XCvcl95+WENwpKwuFwx/5P+Oc6P8CiYGiUA=="

const csTableEUCKR = "" +
	"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"
