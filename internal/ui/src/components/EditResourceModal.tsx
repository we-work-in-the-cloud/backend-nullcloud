import { useState } from 'react'
import Modal from './Modal'
import { TextInput } from './FormField'
import { useAuth } from '../context/AuthContext'
import { useAPI } from '../hooks/useAPI'

interface EditResourceModalProps {
  resource: { id: string; name: string; type: string }
  endpoint: string
  onClose: () => void
  onSuccess: () => void
}

export default function EditResourceModal({ resource, endpoint, onClose, onSuccess }: EditResourceModalProps) {
  const [name, setName] = useState(resource.name)
  const { token } = useAuth()
  const { fetch_api } = useAPI(token)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await fetch_api(`${endpoint}/${resource.id}`, {
        method: 'PATCH',
        body: JSON.stringify({ name }),
      })
      onSuccess()
    } catch (err) {
      alert('Failed to update resource: ' + (err instanceof Error ? err.message : 'Unknown error'))
    }
  }

  return (
    <Modal title="Edit Resource" onClose={onClose} onSubmit={handleSubmit} submitText="Update">
      <TextInput label="Name" id="resource-name" value={name} onChange={setName} required />
    </Modal>
  )
}
