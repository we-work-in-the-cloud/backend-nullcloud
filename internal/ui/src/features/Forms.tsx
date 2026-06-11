import { useState } from 'react'
import Modal from '../components/Modal'
import { TextInput, SelectField, FieldRow } from '../components/FormField'
import { useResources } from '../context/ResourceContext'
import { useAPI } from '../hooks/useAPI'
import { useAuth } from '../context/AuthContext'

interface CreateFormProps {
  onClose: () => void
  onSuccess: () => void
}

export function CreateVPCForm({ onClose, onSuccess }: CreateFormProps) {
  const [name, setName] = useState('')
  const [region, setRegion] = useState('')
  const { token } = useAuth()
  const { fetch_api } = useAPI(token)

  const regions = ['us-east', 'us-west', 'us-central', 'eu-central', 'eu-east', 'eu-west']

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await fetch_api('/v1/vpcs', {
        method: 'POST',
        body: JSON.stringify({ name, region }),
      })
      onSuccess()
      onClose()
    } catch (err) {
      alert('Failed to create VPC: ' + (err instanceof Error ? err.message : 'Unknown error'))
    }
  }

  return (
    <Modal title="Create VPC" onClose={onClose} onSubmit={handleSubmit}>
      <TextInput label="VPC Name" id="vpc-name" value={name} onChange={setName} placeholder="my-vpc" required />
      <SelectField label="Region" id="region" value={region} onChange={setRegion} options={regions.map((r) => ({ value: r, label: r }))} required />
    </Modal>
  )
}

export function CreateSubnetForm({ onClose, onSuccess }: CreateFormProps) {
  const [name, setName] = useState('')
  const [vpcId, setVpcId] = useState('')
  const [zone, setZone] = useState('')
  const [cidrBlock, setCidrBlock] = useState('')
  const { vpcs } = useResources()
  const { token } = useAuth()
  const { fetch_api } = useAPI(token)

  const zones = ['us-east-1', 'us-east-2', 'us-east-3', 'us-west-1', 'us-west-2', 'us-west-3']

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await fetch_api('/v1/subnets', {
        method: 'POST',
        body: JSON.stringify({ name, vpc: { id: vpcId }, zone, cidr_block: cidrBlock }),
      })
      onSuccess()
      onClose()
    } catch (err) {
      alert('Failed to create Subnet: ' + (err instanceof Error ? err.message : 'Unknown error'))
    }
  }

  return (
    <Modal title="Create Subnet" onClose={onClose} onSubmit={handleSubmit}>
      <TextInput label="Subnet Name" id="subnet-name" value={name} onChange={setName} placeholder="my-subnet" required />
      <SelectField label="VPC" id="vpc" value={vpcId} onChange={setVpcId} options={vpcs.map((v) => ({ value: v.id, label: v.name }))} required />
      <FieldRow>
        <SelectField label="Zone" id="zone" value={zone} onChange={setZone} options={zones.map((z) => ({ value: z, label: z }))} required />
        <TextInput label="CIDR Block" id="cidr" value={cidrBlock} onChange={setCidrBlock} placeholder="10.0.1.0/24" required />
      </FieldRow>
    </Modal>
  )
}

export function CreateInstanceForm({ onClose, onSuccess }: CreateFormProps) {
  const [name, setName] = useState('')
  const [subnetId, setSubnetId] = useState('')
  const [profile, setProfile] = useState('')
  const [image, setImage] = useState('')
  const { subnets } = useResources()
  const { token } = useAuth()
  const { fetch_api } = useAPI(token)

  const profiles = ['bx2.2x8', 'bx2.4x16', 'bx2.8x32', 'cx2.2x4', 'cx2.4x8']
  const images = ['ubuntu-22.04', 'centos-8', 'debian-11', 'windows-2022']

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await fetch_api('/v1/instances', {
        method: 'POST',
        body: JSON.stringify({ name, subnet: { id: subnetId }, profile: { name: profile }, image: { id: image } }),
      })
      onSuccess()
      onClose()
    } catch (err) {
      alert('Failed to create Instance: ' + (err instanceof Error ? err.message : 'Unknown error'))
    }
  }

  return (
    <Modal title="Create Instance" onClose={onClose} onSubmit={handleSubmit}>
      <TextInput label="Instance Name" id="instance-name" value={name} onChange={setName} placeholder="my-instance" required />
      <SelectField label="Subnet" id="subnet" value={subnetId} onChange={setSubnetId} options={subnets.map((s) => ({ value: s.id, label: s.name }))} required />
      <FieldRow>
        <SelectField label="Profile" id="profile" value={profile} onChange={setProfile} options={profiles.map((p) => ({ value: p, label: p }))} required />
        <SelectField label="Image" id="image" value={image} onChange={setImage} options={images.map((i) => ({ value: i, label: i }))} required />
      </FieldRow>
    </Modal>
  )
}

