package tidbexplainparser

import (
	"database/sql"
	"encoding/json"
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

	var conx context = context{}
	for _, val := range t {
		layer := GetRowLayer(val.id)
		info := Info{
			Task:         val.task,
			Count:        val.count,
			Operatorinfo: val.operatorinfo,
			Executeinfo:  val.executeinfo,
		}
		n := Node{Name: val.id, Layer: layer, Info: info}
		if layer == 0 {
			conx.layer = 0
			n.parent = nil
			conx.root = &n
			conx.current = conx.root
		} else {
			conx.HandleNode(&n)
		}
	}

	result, err := json.Marshal(&conx.root)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(result))
}

type Info struct {
	Task         string
	Count        string
	Operatorinfo string
	Executeinfo  string
}

type Node struct {
	Name   string
	Layer  int
	Nodes  []*Node
	parent *Node
	Info   Info
}

type context struct {
	root    *Node
	current *Node
	layer   int
}

func (c *context) HandleNode(new *Node) {
	if c.layer < new.Layer {
		new.parent = c.current
		c.current.Nodes = []*Node{}
		c.current.Nodes = append(c.current.Nodes, new)
		c.current = new
		c.layer = new.Layer
	} else {
		loop := c.layer - new.Layer + 1
		for i := 0; i < loop; i++ {
			c.current = c.current.parent
		}
		new.parent = c.current
		c.current.Nodes = append(c.current.Nodes, new)
		c.layer = new.Layer
		c.current = new
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
