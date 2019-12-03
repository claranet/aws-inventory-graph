package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type instanceList struct{ *ec2.DescribeInstancesOutput }
type instanceNodes []instanceNode

type instanceNode struct {
	UID                string                `json:"uid,omitempty"`
	Type               []string              `json:"dgraph.type,omitempty"`
	Name               string                `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID            string                `json:"OwnerId,omitempty"`
	OwnerName          string                `json:"OwnerName,omitempty"`
	Region             string                `json:"Region,omitempty"`
	Service            string                `json:"Service,omitempty"`
	InstanceID         string                `json:"InstanceId,omitempty"`
	InstanceType       string                `json:"InstanceType,omitempty"`
	State              string                `json:"State,omitempty"`
	EbsOptimized       bool                  `json:"EbsOptimized,omitempty"`
	Hypervisor         string                `json:"Hypervisor,omitempty"`
	VirtualizationType string                `json:"VirtualizationType,omitempty"`
	PrivateIPAddress   string                `json:"PrivateIpAddress,omitempty"`
	PublicIPAddress    string                `json:"PublicIpAddress,omitempty"`
	KeyName            keyPairNode           `json:"_KeyName,omitempty"`
	AvailabilityZone   availabilityZoneNodes `json:"_AvailabilityZone,omitempty"`
	Vpc                vpcNode               `json:"_Vpc,omitempty"`
	InstanceProfile    instanceProfileNode   `json:"_InstanceProfile,omitempty"`
	Image              imageNode             `json:"_Image,omitempty"`
	Subnet             subnetNodes           `json:"_Subnet,omitempty"`
	SecurityGroup      securityGroupNodes    `json:"_SecurityGroup,omitempty"`
}

func (c *connector) listInstances() instanceList {
	defer c.waitGroup.Done()

	log.Println("List Instances")
	response, err := ec2.New(c.awsSession).DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		log.Fatal(err)
	}
	return instanceList{response}
}

func (list instanceList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Reservations) == 0 {
		return
	}
	log.Println("Add Instance Nodes")
	a := make(instanceNodes, 0, len(list.Reservations))

	for _, r := range list.Reservations {
		for _, i := range r.Instances {
			var b instanceNode
			b.Service = "ec2"
			b.Type = []string{"Instance"}
			b.OwnerID = c.awsAccountID
			b.OwnerName = c.awsAccountName
			b.Region = c.awsRegion
			b.Name = *i.InstanceId
			for _, tag := range i.Tags {
				if *tag.Key == "Name" {
					b.Name = *tag.Value
				}
			}
			b.InstanceID = *i.InstanceId
			b.InstanceType = *i.InstanceType
			b.State = *i.State.Name
			b.VirtualizationType = *i.VirtualizationType
			b.Hypervisor = *i.Hypervisor
			b.EbsOptimized = *i.EbsOptimized
			if i.PrivateIpAddress != nil {
				b.PrivateIPAddress = *i.PrivateIpAddress
			}
			if i.PublicIpAddress != nil {
				b.PublicIPAddress = *i.PublicIpAddress
			}
			a = append(a, b)
		}
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]instanceNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("Instance"), &m)
	for _, i := range m["list"] {
		n[i.InstanceID] = i.UID
	}
	c.ressources["Instances"] = n
}

func (list instanceList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Reservations) == 0 {
		return
	}
	log.Println("Add Instance Edges")
	a := instanceNodes{}
	for _, r := range list.Reservations {
		for _, i := range r.Instances {
			b := instanceNode{
				UID: c.ressources["Instances"][*i.InstanceId],

				Image:            imageNode{UID: c.ressources["Images"][*i.ImageId]},
				AvailabilityZone: availabilityZoneNodes{availabilityZoneNode{UID: c.ressources["AvailabilityZones"][*i.Placement.AvailabilityZone]}},
			}
			if i.KeyName != nil {
				b.KeyName = keyPairNode{UID: c.ressources["KeyPairs"][*i.KeyName]}
			}
			if i.SubnetId != nil {
				b.Subnet = subnetNodes{subnetNode{UID: c.ressources["Subnets"][*i.SubnetId]}}
			}
			if i.IamInstanceProfile != nil {
				b.InstanceProfile = instanceProfileNode{UID: c.ressources["InstanceProfiles"][*i.IamInstanceProfile.Id]}
			}
			if i.VpcId != nil {
				b.Vpc = vpcNode{UID: c.ressources["Vpcs"][*i.VpcId]}
			}
			for _, j := range i.SecurityGroups {
				b.SecurityGroup = append(b.SecurityGroup, securityGroupNode{UID: c.ressources["SecurityGroups"][*j.GroupId]})
			}
			a = append(a, b)
		}
	}
	c.dgraphAddNodes(a)
}
