package helpers

import (
	"container/list"
	"errors"
	"math"
)

type Value interface {
	Tag() interface{}
	Key() float64
}

type FibHeap struct {
	roots       *list.List
	index       map[interface{}]*node
	treeDegrees map[uint]*list.Element
	min         *node
	num         uint
}

type node struct {
	self     *list.Element
	parent   *node
	children *list.List
	marked   bool
	degree   uint
	position uint
	tag      interface{}
	key      float64
	value    Value
}

// NewFibHeap creates an initialized Fibonacci Heap.
func NewFibHeap() *FibHeap {
	heap := new(FibHeap)
	heap.roots = list.New()
	heap.index = make(map[interface{}]*node)
	heap.treeDegrees = make(map[uint]*list.Element)
	heap.num = 0
	heap.min = nil

	return heap
}

// Num returns the total number of values in the heap.
func (heap *FibHeap) Num() uint {
	return heap.num
}

func (heap *FibHeap) InsertValue(value Value) error {
	if value == nil {
		return errors.New("Input value is nil ")
	}

	return heap.insert(value.Tag(), value.Key(), value)
}

func (heap *FibHeap) MinimumValue() Value {
	if heap.num == 0 {
		return nil
	}

	return heap.min.value
}

func (heap *FibHeap) ExtractMinValue() Value {
	if heap.num == 0 {
		return nil
	}

	min := heap.extractMin()

	return min.value
}

func (heap *FibHeap) DecreaseKeyValue(value Value) error {
	if value == nil {
		return errors.New("Input value is nil ")
	}

	if math.IsInf(value.Key(), -1) {
		return errors.New("Negative infinity key is reserved for internal usage ")
	}

	if node, exists := heap.index[value.Tag()]; exists {
		return heap.decreaseKey(node, value, value.Key())
	}

	return errors.New("Value is not found ")
}

func (heap *FibHeap) IncreaseKeyValue(value Value) error {
	if value == nil {
		return errors.New("Input value is nil ")
	}

	if math.IsInf(value.Key(), -1) {
		return errors.New("Negative infinity key is reserved for internal usage ")
	}

	if node, exists := heap.index[value.Tag()]; exists {
		return heap.increaseKey(node, value, value.Key())
	}

	return errors.New("Value is not found ")
}

func (heap *FibHeap) Delete(tag interface{}) error {
	if tag == nil {
		return errors.New("Input tag is nil ")
	}

	if _, exists := heap.index[tag]; !exists {
		return errors.New("Tag is not found ")
	}

	heap.ExtractValue(tag)

	return nil
}

func (heap *FibHeap) GetValue(tag interface{}) (value Value) {
	if node, exists := heap.index[tag]; exists {
		value = node.value
	}

	return
}

func (heap *FibHeap) ExtractValue(tag interface{}) (value Value) {
	if node, exists := heap.index[tag]; exists {
		value = node.value
		heap.deleteNode(node)
		return
	}

	return nil
}

func (heap *FibHeap) consolidate() {
	for tree := heap.roots.Front(); tree != nil; tree = tree.Next() {
		heap.treeDegrees[tree.Value.(*node).position] = nil
	}

	for tree := heap.roots.Front(); tree != nil; {
		if heap.treeDegrees[tree.Value.(*node).degree] == nil {
			heap.treeDegrees[tree.Value.(*node).degree] = tree
			tree.Value.(*node).position = tree.Value.(*node).degree
			tree = tree.Next()
			continue
		}

		if heap.treeDegrees[tree.Value.(*node).degree] == tree {
			tree = tree.Next()
			continue
		}

		for heap.treeDegrees[tree.Value.(*node).degree] != nil {
			anotherTree := heap.treeDegrees[tree.Value.(*node).degree]
			heap.treeDegrees[tree.Value.(*node).degree] = nil
			if tree.Value.(*node).key <= anotherTree.Value.(*node).key {
				heap.roots.Remove(anotherTree)
				heap.link(tree.Value.(*node), anotherTree.Value.(*node))
			} else {
				heap.roots.Remove(tree)
				heap.link(anotherTree.Value.(*node), tree.Value.(*node))
				tree = anotherTree
			}
		}
		heap.treeDegrees[tree.Value.(*node).degree] = tree
		tree.Value.(*node).position = tree.Value.(*node).degree
	}

	heap.resetMin()
}

func (heap *FibHeap) insert(tag interface{}, key float64, value Value) error {
	if math.IsInf(key, -1) {
		return errors.New("Negative infinity key is reserved for internal usage ")
	}

	if _, exists := heap.index[tag]; exists {
		return errors.New("Duplicate tag is not allowed ")
	}

	node := new(node)
	node.children = list.New()
	node.tag = tag
	node.key = key
	node.value = value

	node.self = heap.roots.PushBack(node)
	heap.index[node.tag] = node
	heap.num++

	if heap.min == nil || heap.min.key > node.key {
		heap.min = node
	}

	return nil
}

func (heap *FibHeap) extractMin() *node {
	min := heap.min

	children := heap.min.children
	if children != nil {
		for e := children.Front(); e != nil; e = e.Next() {
			e.Value.(*node).parent = nil
			e.Value.(*node).self = heap.roots.PushBack(e.Value.(*node))
		}
	}

	heap.roots.Remove(heap.min.self)
	heap.treeDegrees[min.position] = nil
	delete(heap.index, heap.min.tag)
	heap.num--

	if heap.num == 0 {
		heap.min = nil
	} else {
		heap.consolidate()
	}

	return min
}

func (heap *FibHeap) deleteNode(n *node) {
	heap.decreaseKey(n, n.value, math.Inf(-1))
	heap.ExtractMinValue()
}

func (heap *FibHeap) link(parent, child *node) {
	child.marked = false
	child.parent = parent
	child.self = parent.children.PushBack(child)
	parent.degree++
}

func (heap *FibHeap) resetMin() {
	heap.min = heap.roots.Front().Value.(*node)
	for tree := heap.min.self.Next(); tree != nil; tree = tree.Next() {
		if tree.Value.(*node).key < heap.min.key {
			heap.min = tree.Value.(*node)
		}
	}
}

func (heap *FibHeap) decreaseKey(n *node, value Value, key float64) error {
	if key >= n.key {
		return errors.New("New key is not smaller than current key ")
	}

	n.key = key
	n.value = value
	if n.parent != nil {
		parent := n.parent
		if n.key < n.parent.key {
			heap.cut(n)
			heap.cascadingCut(parent)
		}
	}

	if n.parent == nil && n.key < heap.min.key {
		heap.min = n
	}

	return nil
}

func (heap *FibHeap) increaseKey(n *node, value Value, key float64) error {
	if key <= n.key {
		return errors.New("New key is not larger than current key ")
	}

	n.key = key
	n.value = value

	child := n.children.Front()
	for child != nil {
		childNode := child.Value.(*node)
		child = child.Next()
		if childNode.key < n.key {
			heap.cut(childNode)
			heap.cascadingCut(n)
		}
	}

	if heap.min == n {
		heap.resetMin()
	}

	return nil
}

func (heap *FibHeap) cut(n *node) {
	n.parent.children.Remove(n.self)
	n.parent.degree--
	n.parent = nil
	n.marked = false
	n.self = heap.roots.PushBack(n)
}

func (heap *FibHeap) cascadingCut(n *node) {
	if n.parent != nil {
		if !n.marked {
			n.marked = true
		} else {
			parent := n.parent
			heap.cut(n)
			heap.cascadingCut(parent)
		}
	}
}
