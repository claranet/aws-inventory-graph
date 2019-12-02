package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type imageList struct{ *ec2.DescribeImagesOutput }
type imageNodes []imageNode

type imageNode struct {
	UID                string   `json:"uid,omitempty"`
	Type               []string `json:"dgraph.type,omitempty"`
	Name               string   `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID            string   `json:"OwnerId,omitempty"`
	OwnerName          string   `json:"OwnerName,omitempty"`
	Region             string   `json:"Region,omitempty"`
	Service            string   `json:"Service,omitempty"`
	ImageID            string   `json:"ImageId,omitempty"`
	VirtualizationType string   `json:"VirtualizationType,omitempty"`
	Hypervisor         string   `json:"Hypervisor,omitempty"`
	EnaSupport         bool     `json:"EnaSupport,omitempty"`
	SriovNetSupport    string   `json:"SriovNetSupport,omitempty"`
	State              string   `json:"State,omitempty"`
	Architecture       string   `json:"Architecture,omitempty"`
	ImageLocation      string   `json:"ImageLocation,omitempty"`
	RootDeviceType     string   `json:"RootDeviceType,omitempty"`
	RootDeviceName     string   `json:"RootDeviceName,omitempty"`
	Public             bool     `json:"Public,omitempty"`
	ImageType          string   `json:"ImageType,omitempty"`
}

func (c *connector) listImages() imageList {
	defer c.waitGroup.Done()

	log.Println("List Images")
	response, err := ec2.New(c.awsSession).DescribeImages(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("owner-id"),
				Values: []*string{aws.String(c.awsAccountID)},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return imageList{response}
}

func (list imageList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Images) == 0 {
		return
	}
	log.Println("Add Image Nodes")
	a := make(imageNodes, 0, len(list.Images))

	for _, i := range list.Images {
		var b imageNode
		b.Service = "ec2"
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Type = []string{"Image"}
		b.Name = *i.ImageId
		for _, tag := range i.Tags {
			if *tag.Key == "Name" {
				b.Name = *tag.Value
			}
		}
		b.ImageID = *i.ImageId
		b.VirtualizationType = *i.VirtualizationType
		b.Hypervisor = *i.Hypervisor
		if i.EnaSupport != nil {
			b.EnaSupport = *i.EnaSupport
		}
		if i.SriovNetSupport != nil {
			b.SriovNetSupport = *i.SriovNetSupport
		}
		b.State = *i.State
		b.Architecture = *i.Architecture
		b.ImageLocation = *i.ImageLocation
		b.RootDeviceType = *i.RootDeviceType
		b.RootDeviceName = *i.RootDeviceName
		b.Public = *i.Public
		b.ImageType = *i.ImageType
		a = append(a, b)
	}

	c.dgraphAddNodes(a)

	m := make(map[string]imageNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("Image"), &m)
	for _, i := range m["list"] {
		n[i.ImageID] = i.UID
	}
	c.ressources["Images"] = n
}
