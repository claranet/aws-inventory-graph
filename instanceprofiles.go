package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/iam"
)

type instanceProfileList struct {
	*iam.ListInstanceProfilesOutput
}
type instanceProfileNodes []instanceProfileNode

type instanceProfileNode struct {
	UID       string   `json:"uid,omitempty"`
	Type      []string `json:"dgraph.type,omitempty"`
	Name      string   `json:"name,omitempty"` // This field is only for Ratel Viz
	Region    string   `json:"Region,omitempty"`
	OwnerID   string   `json:"OwnerId,omitempty"`
	OwnerName string   `json:"OwnerName,omitempty"`

	Service             string `json:"Service,omitempty"`
	InstanceProfileID   string `json:"InstanceProfileId,omitempty"`
	InstanceProfileName string `json:"InstanceProfileName,omitempty"`
	Arn                 string `json:"Arn,omitempty"`
}

func (c *connector) listInstanceProfiles() instanceProfileList {
	defer c.waitGroup.Done()

	log.Println("List InstanceProfiles")
	response, err := iam.New(c.awsSession).ListInstanceProfiles(&iam.ListInstanceProfilesInput{})
	if err != nil {
		log.Fatal(err)
	}
	return instanceProfileList{response}
}

func (list instanceProfileList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.InstanceProfiles) == 0 {
		return
	}
	log.Println("Add InstanceProfiles Nodes")
	a := make(instanceProfileNodes, 0, len(list.InstanceProfiles))

	for _, i := range list.InstanceProfiles {
		var b instanceProfileNode
		b.Service = "iam"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"InstanceProfile"}
		b.Name = *i.InstanceProfileName
		b.InstanceProfileName = *i.InstanceProfileName
		b.InstanceProfileID = *i.InstanceProfileId
		b.Arn = *i.Arn
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]instanceProfileNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("InstanceProfile"), &m)
	for _, i := range m["list"] {
		n[i.InstanceProfileID] = i.UID
		n[i.InstanceProfileName] = i.UID
	}
	c.ressources["InstanceProfiles"] = n
}
