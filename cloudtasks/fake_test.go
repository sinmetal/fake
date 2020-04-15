package cloudtasks

import (
	"context"
	"fmt"
	"testing"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/golang/protobuf/proto"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

func TestCreateTask(t *testing.T) {
	var name string = "name3373707"
	var dispatchCount int32 = 1217252086
	var responseCount int32 = 424727441
	var expectedResponse = &taskspb.Task{
		Name:          name,
		DispatchCount: dispatchCount,
		ResponseCount: responseCount,
	}

	faker := NewFaker(t)

	faker.mock.err = nil
	faker.mock.resps = nil

	faker.mock.resps = append(faker.mock.resps[:0], expectedResponse)

	var formattedParent string = fmt.Sprintf("projects/%s/locations/%s/queues/%s", "[PROJECT]", "[LOCATION]", "[QUEUE]")
	var task *taskspb.Task = &taskspb.Task{}
	var request = &taskspb.CreateTaskRequest{
		Parent: formattedParent,
		Task:   task,
	}

	c, err := cloudtasks.NewClient(context.Background(), faker.ClientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateTask(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, faker.mock.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}
