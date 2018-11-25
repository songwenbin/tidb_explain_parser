package tidbexplainparser

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func explainReader(input string) {
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type QueryResult struct {
	id           string
	task         string
	count        string
	operatorinfo string
	executeinfo  string
}

func Query(query string) []QueryResult {
	fmt.Println("test")
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:4000)/test?charset=utf8")
	checkErr(err)
	defer db.Close()

	fmt.Println(db)

	rows, err := db.Query(query)
	checkErr(err)

	var result []QueryResult = []QueryResult{}

	for rows.Next() {
		r := QueryResult{}
		err = rows.Scan(&r.id, &r.count, &r.task, &r.operatorinfo, &r.executeinfo)
		checkErr(err)
		result = append(result, r)
	}

	return result
}

var conx context = context{}

func ExplainReader(t []QueryResult) {

	for _, val := range t {
		layer := GetRowLayer(val.id)
		n := Node{name: val.id, layer: layer}
		conx.HandleNode(&n)
	}
}

type Node struct {
	name   string
	layer  int
	nodes  []*Node
	parent *Node
}

type context struct {
	root    *Node
	current *Node
	layer   int
}

func (c *context) HandleNode(new *Node) {

	if c.layer < new.layer {
		new.parent = c.current
		c.current.nodes = []*Node{}
		c.current.nodes = append(c.current.nodes, new)
		c.current = new
		c.layer = new.layer
	} else {
		loop := c.layer - new.layer + 1
		for i := 0; i < loop; i++ {
			c.current = c.current.parent
		}
		new.parent = c.current
		c.current.nodes = append(c.current.nodes, new)
		c.layer = new.layer
	}
}

func GetRowLayer(content string) int {
	var total int = 0
	for i, _ := range content {
		if content[i] < 'A' || content[i] > 'Z' {
			total++
		} else {
			break
		}
	}

	return total / 2
}
