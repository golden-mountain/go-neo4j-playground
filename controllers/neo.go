package controllers

import (
	// "fmt"
	"github.com/jmcvetta/neoism"
	"github.com/astaxie/beego"
)

// NeoController operations for Neo
type NeoController struct {
	beego.Controller
}

// URLMapping ...
func (c *NeoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Neo
// @Param	body		body 	models.Neo	true		"body for Neo content"
// @Success 201 {object} models.Neo
// @Failure 403 body is empty
// @router / [post]
func (c *NeoController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Neo by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Neo
// @Failure 403 :id is empty
// @router /:id [get]
func (c *NeoController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Neo
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Neo
// @Failure 403
// @router / [get]
func (c *NeoController) GetAll() {
	// No error handling in this example - bad, bad, bad!
	//
	// Connect to the Neo4j server
	//
	db, _ := neoism.Connect("http://neo4j:123456@192.168.95.104:7474/db/data")
	kirk := "Captain Kirk"
	mccoy := "Dr McCoy"
	//
	// Create a node
	//
	n0, _ := db.CreateNode(neoism.Props{"name": kirk})
	defer n0.Delete()  // Deferred clean up
	n0.AddLabel("Person") // Add a label
	//
	// Create a node with a Cypher query
	//
	res0 := []struct {
		N neoism.Node // Column "n" gets automagically unmarshalled into field N
	}{}
	cq0 := neoism.CypherQuery{
		Statement: "CREATE (n:Person {name: {name}}) RETURN n",
		// Use parameters instead of constructing a query string
		Parameters: neoism.Props{"name": mccoy},
		Result:     &res0,
	}
	db.Cypher(&cq0)
	n1 := res0[0].N // Only one row of data returned
	n1.Db = db // Must manually set Db with objects returned from Cypher query
	//
	// Create a relationship
	//
	n1.Relate("reports to", n0.Id(), neoism.Props{}) // Empty Props{} is okay
	//
	// Issue a query
	//
	res1 := []struct {
		A   string `json:"a.name"` // `json` tag matches column name in query
		Rel string `json:"type(r)"`
		B   string `json:"b.name"`
	}{}
	cq1 := neoism.CypherQuery{
		// Use backticks for long statements - Cypher is whitespace indifferent
		Statement: `
			MATCH (a:Person)-[r]->(b)
			WHERE a.name = {name}
			RETURN a.name, type(r), b.name
		`,
		Parameters: neoism.Props{"name": mccoy},
		Result:     &res1,
	}
	db.Cypher(&cq1)
	r := res1[0]
	beego.Debug(r.A, r.Rel, r.B)
	//
	// Clean up using a transaction
	//
	qs := []*neoism.CypherQuery{
		&neoism.CypherQuery{
			Statement: `
				MATCH (n:Person)-[r]->()
				WHERE n.name = {name}
				DELETE r
			`,
			Parameters: neoism.Props{"name": mccoy},
		},
		&neoism.CypherQuery{
			Statement: `
				MATCH (n:Person)
				WHERE n.name = {name}
				DELETE n
			`,
			Parameters: neoism.Props{"name": mccoy},
			
		},
	}
	tx, _ := db.Begin(qs)
	tx.Commit()
}

// Put ...
// @Title Put
// @Description update the Neo
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Neo	true		"body for Neo content"
// @Success 200 {object} models.Neo
// @Failure 403 :id is not int
// @router /:id [put]
func (c *NeoController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Neo
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *NeoController) Delete() {

}
