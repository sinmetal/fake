package hook

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func CompareHookRequest(t *testing.T, goldenPath string, req *HookRequest) {
	t.Helper()

	got, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	fn := filepath.Join("testdata", goldenPath)
	if *Update {
		t.Logf("update %s", goldenPath)
		if err := ioutil.WriteFile(fn, got, 0644); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}
	golden, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if !bytes.Equal(got, golden) {
		t.Errorf("want HookRequest : %s but got %s", golden, got)
	}
}

func CompareHookResponse(t *testing.T, goldenPath string, resp *HookResponse) {
	t.Helper()

	got, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	fn := filepath.Join("testdata", goldenPath)
	if *Update {
		t.Logf("update %s", goldenPath)
		if err := ioutil.WriteFile(fn, got, 0644); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}
	golden, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	var want HookResponse
	if err := json.Unmarshal(golden, &want); err != nil {
		t.Fatal("unexpected error:", err)
	}

	compareHeader(t, want.Header, resp.Header)
	if !cmp.Equal(want.Body, resp.Body) {
		t.Errorf("want HookResponse.Body : %+v but got %+v", string(want.Body), string(resp.Body))
	}
}

func compareHeader(t *testing.T, want, got map[string][]string) {
	for k, v := range want {
		if IgnoreHeaderKey(k) {
			continue
		}
		gv, ok := got[k]
		if !ok {
			t.Errorf("want %s but notfound", k)
		}
		if !cmp.Equal(v, gv) {
			t.Errorf("want %s is %v but got %v", k, v, gv)
		}
	}
}
