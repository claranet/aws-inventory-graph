package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/rds"
)

type optionGroupList struct {
	*rds.DescribeOptionGroupsOutput
}
type optionGroupNodes []optionGroupNode

type optionGroupNode struct {
	UID                                   string   `json:"uid,omitempty"`
	Type                                  []string `json:"dgraph.type,omitempty"`
	Name                                  string   `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID                               string   `json:"OwnerId,omitempty"`
	OwnerName                             string   `json:"OwnerName,omitempty"`
	Region                                string   `json:"Region,omitempty"`
	Service                               string   `json:"Service,omitempty"`
	OptionGroupArn                        string   `json:"OptionGroupArn,omitempty"`
	OptionGroupName                       string   `json:"OptionGroupName,omitempty"`
	Description                           string   `json:"Description,omitempty"`
	MajorEngineVersion                    string   `json:"MajorEngineVersion,omitempty"`
	EngineName                            string   `json:"EngineName,omitempty"`
	AllowsVpcAndNonVpcInstanceMemberships bool     `json:"AllowsVpcAndNonVpcInstanceMemberships,omitempty"`
}

func (c *connector) listOptionGroups() optionGroupList {
	defer c.waitGroup.Done()

	log.Println("List OptionGroups")
	response, err := rds.New(c.awsSession).DescribeOptionGroups(&rds.DescribeOptionGroupsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return optionGroupList{response}
}

func (list optionGroupList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.OptionGroupsList) == 0 {
		return
	}
	log.Println("Add OptionGroup Nodes")
	a := make(optionGroupNodes, 0, len(list.OptionGroupsList))

	for _, i := range list.OptionGroupsList {
		var b optionGroupNode
		b.Service = "rds"
		b.Type = []string{"OptionGroup"}
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Name = *i.OptionGroupName
		b.OptionGroupName = *i.OptionGroupName
		b.OptionGroupArn = *i.OptionGroupArn
		b.Description = *i.OptionGroupDescription
		b.MajorEngineVersion = *i.MajorEngineVersion
		b.EngineName = *i.EngineName
		b.AllowsVpcAndNonVpcInstanceMemberships = *i.AllowsVpcAndNonVpcInstanceMemberships
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]optionGroupNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("OptionGroup"), &m)
	for _, i := range m["list"] {
		n[i.OptionGroupName] = i.UID
	}
	c.ressources["OptionGroups"] = n
}
