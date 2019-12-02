package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type vpcPeeringConnectionList struct {
	*ec2.DescribeVpcPeeringConnectionsOutput
}
type vpcPeeringConnectionNodes []vpcPeeringConnectionNode

type vpcPeeringConnectionNode struct {
	UID                    string   `json:"uid,omitempty"`
	Type                   []string `json:"dgraph.type,omitempty"`
	Name                   string   `json:"name,omitempty"` // This field is only for Ratel Viz
	Region                 string   `json:"Region,omitempty"`
	OwnerName              string   `json:"OwnerName,omitempty"`
	OwnerID                string   `json:"OwnerId,omitempty"`
	Service                string   `json:"Service,omitempty"`
	VpcPeeringConnectionID string   `json:"VpcPeeringConnectionId,omitempty"`
	Status                 string   `json:"Status,omitempty"`
	AccepterVpc            vpcNode  `json:"_AccepterVpc,omitempty"`
	RequesterVpc           vpcNode  `json:"_RequesterVpc,omitempty"`
}

func (c *connector) listVpcPeeringConnections() vpcPeeringConnectionList {
	defer c.waitGroup.Done()

	log.Println("List VpcPeeringConnections")
	response, err := ec2.New(c.awsSession).DescribeVpcPeeringConnections(&ec2.DescribeVpcPeeringConnectionsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return vpcPeeringConnectionList{response}
}

func (list vpcPeeringConnectionList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.VpcPeeringConnections) == 0 {
		return
	}
	log.Println("Add VpcPeeringConnection Nodes")
	a := make(vpcPeeringConnectionNodes, 0, len(list.VpcPeeringConnections))

	for _, i := range list.VpcPeeringConnections {
		var b vpcPeeringConnectionNode
		b.Service = "ec2"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"VpcPeeringConnection"}
		b.Name = *i.VpcPeeringConnectionId
		for _, tag := range i.Tags {
			if *tag.Key == "Name" {
				b.Name = *tag.Value
			}
		}
		b.VpcPeeringConnectionID = *i.VpcPeeringConnectionId
		b.Status = *i.Status.Code
		a = append(a, b)
	}
	c.dgraphAddNodes(a)

	m := make(map[string]vpcPeeringConnectionNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("VpcPeeringConnection"), &m)
	for _, i := range m["list"] {
		n[i.VpcPeeringConnectionID] = i.UID
	}
	c.ressources["VpcPeeringConnections"] = n
}

func (list vpcPeeringConnectionList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.VpcPeeringConnections) == 0 {
		return
	}
	log.Println("Add VpcPeeringConnection Edges")
	a := vpcPeeringConnectionNodes{}
	for _, i := range list.VpcPeeringConnections {
		b := vpcPeeringConnectionNode{
			UID:          c.ressources["VpcPeeringConnections"][*i.VpcPeeringConnectionId],
			AccepterVpc:  vpcNode{UID: c.ressources["Vpcs"][*i.AccepterVpcInfo.VpcId]},
			RequesterVpc: vpcNode{UID: c.ressources["Vpcs"][*i.RequesterVpcInfo.VpcId]},
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
