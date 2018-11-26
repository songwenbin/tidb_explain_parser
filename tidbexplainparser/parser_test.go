package tidbexplainparser

import (
	"fmt"
	"strconv"
	"testing"
)

func TestQuery(t *testing.T) {
	s := Query("explain ANALYZE select * from testlib join t1 on testlib.id = t1.id;")
	fmt.Println(s)
	/*
		for _, val := range s {
			fmt.Println(val.id)
			fmt.Println(val.task)
			fmt.Println(val.count)
			fmt.Println(val.operatorinfo)
			fmt.Println(val.executeinfo)
		}
		Assert(t, "ab", "cd")
	*/
}

func TestExplainReader(t *testing.T) {
	s := Query("explain ANALYZE select * from testlib join t1 on testlib.id = t1.id;")

	ExplainReader(s)
}

func Assert(t *testing.T, expect string, actual string) {
	if expect != actual {
		t.Fatal("expect: " + expect + " != actual: " + actual)
	}
}

func AssertInt(t *testing.T, expect int, actual int) {
	if expect != actual {
		t.Fatal("expect: " + strconv.Itoa(expect) + " != actual: " + strconv.Itoa(actual))
	}
}

func TestAssert(t *testing.T) {
	Assert(t, "ab", "cd")
}

func TestGetRowLayer(t *testing.T) {
	result0 := GetRowLayer("HashLeftJoin_6")
	AssertInt(t, 0, result0)

	result1 := GetRowLayer("├─TableReader_9")
	AssertInt(t, 1, result1)

	result2 := GetRowLayer("│ └─TableScan_8")
	AssertInt(t, 2, result2)
}

func TestHandleNodeForNextLayer(t *testing.T) {
	// 下一层测试
	var root Node = Node{Name: "root", Layer: 0}
	var new Node = Node{Name: "new", Layer: 1}
	var conx context = context{root: &root, current: &root, layer: 0}
	conx.HandleNode(&new)
	Assert(t, "new", root.Nodes[0].Name)
	Assert(t, "root", root.Nodes[0].parent.Name)
}

func TestHandleNodeForSameLayer(t *testing.T) {
	var root Node = Node{Name: "root", Layer: 0}
	var node1 Node = Node{Name: "node1", Layer: 1}
	var node2 Node = Node{Name: "node2", Layer: 2}
	var conx context = context{root: &root, current: &root, layer: 0}
	conx.HandleNode(&node1)
	conx.HandleNode(&node2)
	var new Node = Node{Name: "new", Layer: 1}
	conx.HandleNode(&new)
	Assert(t, "new", root.Nodes[1].Name)
	Assert(t, "root", root.Nodes[1].parent.Name)
}

func TestHandleNodeForPreLayer(t *testing.T) {
	var root Node = Node{Name: "root", Layer: 0}
	var node1 Node = Node{Name: "node1", Layer: 1}
	var node2 Node = Node{Name: "node2", Layer: 2}
	var node3 Node = Node{Name: "node3", Layer: 3}
	var conx context = context{root: &root, current: &root, layer: 0}
	conx.HandleNode(&node1)
	conx.HandleNode(&node2)
	conx.HandleNode(&node3)
	var new Node = Node{Name: "new", Layer: 1}
	conx.HandleNode(&new)
	Assert(t, "new", root.Nodes[1].Name)
	Assert(t, "node1", root.Nodes[0].Name)
	Assert(t, "root", root.Nodes[1].parent.Name)
}

func TestHandleNodeForPreLayer2(t *testing.T) {
	var root Node = Node{Name: "root", Layer: 0}
	var node1 Node = Node{Name: "node1", Layer: 1}
	var node2 Node = Node{Name: "node2", Layer: 2}
	var node3 Node = Node{Name: "node3", Layer: 1}
	var conx context = context{root: &root, current: &root, layer: 0}
	conx.HandleNode(&node1)
	conx.HandleNode(&node2)
	conx.HandleNode(&node3)

	var new Node = Node{Name: "new", Layer: 2}
	conx.HandleNode(&new)
	Assert(t, "node1", root.Nodes[0].Name)
	Assert(t, "node2", root.Nodes[0].Nodes[0].Name)
	Assert(t, "node3", root.Nodes[1].Name)
	Assert(t, "new", root.Nodes[1].Nodes[0].Name)
}
