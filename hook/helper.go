package hook

import (
	"flag"
	"strings"
)

var Update = flag.Bool("update", false, "update golden file")
var Generate = flag.Bool("gen", false, "code generate")

func IgnoreHeaderKey(key string) bool {
	keys := []string{"Expires", "Age", "X-GUploader-UploadID", "Alt-Svc", "Date", "X-Goog-Api-Client"}
	for _, v := range keys {
		if strings.ToLower(key) == strings.ToLower(v) {
			return true
		}
	}
	return false
}
