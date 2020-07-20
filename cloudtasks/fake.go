package cloudtasks

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/api/option"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Faker struct {
	mock      *mockCloudTasksServer
	ClientOpt option.ClientOption
}

func NewFaker(t *testing.T) *Faker {
	t.Helper()

	mockCloudTasks := mockCloudTasksServer{}

	serv := grpc.NewServer()
	taskspb.RegisterCloudTasksServer(serv, &mockCloudTasks)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	go serv.Serve(lis)

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return &Faker{
		mock:      &mockCloudTasks,
		ClientOpt: option.WithGRPCConn(conn),
	}
}

// AddMockResponse is Set Response when CreateTask is called
//
// Setting err returns that err.
// When returning a normal Response, err specifies nil.
func (f *Faker) AddMockResponse(err error, resp ...proto.Message) {
	f.mock.tasks = append(f.mock.tasks, &mockTaskContainer{
		err:  err,
		resp: resp,
	})
}

// GetCreateTaskCallCount is Returns the number of times CreateTask was called
func (f *Faker) GetCreateTaskCallCount() int {
	return f.mock.createTaskCallCount
}

// GetCreateTaskRequest is Returns the request passed to CreateTask
func (f *Faker) GetCreateTaskRequest(i int) ([]proto.Message, error) {
	if i > len(f.mock.tasks)-1 {
		return nil, fmt.Errorf("GetCreateTaskRequest out of range. len=%d", len(f.mock.tasks))
	}
	return f.mock.tasks[i].reqs, nil
}

type mockTaskContainer struct {
	reqs []proto.Message
	err  error
	resp []proto.Message
}

type mockCloudTasksServer struct {
	// Embed for forward compatibility.
	// Tests will keep working if more methods are added
	// in the future.
	taskspb.CloudTasksServer

	tasks []*mockTaskContainer

	tasksIndex int

	createTaskCallCount int
}

func (s *mockCloudTasksServer) CreateTask(ctx context.Context, req *taskspb.CreateTaskRequest) (*taskspb.Task, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}

	task := s.tasks[s.tasksIndex]
	s.tasksIndex++
	s.createTaskCallCount++

	task.reqs = append(task.reqs, req)
	if task.err != nil {
		return nil, task.err
	}

	return task.resp[0].(*taskspb.Task), nil
}
