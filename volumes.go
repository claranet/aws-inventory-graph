package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type volumeList struct{ *ec2.DescribeVolumesOutput }
type volumeNodes []volumeNode

type volumeNode struct {
	UID                 string                `json:"uid,omitempty"`
	Type                []string              `json:"dgraph.type,omitempty"`
	Name                string                `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID             string                `json:"OwnerId,omitempty"`
	OwnerName           string                `json:"OwnerName,omitempty"`
	Region              string                `json:"Region,omitempty"`
	Service             string                `json:"Service,omitempty"`
	Encrypted           bool                  `json:"Encrypted,omitempty"`
	VolumeID            string                `json:"VolumeId,omitempty"`
	VolumeType          string                `json:"VolumeType,omitempty"`
	Size                int64                 `json:"Size,omitempty"`
	Iops                int64                 `json:"Iops,omitempty"`
	State               string                `json:"State,omitempty"`
	Device              string                `json:"Device,omitempty"`
	DeleteOnTermination bool                  `json:"DeleteOnTermination,omitempty"`
	AvailabilityZone    availabilityZoneNodes `json:"_AvailabilityZone,omitempty"`
	Instance            instanceNodes         `json:"_Instance,omitempty"`
}

func (c *connector) listVolumes() volumeList {
	defer c.waitGroup.Done()

	log.Println("List Volumes")
	response, err := ec2.New(c.awsSession).DescribeVolumes(&ec2.DescribeVolumesInput{})
	if err != nil {
		log.Fatal(err)
	}
	return volumeList{response}
}

func (list volumeList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Volumes) == 0 {
		return
	}
	log.Println("Add Volume Nodes")
	a := make(volumeNodes, 0, len(list.Volumes))

	for _, i := range list.Volumes {
		var b volumeNode
		b.Service = "ec2"
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Type = []string{"Volume"}
		b.Name = *i.VolumeId
		for _, tag := range i.Tags {
			if *tag.Key == "Name" {
				b.Name = *tag.Value
			}
		}
		b.VolumeID = *i.VolumeId
		b.VolumeType = *i.VolumeType
		b.State = *i.State
		b.Size = *i.Size
		if *i.VolumeType == "gp2" || *i.VolumeType == "iops" {
			b.Iops = *i.Iops
		}
		if len(i.Attachments) != 0 {
			b.Device = *i.Attachments[0].Device
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]volumeNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("Volume"), &m)
	for _, i := range m["list"] {
		n[i.VolumeID] = i.UID
	}
	c.ressources["Volumes"] = n
}

func (list volumeList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Volumes) == 0 {
		return
	}
	log.Println("Add Volume Edges")
	a := volumeNodes{}
	for _, i := range list.Volumes {
		b := volumeNode{
			UID:              c.ressources["Volumes"][*i.VolumeId],
			AvailabilityZone: availabilityZoneNodes{availabilityZoneNode{UID: c.ressources["AvailabilityZones"][*i.AvailabilityZone]}},
		}
		if *i.State == "in-use" {
			b.Instance = instanceNodes{instanceNode{UID: c.ressources["Instances"][*i.Attachments[0].InstanceId]}}
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
