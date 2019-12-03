package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type addressList struct{ *ec2.DescribeAddressesOutput }
type addressNodes []addressNode

type addressNode struct {
	UID              string          `json:"uid,omitempty"`
	Type             []string        `json:"dgraph.type,omitempty"`
	Name             string          `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID          string          `json:"OwnerId,omitempty"`
	OwnerName        string          `json:"OwnerName,omitempty"`
	Region           string          `json:"Region,omitempty"`
	Service          string          `json:"Service,omitempty"`
	Domain           string          `json:"Domain,omitempty"`
	PrivateIPAddress string          `json:"PrivateIpAddress,omitempty"`
	PublicIP         string          `json:"PublicIp,omitempty"`
	AllocationID     string          `json:"AllocationID,omitempty"`
	Instance         instanceNodes   `json:"_Instance,omitempty"`
	NatGateway       natGatewayNodes `json:"_NatGateway,omitempty"`
}

func (c *connector) listAddresses() addressList {
	defer c.waitGroup.Done()

	log.Println("List Addresses")
	response, err := ec2.New(c.awsSession).DescribeAddresses(&ec2.DescribeAddressesInput{})
	if err != nil {
		log.Fatal(err)
	}
	return addressList{response}
}

func (list addressList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Addresses) == 0 {
		return
	}
	log.Println("Add Address Nodes")
	a := make(addressNodes, 0, len(list.Addresses))

	for _, i := range list.Addresses {
		var b addressNode
		b.Service = "ec2"
		b.Type = []string{"Address"}
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Name = *i.PublicIp
		b.PublicIP = *i.PublicIp
		b.AllocationID = *i.AllocationId
		b.Domain = *i.Domain
		if i.PrivateIpAddress != nil {
			b.PrivateIPAddress = *i.PrivateIpAddress
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]addressNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("Address"), &m)
	for _, i := range m["list"] {
		n[i.PublicIP] = i.UID
	}
	c.ressources["Addresses"] = n
}

func (list addressList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Addresses) == 0 {
		return
	}
	log.Println("Add Address Edges")

	a := addressNodes{}
	for _, i := range list.Addresses {
		if i.AssociationId != nil {
			b := addressNode{
				UID:        c.ressources["Addresses"][*i.PublicIp],
				NatGateway: natGatewayNodes{natGatewayNode{UID: c.ressources["NatGateways"][*i.PublicIp]}},
			}
			if i.InstanceId != nil {
				b.Instance = instanceNodes{instanceNode{UID: c.ressources["Instances"][*i.InstanceId]}}
			}
			a = append(a, b)
		}
	}
	if len(a) != 0 {
		c.dgraphAddNodes(a)
	}
}
