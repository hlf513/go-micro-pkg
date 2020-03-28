package handler

import (
	"context"

	"github.com/micro/go-micro"

	// uuid "git.com/protobuf/uuid"
)

func Register(s micro.Service) {
	// _ = uuid.RegisterUuidHandler(s.Server(), Uuid{s})
}

type Uuid struct {
	s micro.Service
}

func (u Uuid) Generate() {
	// handler something
	
	
	// publish message
	pub := micro.NewPublisher("test", u.s.Client())
	pub.Publish(context.Background(), common.Event{
		Id:      1,
		Message: "msg",
	})
}
