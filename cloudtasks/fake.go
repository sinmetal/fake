package cloudtasks

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
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

	var mockCloudTasks mockCloudTasksServer

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

type mockCloudTasksServer struct {
	// Embed for forward compatibility.
	// Tests will keep working if more methods are added
	// in the future.
	taskspb.CloudTasksServer

	reqs []proto.Message

	// If set, all calls return this error.
	err error

	// responses to return if err == nil
	resps []proto.Message
}

func (s *mockCloudTasksServer) CreateTask(ctx context.Context, req *taskspb.CreateTaskRequest) (*taskspb.Task, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*taskspb.Task), nil
}
