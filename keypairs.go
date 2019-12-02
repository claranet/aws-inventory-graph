package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type keyPairList struct{ *ec2.DescribeKeyPairsOutput }
type keyPairNodes []keyPairNode

type keyPairNode struct {
	UID            string   `json:"uid,omitempty"`
	Type           []string `json:"dgraph.type,omitempty"`
	Name           string   `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID        string   `json:"OwnerId,omitempty"`
	OwnerName      string   `json:"OwnerName,omitempty"`
	Region         string   `json:"Region,omitempty"`
	Service        string   `json:"Service,omitempty"`
	KeyName        string   `json:"KeyName,omitempty"`
	KeyFingerprint string   `json:"KeyFingerPrint,omitempty"`
}

func (c *connector) listKeyPairs() keyPairList {
	defer c.waitGroup.Done()

	log.Println("List KeyPair Nodes")
	response, err := ec2.New(c.awsSession).DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return keyPairList{response}
}

func (list keyPairList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.KeyPairs) == 0 {
		return
	}
	log.Println("Add KeyPair Nodes")
	a := make(keyPairNodes, 0, len(list.KeyPairs))

	for _, i := range list.KeyPairs {
		var b keyPairNode
		b.Service = "ec2"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"KeyPair"}
		b.Name = *i.KeyName
		b.KeyName = *i.KeyName
		b.KeyFingerprint = *i.KeyFingerprint
		a = append(a, b)
	}
	c.dgraphAddNodes(a)

	m := make(map[string]keyPairNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("KeyPair"), &m)
	for _, i := range m["list"] {
		n[i.KeyName] = i.UID
	}
	c.ressources["KeyPairs"] = n
}
