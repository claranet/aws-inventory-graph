package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type autoScalingGroupList struct {
	*autoscaling.DescribeAutoScalingGroupsOutput
}
type autoScalingGroupNodes []autoScalingGroupNode

type autoScalingGroupNode struct {
	UID                              string                   `json:"uid,omitempty"`
	Type                             []string                 `json:"dgraph.type,omitempty"`
	Name                             string                   `json:"name,omitempty"` // This field is only for Ratel Viz
	Region                           string                   `json:"Region,omitempty"`
	OwnerID                          string                   `json:"OwnerId,omitempty"`
	OwnerName                        string                   `json:"OwnerName,omitempty"`
	Service                          string                   `json:"Service,omitempty"`
	AutoScalingGroupArn              string                   `json:"AutoScalingGroupArn,omitempty"`
	AutoScalingGroupName             string                   `json:"AutoScalingGroupName,omitempty"`
	DesiredCapacity                  int64                    `json:"DesiredCapacity,omitempty"`
	MinSize                          int64                    `json:"MinSize,omitempty"`
	MaxSize                          int64                    `json:"MaxSize,omitempty"`
	DefaultCooldown                  int64                    `json:"DefaultCooldown,omitempty"`
	HealthCheckGracePeriod           int64                    `json:"HealthCheckGracePeriod,omitempty"`
	HealthCheckType                  string                   `json:"HealthCheckType,omitempty"`
	TerminationPolicy                string                   `json:"TerminationPolicy,omitempty"`
	NewInstancesProtectedFromScaleIn bool                     `json:"NewInstancesProtectedFromScaleIn,omitempty"`
	AvailabilityZone                 availabilityZoneNodes    `json:"_AvailabilityZone,omitempty"`
	Instance                         instanceNodes            `json:"_Instance,omitempty"`
	TargetGroup                      targetGroupNodes         `json:"_TargetGroup,omitempty"`
	LoadBalancer                     loadBalancerNodes        `json:"_LoadBalancer,omitempty"`
	LaunchConfiguration              launchConfigurationNodes `json:"_LaunchConfiguration,omitempty"`
	LaunchTemplate                   launchTemplateNodes      `json:"_LaunchTemplate,omitempty"`
	// Subnet                           subnetNodes             `json:"_Subnet<,omitempty"`
}

func (c *connector) listAutoScalingGroups() autoScalingGroupList {
	defer c.waitGroup.Done()

	log.Println("List AutoScalingGroups")
	response, err := autoscaling.New(c.awsSession).DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return autoScalingGroupList{response}
}

func (list autoScalingGroupList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.AutoScalingGroups) == 0 {
		return
	}
	log.Println("Add AutoScalingGroup Nodes")
	a := make(autoScalingGroupNodes, 0, len(list.AutoScalingGroups))

	for _, i := range list.AutoScalingGroups {
		var b autoScalingGroupNode
		b.Service = "autoscaling"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"AutoScalingGroup"}
		b.Name = *i.AutoScalingGroupName
		for _, tag := range i.Tags {
			if *tag.Key == "Name" {
				b.Name = *tag.Value
			}
		}
		b.AutoScalingGroupArn = *i.AutoScalingGroupARN
		b.MinSize = *i.MinSize
		b.MaxSize = *i.MaxSize
		b.DesiredCapacity = *i.DesiredCapacity
		b.TerminationPolicy = *i.TerminationPolicies[0]
		b.DefaultCooldown = *i.DefaultCooldown
		b.HealthCheckGracePeriod = *i.HealthCheckGracePeriod
		b.HealthCheckType = *i.HealthCheckType
		b.NewInstancesProtectedFromScaleIn = *i.NewInstancesProtectedFromScaleIn
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]autoScalingGroupNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("AutoScalingGroup"), &m)
	for _, i := range m["list"] {
		n[i.AutoScalingGroupArn] = i.UID
	}
	c.ressources["AutoScalingGroups"] = n
}

func (list autoScalingGroupList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.AutoScalingGroups) == 0 {
		return
	}
	log.Println("Add AutoScalingGroup Edges")
	a := autoScalingGroupNodes{}
	for _, i := range list.AutoScalingGroups {
		b := autoScalingGroupNode{
			UID: c.ressources["AutoScalingGroups"][*i.AutoScalingGroupARN],
		}
		if len(i.Instances) != 0 {
			for _, j := range i.Instances {
				b.Instance = append(b.Instance, instanceNode{UID: c.ressources["Instances"][*j.InstanceId]})
			}
		}
		if len(i.TargetGroupARNs) != 0 {
			for _, j := range i.TargetGroupARNs {
				b.TargetGroup = append(b.TargetGroup, targetGroupNode{UID: c.ressources["TargetGroups"][*j]})
			}
		}
		for _, i := range i.AvailabilityZones {
			b.AvailabilityZone = append(b.AvailabilityZone, availabilityZoneNode{UID: c.ressources["AvailabilityZones"][*i]})
		}
		for _, i := range i.LoadBalancerNames {
			b.LoadBalancer = append(b.LoadBalancer, loadBalancerNode{UID: c.ressources["LoadBalancers"][*i]})
		}
		if i.LaunchConfigurationName != nil {
			b.LaunchConfiguration = launchConfigurationNodes{launchConfigurationNode{UID: c.ressources["LaunchConfigurations"][*i.LaunchConfigurationName]}}
		}
		if i.LaunchTemplate != nil {
			b.LaunchTemplate = launchTemplateNodes{launchTemplateNode{UID: c.ressources["LaunchTemplates"][*i.LaunchTemplate.LaunchTemplateId]}}
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
