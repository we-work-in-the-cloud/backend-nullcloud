import { useState } from 'react'
import { useResources } from './context/ResourceContext'
import { CreateVPCForm, CreateSubnetForm, CreateInstanceForm, CreateLoadBalancerForm, CreateBucketForm, CreateDatabaseForm, CreateClusterForm } from './features/Forms'
import EditResourceModal from './components/EditResourceModal'
import HierarchyView from './components/HierarchyView'
import { useAPI } from './hooks/useAPI'
import { useAuth } from './context/AuthContext'

type TabName = 'vpcs' | 'subnets' | 'instances' | 'loadbalancers' | 'buckets' | 'databases' | 'clusters' | 'hierarchy'

interface EditingResource {
  id: string
  name: string
  type: string
}

interface DeletingResource {
  id: string
  type: string
}

export default function MainView() {
  const { vpcs, subnets, instances, loadbalancers, buckets, databases, clusters, loading, refresh } = useResources()
  const { token } = useAuth()
  const { fetch_api } = useAPI(token)
  const [activeTab, setActiveTab] = useState<TabName>('vpcs')
  const [createModalTab, setCreateModalTab] = useState<TabName | null>(null)
  const [editingResource, setEditingResource] = useState<EditingResource | null>(null)
  const [deletingResource, setDeletingResource] = useState<DeletingResource | null>(null)

  const getTabContent = () => {
    const handlers = {
      onCreateClick: () => setCreateModalTab(activeTab as TabName),
      onRefresh: () => refresh(),
      onEdit: (id: string, name: string, type: string) => setEditingResource({ id, name, type }),
      onDelete: (id: string) => setDeletingResource({ id, type: activeTab as string }),
    }

    switch (activeTab) {
      case 'vpcs':
        return <ResourceTab title="Virtual Private Clouds" description="Isolated network environments for your resources." items={vpcs} type="vpcs" {...handlers} />
      case 'subnets':
        return <ResourceTab title="Subnets" description="Network segments within your VPCs." items={subnets} type="subnets" {...handlers} />
      case 'instances':
        return <ResourceTab title="Virtual Server Instances" description="Compute resources running in your subnets." items={instances} type="instances" {...handlers} />
      case 'loadbalancers':
        return <ResourceTab title="Load Balancers" description="Distribute traffic across your instances." items={loadbalancers} type="loadbalancers" {...handlers} />
      case 'buckets':
        return <ResourceTab title="Object Storage Buckets" description="Scalable object storage for your data." items={buckets} type="buckets" {...handlers} />
      case 'databases':
        return <ResourceTab title="Managed Databases" description="Fully managed relational database instances." items={databases} type="databases" {...handlers} />
      case 'clusters':
        return <ResourceTab title="Kubernetes Clusters" description="Managed Kubernetes clusters for container workloads." items={clusters} type="clusters" {...handlers} />
      case 'hierarchy':
        return <HierarchyView vpcs={vpcs} subnets={subnets} databases={databases} clusters={clusters} instances={instances} buckets={buckets} loadbalancers={loadbalancers} />
      default:
        return <div>Coming soon</div>
    }
  }

  const getCreateModal = () => {
    if (!createModalTab) return null
    switch (createModalTab) {
      case 'vpcs':
        return <CreateVPCForm onClose={() => setCreateModalTab(null)} onSuccess={() => { refresh(); setCreateModalTab(null) }} />
      case 'subnets':
        return <CreateSubnetForm onClose={() => setCreateModalTab(null)} onSuccess={() => { refresh(); setCreateModalTab(null) }} />
      case 'instances':
        return <CreateInstanceForm onClose={() => setCreateModalTab(null)} onSuccess={() => { refresh(); setCreateModalTab(null) }} />
      case 'loadbalancers':
        return <CreateLoadBalancerForm onClose={() => setCreateModalTab(null)} onSuccess={() => { refresh(); setCreateModalTab(null) }} />
      case 'buckets':
        return <CreateBucketForm onClose={() => setCreateModalTab(null)} onSuccess={() => { refresh(); setCreateModalTab(null) }} />
      case 'databases':
        return <CreateDatabaseForm onClose={() => setCreateModalTab(null)} onSuccess={() => { refresh(); setCreateModalTab(null) }} />
      case 'clusters':
        return <CreateClusterForm onClose={() => setCreateModalTab(null)} onSuccess={() => { refresh(); setCreateModalTab(null) }} />
      default:
        return null
    }
  }

  const getApiEndpoint = (type: string) => {
    const map: { [key: string]: string } = {
      vpcs: '/v1/vpcs',
      subnets: '/v1/subnets',
      instances: '/v1/instances',
      loadbalancers: '/v1/loadbalancers',
      buckets: '/v1/buckets',
      databases: '/v1/databases',
      clusters: '/v1/clusters',
    }
    return map[type]
  }

  const navItems = [
    { id: 'vpcs', label: 'VPCs', icon: '🌐', count: vpcs.length },
    { id: 'subnets', label: 'Subnets', icon: '📊', count: subnets.length },
    { id: 'instances', label: 'Instances', icon: '💻', count: instances.length },
    { id: 'loadbalancers', label: 'Load Balancers', icon: '⚖️', count: loadbalancers.length },
    { id: 'buckets', label: 'Buckets', icon: '🪣', count: buckets.length },
    { id: 'databases', label: 'Databases', icon: '🗄️', count: databases.length },
    { id: 'clusters', label: 'Clusters', icon: '☸️', count: clusters.length },
    { id: 'hierarchy', label: 'Hierarchy', icon: '🗺️', count: null },
  ] as const

  return (
    <div style={{ display: 'flex', height: '100vh', flexDirection: 'column' }}>
      {/* Sidebar */}
      <div
        style={{
          display: 'flex',
          flex: '1',
          overflow: 'hidden',
          backgroundColor: 'var(--bg)',
        }}
      >
        <div
          style={{
            width: '200px',
            backgroundColor: 'var(--surface)',
            borderRight: '1px solid var(--border)',
            overflowY: 'auto',
            display: 'flex',
            flexDirection: 'column',
          }}
        >
          <div style={{ padding: '16px', borderBottom: '1px solid var(--border)' }}>
            <div style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text-2)', textTransform: 'uppercase', letterSpacing: '0.5px' }}>
              Resources
            </div>
          </div>
          <div style={{ flex: 1, display: 'flex', flexDirection: 'column', padding: '8px' }}>
            {navItems.map((item) => (
              <button
                key={item.id}
                onClick={() => setActiveTab(item.id as TabName)}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '12px',
                  padding: '12px 12px',
                  marginBottom: '4px',
                  border: 'none',
                  borderRadius: 'var(--r-sm)',
                  backgroundColor: activeTab === item.id ? 'var(--brand-subtle)' : 'transparent',
                  color: activeTab === item.id ? 'var(--brand)' : 'var(--text-2)',
                  cursor: 'pointer',
                  fontSize: '13px',
                  fontWeight: activeTab === item.id ? 600 : 500,
                  transition: 'all var(--t)',
                  textAlign: 'left',
                }}
                onMouseEnter={(e) => {
                  if (activeTab !== item.id) {
                    (e.currentTarget as HTMLButtonElement).style.backgroundColor = 'var(--surface-2)'
                  }
                }}
                onMouseLeave={(e) => {
                  if (activeTab !== item.id) {
                    (e.currentTarget as HTMLButtonElement).style.backgroundColor = 'transparent'
                  }
                }}
              >
                <span style={{ fontSize: '18px' }}>{item.icon}</span>
                <div style={{ flex: 1, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <span>{item.label}</span>
                  {item.count !== null && (
                    <span
                      style={{
                        fontSize: '11px',
                        padding: '2px 6px',
                        backgroundColor: activeTab === item.id ? 'var(--brand)' : 'var(--border)',
                        color: activeTab === item.id ? 'white' : 'var(--text-2)',
                        borderRadius: 'var(--r-xs)',
                        fontWeight: 600,
                      }}
                    >
                      {item.count}
                    </span>
                  )}
                </div>
              </button>
            ))}
          </div>
        </div>

        {/* Content */}
        <main style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
          {loading && <div style={{ padding: '20px', textAlign: 'center' }}>Loading resources...</div>}
          {!loading && (
            <div style={{ flex: 1, overflow: 'auto' }}>
              {getTabContent()}
            </div>
          )}
          {getCreateModal()}
        </main>
      </div>

      {editingResource && (
        <EditResourceModal
          resource={editingResource}
          endpoint={getApiEndpoint(editingResource.type)}
          onClose={() => setEditingResource(null)}
          onSuccess={() => {
            refresh()
            setEditingResource(null)
          }}
        />
      )}
      {deletingResource && (
        <div className="overlay" onClick={() => setDeletingResource(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Delete Resource</h3>
              <button className="btn-icon" onClick={() => setDeletingResource(null)} aria-label="Close">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </button>
            </div>
            <div className="modal-body">
              <p>Are you sure you want to delete this resource? This action cannot be undone.</p>
            </div>
            <div className="modal-footer">
              <button className="btn btn-ghost" onClick={() => setDeletingResource(null)}>
                Cancel
              </button>
              <button
                className="btn btn-danger"
                onClick={async () => {
                  try {
                    const endpoint = getApiEndpoint(deletingResource.type)
                    await fetch_api(`${endpoint}/${deletingResource.id}`, {
                      method: 'DELETE',
                    })
                    refresh()
                    setDeletingResource(null)
                  } catch (err) {
                    alert('Failed to delete resource: ' + (err instanceof Error ? err.message : 'Unknown error'))
                  }
                }}
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

function ResourceTab({
  title,
  description,
  items,
  type,
  onCreateClick,
  onRefresh,
  onEdit,
  onDelete,
}: {
  title: string
  description: string
  items: any[]
  type: string
  onCreateClick: () => void
  onRefresh: () => void
  onEdit: (id: string, name: string, type: string) => void
  onDelete: (id: string) => void
}) {
  const resourceType = title.split(' ')[0]

  return (
    <div id={`tab-${title.toLowerCase().replace(/\s+/g, '')}`} className="tab-panel">
      <div className="panel-header">
        <div>
          <h2>{title}</h2>
          <p className="panel-desc">{description}</p>
        </div>
        <div className="panel-actions">
          <button className="btn btn-ghost" title="Refresh resources" onClick={onRefresh}>
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"><path d="M23 4v6h-6" /><path d="M1 20v-6h6" /><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15" /></svg>
            Refresh
          </button>
          <button className="btn btn-primary" title={`Create new ${resourceType}`} onClick={onCreateClick}>
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><line x1="12" y1="5" x2="12" y2="19" /><line x1="5" y1="12" x2="19" y2="12" /></svg>
            Create {resourceType}
          </button>
        </div>
      </div>
      <div className="resource-card">
        {items.length === 0 ? (
          <div className="empty">
            <div className="empty-icon">📭</div>
            <strong>No resources</strong>
            <p>Create one to get started</p>
          </div>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Status</th>
                <th>ID</th>
                <th>Created</th>
                <th style={{ width: '100px', textAlign: 'center' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {items.map((item: any) => (
                <tr key={item.id}>
                  <td><div className="rname">{item.name}</div></td>
                  <td><span className={`badge badge-${item.status}`}>{item.status}</span></td>
                  <td><code className="rid">{item.id.substring(0, 8)}...</code></td>
                  <td className="crn">{new Date(item.created_at).toLocaleDateString()}</td>
                  <td style={{ textAlign: 'center' }}>
                    <button
                      className="btn btn-sm"
                      style={{ padding: '4px 8px', marginRight: '4px', fontSize: '12px' }}
                      onClick={() => onEdit(item.id, item.name, type)}
                      title="Edit"
                    >
                      ✎
                    </button>
                    <button
                      className="btn btn-sm"
                      style={{ padding: '4px 8px', fontSize: '12px', color: 'var(--error)' }}
                      onClick={() => onDelete(item.id)}
                      title="Delete"
                    >
                      ✕
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
