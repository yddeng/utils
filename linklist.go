package dutil

type Node struct {
	data interface{}
	next *Node
}

type LinkList struct {
	size int
	head *Node
	end  *Node
}

func NewLinkList() *LinkList {
	return &LinkList{}
}

func (l *LinkList) Size() int {
	return l.size
}

func (l *LinkList) Add(data interface{}) *Node {
	node := &Node{
		data: data,
	}

	if l.end == nil {
		l.head = node
		l.end = node
	} else {
		l.end.next = node
		l.end = node
	}
	l.size++

	return node
}

func (l *LinkList) Delete(node *Node) {

	cNode := l.head
	befNode := cNode
	for cNode != nil && node != cNode {
		befNode = cNode
		cNode = cNode.next
	}
	if cNode != nil {
		befNode.next = cNode.next
		cNode.next = nil
		l.size--
	}

}

func (l *LinkList) Get() *Node {
	if l.size == 0 {
		return nil
	} else {
		node := l.head
		l.head = node.next
		l.size--
		if l.Size() == 0 {
			l.Clear()
		}
		return node
	}
}

func (l *LinkList) Clear() {
	l.size = 0
	l.head = nil
	l.end = nil
}
