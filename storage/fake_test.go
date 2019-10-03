package storage_test

import (
	"context"
	"fmt"
	"github.com/sinmetal/fake/hook/hars"
	"io/ioutil"
	"net/http"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/sinmetal/fake/hook"
	"github.com/vvakame/go-harlog"
	"google.golang.org/api/option"

	. "github.com/sinmetal/fake/storage"
)

func TestGetObject(t *testing.T) {
	ctx := context.Background()

	faker := NewFaker(t)

	stg, err := storage.NewClient(ctx, option.WithHTTPClient(faker.Client))
	if err != nil {
		t.Fatal(err)
	}

	reader, err := stg.Bucket("hoge").Object("hoge.txt").NewReader(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", string(body))
}

// TestRealGetObject is 実際にCloud StorageにGetを投げてResponseの内容に変更がないかをチェックする
func TestRealGetObject(t *testing.T) {
	ctx := context.Background()

	hooker := hook.NewHooker(t)
	stg, err := storage.NewClient(ctx, option.WithHTTPClient(hooker.Client))
	if err != nil {
		t.Fatal(err)
	}

	_, err = stg.Bucket("hoge").Object("hoge.txt").NewReader(ctx)
	if err != nil {
		t.Fatal(err)
	}

	req := hooker.GetRequest()
	hook.CompareHookRequest(t, "object.get.request.golden", req)

	resp := hooker.GetResponse()
	hook.CompareHookResponse(t, "object.get.response.golden", resp)
}

func TestRealGetObjectHar(t *testing.T) {
	ctx := context.Background()

	har := &harlog.Transport{}
	hc := &http.Client{
		Transport: har,
	}

	stg, err := storage.NewClient(ctx, option.WithHTTPClient(hc))
	if err != nil {
		t.Fatal(err)
	}

	_, err = stg.Bucket("hoge").Object("hoge.txt").NewReader(ctx)
	if err != nil {
		t.Fatal(err)
	}

	hars.Compare(t, "object.get.har.golden", har.HAR())
}
