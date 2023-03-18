package utils

import (
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
)

var Json = func() jsoniter.API {
	// https://github.com/json-iterator/go/issues/244
	jsoniter.RegisterExtension(&BinaryAsStringExtension{})
	return jsoniter.ConfigCompatibleWithStandardLibrary
}()

func Debug(l ...interface{}) {
	if os.Getenv("DEBUG") == "1" {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)

		prefix := "[shadiaosocketio]"

		log.Println(append([]interface{}{prefix}, l...)...)
	}
}
