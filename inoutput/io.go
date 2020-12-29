package inoutput

/*
 * 多对一
 * 给每一个资源一个uint16 ID，根据合成规则，将消耗的资源按照生序生成索引。
 * 树形结构
 */

type Resource struct {
	Name string
}

type IORule struct {
	Name string
	Src  []string
	Out  string
}

func NewIORule(name string, src []string, out string) *IORule {
	if len(src) == 0 {
		panic("src length is 0")
	}
	return &IORule{
		Name: name,
		Src:  src,
		Out:  out,
	}
}

func (this *IORule) check(in map[string]string) bool {
	for _, v := range this.Src {
		if _, ok := in[v.Name]; !ok {
			return false
		}
	}
	return true
}

type ctNode struct {
	Children map[int]*ctNode
	End      map[int]struct{}
}

func newCtNode() *ctNode {
	return &ctNode{
		Children: map[int]*ctNode{},
		End:      map[int]struct{}{},
	}
}

type outNode struct {
	out int
	Src [][]int
}

// IO 工具台
type IOCraftingTable struct {
	genID  int
	id2Res map[int]string
	res2ID map[string]int

	rules map[string]*IORule

	inRoot *ctNode
	outMap map[int]*outNode
}

func NewIOCraftingTable() *IOCraftingTable {
	return &IOCraftingTable{
		genID:  0,
		id2Res: map[int]string{},
		res2ID: map[string]int{},
		rules:  map[string]*IORule{},
		inRoot: newCtNode(),
		outMap: map[int]*outNode{},
	}
}

func (this *IOCraftingTable) foundOutSrc(out int) [][]int {
	if n, ok := this.outMap[out]; ok {
		return n.Src
	}
	return nil
}

func (this *IOCraftingTable)makeSrc(out int,srcSlice []int)*outNode{
	node := &outNode{
		out: out,
		Src: [][]int{srcSlice},
	}
	for _,src := range srcSlice{
		if this.foundOutSrc(src)
	}
}

func (this *IOCraftingTable) Register(rule *IORule) {
	//if _, ok := this.rules[rule.Name]; ok {
	//	panic(fmt.Sprintf("rule name %s is already register", rule.Name))
	//}
	//this.rules[rule.Name] = rule

	srcSlice := make([]int, 0, len(rule.Src))
	for _, name := range rule.Src {
		id, ok := this.res2ID[name]
		if !ok {
			this.genID++
			id = this.genID

			this.res2ID[name] = id
			this.id2Res[id] = name
		}
		srcSlice = append(srcSlice, id)
	}

	for _, src := range rule.Src {
		if node, ok := this.src[src.Name]; ok {
			if _, ok2 := node.Rules[rule.Name]; !ok2 {
				node.Rules[rule.Name] = rule
			}
		} else {
			this.src[src.Name] = &srcNode{
				Current: src,
				Rules: map[string]*IORule{
					rule.Name: rule,
				},
			}
		}
	}
}

func (this *IOCraftingTable) Compound(in ...*Resource) (map[string]*IORule, error) {
	items := map[string]*Resource{}
	for _, v := range in {
		items[v.Name] = v
	}

	ret := map[string]*IORule{}
	for _, v := range this.rules {
		if v.check(items) {
			ret[v.Name] = v
		}
	}
	return ret, nil
}

func dp() {

}

func intersectionM(m []map[string]*IORule) map[string]*IORule {
	if len(m) == 1 {
		return m[0]
	}

	ret := intersection(m[0], m[1])
	for i := 2; i < len(m); i++ {
		ret = intersection(ret, m[i])
	}
	return ret
}

func intersection(m1, m2 map[string]*IORule) (ret map[string]*IORule) {
	ret = map[string]*IORule{}
	for _, v := range m1 {
		if _, ok := m2[v.Name]; ok {
			ret[v.Name] = v
		}
	}
	return
}
