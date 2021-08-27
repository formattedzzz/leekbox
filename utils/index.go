package utils

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
)

var FuncMapUnion = template.FuncMap{
	"prefix": func(str string) string {
		return "prefix-something-" + str
	},
}

func MD5(raw string) string {
	d := []byte(raw)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func If(ok bool, a, b interface{}) interface{} {
	if ok {
		return a
	} else {
		return b
	}
}
