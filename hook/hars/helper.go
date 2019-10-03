package hars

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/sinmetal/fake/hook"
	"github.com/vvakame/go-harlog"
)

// LogFakeResponseCode is 指定したharのgolden fileを読み込んで、fakeなhttp.Responseを作るコードをログに吐き出す
func LogFakeResponseCode(t *testing.T, goldenPath string) {
	if !*hook.Generate {
		return
	}
	fn := filepath.Join("testdata", goldenPath)
	golden, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	var got harlog.HARContainer
	if err := json.Unmarshal(golden, &got); err != nil {
		t.Fatal("unexpected error:", err)
	}

	for _, entry := range got.Log.Entries {
		fmt.Println(`header := make(map[string][]string)`)
		for _, header := range entry.Response.Headers {
			fmt.Printf("header[\"%v\"] = []string{\"%v\"}\n", header.Name, header.Value)
		}

		fmt.Printf("r := ioutil.NopCloser(strings.NewReader(`%s`))", entry.Response.Content.Text)

		fmt.Printf(`
return &http.Response{
			Status:        "%v",
			StatusCode:    nil,
			Header:        header,
			Body:          r,
			ContentLength: %v,
		}
`, entry.Response.Status, entry.Response.Content.Size)

	}
}
