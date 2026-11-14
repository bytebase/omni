package deparse

// DefaultCollations holds the default collation for each charset
// MariaDB has. Measured on mariadb:11.8.8 through
// information_schema.CHARACTER_SETS: the Unicode charsets default to the
// UCA 14.0.0 collations, where MySQL 8.0 defaults utf8mb4 to
// utf8mb4_0900_ai_ci and the rest to <charset>_general_ci. MariaDB has no
// gb18030.
var DefaultCollations = map[string]string{
	"utf8mb4":  "utf8mb4_uca1400_ai_ci",
	"utf8mb3":  "utf8mb3_uca1400_ai_ci",
	"utf8":     "utf8mb3_uca1400_ai_ci",
	"ucs2":     "ucs2_uca1400_ai_ci",
	"utf16":    "utf16_uca1400_ai_ci",
	"utf32":    "utf32_uca1400_ai_ci",
	"utf16le":  "utf16le_general_ci",
	"latin1":   "latin1_swedish_ci",
	"latin2":   "latin2_general_ci",
	"latin5":   "latin5_turkish_ci",
	"latin7":   "latin7_general_ci",
	"ascii":    "ascii_general_ci",
	"binary":   "binary",
	"gbk":      "gbk_chinese_ci",
	"big5":     "big5_chinese_ci",
	"euckr":    "euckr_korean_ci",
	"gb2312":   "gb2312_chinese_ci",
	"sjis":     "sjis_japanese_ci",
	"ujis":     "ujis_japanese_ci",
	"cp932":    "cp932_japanese_ci",
	"eucjpms":  "eucjpms_japanese_ci",
	"cp1250":   "cp1250_general_ci",
	"cp1251":   "cp1251_general_ci",
	"cp1256":   "cp1256_general_ci",
	"cp1257":   "cp1257_general_ci",
	"cp850":    "cp850_general_ci",
	"cp852":    "cp852_general_ci",
	"cp866":    "cp866_general_ci",
	"dec8":     "dec8_swedish_ci",
	"hp8":      "hp8_english_ci",
	"keybcs2":  "keybcs2_general_ci",
	"koi8r":    "koi8r_general_ci",
	"koi8u":    "koi8u_general_ci",
	"macce":    "macce_general_ci",
	"macroman": "macroman_general_ci",
	"swe7":     "swe7_swedish_ci",
	"geostd8":  "geostd8_general_ci",
	"tis620":   "tis620_thai_ci",
	"hebrew":   "hebrew_general_ci",
	"greek":    "greek_general_ci",
	"armscii8": "armscii8_general_ci",
}
