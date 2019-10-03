package hars

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sinmetal/fake/hook"
	"github.com/vvakame/go-harlog"
)

func Compare(t *testing.T, goldenPath string, har *harlog.HARContainer) {
	t.Helper()

	got, err := json.MarshalIndent(har, "", "  ")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	fn := filepath.Join("testdata", goldenPath)
	if *hook.Update {
		t.Logf("update %s", goldenPath)
		if err := ioutil.WriteFile(fn, got, 0644); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	golden, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	var want harlog.HARContainer
	if err := json.Unmarshal(golden, &want); err != nil {
		t.Fatal("unexpected error:", err)
	}

	compare(t, &want, har)
}

func compare(t *testing.T, want, got *harlog.HARContainer) {
	if e, g := len(want.Log.Entries), len(got.Log.Entries); e != g {
		t.Fatalf("want Entries.len is %v but got %v", e, g)
	}

	for i, v := range want.Log.Entries {
		v2 := got.Log.Entries[i]
		if e, g := v.Request.URL, v2.Request.URL; e != g {
			t.Errorf("want Request.URL is %v but got %v", e, g)
		}
		if e, g := v.Request.Method, v2.Request.Method; e != g {
			t.Errorf("want Request.Method is %v but got %v", e, g)
		}
		compareNVPArray(t, "Request.Headers", v.Request.Headers, v2.Request.Headers)
		compareNVPArray(t, "Request.QueryString", v.Request.QueryString, v2.Request.QueryString)

		cmp.Equal(v.Request.PostData, v2.Request.PostData)

		if e, g := v.Response.Status, v2.Response.Status; e != g {
			t.Errorf("want Response.Status is %v but got %v", e, g)
		}
		cmp.Equal(v.Response.Content, v2.Response.Content)
	}
}

func compareNVPArray(t *testing.T, title string, want, got []*harlog.NVP) {
	t.Helper()

	gm := nvpArrayToMap(got)

	for _, w := range want {
		if hook.IgnoreHeaderKey(w.Name) {
			continue
		}
		gv, ok := gm[w.Name]
		if !ok {
			t.Errorf("want %s %s but notfound", title, w.Name)
		}
		if w.Value != gv {
			t.Errorf("want %s[%s] is %v but got %v", title, w.Name, w.Value, gv)
		}
	}
}

func nvpArrayToMap(array []*harlog.NVP) map[string]string {
	m := make(map[string]string)
	for _, v := range array {
		m[v.Name] = v.Value
	}
	return m
}
