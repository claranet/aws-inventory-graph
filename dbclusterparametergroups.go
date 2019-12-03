package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/rds"
)

type dbClusterParameterGroupList struct {
	*rds.DescribeDBClusterParameterGroupsOutput
}
type dbClusterParameterGroupNodes []dbClusterParameterGroupNode

type dbClusterParameterGroupNode struct {
	UID                         string   `json:"uid,omitempty"`
	Type                        []string `json:"dgraph.type,omitempty"`
	Name                        string   `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID                     string   `json:"OwnerId,omitempty"`
	OwnerName                   string   `json:"OwnerName,omitempty"`
	Region                      string   `json:"Region,omitempty"`
	Service                     string   `json:"Service,omitempty"`
	DBClusterParameterGroupArn  string   `json:"DbClusterParameterGroupArn,omitempty"`
	DBClusterParameterGroupName string   `json:"DbClusterParameterGroupName,omitempty"`
	DBParameterGroupFamily      string   `json:"DbClusterParameterGroupFamily,omitempty"`
	Description                 string   `json:"Description,omitempty"`
}

func (c *connector) listDbClusterParameterGroups() dbClusterParameterGroupList {
	defer c.waitGroup.Done()

	log.Println("List DbClusterParameterGroups")
	response, err := rds.New(c.awsSession).DescribeDBClusterParameterGroups(&rds.DescribeDBClusterParameterGroupsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return dbClusterParameterGroupList{response}
}

func (list dbClusterParameterGroupList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.DBClusterParameterGroups) == 0 {
		return
	}
	log.Println("Add ParameterGroup Nodes")
	a := make(dbClusterParameterGroupNodes, 0, len(list.DBClusterParameterGroups))

	for _, i := range list.DBClusterParameterGroups {
		var b dbClusterParameterGroupNode
		b.Service = "rds"
		b.Type = []string{"DbClusterParameterGroup"}
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Name = *i.DBClusterParameterGroupName
		b.DBClusterParameterGroupName = *i.DBClusterParameterGroupName
		b.DBClusterParameterGroupArn = *i.DBClusterParameterGroupArn
		b.DBParameterGroupFamily = *i.DBParameterGroupFamily
		b.Description = *i.Description
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]dbClusterParameterGroupNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("DbClusterParameterGroup"), &m)
	for _, i := range m["list"] {
		n[i.DBClusterParameterGroupName] = i.UID
	}
	c.ressources["DbClusterParameterGroups"] = n
}
