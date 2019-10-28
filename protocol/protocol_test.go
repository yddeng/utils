package protocol_test

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dutil/protocol"
	"github.com/yddeng/dutil/protocol/pb"
	"github.com/yddeng/dutil/protocol/protobuf"
	"testing"
)

func TestInitProtocol(t *testing.T) {
	protocol.InitProtocol(protobuf.Protobuf{})

	protocol.Register(1, &pb.EchoToS{})
	protocol.Register(2, &pb.EchoToC{})

	id, data, err := protocol.Marshal(&pb.EchoToS{
		Msg: proto.String("test"),
	})
	fmt.Println(id, data, err)

	msg, err2 := protocol.Unmarshal(id, data)
	fmt.Println(msg, err2)
}
