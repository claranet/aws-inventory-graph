package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type loadBalancerList struct {
	*elb.DescribeLoadBalancersOutput
}

type loadBalancerV2List struct {
	*elbv2.DescribeLoadBalancersOutput
}

type loadBalancerNodes []loadBalancerNode

type loadBalancerNode struct {
	UID                       string                `json:"uid,omitempty"`
	Type                      []string              `json:"dgraph.type,omitempty"`
	Name                      string                `json:"name,omitempty"` // This field is only for Ratel Viz
	Region                    string                `json:"Region,omitempty"`
	OwnerID                   string                `json:"OwnerId,omitempty"`
	OwnerName                 string                `json:"OwnerName,omitempty"`
	Service                   string                `json:"Service,omitempty"`
	LoadBalancerName          string                `json:"LoadBalancerName,omitempty"`
	LoadBalancerArn           string                `json:"LoadBalancerArn,omitempty"`
	CanonicalHostedZoneID     string                `json:"CanonicalHostedZoneId,omitempty"`
	CanonicalHostedZoneNameID string                `json:"CanonicalHostedZoneNameId,omitempty"`
	CanonicalHostedZoneName   string                `json:"CanonicalHostedZoneName,omitempty"`
	DNSName                   string                `json:"DNSName,omitempty"`
	Scheme                    string                `json:"Scheme,omitempty"`
	LoadBalancerType          string                `json:"LoadBalancerType,omitempty"`
	IPAddressType             string                `json:"IpAddressType,omitempty"`
	State                     string                `json:"State,omitempty"`
	Instance                  instanceNodes         `json:"_Instance,omitempty"`
	Vpc                       vpcNode               `json:"_Vpc,omitempty"`
	AvailabilityZone          availabilityZoneNodes `json:"_AvailabilityZone,omitempty"`
	SecurityGroup             securityGroupNodes    `json:"_SecurityGroup,omitempty"`
	Subnet                    subnetNodes           `json:"_Subnet,omitempty"`
}

func (c *connector) listLoadBalancers() loadBalancerList {
	defer c.waitGroup.Done()

	log.Println("List LoadBalancers")
	response, err := elb.New(c.awsSession).DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{})
	if err != nil {
		log.Fatal(err)
	}
	return loadBalancerList{response}
}

func (c *connector) listLoadBalancersV2() loadBalancerV2List {
	defer c.waitGroup.Done()

	log.Println("List LoadBalancersV2")
	response, err := elbv2.New(c.awsSession).DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		log.Fatal(err)
	}
	return loadBalancerV2List{response}
}

func (list loadBalancerList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.LoadBalancerDescriptions) == 0 {
		return
	}
	log.Println("Add LoadBalancer Nodes")
	a := make(loadBalancerNodes, 0, len(list.LoadBalancerDescriptions))

	for _, i := range list.LoadBalancerDescriptions {
		var b loadBalancerNode
		b.Service = "elb"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"LoadBalancer"}
		b.LoadBalancerType = "classic"
		b.Name = *i.LoadBalancerName
		b.LoadBalancerName = *i.LoadBalancerName
		b.CanonicalHostedZoneNameID = *i.CanonicalHostedZoneNameID
		if i.CanonicalHostedZoneName != nil {
			b.CanonicalHostedZoneName = *i.CanonicalHostedZoneName
		}
		b.DNSName = *i.DNSName
		b.Scheme = *i.Scheme
		a = append(a, b)
	}
	c.dgraphAddNodes(a)

	m := make(map[string]loadBalancerNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("LoadBalancer"), &m)
	for _, i := range m["list"] {
		n[i.LoadBalancerName] = i.UID
	}
	c.ressources["LoadBalancers"] = n
}

func (list loadBalancerV2List) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.LoadBalancers) == 0 {
		return
	}
	log.Println("Add LoadBalancerV2 Nodes")
	a := make(loadBalancerNodes, 0, len(list.LoadBalancers))

	for _, i := range list.LoadBalancers {
		var b loadBalancerNode
		b.Service = "elbv2"
		b.Region = c.awsRegion
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Type = []string{"LoadBalancer"}
		b.LoadBalancerType = *i.Type
		b.Name = *i.LoadBalancerName
		b.LoadBalancerName = *i.LoadBalancerName
		b.LoadBalancerArn = *i.LoadBalancerArn
		b.CanonicalHostedZoneID = *i.CanonicalHostedZoneId
		b.IPAddressType = *i.IpAddressType
		b.State = *i.State.Code
		b.Name = *i.DNSName
		b.Scheme = *i.Scheme
		a = append(a, b)
	}
	c.dgraphAddNodes(a)

	m := make(map[string]loadBalancerNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("LoadBalancer"), &m)
	for _, i := range m["list"] {
		n[i.LoadBalancerArn] = i.UID
	}
	c.ressources["LoadBalancersV2"] = n
}

func (list loadBalancerList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.LoadBalancerDescriptions) == 0 {
		return
	}
	log.Println("Add LoadBalancer Edges")
	a := loadBalancerNodes{}
	for _, i := range list.LoadBalancerDescriptions {
		b := loadBalancerNode{
			UID: c.ressources["LoadBalancers"][*i.LoadBalancerName],
			Vpc: vpcNode{UID: c.ressources["Vpcs"][*i.VPCId]},
		}
		if len(i.Instances) != 0 {
			for _, j := range i.Instances {
				b.Instance = append(b.Instance, instanceNode{UID: c.ressources["Instances"][*j.InstanceId]})
			}
		}
		for _, i := range i.AvailabilityZones {
			b.AvailabilityZone = append(b.AvailabilityZone, availabilityZoneNode{UID: c.ressources["AvailabilityZones"][*i]})
		}
		for _, i := range i.Subnets {
			b.Subnet = append(b.Subnet, subnetNode{UID: c.ressources["Subnets"][*i]})
		}
		for _, j := range i.SecurityGroups {
			b.SecurityGroup = append(b.SecurityGroup, securityGroupNode{UID: c.ressources["SecurityGroups"][*j]})
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}

func (list loadBalancerV2List) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.LoadBalancers) == 0 {
		return
	}
	log.Println("Add LoadBalancerV2 Edges")
	a := loadBalancerNodes{}
	for _, i := range list.LoadBalancers {
		b := loadBalancerNode{
			UID: c.ressources["LoadBalancersV2"][*i.LoadBalancerArn],
			Vpc: vpcNode{UID: c.ressources["Vpcs"][*i.VpcId]},
		}
		for _, i := range i.AvailabilityZones {
			b.AvailabilityZone = append(b.AvailabilityZone, availabilityZoneNode{UID: c.ressources["AvailabilityZones"][*i.ZoneName]})
			b.Subnet = append(b.Subnet, subnetNode{UID: c.ressources["Subnets"][*i.SubnetId]})
		}
		for _, j := range i.SecurityGroups {
			b.SecurityGroup = append(b.SecurityGroup, securityGroupNode{UID: c.ressources["SecurityGroups"][*j]})
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
