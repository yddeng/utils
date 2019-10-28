package tool_test

import (
	"fmt"
	"github.com/yddeng/dutil/protocol/tool"
	"os"
	"testing"
)

func TestMake(t *testing.T) {
	os.MkdirAll("../message", os.ModePerm)
	tool.GenProto("../proto/message", "echo")
	fmt.Println("gen proto ok")
}
