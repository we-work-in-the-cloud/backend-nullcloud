package model

import "time"

type VPC struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CRN       string    `json:"crn"`
	CreatedAt time.Time `json:"created_at"`
}

type Subnet struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CRN       string    `json:"crn"`
	VPCID     string    `json:"vpc_id"`
	CIDRBlock string    `json:"cidr_block"`
	CreatedAt time.Time `json:"created_at"`
}

type VSI struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CRN       string    `json:"crn"`
	SubnetID  string    `json:"subnet_id"`
	VPCID     string    `json:"vpc_id"`
	Profile   string    `json:"profile"`
	Image     string    `json:"image"`
	PrimaryIP string    `json:"primary_ip"`
	CreatedAt time.Time `json:"created_at"`
}

type LoadBalancerTarget struct {
	Type string `json:"type"` // "cluster" or "vsi"
	ID   string `json:"id"`
}

type LoadBalancer struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Status    string               `json:"status"`
	CRN       string               `json:"crn"`
	Protocol  string               `json:"protocol"`
	Port      int                  `json:"port"`
	Targets   []LoadBalancerTarget `json:"targets"`
	CreatedAt time.Time            `json:"created_at"`
}

type Bucket struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CRN       string    `json:"crn"`
	Region    string    `json:"region"`
	CreatedAt time.Time `json:"created_at"`
}

type Database struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CRN       string    `json:"crn"`
	Engine    string    `json:"engine"`
	Version   string    `json:"version"`
	Plan      string    `json:"plan"`
	SubnetIDs []string  `json:"subnet_ids"`
	CreatedAt time.Time `json:"created_at"`
}

type KubernetesCluster struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CRN       string    `json:"crn"`
	Version   string    `json:"version"`
	NodeCount int       `json:"node_count"`
	SubnetIDs []string  `json:"subnet_ids"`
	CreatedAt time.Time `json:"created_at"`
}
