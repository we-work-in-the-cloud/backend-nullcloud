export interface VPC {
  id: string
  name: string
  status: string
  crn: string
  region: string
  created_at: string
}

export interface Subnet {
  id: string
  name: string
  status: string
  crn: string
  vpc_id: string
  zone: string
  cidr_block: string
  created_at: string
}

export interface VSI {
  id: string
  name: string
  status: string
  crn: string
  subnet_id: string
  profile: string
  image: string
  primary_ip: string
  created_at: string
}

export interface LoadBalancer {
  id: string
  name: string
  status: string
  crn: string
  protocol: string
  port: number
  targets: LoadBalancerTarget[]
  created_at: string
}

export interface LoadBalancerTarget {
  type: string
  id: string
}

export interface Bucket {
  id: string
  name: string
  status: string
  crn: string
  region: string
  created_at: string
}

export interface Database {
  id: string
  name: string
  status: string
  crn: string
  engine: string
  version: string
  plan: string
  subnet_ids: string[]
  created_at: string
  endpoint: string
}

export interface KubernetesCluster {
  id: string
  name: string
  status: string
  crn: string
  version: string
  node_count: number
  subnet_ids: string[]
  created_at: string
}

export interface APIResponse<T> {
  [key: string]: T | T[]
}

export interface APIError {
  code: string
  message: string
}
