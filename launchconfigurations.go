package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type launchConfigurationList struct {
	*autoscaling.DescribeLaunchConfigurationsOutput
}
type launchConfigurationNodes []launchConfigurationNode

type launchConfigurationNode struct {
	UID                      string              `json:"uid,omitempty"`
	Type                     []string            `json:"dgraph.type,omitempty"`
	Name                     string              `json:"name,omitempty"` // This field is only for Ratel Viz
	Region                   string              `json:"Region,omitempty"`
	OwnerID                  string              `json:"OwnerId,omitempty"`
	OwnerName                string              `json:"OwnerName,omitempty"`
	Service                  string              `json:"Service,omitempty"`
	LaunchConfigurationArn   string              `json:"LaunchConfigurationArn,omitempty"`
	LaunchConfigurationName  string              `json:"LaunchConfigurationName,omitempty"`
	InstanceType             string              `json:"InstanceType,omitempty"`
	UserData                 string              `json:"UserData,omitempty"`
	EbsOptimized             bool                `json:"EbsOptimized,omitempty"`
	AssociatePublicIPAddress bool                `json:"AssociatePublicIpAddress,omitempty"`
	KeyName                  keyPairNode         `json:"_KeyName,omitempty"`
	InstanceProfile          instanceProfileNode `json:"_InstanceProfile,omitempty"`
	Image                    imageNode           `json:"_Image,omitempty"`
	SecurityGroup            securityGroupNodes  `json:"_SecurityGroup,omitempty"`
}

func (c *connector) listLaunchConfigurations() launchConfigurationList {
	defer c.waitGroup.Done()

	log.Println("List LaunchConfigurations")
	response, err := autoscaling.New(c.awsSession).DescribeLaunchConfigurations(&autoscaling.DescribeLaunchConfigurationsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return launchConfigurationList{response}
}

func (list launchConfigurationList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.LaunchConfigurations) == 0 {
		return
	}
	log.Println("Add LaunchConfiguration Nodes")
	a := make(launchConfigurationNodes, 0, len(list.LaunchConfigurations))

	for _, i := range list.LaunchConfigurations {
		var b launchConfigurationNode
		b.Service = "autoscaling"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"LaunchConfiguration"}
		b.Name = *i.LaunchConfigurationName
		b.LaunchConfigurationArn = *i.LaunchConfigurationARN
		b.LaunchConfigurationName = *i.LaunchConfigurationName
		b.InstanceType = *i.InstanceType
		b.UserData = *i.UserData
		b.EbsOptimized = *i.EbsOptimized
		if i.AssociatePublicIpAddress != nil {
			b.AssociatePublicIPAddress = *i.AssociatePublicIpAddress
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)

	m := make(map[string]launchConfigurationNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("LaunchConfiguration"), &m)
	for _, i := range m["list"] {
		n[i.LaunchConfigurationName] = i.UID
	}
	c.ressources["LaunchConfigurations"] = n
}

func (list launchConfigurationList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.LaunchConfigurations) == 0 {
		return
	}
	log.Println("Add LaunchConfiguration Edges")
	a := launchConfigurationNodes{}
	for _, i := range list.LaunchConfigurations {
		b := launchConfigurationNode{
			UID:     c.ressources["LaunchConfigurations"][*i.LaunchConfigurationName],
			KeyName: keyPairNode{UID: c.ressources["KeyPairs"][*i.KeyName]},
			Image:   imageNode{UID: c.ressources["Images"][*i.ImageId]},
		}
		if i.IamInstanceProfile != nil {
			b.InstanceProfile = instanceProfileNode{UID: c.ressources["InstanceProfiles"][*i.IamInstanceProfile]}
		}
		for _, j := range i.SecurityGroups {
			b.SecurityGroup = append(b.SecurityGroup, securityGroupNode{UID: c.ressources["SecurityGroups"][*j]})
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
