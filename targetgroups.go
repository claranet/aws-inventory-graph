package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

type targetGroupList struct {
	*elbv2.DescribeTargetGroupsOutput
}
type targetGroupNodes []targetGroupNode

type targetGroupNode struct {
	UID                     string            `json:"uid,omitempty"`
	Type                    []string          `json:"dgraph.type,omitempty"`
	Name                    string            `json:"name,omitempty"` // This field is only for Ratel Viz
	Region                  string            `json:"Region,omitempty"`
	OwnerID                 string            `json:"OwnerId,omitempty"`
	OwnerName               string            `json:"OwnerName,omitempty"`
	Service                 string            `json:"Service,omitempty"`
	TargetGroupArn          string            `json:"TargetGroupArn,omitempty"`
	TargetGroupName         string            `json:"TargetGroupName,omitempty"`
	TargetType              string            `json:"TargetType,omitempty"`
	HealthCheckPath         string            `json:"HealthCheckPath,omitempty"`
	HealthCheckProtocol     string            `json:"HealthCheckProtocol,omitempty"`
	HealthCheckPort         string            `json:"HealthCheckPort,omitempty"`
	HealthyThresholdCount   int64             `json:"HealthCheckIntervalSeconds,omitempty"`
	Port                    int64             `json:"Port,omitempty"`
	UnhealthyThresholdCount int64             `json:"UnhealthyThresholdCount,omitempty"`
	Protocol                string            `json:"Protocol,omitempty"`
	Vpc                     vpcNode           `json:"_Vpc,omitempty"`
	LoadBalancer            loadBalancerNodes `json:"_LoadBalancer,omitempty"`
}

func (c *connector) listTargetGroups() targetGroupList {
	defer c.waitGroup.Done()

	log.Println("List TargetGroups")
	response, err := elbv2.New(c.awsSession).DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return targetGroupList{response}
}

func (list targetGroupList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.TargetGroups) == 0 {
		return
	}
	log.Println("Add TargetGroup Nodes")
	a := make(targetGroupNodes, 0, len(list.TargetGroups))

	for _, i := range list.TargetGroups {
		var b targetGroupNode
		b.Service = "elbv2"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"TargetGroup"}
		b.Name = *i.TargetGroupName
		b.TargetGroupName = *i.TargetGroupName
		b.TargetGroupArn = *i.TargetGroupArn
		b.TargetType = *i.TargetType
		if i.HealthCheckPath != nil {
			b.HealthCheckPath = *i.HealthCheckPath
			b.HealthCheckProtocol = *i.HealthCheckProtocol
			b.HealthCheckPort = *i.HealthCheckPort
			b.HealthyThresholdCount = *i.HealthyThresholdCount
			b.UnhealthyThresholdCount = *i.UnhealthyThresholdCount
		}
		b.Port = *i.Port
		b.Protocol = *i.Protocol
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	m := make(map[string]targetGroupNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("TargetGroup"), &m)
	for _, i := range m["list"] {
		n[i.TargetGroupArn] = i.UID
	}
	c.ressources["TargetGroups"] = n
}

func (list targetGroupList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.TargetGroups) == 0 {
		return
	}
	log.Println("Add TargetGroup Edges")
	a := targetGroupNodes{}
	for _, i := range list.TargetGroups {
		b := targetGroupNode{
			UID: c.ressources["TargetGroups"][*i.TargetGroupArn],
			Vpc: vpcNode{UID: c.ressources["Vpcs"][*i.VpcId]},
		}
		for _, i := range i.LoadBalancerArns {
			b.LoadBalancer = append(b.LoadBalancer, loadBalancerNode{UID: c.ressources["LoadBalancersV2"][*i]})
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
