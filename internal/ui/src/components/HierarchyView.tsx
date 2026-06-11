import type { VPC, Subnet, Database, KubernetesCluster, VSI, Bucket, LoadBalancer } from '../types/api'

interface HierarchyViewProps {
  vpcs: VPC[]
  subnets: Subnet[]
  databases: Database[]
  clusters: KubernetesCluster[]
  instances: VSI[]
  buckets: Bucket[]
  loadbalancers: LoadBalancer[]
}

export default function HierarchyView({ vpcs, subnets, databases, clusters, instances, buckets, loadbalancers }: HierarchyViewProps) {
  const groupByRegion = () => {
    const regions: { [key: string]: VPC[] } = {}
    vpcs.forEach((vpc) => {
      if (!regions[vpc.region]) {
        regions[vpc.region] = []
      }
      regions[vpc.region].push(vpc)
    })
    return regions
  }

  const getSubnetsForVpc = (vpcId: string) => subnets.filter((s) => s.vpc_id === vpcId)
  const getSubnetsByZone = (vpcSubnets: Subnet[]) => {
    const zones: { [key: string]: Subnet[] } = {}
    vpcSubnets.forEach((subnet) => {
      if (!zones[subnet.zone]) {
        zones[subnet.zone] = []
      }
      zones[subnet.zone].push(subnet)
    })
    return zones
  }
  const getInstancesForSubnet = (subnetId: string) => instances.filter((i) => i.subnet_id === subnetId)
  const getDatabasesForSubnet = (subnetId: string) =>
    databases.filter((d) => d.subnet_ids && d.subnet_ids.includes(subnetId))
  const getClustersForSubnet = (subnetId: string) =>
    clusters.filter((c) => c.subnet_ids && c.subnet_ids.includes(subnetId))
  const getBucketsForRegion = (region: string) => buckets.filter((b) => b.region === region)

  const regions = groupByRegion()

  if (vpcs.length === 0 && buckets.length === 0 && loadbalancers.length === 0) {
    return (
      <div className="tab-panel">
        <div style={{ padding: '60px 20px', textAlign: 'center' }}>
          <div style={{ marginBottom: '16px', fontSize: '48px' }}>🗺️</div>
          <strong>No infrastructure yet</strong>
          <p style={{ color: 'var(--text-2)', marginTop: '8px' }}>Create VPCs and buckets to visualize your infrastructure hierarchy</p>
        </div>
      </div>
    )
  }

  return (
    <div className="tab-panel">
      <div className="panel-header">
        <div>
          <h2>Infrastructure Hierarchy</h2>
          <p className="panel-desc">Visual representation of your regions, zones, VPCs, subnets, and resources.</p>
        </div>
      </div>
      <div style={{ padding: '24px', overflow: 'auto' }}>
        {Object.entries(regions).map(([region, regionVpcs]) => {
          const regionBuckets = getBucketsForRegion(region)
          const regionLBs = loadbalancers.filter((lb) => {
            // Load balancers target instances/clusters which may be in this region
            // For now, associate LBs with all regions if they have targets
            return lb.targets && lb.targets.length > 0
          })
          return (
            <div key={region} style={{ marginBottom: '32px' }}>
              {/* Region Header */}
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '12px',
                  marginBottom: '16px',
                  paddingBottom: '12px',
                  borderBottom: '2px solid var(--brand)',
                }}
              >
                <div style={{ fontSize: '28px' }}>📍</div>
                <div>
                  <div style={{ fontSize: '18px', fontWeight: 600, color: 'var(--text)' }}>{region}</div>
                  <div style={{ fontSize: '12px', color: 'var(--text-3)' }}>
                    {regionVpcs.length} VPC{regionVpcs.length !== 1 ? 's' : ''} • {regionBuckets.length} Bucket{regionBuckets.length !== 1 ? 's' : ''} • {regionLBs.length} LB{regionLBs.length !== 1 ? 's' : ''}
                  </div>
                </div>
              </div>

              {/* VPCs and Buckets Grid */}
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(320px, 1fr))', gap: '16px', marginBottom: '24px' }}>
                {/* VPCs */}
                {regionVpcs.map((vpc) => {
                  const vpcSubnets = getSubnetsForVpc(vpc.id)
                  const subnetsByZone = getSubnetsByZone(vpcSubnets)

                  return (
                    <div
                      key={vpc.id}
                      style={{
                        border: '2px solid var(--border)',
                        borderRadius: 'var(--r)',
                        padding: '16px',
                        backgroundColor: 'var(--surface)',
                        boxShadow: 'var(--sh-sm)',
                      }}
                    >
                      {/* VPC Header */}
                      <div style={{ marginBottom: '12px' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '4px' }}>
                          <span style={{ fontSize: '20px' }}>🌐</span>
                          <div style={{ fontWeight: 600, fontSize: '14px', color: 'var(--text)' }}>{vpc.name}</div>
                        </div>
                        <div style={{ fontSize: '11px', color: 'var(--text-3)', fontFamily: 'monospace' }}>{vpc.id}</div>
                      </div>

                      {/* Zones */}
                      {Object.entries(subnetsByZone).length === 0 ? (
                        <div style={{ padding: '12px', textAlign: 'center', color: 'var(--text-3)', fontSize: '12px' }}>
                          No zones
                        </div>
                      ) : (
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                          {Object.entries(subnetsByZone).map(([zone, zoneSubnets]) => (
                            <div
                              key={zone}
                              style={{
                                border: '1px solid var(--border)',
                                borderRadius: 'var(--r-sm)',
                                padding: '10px',
                                backgroundColor: 'var(--surface-2)',
                              }}
                            >
                              {/* Zone Header */}
                              <div style={{ marginBottom: '8px', paddingBottom: '6px', borderBottom: '1px solid var(--border-sub)' }}>
                                <div style={{ fontSize: '12px', fontWeight: 600, color: 'var(--brand)', display: 'flex', alignItems: 'center', gap: '4px' }}>
                                  <span>🌍</span>
                                  {zone}
                                </div>
                              </div>

                              {/* Subnets in Zone */}
                              <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
                                {zoneSubnets.map((subnet) => {
                                  const subnetInstances = getInstancesForSubnet(subnet.id)
                                  const subnetDbs = getDatabasesForSubnet(subnet.id)
                                  const subnetClusters = getClustersForSubnet(subnet.id)
                                  const hasResources = subnetInstances.length > 0 || subnetDbs.length > 0 || subnetClusters.length > 0

                                  return (
                                    <div
                                      key={subnet.id}
                                      style={{
                                        border: '1px solid var(--border)',
                                        borderRadius: 'var(--r-xs)',
                                        padding: '8px',
                                        backgroundColor: 'var(--surface)',
                                      }}
                                    >
                                      {/* Subnet Header */}
                                      <div style={{ marginBottom: hasResources ? '6px' : '0' }}>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '4px', marginBottom: '2px' }}>
                                          <span style={{ fontSize: '14px' }}>📊</span>
                                          <div style={{ fontWeight: 500, fontSize: '12px', color: 'var(--text)' }}>{subnet.name}</div>
                                        </div>
                                        <div style={{ fontSize: '10px', color: 'var(--text-3)', marginLeft: '18px' }}>
                                          {subnet.cidr_block}
                                        </div>
                                      </div>

                                      {/* Resources in Subnet */}
                                      {hasResources && (
                                        <div style={{ marginTop: '6px', paddingTop: '6px', borderTop: '1px solid var(--border-sub)' }}>
                                          {subnetInstances.length > 0 && (
                                            <div style={{ marginBottom: subnetDbs.length > 0 || subnetClusters.length > 0 ? '4px' : '0' }}>
                                              <div style={{ fontSize: '10px', fontWeight: 500, color: 'var(--text-2)', marginBottom: '2px' }}>
                                                Instances ({subnetInstances.length})
                                              </div>
                                              {subnetInstances.map((inst) => (
                                                <div
                                                  key={inst.id}
                                                  style={{
                                                    fontSize: '11px',
                                                    padding: '2px 4px',
                                                    backgroundColor: 'rgba(59, 130, 246, 0.1)',
                                                    borderRadius: 'var(--r-xs)',
                                                    marginBottom: '1px',
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    gap: '4px',
                                                  }}
                                                >
                                                  <span>💻</span>
                                                  <span style={{ color: 'var(--text)', fontWeight: 500 }}>{inst.name}</span>
                                                </div>
                                              ))}
                                            </div>
                                          )}
                                          {subnetDbs.length > 0 && (
                                            <div style={{ marginBottom: subnetClusters.length > 0 ? '4px' : '0' }}>
                                              <div style={{ fontSize: '10px', fontWeight: 500, color: 'var(--text-2)', marginBottom: '2px' }}>
                                                Databases ({subnetDbs.length})
                                              </div>
                                              {subnetDbs.map((db) => (
                                                <div
                                                  key={db.id}
                                                  style={{
                                                    fontSize: '11px',
                                                    padding: '2px 4px',
                                                    backgroundColor: 'rgba(20, 184, 166, 0.1)',
                                                    borderRadius: 'var(--r-xs)',
                                                    marginBottom: '1px',
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    gap: '4px',
                                                  }}
                                                >
                                                  <span>🗄️</span>
                                                  <span style={{ color: 'var(--text)', fontWeight: 500 }}>{db.name}</span>
                                                </div>
                                              ))}
                                            </div>
                                          )}
                                          {subnetClusters.length > 0 && (
                                            <div>
                                              <div style={{ fontSize: '10px', fontWeight: 500, color: 'var(--text-2)', marginBottom: '2px' }}>
                                                Clusters ({subnetClusters.length})
                                              </div>
                                              {subnetClusters.map((cluster) => (
                                                <div
                                                  key={cluster.id}
                                                  style={{
                                                    fontSize: '11px',
                                                    padding: '2px 4px',
                                                    backgroundColor: 'rgba(168, 85, 247, 0.1)',
                                                    borderRadius: 'var(--r-xs)',
                                                    marginBottom: '1px',
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    gap: '4px',
                                                  }}
                                                >
                                                  <span>☸️</span>
                                                  <span style={{ color: 'var(--text)', fontWeight: 500 }}>{cluster.name}</span>
                                                </div>
                                              ))}
                                            </div>
                                          )}
                                        </div>
                                      )}
                                    </div>
                                  )
                                })}
                              </div>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  )
                })}

                {/* Buckets */}
                {regionBuckets.map((bucket) => (
                  <div
                    key={bucket.id}
                    style={{
                      border: '2px solid var(--border)',
                      borderRadius: 'var(--r)',
                      padding: '16px',
                      backgroundColor: 'var(--surface)',
                      boxShadow: 'var(--sh-sm)',
                      display: 'flex',
                      flexDirection: 'column',
                      justifyContent: 'center',
                      alignItems: 'center',
                      textAlign: 'center',
                      minHeight: '140px',
                    }}
                  >
                    <div style={{ fontSize: '40px', marginBottom: '8px' }}>🪣</div>
                    <div style={{ fontWeight: 600, fontSize: '14px', color: 'var(--text)', marginBottom: '4px' }}>
                      {bucket.name}
                    </div>
                    <div style={{ fontSize: '11px', color: 'var(--text-3)', fontFamily: 'monospace' }}>
                      {bucket.id}
                    </div>
                  </div>
                ))}

                {/* Load Balancers */}
                {regionLBs.map((lb) => (
                  <div
                    key={lb.id}
                    style={{
                      border: '2px solid var(--border)',
                      borderRadius: 'var(--r)',
                      padding: '16px',
                      backgroundColor: 'var(--surface)',
                      boxShadow: 'var(--sh-sm)',
                    }}
                  >
                    <div style={{ marginBottom: '12px' }}>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '4px' }}>
                        <span style={{ fontSize: '20px' }}>⚖️</span>
                        <div style={{ fontWeight: 600, fontSize: '14px', color: 'var(--text)' }}>{lb.name}</div>
                      </div>
                      <div style={{ fontSize: '11px', color: 'var(--text-3)', fontFamily: 'monospace', marginBottom: '8px' }}>
                        {lb.id}
                      </div>
                      <div style={{ fontSize: '11px', color: 'var(--text-2)', display: 'flex', gap: '12px' }}>
                        <span>{lb.protocol.toUpperCase()}</span>
                        <span>:{lb.port}</span>
                      </div>
                    </div>

                    {lb.targets && lb.targets.length > 0 && (
                      <div style={{ marginTop: '12px', paddingTop: '12px', borderTop: '1px solid var(--border)' }}>
                        <div style={{ fontSize: '11px', fontWeight: 500, color: 'var(--text-2)', marginBottom: '6px' }}>
                          Targets ({lb.targets.length})
                        </div>
                        {lb.targets.map((target, idx) => (
                          <div
                            key={idx}
                            style={{
                              fontSize: '11px',
                              padding: '4px 6px',
                              backgroundColor: 'rgba(251, 146, 60, 0.1)',
                              borderRadius: 'var(--r-xs)',
                              marginBottom: '2px',
                              color: 'var(--text)',
                            }}
                          >
                            {target.type === 'vsi' ? '💻' : '☸️'} {target.id.substring(0, 8)}...
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