export function CreateLoadBalancerForm({ onClose, onSuccess }: CreateFormProps) {
  const [name, setName] = useState('')
  const [protocol, setProtocol] = useState('http')
  const [port, setPort] = useState('80')
  const [selectedTargets, setSelectedTargets] = useState<string[]>([])
  const { instances, clusters } = useResources()
  const { token } = useAuth()
  const { fetch_api } = useAPI(token)

  const protocols = ['http', 'https', 'tcp', 'udp']

  const handleTargetToggle = (id: string, type: string) => {
    const targetId = `${type}:${id}`
    setSelectedTargets((prev) =>
      prev.includes(targetId) ? prev.filter((t) => t !== targetId) : [...prev, targetId]
    )
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const targets = selectedTargets.map((t) => {
        const [type, id] = t.split(':')
        return { type: type === 'instance' ? 'vsi' : type, id }
      })
      await fetch_api('/v1/loadbalancers', {
        method: 'POST',
        body: JSON.stringify({ name, protocol, port: parseInt(port), targets }),
      })
      onSuccess()
      onClose()
    } catch (err) {
      alert('Failed to create Load Balancer: ' + (err instanceof Error ? err.message : 'Unknown error'))
    }
  }

  return (
    <Modal title="Create Load Balancer" onClose={onClose} onSubmit={handleSubmit}>
      <TextInput label="Load Balancer Name" id="lb-name" value={name} onChange={setName} placeholder="my-lb" required />
      <FieldRow>
        <SelectField label="Protocol" id="protocol" value={protocol} onChange={setProtocol} options={protocols.map((p) => ({ value: p, label: p }))} required />
        <TextInput label="Port" id="port" value={port} onChange={setPort} type="number" required />
      </FieldRow>
      <div className="field">
        <label>Targets</label>
        <div style={{ maxHeight: '200px', overflowY: 'auto', border: '1px solid var(--border)', borderRadius: 'var(--r-sm)', padding: '8px' }}>
          {instances.length > 0 && (
            <div style={{ marginBottom: '8px' }}>
              <strong style={{ fontSize: '12px', color: 'var(--text-2)' }}>Instances</strong>
              {instances.map((i) => (
                <div key={i.id} style={{ padding: '4px 0' }}>
                  <label style={{ display: 'flex', alignItems: 'center', gap: '6px', fontSize: '13px', cursor: 'pointer' }}>
                    <input
                      type="checkbox"
                      checked={selectedTargets.includes(`instance:${i.id}`)}
                      onChange={() => handleTargetToggle(i.id, 'instance')}
                    />
                    {i.name}
                  </label>
                </div>
              ))}
            </div>
          )}
          {clusters.length > 0 && (
            <div>
              <strong style={{ fontSize: '12px', color: 'var(--text-2)' }}>Clusters</strong>
              {clusters.map((c) => (
                <div key={c.id} style={{ padding: '4px 0' }}>
                  <label style={{ display: 'flex', alignItems: 'center', gap: '6px', fontSize: '13px', cursor: 'pointer' }}>
                    <input
                      type="checkbox"
                      checked={selectedTargets.includes(`cluster:${c.id}`)}
                      onChange={() => handleTargetToggle(c.id, 'cluster')}
                    />
                    {c.name}
                  </label>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </Modal>
  )
}

export function CreateBucketForm({ onClose, onSuccess }: CreateFormProps) {
  const [name, setName] = useState('')
  const [region, setRegion] = useState('')
  const { token } = useAuth()
  const { fetch_api } = useAPI(token)

  const regions = ['us-east', 'us-west', 'us-central', 'eu-central', 'eu-east', 'eu-west']

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await fetch_api('/v1/buckets', {
        method: 'POST',
        body: JSON.stringify({ name, region }),
      })
      onSuccess()
      onClose()
    } catch (err) {
      alert('Failed to create Bucket: ' + (err instanceof Error ? err.message : 'Unknown error'))
    }
  }

  return (
    <Modal title="Create Bucket" onClose={onClose} onSubmit={handleSubmit}>
      <TextInput label="Bucket Name" id="bucket-name" value={name} onChange={setName} placeholder="my-bucket" required />
      <SelectField label="Region" id="region" value={region} onChange={setRegion} options={regions.map((r) => ({ value: r, label: r }))} required />
    </Modal>
  )
}

export function CreateDatabaseForm({ onClose, onSuccess }: CreateFormProps) {
  const [name, setName] = useState('')
  const [engine, setEngine] = useState('postgres')
  const [version, setVersion] = useState('13')
  const [plan, setPlan] = useState('small')
  const [selectedSubnets, setSelectedSubnets] = useState<string[]>([])
  const { subnets } = useResources()
  const { token } = useAuth()
  const { fetch_api } = useAPI(token)

  const engines = [
    { value: 'postgres', label: 'PostgreSQL' },
    { value: 'mysql', label: 'MySQL' },
    { value: 'mariadb', label: 'MariaDB' },
  ]
  const versions = ['9.6', '10', '11', '12', '13', '14', '15']
  const plans = ['small', 'medium', 'large']

  const handleSubnetToggle = (id: string) => {
    setSelectedSubnets((prev) => (prev.includes(id) ? prev.filter((s) => s !== id) : [...prev, id]))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await fetch_api('/v1/databases', {
        method: 'POST',
        body: JSON.stringify({ name, engine, version, plan, subnet_ids: selectedSubnets }),
      })
      onSuccess()
      onClose()
    } catch (err) {
      alert('Failed to create Database: ' + (err instanceof Error ? err.message : 'Unknown error'))
    }
  }

  return (
    <Modal title="Create Database" onClose={onClose} onSubmit={handleSubmit}>
      <TextInput label="Database Name" id="db-name" value={name} onChange={setName} placeholder="my-db" required />
      <FieldRow>
        <SelectField label="Engine" id="engine" value={engine} onChange={setEngine} options={engines} required />
        <SelectField label="Version" id="version" value={version} onChange={setVersion} options={versions.map((v) => ({ value: v, label: v }))} required />
      </FieldRow>
      <SelectField label="Plan" id="plan" value={plan} onChange={setPlan} options={plans.map((p) => ({ value: p, label: p }))} required />
      <div className="field">
        <label>Subnets</label>
        <div style={{ maxHeight: '150px', overflowY: 'auto', border: '1px solid var(--border)', borderRadius: 'var(--r-sm)', padding: '8px' }}>
          {subnets.map((s) => (
            <div key={s.id} style={{ padding: '4px 0' }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: '6px', fontSize: '13px', cursor: 'pointer' }}>
                <input
                  type="checkbox"
                  checked={selectedSubnets.includes(s.id)}
                  onChange={() => handleSubnetToggle(s.id)}
                />
                {s.name}
              </label>
            </div>
          ))}
        </div>
      </div>
    </Modal>
  )
}

