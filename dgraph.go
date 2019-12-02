package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

func (c *connector) dgraphDropAll() {
	log.Println("Drop all previous data")
	op := api.Operation{DropAll: true}
	if err := c.dgraphClient.Alter(*c.context, &op); err != nil {
		log.Fatal(err)
	}
}

func (c *connector) dgraphDropPrevious() {
	log.Println("List previous data")
	txn := c.dgraphClient.NewTxn()
	defer txn.Discard(*c.context)

	q := `query query($owner: string, $region: string){
		list(func: eq(OwnerId, $owner)) @filter(eq(Region, $region)) {
			uid
		}
		}`

	res, err := txn.QueryWithVars(*c.context, q, map[string]string{"$owner": c.awsAccountID, "$region": c.awsRegion})
	if err != nil {
		log.Println(err.Error())
	}

	m := make(map[string]cidrNodes)
	json.Unmarshal(res.Json, &m)

	if len(m["list"]) != 0 {
		log.Println("Drop previous data")
		n, _ := json.Marshal(m["list"])
		mu := &api.Mutation{
			CommitNow:  true,
			DeleteJson: n,
		}

		_, err = txn.Mutate(*c.context, mu)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (c *connector) dgraphAddSchema() {
	log.Println("Add schema")
	op := &api.Operation{Schema: getdgraphSchema()}
	err := c.dgraphClient.Alter(*c.context, op)
	if err != nil {
		log.Fatal(err)
	}
}

func dgraphDisplaySchema(dgraph, dtype *string) {
	connexion, err := grpc.Dial(*dgraph, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}

	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(connexion))
	ctx := context.Background()
	resp, err := dgraphClient.NewReadOnlyTxn().Query(ctx, `schema(type: [`+*dtype+`]) {type}`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resp.Json))
}

func (c *connector) dgraphAddNodes(nodes interface{}) {
	txn := c.dgraphClient.NewTxn()
	defer txn.Discard(*c.context)
	n, err := json.Marshal(nodes)
	if err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   n,
	}
	_, err = txn.Mutate(*c.context, mu)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *connector) dgraphQuery(nodeType string) []byte {
	txn := c.dgraphClient.NewReadOnlyTxn()
	defer txn.Discard(*c.context)

	q := `query query($type: string){
			list(func: type($type)) {
				uid
				dgraph.type
				expand(_all_)
			}
	  	}`

	res, err := txn.QueryWithVars(*c.context, q, map[string]string{"$type": nodeType})
	if err != nil {
		log.Println(err.Error())
	}
	return res.Json
}
