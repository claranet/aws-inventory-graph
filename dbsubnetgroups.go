package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/rds"
)

type dbSubnetGroupList struct {
	*rds.DescribeDBSubnetGroupsOutput
}
type dbSubnetGroupNodes []dbSubnetGroupNode

type dbSubnetGroupNode struct {
	UID               string      `json:"uid,omitempty"`
	Type              []string    `json:"dgraph.type,omitempty"`
	Name              string      `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID           string      `json:"OwnerId,omitempty"`
	OwnerName         string      `json:"OwnerName,omitempty"`
	Region            string      `json:"Region,omitempty"`
	Service           string      `json:"Service,omitempty"`
	DBSubnetGroupArn  string      `json:"DbSubnetGroupArn,omitempty"`
	DBSubnetGroupName string      `json:"DbSubnetGroupName,omitempty"`
	SubnetGroupStatus string      `json:"SubnetGroupStatus,omitempty"`
	Description       string      `json:"Description,omitempty"`
	Vpc               vpcNode     `json:"_Vpc,omitempty"`
	Subnet            subnetNodes `json:"_Subnet,omitempty"`
}

func (c *connector) listDbSubnetGroups() dbSubnetGroupList {
	defer c.waitGroup.Done()

	log.Println("List DbSubnetGroups")
	response, err := rds.New(c.awsSession).DescribeDBSubnetGroups(&rds.DescribeDBSubnetGroupsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return dbSubnetGroupList{response}
}

func (list dbSubnetGroupList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.DBSubnetGroups) == 0 {
		return
	}
	log.Println("Add DbSubnetGroup Nodes")
	a := make(dbSubnetGroupNodes, 0, len(list.DBSubnetGroups))

	for _, i := range list.DBSubnetGroups {
		var b dbSubnetGroupNode
		b.Service = "rds"
		b.Type = []string{"DbSubnetGroup"}
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Name = *i.DBSubnetGroupName
		b.DBSubnetGroupName = *i.DBSubnetGroupName
		b.DBSubnetGroupArn = *i.DBSubnetGroupArn
		b.SubnetGroupStatus = *i.SubnetGroupStatus
		b.Description = *i.DBSubnetGroupDescription
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]dbSubnetGroupNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("DbSubnetGroup"), &m)
	for _, i := range m["list"] {
		n[i.DBSubnetGroupName] = i.UID
	}
	c.ressources["DbSubnetGroups"] = n
}

func (list dbSubnetGroupList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.DBSubnetGroups) == 0 {
		return
	}
	log.Println("Add DbSubnetGroup Edges")
	a := dbSubnetGroupNodes{}
	for _, i := range list.DBSubnetGroups {
		b := dbSubnetGroupNode{
			UID: c.ressources["DbSubnetGroups"][*i.DBSubnetGroupName],
		}
		for _, j := range i.Subnets {
			b.Subnet = append(b.Subnet, subnetNode{UID: c.ressources["Subnets"][*j.SubnetIdentifier]})
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
