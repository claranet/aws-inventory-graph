package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type vpcList struct {
	*ec2.DescribeVpcsOutput
}
type vpcNodes []vpcNode

type vpcNode struct {
	UID       string   `json:"uid,omitempty"`
	Type      []string `json:"dgraph.type,omitempty"`
	Name      string   `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID   string   `json:"OwnerId,omitempty"`
	OwnerName string   `json:"OwnerName,omitempty"`
	Region    string   `json:"Region,omitempty"`
	Service   string   `json:"Service,omitempty"`
	VpcID     string   `json:"VpcId,omitempty"`
	CidrBlock string   `json:"CidrBlock,omitempty"`
}

func (c *connector) listVpcs() vpcList {
	defer c.waitGroup.Done()

	log.Println("List VPCs")
	response, err := ec2.New(c.awsSession).DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return vpcList{response}
}

func (list vpcList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Vpcs) == 0 {
		return
	}
	log.Println("Add VPC Nodes")
	a := make(vpcNodes, 0, len(list.Vpcs))

	for _, i := range list.Vpcs {
		var b vpcNode
		b.Service = "ec2"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"Vpc"}
		b.Name = *i.VpcId
		for _, tag := range i.Tags {
			if *tag.Key == "Name" {
				b.Name = *tag.Value
			}
		}
		b.VpcID = *i.VpcId
		a = append(a, b)
	}
	c.dgraphAddNodes(a)

	m := make(map[string]vpcNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("Vpc"), &m)
	for _, i := range m["list"] {
		n[i.VpcID] = i.UID
	}
	c.ressources["Vpcs"] = n
}