export function CreateClusterForm({ onClose, onSuccess }: CreateFormProps) {
  const [name, setName] = useState('')
  const [version, setVersion] = useState('1.24')
  const [nodeCount, setNodeCount] = useState('3')
  const [selectedSubnets, setSelectedSubnets] = useState<string[]>([])
  const { subnets } = useResources()
  const { token } = useAuth()
  const { fetch_api } = useAPI(token)

  const versions = ['1.20', '1.21', '1.22', '1.23', '1.24', '1.25']

  const handleSubnetToggle = (id: string) => {
    setSelectedSubnets((prev) => (prev.includes(id) ? prev.filter((s) => s !== id) : [...prev, id]))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await fetch_api('/v1/clusters', {
        method: 'POST',
        body: JSON.stringify({ name, version, node_count: parseInt(nodeCount), subnet_ids: selectedSubnets }),
      })
      onSuccess()
      onClose()
    } catch (err) {
      alert('Failed to create Cluster: ' + (err instanceof Error ? err.message : 'Unknown error'))
    }
  }

  return (
    <Modal title="Create Kubernetes Cluster" onClose={onClose} onSubmit={handleSubmit}>
      <TextInput label="Cluster Name" id="cluster-name" value={name} onChange={setName} placeholder="my-cluster" required />
      <FieldRow>
        <SelectField label="Kubernetes Version" id="version" value={version} onChange={setVersion} options={versions.map((v) => ({ value: v, label: v }))} required />
        <TextInput label="Node Count" id="nodes" value={nodeCount} onChange={setNodeCount} type="number" required />
      </FieldRow>
      <div className="field">
        <label>Subnets</label>
        <div style={{ maxHeight: '150px', overflowY: 'auto', border: '1px solid var(--border)', borderRadius: 'var(--r-sm)', padding: '8px' }}>
          {subnets.map((s) => (
            <div key={s.id} style={{ padding: '4px 0' }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: '6px', fontSize: '13px', cursor: 'pointer' }}>
                <input
                  type="checkbox"
                  checked={selectedSubnets.includes(s.id)}
                  onChange={() => handleSubnetToggle(s.id)}
                />
                {s.name}
              </label>
            </div>
          ))}
        </div>
      </div>
    </Modal>
  )
}
