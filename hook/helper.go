package hook

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var Update = flag.Bool("update", false, "update golden file")

func IgnoreHeaderKey(key string) bool {
	keys := []string{"Expires", "Age", "X-GUploader-UploadID"}
	for _, v := range keys {
		if key == v {
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

	compResp := HookResponse{
		Header: copyUnstableHeader(t, want.Header, resp.Header),
		Body:   resp.Body,
	}

	if !cmp.Equal(want, compResp) {
		t.Errorf("want HookResponse : %+v but got %+v", want, resp)
	}
}

// copyUnstableHeader is 動的に変わる値を期待する値で埋めてしまう
func copyUnstableHeader(t *testing.T, want, got map[string][]string) map[string][]string {
	copyKeys := []string{"Expires", "Age", "X-GUploader-UploadID"}
	for _, key := range copyKeys {
		_, ok := want[key]
		if ok {
			got[key] = want[key]
		}
	}
	return got
}
