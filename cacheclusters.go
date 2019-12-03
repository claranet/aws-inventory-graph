package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/elasticache"
)

type cacheClusterList struct {
	*elasticache.DescribeCacheClustersOutput
}
type cacheClusterNodes []cacheClusterNode

type cacheClusterNode struct {
	UID                        string                `json:"uid,omitempty"`
	Type                       []string              `json:"dgraph.type,omitempty"`
	Name                       string                `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID                    string                `json:"OwnerId,omitempty"`
	OwnerName                  string                `json:"OwnerName,omitempty"`
	Region                     string                `json:"Region,omitempty"`
	Service                    string                `json:"Service,omitempty"`
	CacheClusterID             string                `json:"CacheClusterId,omitempty"`
	Engine                     string                `json:"Engine,omitempty"`
	EngineVersion              string                `json:"EngineVersion,omitempty"`
	AuthTokenEnabled           bool                  `json:"AuthTokenEnabled,omitempty"`
	Endpoint                   string                `json:"Endpoint,omitempty"`
	Port                       int64                 `json:"Port,omitempty"`
	AutoMinorVersionUpgrade    bool                  `json:"AutoMinorVersionUpgrade,omitempty"`
	Status                     string                `json:"CacheClusterStatus,omitempty"`
	NumCacheNodes              int64                 `json:"NumCacheNodes,omitempty"`
	TransitEncryptionEnabled   bool                  `json:"TransitEncryptionEnabled,omitempty"`
	PreferredMaintenanceWindow string                `json:"PreferredMaintenanceWindow,omitempty"`
	CacheNodeType              string                `json:"CacheNodeType,omitempty"`
	AvailabilityZone           availabilityZoneNodes `json:"_AvailabilityZone,omitempty"`
	SecurityGroup              securityGroupNodes    `json:"_SecurityGroup,omitempty"`
	CacheSubnetGroup           cacheSubnetGroupNodes `json:"_CacheSubnetGroup,omitempty"`
	// CacheParameterGroup        cacheParameterGroupNode `json:"_CacheParameterGroup,omitempty"`
}

func (c *connector) listCacheClusters() cacheClusterList {
	defer c.waitGroup.Done()

	log.Println("List CacheClusters")
	response, err := elasticache.New(c.awsSession).DescribeCacheClusters(&elasticache.DescribeCacheClustersInput{})
	if err != nil {
		log.Fatal(err)
	}
	return cacheClusterList{response}
}

func (list cacheClusterList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.CacheClusters) == 0 {
		return
	}
	log.Println("Add CacheCluster Nodes")
	a := make(cacheClusterNodes, 0, len(list.CacheClusters))

	for _, i := range list.CacheClusters {
		var b cacheClusterNode
		b.Service = "elasticache"
		b.Type = []string{"CacheCluster"}
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Name = *i.CacheClusterId
		b.CacheClusterID = *i.CacheClusterId
		b.Engine = *i.Engine
		b.EngineVersion = *i.EngineVersion
		b.AuthTokenEnabled = *i.AuthTokenEnabled
		if i.ConfigurationEndpoint != nil {
			b.Endpoint = *i.ConfigurationEndpoint.Address
			b.Port = *i.ConfigurationEndpoint.Port
		}
		b.AutoMinorVersionUpgrade = *i.AutoMinorVersionUpgrade
		b.Status = *i.CacheClusterStatus
		b.NumCacheNodes = *i.NumCacheNodes
		b.TransitEncryptionEnabled = *i.TransitEncryptionEnabled
		b.PreferredMaintenanceWindow = *i.PreferredMaintenanceWindow
		b.CacheNodeType = *i.CacheNodeType
		a = append(a, b)
	}
	c.dgraphAddNodes(a)

	m := make(map[string]cacheClusterNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("CacheCluster"), &m)
	for _, i := range m["list"] {
		n[i.CacheClusterID] = i.UID
	}
	c.ressources["CacheClusters"] = n
}

func (list cacheClusterList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.CacheClusters) == 0 {
		return
	}
	log.Println("Add CacheCluster Edges")
	a := cacheClusterNodes{}
	for _, i := range list.CacheClusters {
		b := cacheClusterNode{
			UID:              c.ressources["CacheClusters"][*i.CacheClusterId],
			CacheSubnetGroup: cacheSubnetGroupNodes{cacheSubnetGroupNode{UID: c.ressources["CacheSubnetGroups"][*i.CacheSubnetGroupName]}},
			AvailabilityZone: availabilityZoneNodes{availabilityZoneNode{UID: c.ressources["AvailabilityZones"][*i.PreferredAvailabilityZone]}},
		}
		for _, j := range i.SecurityGroups {
			b.SecurityGroup = append(b.SecurityGroup, securityGroupNode{UID: c.ressources["SecurityGroups"][*j.SecurityGroupId]})
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)
}
