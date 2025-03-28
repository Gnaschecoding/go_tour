package word

import (
	"strings"
	"unicode"
)

func ToUpper(str string) string {
	return strings.ToUpper(str)
}

func ToLower(str string) string {
	return strings.ToLower(str)
}

func UnderscoreToUpperCamelCase(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	str = strings.Replace(str, " ", "", -1)
	return str
}

func UnderscoreToLowerCamelCase(str string) string {
	str = UnderscoreToUpperCamelCase(str)
	return string(unicode.ToLower(rune(str[0]))) + str[1:]
}

func CamelCaseToUnderscore(str string) string {
	var out []rune
	for k, v := range str {
		if k == 0 {
			out = append(out, unicode.ToLower(v))
			continue
		}

		if unicode.IsUpper(v) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(v))
	}
	return string(out)
}
