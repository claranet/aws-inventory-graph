package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type natGatewayList struct {
	*ec2.DescribeNatGatewaysOutput
}
type natGatewayNodes []natGatewayNode

type natGatewayNode struct {
	UID              string      `json:"uid,omitempty"`
	Type             []string    `json:"dgraph.type,omitempty"`
	Name             string      `json:"name,omitempty"` // This field is only for Ratel Viz
	Region           string      `json:"Region,omitempty"`
	OwnerName        string      `json:"OwnerName,omitempty"`
	OwnerID          string      `json:"OwnerId,omitempty"`
	Service          string      `json:"Service,omitempty"`
	NatGatewayID     string      `json:"NatGatewayId,omitempty"`
	PublicIPAddress  string      `json:"PublicIpAddress,omitempty"`
	PrivateIPAddress string      `json:"PrivateIpAddress,omitempty"`
	State            string      `json:"State,omitempty"`
	Vpc              vpcNode     `json:"_Vpc,omitempty"`
	Subnet           subnetNodes `json:"_Subnet,omitempty"`
}

func (c *connector) listNatGateways() natGatewayList {
	defer c.waitGroup.Done()

	log.Println("List NatGateways")
	response, err := ec2.New(c.awsSession).DescribeNatGateways(&ec2.DescribeNatGatewaysInput{})
	if err != nil {
		log.Fatal(err)
	}
	return natGatewayList{response}
}

func (list natGatewayList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.NatGateways) == 0 {
		return
	}
	log.Println("Add NatGateway Nodes")
	a := make(natGatewayNodes, 0, len(list.NatGateways))

	for _, i := range list.NatGateways {
		var b natGatewayNode
		b.Service = "ec2"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"NatGateway"}
		b.Name = *i.NatGatewayId
		for _, tag := range i.Tags {
			if *tag.Key == "Name" {
				b.Name = *tag.Value
			}
		}
		b.NatGatewayID = *i.NatGatewayId
		b.State = *i.State
		b.PublicIPAddress = *i.NatGatewayAddresses[0].PublicIp
		b.PrivateIPAddress = *i.NatGatewayAddresses[0].PrivateIp
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]natGatewayNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("NatGateway"), &m)
	for _, i := range m["list"] {
		n[i.PublicIPAddress] = i.UID
	}
	c.ressources["NatGateways"] = n
}

func (list natGatewayList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.NatGateways) == 0 {
		return
	}
	log.Println("Add NatGateway Edges")
	a := natGatewayNodes{}
	for _, i := range list.NatGateways {
		b := natGatewayNode{
			UID:    c.ressources["NatGateways"][*i.NatGatewayAddresses[0].PublicIp],
			Subnet: subnetNodes{subnetNode{UID: c.ressources["Subnets"][*i.SubnetId]}},
			Vpc:    vpcNode{UID: c.ressources["Vpcs"][*i.VpcId]},
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
