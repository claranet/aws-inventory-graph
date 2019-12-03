package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/elasticache"
)

type cacheSubnetGroupList struct {
	*elasticache.DescribeCacheSubnetGroupsOutput
}
type cacheSubnetGroupNodes []cacheSubnetGroupNode

type cacheSubnetGroupNode struct {
	UID                  string      `json:"uid,omitempty"`
	Type                 []string    `json:"dgraph.type,omitempty"`
	Name                 string      `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID              string      `json:"OwnerId,omitempty"`
	OwnerName            string      `json:"OwnerName,omitempty"`
	Region               string      `json:"Region,omitempty"`
	Service              string      `json:"Service,omitempty"`
	CacheSubnetGroupName string      `json:"CacheSubnetGroupName,omitempty"`
	Description          string      `json:"Description,omitempty"`
	Vpc                  vpcNode     `json:"_Vpc,omitempty"`
	Subnet               subnetNodes `json:"_Subnet,omitempty"`
}

func (c *connector) listCacheSubnetGroups() cacheSubnetGroupList {
	defer c.waitGroup.Done()

	log.Println("List CacheSubnetGroups")
	response, err := elasticache.New(c.awsSession).DescribeCacheSubnetGroups(&elasticache.DescribeCacheSubnetGroupsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return cacheSubnetGroupList{response}
}

func (list cacheSubnetGroupList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.CacheSubnetGroups) == 0 {
		return
	}
	log.Println("Add CacheSubnetGroup Nodes")
	a := make(cacheSubnetGroupNodes, 0, len(list.CacheSubnetGroups))

	for _, i := range list.CacheSubnetGroups {
		var b cacheSubnetGroupNode
		b.Service = "elasticache"
		b.Type = []string{"CacheSubnetGroup"}
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Name = *i.CacheSubnetGroupName
		b.CacheSubnetGroupName = *i.CacheSubnetGroupName
		b.Description = *i.CacheSubnetGroupDescription
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]cacheSubnetGroupNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("CacheSubnetGroup"), &m)
	for _, i := range m["list"] {
		n[i.CacheSubnetGroupName] = i.UID
	}
	c.ressources["CacheSubnetGroups"] = n
}

func (list cacheSubnetGroupList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.CacheSubnetGroups) == 0 {
		return
	}
	log.Println("Add CacheSubnetGroup Edges")
	a := cacheSubnetGroupNodes{}
	for _, i := range list.CacheSubnetGroups {
		b := cacheSubnetGroupNode{
			UID: c.ressources["CacheSubnetGroups"][*i.CacheSubnetGroupName],
			Vpc: vpcNode{UID: c.ressources["Vpcs"][*i.VpcId]},
		}
		for _, j := range i.Subnets {
			b.Subnet = append(b.Subnet, subnetNode{UID: c.ressources["Subnets"][*j.SubnetIdentifier]})
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
