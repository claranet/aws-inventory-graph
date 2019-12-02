package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type launchTemplateList struct {
	*ec2.DescribeLaunchTemplatesOutput
}
type launchTemplateNodes []launchTemplateNode

type launchTemplateNode struct {
	UID                  string   `json:"uid,omitempty"`
	Type                 []string `json:"dgraph.type,omitempty"`
	Name                 string   `json:"name,omitempty"` // This field is only for Ratel Viz
	Region               string   `json:"Region,omitempty"`
	OwnerID              string   `json:"OwnerId,omitempty"`
	OwnerName            string   `json:"OwnerName,omitempty"`
	Service              string   `json:"Service,omitempty"`
	LaunchTemplateID     string   `json:"LaunchTemplateId,omitempty"`
	LaunchTemplateName   string   `json:"LaunchTemplateName,omitempty"`
	LatestVersionNumber  int64    `json:"LatestVersionNumber,omitempty"`
	DefaultVersionNumber int64    `json:"DefaultVersionNumber,omitempty"`
	CreatedBy            string   `json:"CreatedBy,omitempty"`
}

func (c *connector) listLaunchTemplates() launchTemplateList {
	defer c.waitGroup.Done()

	log.Println("List LaunchTemplates")
	response, err := ec2.New(c.awsSession).DescribeLaunchTemplates(&ec2.DescribeLaunchTemplatesInput{})
	if err != nil {
		log.Fatal(err)
	}
	return launchTemplateList{response}
}

func (list launchTemplateList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.LaunchTemplates) == 0 {
		return
	}
	log.Println("Add LaunchTemplate Nodes")
	a := make(launchTemplateNodes, 0, len(list.LaunchTemplates))

	for _, i := range list.LaunchTemplates {
		var b launchTemplateNode
		b.Service = "ec2"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"LaunchTemplate"}
		b.Name = *i.LaunchTemplateName
		b.LaunchTemplateName = *i.LaunchTemplateName
		b.LaunchTemplateID = *i.LaunchTemplateId
		b.LatestVersionNumber = *i.LatestVersionNumber
		b.DefaultVersionNumber = *i.DefaultVersionNumber
		b.CreatedBy = *i.CreatedBy
		a = append(a, b)
	}
	c.dgraphAddNodes(a)

	m := make(map[string]launchTemplateNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("LaunchTemplate"), &m)
	for _, i := range m["list"] {
		n[i.LaunchTemplateID] = i.UID
	}
	c.ressources["LaunchTemplates"] = n
}
