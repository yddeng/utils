package wordsFilter

import (
	"strings"
	"unicode/utf8"
)

var wordsFilter *filter

type filter struct {
	Root   *filterNode
	Ignore string
}

type filterNode struct {
	Children map[rune]*filterNode
	End      bool
}

func NewFilter(ignore, words string) *filter {
	filter := &filter{
		Root:   newFilterNode(),
		Ignore: ignore,
	}

	filter.newFilter(strings.Split(words, ";"))

	return filter
}

func newFilterNode() *filterNode {
	n := &filterNode{}
	n.Children = map[rune]*filterNode{}
	return n
}

func (this *filter) newFilter(s []string) {
	for _, _s := range s {
		if _s != "" {
			this.inster(_s)
		}
	}
}

func (this *filter) inster(txt string) {
	if len(txt) < 1 {
		return
	}
	node := this.Root
	key := []rune(txt)
	for i := 0; i < len(key); i++ {
		if _, b := node.Children[key[i]]; !b {
			node.Children[key[i]] = newFilterNode()
		}
		node = node.Children[key[i]]
	}

	node.End = true
}

func eliminate(str string, e string) string {
	str_ := []byte(str)
	e_ := []byte(e)
	out := []byte{}

	f := func(c byte) bool {
		for _, v := range e_ {
			if c == v {
				return false
			}
		}
		return true
	}

	for _, v := range str_ {
		if f(v) {
			out = append(out, v)
		}
	}
	return string(out)
}

func (this *filter) check(txt string) bool {
	if len(txt) < 1 {
		return false
	}

	txt = eliminate(txt, this.Ignore)
	node := this.Root
	key := []rune(txt)

	for i := 0; i < len(key); i++ {
		for j := i; j < len(key); j++ {
			if node, b := node.Children[key[j]]; !b {
				break
			} else {
				if node.End == true {
					return true
				}
			}
			node = node.Children[key[j]]
		}
		node = this.Root
	}
	return false
}

func Check(txt string) bool {
	return wordsFilter.check(txt)
}

func (this *filter) replace(txt string) string {
	if len(txt) < 1 {
		return txt
	}

	node := this.Root
	key := []rune(txt)
	var chars []rune = key
	c, _ := utf8.DecodeRuneInString("*")

	for i := 0; i < len(key); i++ {
		if node1_, b := node.Children[key[i]]; b {
			if node1_.End {
				chars[i] = c
				continue
			}
			node = node1_
			for j := i + 1; j < len(key); j++ {
				key_ := eliminate(string(key[j]), this.Ignore)
				if key_ == "" {
					continue
				} else {
					if node2_, bb := node.Children[[]rune(key_)[0]]; bb {
						if node2_.End {
							for t := i; t <= j; t++ {
								chars[t] = c
							}
							i = j
							//break
						}
						node = node2_
					} else {
						break
					}
				}
			}
			node = this.Root
		}
	}

	return string(chars)
}

func Replace(txt string) string {
	return wordsFilter.replace(txt)
}

func init() {
	ignore := "!@#$%^&*()_+/*-="
	words := "毛泽东;毛主席;狗屎;臭狗屎;fuck;fuck you;"
	wordsFilter = NewFilter(ignore, words)
}
