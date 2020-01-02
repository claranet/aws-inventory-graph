package main

func getdgraphSchema() string {
	s := `
	name: string @index(exact) .
	OwnerId: string @index(term) .
	OwnerName: string @index(term) .
	Region: string @index(term) .
	Service: string @index(term) .
	AvailabilityZone: string @index(term) .
	State: string @index(term) .
	Status: string @index(term) .
	PrivateIpAddress: string @index(term) .
	PublicIpAddress: string @index(term) .
	Arn: string @index(term) .
	Description: string @index(term) .
	CidrBlock: string @index(term) .
	_Vpc: uid @reverse .
	_Instance: [uid] @reverse .
	_AvailabilityZone: [uid] @reverse .
	_TargetGroup: [uid] @reverse .
	_LoadBalancer: [uid] @reverse .
	_Image: uid @reverse .
	_Subnet: [uid] @reverse .
	_SecurityGroup: [uid] @reverse .
	_Snapshot: [uid] @reverse .
	
	InstanceId: string @index(term) .
	InstanceType: string @index(term) .
	EbsOptimized: bool .
	Hypervisor: string @index(term) .
	VirtualizationType: string @index(term) .
	_KeyName: uid @reverse .
	_InstanceProfile: uid @reverse .

	type Instance {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		InstanceId: string
		InstanceType: string
		State: string
		EbsOptimized: bool
		Hypervisor: string
		VirtualizationType: string
		PrivateIpAddress: string
		PublicIpAddress: string
		_KeyName: KeyPair
		_AvailabilityZone: AvailabilityZone
		_Vpc: Vpc
		_InstanceProfile: InstanceProfile
		_Image: Image
		_Subnet: Subnet
	}

	Encrypted: string @index(term) .
	VolumeId: string @index(term) .
	VolumeType: string @index(term) .
	AvailabilityZone: string @index(term) .
	Iops: int .
	Size: int .
	Device: string @index(term) .
	DeleteOnTermination: bool .

	type Volume {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		VolumeId: string
		DeleteOnTermination: bool
		Device: string
		Encrypted: bool
		VolumeType: string
		Size: int
		Iops: int
		State: string
		_AvailabilityZone: AvailabilityZone
		_Instance: Instance
		_Snapshot: Snapshot
	}

	PublicIp: string @index(term) .
	Domain: string @index(term) .
	AllocationId: string @index(term) . 
	_NatGateway: [uid] @reverse .

	type Address {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		PrivateIpAddress: string
		PublicIp: string
		Domain: string
		AllocationId: string
		_Instance: Instance
		_NatGateway: NatGateway
	}

	ZoneName: string @index(term) .
	ZoneId: string @index(term) . 

	type AvailabilityZone {
		name: string
		Service: string
		Region: string
		State: string
		ZoneName: string
		ZoneId: string
	}

	KeyName: string @index(term) .
	KeyFingerprint: string @index(term) . 

	type KeyPair {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		State: string
		KeyName: string
		KeyFingerprint: string
	}

	VpcId: string @index(term) .

	type Vpc {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		VpcId: string
		CidrBlock: string
	}

	InstanceProfileId: string @index(term) .
	InstanceProfileName: string @index(term) .

	type InstanceProfile {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		InstanceProfileId: string
		InstanceProfileName: string
		Arn: string
	}

	AutoScalingGroupArn: string @index(term) .
	AutoScalingGroupName: string @index(term) .
	MinSize: int .
	MaxSize: int .
	DesiredCapacity: int .
	DefaultCooldown: int .
	HealthCheckGracePeriod: int .
	HealthCheckType: string @index(term) .
	NewInstancesProtectedFromScaleIn: bool .
	TerminationPolicy: string @index(term) .
	_LaunchConfiguration: uid @reverse .
	_LaunchTemplate: uid @reverse .
	
	type AutoScalingGroup {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		AutoScalingGroupArn: string
		AutoScalingGroupName: string
		TerminationPolicy: string
		MinSize: int
		MaxSize: int
		DesiredCapacity: int
		DefaultCooldown: int
		HealthCheckGracePeriod: int
		HealthCheckType: string
		NewInstancesProtectedFromScaleIn: bool
		_Instance: Instance
		_TargetGroup: TargetGroup
		_AvailabilityZone: AvailabilityZone
		_LoadBalancer: LoadBalancer
		_LaunchConfiguration: LaunchConfiguration
		_LaunchTemplate: LaunchTemplate
	}

	LaunchConfigurationArn: string @index(term) .
	LaunchConfigurationName: string @index(term) .
	UserData: string @index(term) .
	AssociatePublicIpAddress: bool .

	type LaunchConfiguration {
		name: string
		Service: string
		Region: string
		OwnerId: string
		LaunchConfigurationArn: string
		LaunchConfigurationName: string
		UserData: string
		AssociatePublicIpAddress: bool
		InstanceType: string
		EbsOptimized: bool
		_KeyName: KeyPair
		_InstanceProfile: InstanceProfile
		_Image: Image
		_SecurityGroup: SecurityGroup
	}

	LaunchTemplateId: string @index(term) .
	LaunchTemplateName: string @index(term) .
	LatestVersionNumber: int .
	DefaultVersionNumber: int .
	CreatedBy: string @index(term) .

	type LaunchTemplate {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		LaunchTemplateId: string
		LaunchTemplateName: string
		LatestVersionNumber: int
		DefaultVersionNumber: int
		CreatedBy: string
	}

	TargetGroupArn: string @index(term) .
	TargetGroupName: string @index(term) .
	TargetType: string @index(term) .
	HealthCheckPath: string @index(term) .
	HealthCheckPort: string @index(term) .
	HealthyThresholdCount: int .
	Port: int .
	UnhealthyThresholdCount: int .
	Protocol: string @index(term) .

	type TargetGroup {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		TargetGroupArn: string
		TargetGroupName: string
		TargetType: string
		HealthCheckPath: string
		HealthCheckPort: string
		HealthyThresholdCount: int
		Port: int
		UnhealthyThresholdCount: int
		Protocol: string
		_Vpc: Vpc
		_LoadBalancer: LoadBalancer
	}

	LoadBalancerName: string @index(term) .
	LoadBalancerArn: string @index(term) .
	CanonicalHostedZoneNameId: string @index(term) .
	CanonicalHostedZoneId: string @index(term) .
	CanonicalHostedZoneName: string @index(term) .
	DNSName: string @index(term) .
	Scheme: string @index(term) .
	LoadBalancerType: string @index(term) .
	IPAddressType: string @index(term) .

	type LoadBalancer {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		LoadBalancerName: string
		LoadBalancerArn: string
		CanonicalHostedZoneNameId: string
		CanonicalHostedZoneId: string
		CanonicalHostedZoneName: string
		DNSName: string
		Scheme: string
		LoadBalancerType: string
		IPAddressType: string
		State: string
		_Vpc: Vpc
		_Instance: Instance
		_AvailabilityZone: AvailabilityZone
	}

	SnapshotId: string @index(term) .
	VolumeSize: int .
	Encrypted: bool .
	Progress: string @index(term) .
	_Volume: uid @reverse .
	DeviceName: string @index(term) .

	type Snapshot {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		VpcId: string
		SnapshotId: string
		Description: string
		VolumeSize: int
		Encrypted: bool
		Progress: string
		DeviceName: string
		_Volume: Volume
	}

	ImageId: string @index(term) .
	VirtualizationType: string @index(term) .
	Hypervisor: string @index(term) .
	EnaSupport: bool .
	SriovNetSupport: string @index(term) .
	State: string @index(term) .
	Architecture: string @index(term) .
	ImageLocation: string @index(term) .
	RootDeviceType: string @index(term) .
	RootDeviceName: string @index(term) .
	Public: bool .
	ImageType: string @index(term) .

	type Image {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		ImageId: string
		VirtualizationType: string
		Hypervisor: string
		EnaSupport: bool
		SriovNetSupport: string
		State: string
		Architecture: string
		ImageLocation: string
		RootDeviceType: string
		RootDeviceName: string
		Public: bool
		ImageType: string
		_Snapshot: Snapshot
	}

	SubnetId: string @index(term) .
	MapPublicIPOnLaunch: bool .
	DefaultForAz: bool .
	AssignIPv6AddressOnCreation: bool .
	AvailableIPAddressCount: string @index(term) .
	PortTcp: string @index(term) .
	PortUdp: string @index(term) .

	type Subnet {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		SubnetId: string
		MapPublicIPOnLaunch: bool
		DefaultForAz: bool
		AssignIPv6AddressOnCreation: bool
		AvailableIPAddressCount: string
		State: string
		CidrBlock: string
		_Vpc: Vpc
		_AvailabilityZone: AvailabilityZone
	}

	GroupId: string @index(term) .
	GroupName: string @index(term) .
	PortTcp: string .
	_Cidr: [uid] @reverse .
	
	type SecurityGroup {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		GroupId: string
		GroupName: string
		PortTcp: string
		_Vpc: Vpc
		_SecurityGroup: SecurityGroup
		_Cidr: Cidr
	}
	
	type Cidr {
		name: string
		Service: string
		PortTcp: string
	}

	DbInstanceArn: string @index(term) .
	DbInstanceIdentifier: string @index(term) .
	DBName: string @index(term) .
	DbInstanceStatus: string @index(term) .
	DbInstanceClass: string @index(term) .
	DbiResourceID: string @index(term) .
	CACertificateIdentifier: string @index(term) .
	PubliclyAccessible: bool .
	MasterUsername: string @index(term) .
	LicenseModel: string @index(term) .
	CopyTagsToSnapshot: bool .
	Engine: string @index(term) .
	EngineVersion: string @index(term) .
	Endpoint: string @index(term) .
	Port: int .
	MultiAZ: bool .
	AutoMinorVersionUpgrade: bool .
	IAMDatabaseAuthenticationEnabled: bool .
	DeletionProtection: bool .
	PreferredBackupWindow: string @index(term) .
	PreferredMaintenanceWindow: string @index(term) .
	PromotionTier: int .
	AllocatedStorage: int .
	StorageType: string @index(term) .
	StorageEncrypted: bool .
	BackupRetentionPeriod: int .
	IsClusterWriter: bool .			
	_DbSubnetGroup: [uid] @reverse . 
	_DbParameterGroup: uid @reverse . 
	_OptionGroup: uid @reverse .
	_SecondaryAvailabilityZone: uid @reverse .
	_DbClusterParameterGroup: uid @reverse .

	type DbInstance {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		DbInstanceArn: string
		DbInstanceIdentifier: string
		DBName: string
		DbInstanceStatus: string
		DbInstanceClass: string
		DbiResourceID: string
		CACertificateIdentifier: string
		PubliclyAccessible: bool
		MasterUsername: string
		LicenseModel: string
		CopyTagsToSnapshot: bool
		Engine: string
		EngineVersion: string
		Endpoint: string
		Port: int
		MultiAZ: bool
		AutoMinorVersionUpgrade: bool
		IAMDatabaseAuthenticationEnabled: bool
		DeletionProtection: bool
		PreferredBackupWindow: string
		PreferredMaintenanceWindow: string
		PromotionTier: int
		AllocatedStorage: int
		StorageType: string
		StorageEncrypted: bool
		BackupRetentionPeriod: int
		_AvailabilityZone: AvailabilityZone
		_SecondaryAvailabilityZone: AvailabilityZone
		_SecurityGroup: SecurityGroup
		_Vpc: Vpc
		_DbSubnetGroup: DbSubnetGroup 
		_DbParameterGroup: DbParameterGroup
		_OptionGroup: OptionGroup
		_DbClusterParameterGroup: DbClusterParameterGroup
	}

	DbClusterArn: string @index(term) .
	DbClusterIdentifier: string @index(term) .
	DbClusterResourceID: string @index(term) .
	HostedZoneID: string @index(term) .
	HTTPEndpointEnabled: bool .
	EngineMode: string @index(term) .
	ReaderEndpoint: string @index(term) .
	_DbClusterMember: [uid] @reverse .
	_ReadReplica: [uid] @reverse .

	type DbCluster {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		DbClusterArn: string
		DbClusterIdentifier: string
		DbClusterResourceID: string
		HostedZoneID: string
		HTTPEndpointEnabled: bool
		DBName: string
		DbiResourceID: string
		MasterUsername: string
		Engine: string
		EngineMode: string
		EngineVersion: string
		Endpoint: string
		ReaderEndpoint: string
		Port: int
		MultiAZ: bool
		Status: string
		IAMDatabaseAuthenticationEnabled: bool
		DeletionProtection: bool
		PreferredBackupWindow: string
		PreferredMaintenanceWindow: string
		AllocatedStorage: int
		StorageType: string
		StorageEncrypted: bool
		BackupRetentionPeriod: int
		_AvailabilityZone: AvailabilityZone
		_SecurityGroup: SecurityGroup
		_DbClusterMember: DbInstance
		_ReadReplica: DbInstance
		_DbClusterParameterGroup: DbClusterParameterGroup
		_DbSubnetGroup: DbSubnetGroup 
	}

	OptionGroupArn: string @index(term) .
	OptionGroupName: string @index(term) .
	MajorEngineVersion: string @index(term) .
	EngineName: string @index(term) .
	AllowsVpcAndNonVpcInstanceMemberships: bool .

	type OptionGroup {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		OptionGroupArn: string
		OptionGroupName: string
		MajorEngineVersion: string
		EngineName: string
		AllowsVpcAndNonVpcInstanceMemberships: bool
	}

	DbParameterGroupArn: string @index(term) .
	DbParameterGroupName: string @index(term) .
	DbParameterGroupFamily: string @index(term) .

	type DbParameterGroup {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		DbParameterGroupArn: string
		DbParameterGroupName: string
		DbParameterGroupFamily: string
		Description: string
	}

	DbClusterParameterGroupArn: string @index(term) .
	DbClusterParameterGroupName: string @index(term) .
	DbClusterscription: string @index(term) .

	type DbClusterParameterGroup {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		DbClusterParameterGroupArn: string
		DbClusterParameterGroupName: string
		DbParameterGroupFamily: string
		DbClusterscription: string
	}

	DbSubnetGroupArn: string @index(term) .
	DbSubnetGroupName: string @index(term) .
	SubnetGroupStatus: string @index(term) .

	type DbSubnetGroup {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		DbSubnetGroupArn: string
		DbSubnetGroupName: string
		SubnetGroupStatus: string
		Description: string
	}

	CacheClusterId: string @index(term) .
	AuthTokenEnabled: bool .
	NumCacheNodes: int .
	TransitEncryptionEnabled: bool .
	CacheNodeType: string @index(term) .
	_CacheSubnetGroup: uid @reverse .

	type CacheCluster {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		CacheClusterId: string
		AuthTokenEnabled: bool
		NumCacheNodes: int
		TransitEncryptionEnabled: bool
		CacheNodeType: string
		Engine: string
		Endpoint: string
		Port: string
		EngineVersion: string
		AutoMinorVersionUpgrade: bool
		Status: string
		_AvailabilityZone: AvailabilityZone
		_SecurityGroup: SecurityGroup
		_Subnet: Subnet
	}

	CacheSubnetGroupName: string @index(term) .

	type CacheSubnetGroup {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		CacheSubnetGroupName: string
		Description: string
		_Vpc: Vpc
		_Subnet: Subnet
	}

	VpcPeeringConnectionId: string @index(term) .
	_AccepterVpc: uid @reverse .
	_RequesterVpc: uid @reverse .

	type VpcPeeringConnection {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		VpcPeeringConnectionId: string
		Status: string
		_AccepterVpc: Vpc
		_RequesterVpc: Vpc
	}

	NatGatewayId: string @index(term) .

	type NatGateway {
		name: string
		Service: string
		Region: string
		OwnerId: string
		OwnerName: string
		NatGatewayId: string
		PublicIpAddress: string
		PrivateIpAddress: string
		State: string
		_Vpc: Vpc
		_Subnet: Subnet
	}
`
	return s
}
