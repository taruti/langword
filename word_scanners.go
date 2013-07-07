package langword

import (
	"unicode"
	"unicode/utf8"
	"code.google.com/p/go.text/unicode/norm"
)

func ScanLatinWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
        // Skip leading spaces.
        start := 0
        for width := 0; start < len(data); start += width {
                var r rune
                r, width = utf8.DecodeRune(data[start:])
                if unicode.IsLetter(r) {
                        break
                }
        }
        if atEOF && len(data) == 0 {
                return 0, nil, nil
        }
        // Scan until space, marking end of word.
        for width, i := 0, start; i < len(data); i += width {
                var r rune
                r, width = utf8.DecodeRune(data[i:])
		
			if !unicode.IsLetter(r) {
				return i + width, mutateWord(data[start:i]), nil
			} 
       }
        // If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
        if atEOF && len(data) > start {
			return len(data), mutateWord(data[start:]), nil
        }
        // Request more data.
        return 0, nil, nil
}

func mutateWord(bs []byte) []byte {
	return norm.NFC.Bytes(toLowerBytes(bs))
}

func toLowerBytes(src []byte) []byte {
	dst := make([]byte, 0, len(src))
	for width, i, di := 0,0,0; i<len(src); i+= width {
		var r rune
		r, width = utf8.DecodeRune(src[i:])
		r = unicode.ToLower(r)
		dw:=utf8.RuneLen(r)
		dst = append(dst, make([]byte, dw)...)
		utf8.EncodeRune(dst[di:], r)
		di += dw
	}
	return dst
}
