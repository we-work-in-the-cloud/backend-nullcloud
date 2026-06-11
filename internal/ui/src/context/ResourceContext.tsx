import React, { ReactNode, useState, useEffect } from 'react'
import { VPC, Subnet, VSI, LoadBalancer, Bucket, Database, KubernetesCluster } from '../types'
import { useAPI } from '../hooks/useAPI'
import { useAuth } from './AuthContext'

interface ResourceContextType {
  vpcs: VPC[]
  subnets: Subnet[]
  instances: VSI[]
  loadbalancers: LoadBalancer[]
  buckets: Bucket[]
  databases: Database[]
  clusters: KubernetesCluster[]
  loading: boolean
  error: string | null
  refresh: () => Promise<void>
}

export const ResourceContext = React.createContext<ResourceContextType | undefined>(undefined)

export function ResourceProvider({ isConnected, children }: { isConnected: boolean; children: ReactNode }) {
  const { token } = useAuth()
  const { fetch_api } = useAPI(token)
  const [vpcs, setVpcs] = useState<VPC[]>([])
  const [subnets, setSubnets] = useState<Subnet[]>([])
  const [instances, setInstances] = useState<VSI[]>([])
  const [loadbalancers, setLoadbalancers] = useState<LoadBalancer[]>([])
  const [buckets, setBuckets] = useState<Bucket[]>([])
  const [databases, setDatabases] = useState<Database[]>([])
  const [clusters, setClusters] = useState<KubernetesCluster[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const refresh = async () => {
    if (!isConnected) return

    setLoading(true)
    setError(null)

    try {
      const [vpcsData, subnetsData, instancesData, lbData, bucketData, dbData, clusterData] = await Promise.all([
        fetch_api<{ vpcs: VPC[] }>('/v1/vpcs'),
        fetch_api<{ subnets: Subnet[] }>('/v1/subnets'),
        fetch_api<{ instances: VSI[] }>('/v1/instances'),
        fetch_api<{ load_balancers: LoadBalancer[] }>('/v1/loadbalancers'),
        fetch_api<{ buckets: Bucket[] }>('/v1/buckets'),
        fetch_api<{ databases: Database[] }>('/v1/databases'),
        fetch_api<{ clusters: KubernetesCluster[] }>('/v1/clusters'),
      ])

      setVpcs(vpcsData.vpcs || [])
      setSubnets(subnetsData.subnets || [])
      setInstances(instancesData.instances || [])
      setLoadbalancers(lbData.load_balancers || [])
      setBuckets(bucketData.buckets || [])
      setDatabases(dbData.databases || [])
      setClusters(clusterData.clusters || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load resources')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (isConnected) {
      refresh()
    }
  }, [isConnected])

  return (
    <ResourceContext.Provider
      value={{
        vpcs,
        subnets,
        instances,
        loadbalancers,
        buckets,
        databases,
        clusters,
        loading,
        error,
        refresh,
      }}
    >
      {children}
    </ResourceContext.Provider>
  )
}

export function useResources() {
  const context = React.useContext(ResourceContext)
  if (!context) {
    throw new Error('useResources must be used within ResourceProvider')
  }
  return context
}
