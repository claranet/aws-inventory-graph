package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/rds"
)

type dbClusterList struct{ *rds.DescribeDBClustersOutput }
type dbClusterNodes []dbClusterNode

type dbClusterNode struct {
	UID                              string                      `json:"uid,omitempty"`
	Type                             []string                    `json:"dgraph.type,omitempty"`
	Name                             string                      `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID                          string                      `json:"OwnerId,omitempty"`
	OwnerName                        string                      `json:"OwnerName,omitempty"`
	Region                           string                      `json:"Region,omitempty"`
	Service                          string                      `json:"Service,omitempty"`
	DBClusterArn                     string                      `json:"DbClusterArn,omitempty"`
	DBClusterIdentifier              string                      `json:"DbClusterIdentifier,omitempty"`
	DBClusterResourceID              string                      `json:"DbClusterResourceID,omitempty"`
	DBName                           string                      `json:"DbName,omitempty"`
	MasterUsername                   string                      `json:"MasterUsername,omitempty"`
	Engine                           string                      `json:"Engine,omitempty"`
	EngineMode                       string                      `json:"EngineMode,omitempty"`
	EngineVersion                    string                      `json:"EngineVersion,omitempty"`
	Endpoint                         string                      `json:"Endpoint,omitempty"`
	ReaderEndpoint                   string                      `json:"ReaderEndpoint,omitempty"`
	Port                             int64                       `json:"Port,omitempty"`
	MultiAZ                          bool                        `json:"MultiAz,omitempty"`
	HTTPEndpointEnabled              bool                        `json:"HttpEndpointEnabled,omitempty"`
	IAMDatabaseAuthenticationEnabled bool                        `json:"IamDatabaseAuthenticationEnabled,omitempty"`
	PreferredMaintenanceWindow       string                      `json:"PreferredMaintenanceWindow,omitempty"`
	DeletionProtection               bool                        `json:"DeletionProtection,omitempty"`
	HostedZoneID                     string                      `json:"HostedZoneId,omitempty"`
	Status                           string                      `json:"Status,omitempty"`
	PreferredBackupWindow            string                      `json:"PreferredBackupWindow,omitempty"`
	AllocatedStorage                 int64                       `json:"AllocatedStorage,omitempty"`
	BackupRetentionPeriod            int64                       `json:"BackupRetentionPeriod,omitempty"`
	StorageEncrypted                 bool                        `json:"StorageEncrypted,omitempty"`
	AvailabilityZone                 availabilityZoneNodes       `json:"_AvailabilityZone,omitempty"`
	DBClusterMember                  dbInstanceNodes             `json:"_DbClusterMember,omitempty"`
	ReadReplica                      dbInstanceNodes             `json:"_ReadReplica,omitempty"`
	DBSubnetGroup                    dbSubnetGroupNodes          `json:"_DbSubnetGroup,omitempty"`
	SecurityGroup                    securityGroupNodes          `json:"_SecurityGroup,omitempty"`
	DBClusterParameterGroup          dbClusterParameterGroupNode `json:"_DbClusterParameterGroup,omitempty"`
}

func (c *connector) listDbClusters() dbClusterList {
	defer c.waitGroup.Done()

	log.Println("List DbClusters")
	response, err := rds.New(c.awsSession).DescribeDBClusters(&rds.DescribeDBClustersInput{})
	if err != nil {
		log.Fatal(err)
	}
	return dbClusterList{response}
}

func (list dbClusterList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.DBClusters) == 0 {
		return
	}
	log.Println("Add DbCluster Nodes")
	a := make(dbClusterNodes, 0, len(list.DBClusters))

	for _, i := range list.DBClusters {
		var b dbClusterNode
		b.Service = "rds"
		b.Type = []string{"DbCluster"}
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Name = *i.DBClusterIdentifier
		b.DBClusterArn = *i.DBClusterArn
		b.DBClusterIdentifier = *i.DBClusterIdentifier
		b.DBClusterResourceID = *i.DbClusterResourceId
		b.MasterUsername = *i.MasterUsername
		b.Engine = *i.Engine
		b.EngineMode = *i.EngineMode
		b.EngineVersion = *i.EngineVersion
		b.Endpoint = *i.Endpoint
		b.ReaderEndpoint = *i.ReaderEndpoint
		b.Port = *i.Port
		b.MultiAZ = *i.MultiAZ
		b.HTTPEndpointEnabled = *i.HttpEndpointEnabled
		b.IAMDatabaseAuthenticationEnabled = *i.IAMDatabaseAuthenticationEnabled
		b.PreferredMaintenanceWindow = *i.PreferredMaintenanceWindow
		b.DeletionProtection = *i.IAMDatabaseAuthenticationEnabled
		b.HostedZoneID = *i.HostedZoneId
		b.Status = *i.Status
		b.AllocatedStorage = *i.AllocatedStorage
		b.BackupRetentionPeriod = *i.BackupRetentionPeriod
		b.StorageEncrypted = *i.StorageEncrypted
		b.DBName = *i.DatabaseName
		a = append(a, b)
	}
	c.dgraphAddNodes(a)

	m := make(map[string]dbClusterNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("DbCluster"), &m)
	for _, i := range m["list"] {
		n[i.DBClusterIdentifier] = i.UID
	}
	c.ressources["DbClusters"] = n
}

func (list dbClusterList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.DBClusters) == 0 {
		return
	}
	log.Println("Add DbCluster Edges")
	a := dbClusterNodes{}
	for _, i := range list.DBClusters {
		b := dbClusterNode{
			UID:                     c.ressources["DbClusters"][*i.DBClusterIdentifier],
			DBClusterParameterGroup: dbClusterParameterGroupNode{UID: c.ressources["DbParameterGroups"][*i.DBClusterParameterGroup]},
			DBSubnetGroup:           dbSubnetGroupNodes{dbSubnetGroupNode{UID: c.ressources["DbSubnetGroups"][*i.DBSubnetGroup]}},
		}
		for _, j := range i.AvailabilityZones {
			b.AvailabilityZone = append(b.AvailabilityZone, availabilityZoneNode{UID: c.ressources["AvailabilityZones"][*j]})
		}
		for _, j := range i.VpcSecurityGroups {
			b.SecurityGroup = append(b.SecurityGroup, securityGroupNode{UID: c.ressources["SecurityGroups"][*j.VpcSecurityGroupId]})
		}
		for _, j := range i.DBClusterMembers {
			b.DBClusterMember = append(b.DBClusterMember, dbInstanceNode{UID: c.ressources["DbInstances"][*j.DBInstanceIdentifier], IsClusterWriter: *j.IsClusterWriter})
		}
		for _, j := range i.ReadReplicaIdentifiers {
			s := strings.Split(*j, ":")
			b.ReadReplica = append(b.ReadReplica, dbInstanceNode{UID: c.ressources["DbInstances"][s[6]]})
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
