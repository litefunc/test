package internal

import (
	"fmt"
	"test/logger"
)

type Node struct {
	valid    bool
	B        bool
	hasItems bool
	Items    map[string]*Node
}

func NewNode(hasItems bool) *Node {

	n := Node{
		valid:    true,
		hasItems: hasItems,
	}
	if hasItems {
		n.Items = make(map[string]*Node)
	}
	return &n
}

func (rec *Node) Valid() bool {
	return rec.valid
}

func (rec *Node) get(id string) *Node {
	if !rec.valid {
		logger.Error("invalid")
		return &Node{}
	}
	if !rec.hasItems {
		logger.Error("is leaf")
		return &Node{}
	}

	for k, v := range rec.Items {
		if k == id {
			return v
		}
	}
	logger.Error(id, "not found")
	return &Node{}
}

func (rec *Node) Get(id ...string) *Node {
	n := rec
	for _, v := range id {
		n = n.get(v)
	}
	return n
}

func (rec *Node) Set(key string, v *Node) error {
	if !rec.hasItems {
		err := fmt.Errorf(`is leaf`)
		logger.Error(err)
		return err
	}
	rec.Items[key] = v
	return nil
}
