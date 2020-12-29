package inoutput

import (
	"fmt"
	"testing"
)

func TestIOCraftingTable_Compound(t *testing.T) {
	ct := NewIOCraftingTable()
	ct.Register(NewIORule("1->2", []*Resource{{Name: "1"}}, &Resource{Name: "2"}))
	ct.Register(NewIORule("1->3", []*Resource{{Name: "1"}}, &Resource{Name: "3"}))
	ct.Register(NewIORule("1->4", []*Resource{{Name: "1"}}, &Resource{Name: "4"}))
	ct.Register(NewIORule("1,4->5", []*Resource{{Name: "1"}, {Name: "4"}}, &Resource{Name: "5"}))

	out, _ := ct.Compound(&Resource{Name: "1"})
	for _, v := range out {
		fmt.Println(v.Name, v.Out)
	}

	fmt.Println()
	out2, _ := ct.Compound(&Resource{Name: "1"}, &Resource{Name: "4"})
	for _, v := range out2 {
		fmt.Println(v.Name, v.Out)
	}

	fmt.Println()
	out3, _ := ct.Compound(&Resource{Name: "4"})
	for _, v := range out3 {
		fmt.Println(v.Name, v.Out)
	}

}
