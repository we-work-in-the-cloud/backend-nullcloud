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
	VPCID     string    `json:"vpc_id"`
	CIDRBlock string    `json:"cidr_block"`
	CreatedAt time.Time `json:"created_at"`
}

type VSI struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	SubnetID  string    `json:"subnet_id"`
	VPCID     string    `json:"vpc_id"`
	Profile   string    `json:"profile"`
	Image     string    `json:"image"`
	PrimaryIP string    `json:"primary_ip"`
	CreatedAt time.Time `json:"created_at"`
}
