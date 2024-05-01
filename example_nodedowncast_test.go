package webidl_test

import "fmt"

type Node interface {
	NodeValue() *string
	SetNodeValue(v *string)
	sealedNode()
}

type Element interface {
	Node
	TagName() string
	LocalName() string
	sealedElement()
}

type node struct {
	nodeValue   *string
}

var _ Node = (*node)(nil)

func (n *node) NodeValue() (string, bool) {
	return n.nodeValue, n.nodeValueOk
}

func (n *node) SetNodeValue(v string, ok bool) {
	if ok {
		n.nodeValue = v
		n.nodeValueOk = true
	} else {
		n.nodeValue = ""
		n.nodeValueOk = false
	}
}

func (n *node) sealedNode() {}

type element struct {
	node
	tagName   string
	localName string
}

var _ Element = (*element)(nil)

func (e *element) TagName() string {
	return e.tagName
}

func (e *element) LocalName() string {
	return e.localName
}

func (e *element) sealedElement() {}

func getNode() Node {
	return &element{
		node: node{
			nodeValue:   "",
			nodeValueOk: false,
		},
		tagName:   "div",
		localName: "div",
	}
}

func Example() {
	e := getNode().(Element)
	fmt.Println(e.TagName())
	nodeValue, nodeValueOk := e.NodeValue()
	fmt.Printf("%#+v %#+v", nodeValue, nodeValueOk)
	// Output:
	// div
	// "" false
}
