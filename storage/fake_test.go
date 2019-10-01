package storage_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/sinmetal/fake"
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

	hooker := fake.NewHooker(t)
	stg, err := storage.NewClient(ctx, option.WithHTTPClient(hooker.Client))
	if err != nil {
		t.Fatal(err)
	}

	_, err = stg.Bucket("hoge").Object("hoge.txt").NewReader(ctx)
	if err != nil {
		t.Fatal(err)
	}

	req := hooker.GetRequest()
	fake.CompareHookRequest(t, "object.get.request.golden", req)

	resp := hooker.GetResponse()
	fake.CompareHookResponse(t, "object.get.response.golden", resp)
}
