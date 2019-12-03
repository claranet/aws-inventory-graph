package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type subnetList struct {
	*ec2.DescribeSubnetsOutput
}
type subnetNodes []subnetNode

type subnetNode struct {
	UID                         string                `json:"uid,omitempty"`
	Type                        []string              `json:"dgraph.type,omitempty"`
	Name                        string                `json:"name,omitempty"` // This field is only for Ratel Viz
	Region                      string                `json:"Region,omitempty"`
	OwnerID                     string                `json:"OwnerId,omitempty"`
	OwnerName                   string                `json:"OwnerName,omitempty"`
	Service                     string                `json:"Service,omitempty"`
	SubnetID                    string                `json:"SubnetId,omitempty"`
	CidrBlock                   string                `json:"CidrBlock,omitempty"`
	MapPublicIPOnLaunch         bool                  `json:"MapPublicIpOnLaunch,omitempty"`
	DefaultForAz                bool                  `json:"MapPuDefaultForAzblicIpOnLaunch,omitempty"`
	State                       string                `json:"State,omitempty"`
	AssignIPv6AddressOnCreation bool                  `json:"AssignIpv6AddressOnCreation,omitempty"`
	AvailableIPAddressCount     int64                 `json:"AvailableIpAddressCount,omitempty"`
	AvailabilityZone            availabilityZoneNodes `json:"_AvailabilityZone,omitempty"`
	Vpc                         vpcNode               `json:"_Vpc,omitempty"`
}

func (c *connector) listSubnets() subnetList {
	defer c.waitGroup.Done()

	log.Println("List Subnets")
	response, err := ec2.New(c.awsSession).DescribeSubnets(&ec2.DescribeSubnetsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return subnetList{response}
}

func (list subnetList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Subnets) == 0 {
		return
	}
	log.Println("Add Subnet Nodes")
	a := make(subnetNodes, 0, len(list.Subnets))

	for _, i := range list.Subnets {
		var b subnetNode
		b.Service = "ec2"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"Subnet"}
		b.Name = *i.SubnetId
		for _, tag := range i.Tags {
			if *tag.Key == "Name" {
				b.Name = *tag.Value
			}
		}
		b.SubnetID = *i.SubnetId
		b.MapPublicIPOnLaunch = *i.MapPublicIpOnLaunch
		b.AvailableIPAddressCount = *i.AvailableIpAddressCount
		b.DefaultForAz = *i.DefaultForAz
		b.State = *i.State
		b.CidrBlock = *i.CidrBlock
		b.AssignIPv6AddressOnCreation = *i.AssignIpv6AddressOnCreation
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]subnetNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("Subnet"), &m)
	for _, i := range m["list"] {
		n[i.SubnetID] = i.UID
	}
	c.ressources["Subnets"] = n
}

func (list subnetList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Subnets) == 0 {
		return
	}
	log.Println("Add Subnet Edges")
	a := subnetNodes{}
	for _, i := range list.Subnets {
		b := subnetNode{
			UID:              c.ressources["Subnets"][*i.SubnetId],
			Vpc:              vpcNode{UID: c.ressources["Vpcs"][*i.VpcId]},
			AvailabilityZone: availabilityZoneNodes{availabilityZoneNode{UID: c.ressources["AvailabilityZones"][*i.AvailabilityZone]}},
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
