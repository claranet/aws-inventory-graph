package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type snapshotList struct{ *ec2.DescribeSnapshotsOutput }
type snapshotNodes []snapshotNode

type snapshotNode struct {
	UID         string     `json:"uid,omitempty"`
	Type        []string   `json:"dgraph.type,omitempty"`
	Name        string     `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID     string     `json:"OwnerId,omitempty"`
	OwnerName   string     `json:"OwnerName,omitempty"`
	Region      string     `json:"Region,omitempty"`
	Service     string     `json:"Service,omitempty"`
	SnapshotID  string     `json:"SnapshotId,omitempty"`
	Description string     `json:"Description,omitempty"`
	VolumeSize  int64      `json:"VolumeSize,omitempty"`
	Encrypted   bool       `json:"Encrypted,omitempty"`
	Progress    string     `json:"Progress,omitempty"`
	Volume      volumeNode `json:"_Volume,omitempty"`
	DeviceName  string     `json:"_Snapshot|DeviceName,omitempty"`
}

func (c *connector) listSnapshots() snapshotList {
	defer c.waitGroup.Done()

	log.Println("List Snapshots")
	response, err := ec2.New(c.awsSession).DescribeSnapshots(&ec2.DescribeSnapshotsInput{
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
	return snapshotList{response}
}

func (list snapshotList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Snapshots) == 0 {
		return
	}
	log.Println("Add Snapshot Nodes")
	a := make(snapshotNodes, 0, len(list.Snapshots))

	for _, i := range list.Snapshots {
		var b snapshotNode
		b.Service = "ec2"
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Type = []string{"Snapshot"}
		b.Name = *i.SnapshotId
		for _, tag := range i.Tags {
			if *tag.Key == "Name" {
				b.Name = *tag.Value
			}
		}
		b.SnapshotID = *i.SnapshotId
		b.Description = *i.Description
		b.VolumeSize = *i.VolumeSize
		b.Encrypted = *i.Encrypted
		b.Progress = *i.Progress
		a = append(a, b)
		if len(a) == 100 {
			c.dgraphAddNodes(a)
			c.stats.NumberOfNodes += len(a)
			a = snapshotNodes{}
		}
	}
	if len(a) != 0 {
		c.dgraphAddNodes(a)
		c.stats.NumberOfNodes += len(a)
	}

	m := make(map[string]snapshotNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("Snapshot"), &m)
	for _, i := range m["list"] {
		n[i.SnapshotID] = i.UID
	}
	c.ressources["Snapshots"] = n
}

func (list snapshotList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.Snapshots) == 0 {
		return
	}
	log.Println("Add Snapshot Edges")
	a := snapshotNodes{}
	for _, i := range list.Snapshots {
		b := snapshotNode{
			UID:    c.ressources["Snapshots"][*i.SnapshotId],
			Volume: volumeNode{UID: c.ressources["Volumes"][*i.VolumeId]},
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
