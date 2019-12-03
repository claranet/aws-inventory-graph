package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type securityGroupList struct {
	*ec2.DescribeSecurityGroupsOutput
}
type securityGroupNodes []securityGroupNode
type cidrNodes []cidrNode

type securityGroupNode struct {
	UID                  string             `json:"uid,omitempty"`
	Type                 []string           `json:"dgraph.type,omitempty"`
	Name                 string             `json:"name,omitempty"` // This field is only for Ratel Viz
	OwnerID              string             `json:"OwnerId,omitempty"`
	OwnerName            string             `json:"OwnerName,omitempty"`
	Region               string             `json:"Region,omitempty"`
	Service              string             `json:"Service,omitempty"`
	GroupID              string             `json:"GroupId,omitempty"`
	Description          string             `json:"Description,omitempty"`
	GroupName            string             `json:"GroupName,omitempty"`
	Vpc                  vpcNode            `json:"_Vpc,omitempty"`
	SecurityGroup        securityGroupNodes `json:"_SecurityGroup,omitempty"`
	SecurityGroupPortTCP string             `json:"_SecurityGroup|PortTcp,omitempty"`
	Cidr                 cidrNodes          `json:"_Cidr,omitempty"`
	// SecurityGroupPortUDP string             `json:"_SecurityGroup|PortUdp,omitempty"`
}

type cidrNode struct {
	UID         string   `json:"uid,omitempty"`
	Type        []string `json:"dgraph.type,omitempty"`
	Name        string   `json:"name,omitempty"` // This field is only for Ratel Viz
	Region      string   `json:"Region,omitempty"`
	Service     string   `json:"Service,omitempty"`
	CidrPortTCP string   `json:"_Cidr|PortTcp,omitempty"`
}

func (c *connector) listSecurityGroups() securityGroupList {
	defer c.waitGroup.Done()

	log.Println("List securityGroups")
	response, err := ec2.New(c.awsSession).DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return securityGroupList{response}
}

func (list securityGroupList) addNodes(c *connector) {
	defer c.waitGroup.Done()

	if len(list.SecurityGroups) == 0 {
		return
	}
	log.Println("Add SecurityGroup Nodes")
	a := make(securityGroupNodes, 0, len(list.SecurityGroups))
	z := make(cidrNodes, 0)
	x, y := make(map[string]string), make(map[string]string)

	for _, i := range list.SecurityGroups {
		var b securityGroupNode
		b.Service = "ec2"
		b.OwnerID = c.awsAccountID
		b.OwnerName = c.awsAccountName
		b.Region = c.awsRegion
		b.Type = []string{"SecurityGroup"}
		b.Name = *i.GroupName
		b.GroupID = *i.GroupId
		b.GroupName = *i.GroupName
		b.Description = *i.Description
		for _, j := range i.IpPermissions {
			for _, k := range j.UserIdGroupPairs {
				var p string
				if *j.IpProtocol == "tcp" {
					if *j.FromPort == *j.ToPort {
						p = strconv.FormatInt(*j.FromPort, 10)
					} else {
						p = strconv.FormatInt(*j.FromPort, 10) + "-" + strconv.FormatInt(*j.ToPort, 10)
					}
				}
				if *j.IpProtocol == "-1" {
					p = "0-65535"
				}
				x[*k.GroupId] += *i.GroupId + "#" + p + "!"
			}
			for _, k := range j.IpRanges {
				var p string
				if *j.IpProtocol == "tcp" {
					if *j.FromPort == *j.ToPort {
						p = strconv.FormatInt(*j.FromPort, 10)
					} else {
						p = strconv.FormatInt(*j.FromPort, 10) + "-" + strconv.FormatInt(*j.ToPort, 10)
					}
				}
				if *j.IpProtocol == "-1" {
					p = "0-65535"
				}
				y[*k.CidrIp] += *i.GroupId + "#" + p + "!"
			}
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
	c.stats.NumberOfNodes += len(a)

	if len(y) != 0 {
		for i := range y {
			z = append(z, cidrNode{
				Service: "ec2",
				Type:    []string{"Cidr"},
				Name:    i,
			})
		}
		c.dgraphAddNodes(z)
		c.stats.NumberOfNodes += len(z)
	}

	m := make(map[string]securityGroupNodes)
	n := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("SecurityGroup"), &m)
	for _, i := range m["list"] {
		n[i.GroupID] = i.UID
	}
	c.ressources["SecurityGroups"] = n

	l := make(map[string]string)
	for _, i := range m["list"] {
		l[i.GroupID] = i.UID + "%" + x[i.GroupID]
	}
	c.ressources["SecurityGroupEdges"] = l

	o := make(map[string]cidrNodes)
	p := make(map[string]string)
	json.Unmarshal(c.dgraphQuery("Cidr"), &o)
	for _, i := range o["list"] {
		p[i.Name] = i.UID + "%" + y[i.Name]
	}
	c.ressources["Cidr"] = p
}

func (list securityGroupList) addEdges(c *connector) {
	defer c.waitGroup.Done()

	if len(list.SecurityGroups) == 0 {
		return
	}
	log.Println("Add SecurityGroup Edges")
	a := securityGroupNodes{}
	for _, i := range list.SecurityGroups {
		b := securityGroupNode{
			UID: c.ressources["SecurityGroups"][*i.GroupId],
		}
		if i.VpcId != nil {
			b.Vpc = vpcNode{UID: c.ressources["Vpcs"][*i.VpcId]}
		}
		for _, j := range i.IpPermissions {
			for _, k := range j.UserIdGroupPairs {
				s := strings.Split(c.ressources["SecurityGroupEdges"][*k.GroupId], "%")
				if len(s) >= 2 {
					t := strings.Split(s[1], "!")
					var u string
					for _, l := range t {
						m := strings.Split(l, "#")
						if m[0] == *i.GroupId {
							u += m[1] + ","
						}
					}
					b.SecurityGroup = append(b.SecurityGroup, securityGroupNode{UID: s[0], SecurityGroupPortTCP: u})
				}
			}
			for _, k := range j.IpRanges {
				s := strings.Split(c.ressources["Cidr"][*k.CidrIp], "%")
				if len(s) >= 2 {
					t := strings.Split(s[1], "!")
					var u string
					for _, l := range t {
						m := strings.Split(l, "#")
						if m[0] == *i.GroupId {
							u += m[1] + ","
						}
					}
					b.Cidr = append(b.Cidr, cidrNode{UID: s[0], CidrPortTCP: u})
				}
			}
		}
		a = append(a, b)
	}
	c.dgraphAddNodes(a)
}
