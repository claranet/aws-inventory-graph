package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/rds"
)

type dbParameterGroupList struct {
	*rds.DescribeDBParameterGroupsOutput
}
type dbParameterGroupNodes []dbParameterGroupNode

type dbParameterGroupNode struct {
	UID                    string   `json:"uid,omitempty"`
	Type                   []string `json:"dgraph.type,omitempty"`
	Name                   string   `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID                string   `json:"OwnerId,omitempty"`
	OwnerName              string   `json:"OwnerName,omitempty"`
	Region                 string   `json:"Region,omitempty"`
	Service                string   `json:"Service,omitempty"`
	DBParameterGroupArn    string   `json:"DbParameterGroupArn,omitempty"`
	DBParameterGroupName   string   `json:"DbParameterGroupName,omitempty"`
	DBParameterGroupFamily string   `json:"DbParameterGroupFamily,omitempty"`
	Description            string   `json:"Description,omitempty"`
}

func (c *connector) listDbParameterGroups() dbParameterGroupList {
	defer c.waitGroup.Done()

	log.Println("List DbParameterGroups")
	response, err := rds.New(c.awsSession).DescribeDBParameterGroups(&rds.DescribeDBParameterGroupsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return dbParameterGroupList{response}
}

func (list dbParameterGroupList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.DBParameterGroups) == 0 {
		return
	}
	log.Println("Add DbParameterGroup Nodes")
	a := make(dbParameterGroupNodes, 0, len(list.DBParameterGroups))

	for _, i := range list.DBParameterGroups {
		var b dbParameterGroupNode
		b.Service = "rds"
		b.Type = []string{"DbParameterGroup"}
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Name = *i.DBParameterGroupName
		b.DBParameterGroupName = *i.DBParameterGroupName
		b.DBParameterGroupArn = *i.DBParameterGroupArn
		b.DBParameterGroupFamily = *i.DBParameterGroupFamily
		b.Description = *i.Description
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]dbParameterGroupNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("DbParameterGroup"), &m)
	for _, i := range m["list"] {
		n[i.DBParameterGroupName] = i.UID
	}
	c.ressources["DbParameterGroups"] = n
}
