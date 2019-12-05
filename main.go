package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	profile := flag.String("profile", "default", "Profile from ~/.aws/config")
	region := flag.String("region", "eu-west-1", "AWS Region")
	dgraph := flag.String("dgraph", "127.0.0.1:9080", "Dgraph server (ip:port)")
	dtype := flag.String("type", "", "Get the schema for a type (only after importing some data)")
	list := flag.Bool("list", false, "List available ressource types")
	drop := flag.Bool("drop", false, "Drop all nodes and the schema")
	noschema := flag.Bool("no-schema", false, "Disable the refresh schema at each run")
	flag.Parse()

	if *list == true {
		fmt.Println("Address, AutoScalingGroup, AvailabilityZone, CacheCluster, CacheSubnetGroup, Cidr, DbCluster, DbClusterParameterGroup, DbInstance, DbParameterGroup, DbSubnetGroup, Instance, InstanceProfile, KeyPair, LaunchConfiguration, LaunchTemplate, LoadBalancer, NatGateway, OptionGroup, SecurityGroup, Subnet, TargetGroup, Volume, Vpc, VpcPeeringConnection")
		os.Exit(0)
	}

	connector := newConnector(profile, region, dgraph)
	defer connector.grpcConnexion.Close()

	if *drop == true {
		connector.dgraphDropAll()
		os.Exit(0)
	}

	if *dtype != "" {
		dgraphDisplaySchema(dgraph, dtype)
		os.Exit(0)
	}

	if *noschema != true {
		connector.dgraphAddSchema()
	}

	connector.dgraphDropPrevious()

	var instances instanceList
	var keypairs keyPairList
	var volumes volumeList
	var addresses addressList
	var availabilityzones availabilityZoneList
	var vpcs vpcList
	var instanceprofiles instanceProfileList
	var autoscalinggroups autoScalingGroupList
	var launchconfigurations launchConfigurationList
	var launchtemplates launchTemplateList
	var targetgroups targetGroupList
	var loadbalancers loadBalancerList
	var loadbalancersv2 loadBalancerV2List
	var subnets subnetList
	var securitygroups securityGroupList
	var dbinstances dbInstanceList
	var dbclusters dbClusterList
	var optiongroups optionGroupList
	var dbparametergroups dbParameterGroupList
	var dbclusterparametergroups dbClusterParameterGroupList
	var dbsubnetgroups dbSubnetGroupList
	var cacheclusters cacheClusterList
	var cachesubnetgroups cacheSubnetGroupList
	var vpcpeeringconnections vpcPeeringConnectionList
	var natgateways natGatewayList
	var snapshots snapshotList
	var images imageList

	// List ressources
	connector.waitGroup.Add(27)
	start := time.Now()
	go func() { instances = connector.listInstances() }()
	go func() { keypairs = connector.listKeyPairs() }()
	go func() { volumes = connector.listVolumes() }()
	go func() { addresses = connector.listAddresses() }()
	go func() { availabilityzones = connector.listAvailabilityZones() }()
	go func() { vpcs = connector.listVpcs() }()
	go func() { instanceprofiles = connector.listInstanceProfiles() }()
	go func() { autoscalinggroups = connector.listAutoScalingGroups() }()
	go func() { launchconfigurations = connector.listLaunchConfigurations() }()
	go func() { launchtemplates = connector.listLaunchTemplates() }()
	go func() { targetgroups = connector.listTargetGroups() }()
	go func() { loadbalancers = connector.listLoadBalancers() }()
	go func() { loadbalancersv2 = connector.listLoadBalancersV2() }()
	go func() { subnets = connector.listSubnets() }()
	go func() { securitygroups = connector.listSecurityGroups() }()
	go func() { dbinstances = connector.listDbInstances() }()
	go func() { dbclusters = connector.listDbClusters() }()
	go func() { optiongroups = connector.listOptionGroups() }()
	go func() { dbparametergroups = connector.listDbParameterGroups() }()
	go func() { dbclusterparametergroups = connector.listDbClusterParameterGroups() }()
	go func() { dbsubnetgroups = connector.listDbSubnetGroups() }()
	go func() { cacheclusters = connector.listCacheClusters() }()
	go func() { cachesubnetgroups = connector.listCacheSubnetGroups() }()
	go func() { vpcpeeringconnections = connector.listVpcPeeringConnections() }()
	go func() { natgateways = connector.listNatGateways() }()
	go func() { snapshots = connector.listSnapshots() }()
	go func() { images = connector.listImages() }()

	connector.waitGroup.Wait()

	// Add Nodes
	connector.waitGroup.Add(27)
	instances.addNodes(connector)
	keypairs.addNodes(connector)
	volumes.addNodes(connector)
	addresses.addNodes(connector)
	availabilityzones.addNodes(connector)
	vpcs.addNodes(connector)
	instanceprofiles.addNodes(connector)
	autoscalinggroups.addNodes(connector)
	launchconfigurations.addNodes(connector)
	launchtemplates.addNodes(connector)
	targetgroups.addNodes(connector)
	loadbalancers.addNodes(connector)
	loadbalancersv2.addNodes(connector)
	subnets.addNodes(connector)
	securitygroups.addNodes(connector)
	dbinstances.addNodes(connector)
	dbclusters.addNodes(connector)
	optiongroups.addNodes(connector)
	dbparametergroups.addNodes(connector)
	dbclusterparametergroups.addNodes(connector)
	dbsubnetgroups.addNodes(connector)
	cacheclusters.addNodes(connector)
	cachesubnetgroups.addNodes(connector)
	vpcpeeringconnections.addNodes(connector)
	natgateways.addNodes(connector)
	snapshots.addNodes(connector)
	images.addNodes(connector)

	connector.waitGroup.Wait()

	// Add Edges
	connector.waitGroup.Add(19)
	instances.addEdges(connector)
	addresses.addEdges(connector)
	volumes.addEdges(connector)
	autoscalinggroups.addEdges(connector)
	launchconfigurations.addEdges(connector)
	loadbalancers.addEdges(connector)
	loadbalancersv2.addEdges(connector)
	targetgroups.addEdges(connector)
	subnets.addEdges(connector)
	securitygroups.addEdges(connector)
	dbinstances.addEdges(connector)
	dbsubnetgroups.addEdges(connector)
	dbclusters.addEdges(connector)
	cacheclusters.addEdges(connector)
	cachesubnetgroups.addEdges(connector)
	vpcpeeringconnections.addEdges(connector)
	natgateways.addEdges(connector)
	snapshots.addEdges(connector)
	images.addEdges(connector)

	connector.waitGroup.Wait()

	log.Printf("%v Nodes have been imported in %s\n", connector.stats.NumberOfNodes, time.Since(start))
}
