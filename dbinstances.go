package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/rds"
)

type dbInstanceList struct{ *rds.DescribeDBInstancesOutput }
type dbInstanceNodes []dbInstanceNode

type dbInstanceNode struct {
	UID                              string                `json:"uid,omitempty"`
	Type                             []string              `json:"dgraph.type,omitempty"`
	Name                             string                `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID                          string                `json:"OwnerId,omitempty"`
	OwnerName                        string                `json:"OwnerName,omitempty"`
	Region                           string                `json:"Region,omitempty"`
	Service                          string                `json:"Service,omitempty"`
	DBInstanceArn                    string                `json:"DbInstanceArn,omitempty"`
	DBInstanceIdentifier             string                `json:"DbInstanceIdentifier,omitempty"`
	DBName                           string                `json:"DbName,omitempty"`
	DBInstanceStatus                 string                `json:"DbInstanceStatus,omitempty"`
	DBInstanceClass                  string                `json:"DbInstanceClass,omitempty"`
	DbiResourceID                    string                `json:"DbiResourceId,omitempty"`
	CACertificateIdentifier          string                `json:"CACertificateIdentifier,omitempty"`
	PubliclyAccessible               bool                  `json:"PubliclyAccessible,omitempty"`
	MasterUsername                   string                `json:"MasterUsername,omitempty"`
	LicenseModel                     string                `json:"LicenseModel,omitempty"`
	CopyTagsToSnapshot               bool                  `json:"CopyTagsToSnapshot,omitempty"`
	Engine                           string                `json:"Engine,omitempty"`
	EngineVersion                    string                `json:"EngineVersion,omitempty"`
	Endpoint                         string                `json:"Endpoint,omitempty"`
	Port                             int64                 `json:"Port,omitempty"`
	MultiAZ                          bool                  `json:"MultiAz,omitempty"`
	AutoMinorVersionUpgrade          bool                  `json:"AutoMinorVersionUpgrade,omitempty"`
	IAMDatabaseAuthenticationEnabled bool                  `json:"IamDatabaseAuthenticationEnabled,omitempty"`
	DeletionProtection               bool                  `json:"DeletionProtection,omitempty"`
	PreferredBackupWindow            string                `json:"PreferredBackupWindow,omitempty"`
	PreferredMaintenanceWindow       string                `json:"PreferredMaintenanceWindow,omitempty"`
	PromotionTier                    int64                 `json:"PromotionTier,omitempty"`
	AllocatedStorage                 int64                 `json:"AllocatedStorage,omitempty"`
	StorageType                      string                `json:"StorageType,omitempty"`
	StorageEncrypted                 bool                  `json:"StorageEncrypted,omitempty"`
	BackupRetentionPeriod            int64                 `json:"BackupRetentionPeriod,omitempty"`
	AvailabilityZone                 availabilityZoneNodes `json:"_AvailabilityZone,omitempty"`
	SecondaryAvailabilityZone        availabilityZoneNodes `json:"_SecondaryAvailabilityZone,omitempty"`
	SecurityGroup                    securityGroupNodes    `json:"_SecurityGroup,omitempty"`
	Vpc                              vpcNode               `json:"_Vpc,omitempty"`
	OptionGroup                      optionGroupNode       `json:"_OptionGroup,omitempty"`
	IsClusterWriter                  bool                  `json:"_DbClusterMember|IsClusterWriter,omitempty"`
	DBSubnetGroup                    dbSubnetGroupNodes    `json:"_DbSubnetGroup,omitempty"`
	DBParameterGroup                 dbParameterGroupNode  `json:"_DbParameterGroup,omitempty"`
}

func (c *connector) listDbInstances() dbInstanceList {
	defer c.waitGroup.Done()

	log.Println("List DbInstances")
	response, err := rds.New(c.awsSession).DescribeDBInstances(&rds.DescribeDBInstancesInput{})
	if err != nil {
		log.Fatal(err)
	}
	return dbInstanceList{response}
}

func (list dbInstanceList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.DBInstances) == 0 {
		return
	}
	log.Println("Add Dbinstance Nodes")
	a := make(dbInstanceNodes, 0, len(list.DBInstances))

	for _, i := range list.DBInstances {
		var b dbInstanceNode
		b.Service = "rds"
		b.Type = []string{"DbInstance"}
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Name = *i.DBInstanceIdentifier
		b.DBInstanceArn = *i.DBInstanceArn
		b.DBInstanceIdentifier = *i.DBInstanceIdentifier
		if i.DBName != nil {
			b.DBName = *i.DBName
		}
		b.DBInstanceStatus = *i.DBInstanceStatus
		b.DBInstanceClass = *i.DBInstanceClass
		b.DbiResourceID = *i.DbiResourceId
		b.CACertificateIdentifier = *i.CACertificateIdentifier
		b.PubliclyAccessible = *i.PubliclyAccessible
		b.MasterUsername = *i.MasterUsername
		b.LicenseModel = *i.LicenseModel
		b.CopyTagsToSnapshot = *i.CopyTagsToSnapshot
		b.Engine = *i.Engine
		b.EngineVersion = *i.EngineVersion
		b.Endpoint = *i.Endpoint.Address
		b.Port = *i.Endpoint.Port
		b.MultiAZ = *i.MultiAZ
		b.AutoMinorVersionUpgrade = *i.AutoMinorVersionUpgrade
		b.IAMDatabaseAuthenticationEnabled = *i.IAMDatabaseAuthenticationEnabled
		b.DeletionProtection = *i.IAMDatabaseAuthenticationEnabled
		b.PreferredBackupWindow = *i.PreferredBackupWindow
		b.PreferredMaintenanceWindow = *i.PreferredMaintenanceWindow
		if i.PromotionTier != nil {
			b.PromotionTier = *i.PromotionTier
		}
		b.AllocatedStorage = *i.AllocatedStorage
		b.StorageType = *i.StorageType
		b.StorageEncrypted = *i.StorageEncrypted
		b.BackupRetentionPeriod = *i.BackupRetentionPeriod
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]dbInstanceNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("DbInstance"), &m)
	for _, i := range m["list"] {
		n[i.DBInstanceIdentifier] = i.UID
	}
	c.ressources["DbInstances"] = n
}

func (list dbInstanceList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.DBInstances) == 0 {
		return
	}
	log.Println("Add Dbinstance Edges")
	a := dbInstanceNodes{}
	for _, i := range list.DBInstances {
		b := dbInstanceNode{
			UID:              c.ressources["DbInstances"][*i.DBInstanceIdentifier],
			AvailabilityZone: availabilityZoneNodes{availabilityZoneNode{UID: c.ressources["AvailabilityZones"][*i.AvailabilityZone]}},
			Vpc:              vpcNode{UID: c.ressources["Vpcs"][*i.DBSubnetGroup.VpcId]},
			OptionGroup:      optionGroupNode{UID: c.ressources["OptionGroups"][*i.OptionGroupMemberships[0].OptionGroupName]},
			DBParameterGroup: dbParameterGroupNode{UID: c.ressources["DbParameterGroups"][*i.DBParameterGroups[0].DBParameterGroupName]},
			DBSubnetGroup:    dbSubnetGroupNodes{dbSubnetGroupNode{UID: c.ressources["DbSubnetGroups"][*i.DBSubnetGroup.DBSubnetGroupName]}},
		}
		if i.SecondaryAvailabilityZone != nil {
			b.SecondaryAvailabilityZone = availabilityZoneNodes{availabilityZoneNode{UID: c.ressources["AvailabilityZones"][*i.SecondaryAvailabilityZone]}}
		}
		for _, j := range i.VpcSecurityGroups {
			b.SecurityGroup = append(b.SecurityGroup, securityGroupNode{UID: c.ressources["SecurityGroups"][*j.VpcSecurityGroupId]})
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
