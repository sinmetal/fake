package cloudtasks

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/golang/protobuf/proto"
	_ "github.com/sinmetal/fake/hook"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

func TestCreateTask(t *testing.T) {
	cases := []struct {
		name      string
		callCount int
	}{
		{"one", 1},
		{"two", 2},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			faker := NewFaker(t)
			var expectedResponses []*taskspb.Task
			for i := 0; i < tt.callCount; i++ {
				var name string = fmt.Sprintf("name%d", rand.Int())
				var dispatchCount int32 = 1217252086
				var responseCount int32 = 424727441
				var expectedResponse = &taskspb.Task{
					Name:          name,
					DispatchCount: dispatchCount,
					ResponseCount: responseCount,
				}
				faker.AddMockResponse(nil, expectedResponse)
				expectedResponses = append(expectedResponses, expectedResponse)
			}

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

			for i := 0; i < tt.callCount; i++ {
				resp, err := c.CreateTask(context.Background(), request)
				if err != nil {
					t.Fatal(err)
				}

				if e, g := request, faker.mock.tasks[i].reqs[0]; !proto.Equal(e, g) {
					t.Errorf("request want %q, but got %q", e, g)
				}

				if e, g := expectedResponses[i], resp; !proto.Equal(e, g) {
					t.Errorf("response want %q, but got %q)", e, g)
				}
			}

			if e, g := tt.callCount, faker.mock.createTaskCallCount; e != g {
				t.Errorf("createTaskCallCount want %v but got %v", e, g)
			}
		})
	}
}
