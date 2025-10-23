<script setup lang="ts">
import { useClients, type Client } from '~/composables/useClients'
import { useVueTable, getCoreRowModel, createColumnHelper, FlexRender } from '@tanstack/vue-table'

definePageMeta({
  middleware: 'auth'
})

useHead({
  title: 'OAuth Clients'
})

const { getClients, deleteClient } = useClients()

const showCreateModal = ref(false)
const showSecretModal = ref(false)
const newClientSecret = ref('')
const deletingClientId = ref<string | null>(null)

const { data: clients, error, pending: loading, refresh } = await useAsyncData(
  'oauth-clients',
  async () => {
    const result = await getClients()
    if (result.data.value === undefined) {
      throw result.error.value || new Error('Failed to load clients')
    }
    return result.data.value.data.sort((a, b) =>
      new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    )
  }
)

const columnHelper = createColumnHelper<Client>()

const columns = [
  columnHelper.accessor('name', {
    header: 'Name',
    cell: info => info.getValue()
  }),
  columnHelper.accessor('client_id', {
    header: 'Client ID',
    cell: info => info.getValue()
  }),
  columnHelper.accessor('redirect_uris', {
    header: 'Redirect URIs',
    cell: info => info.getValue().join(', ')
  }),
  columnHelper.accessor('created_at', {
    header: 'Created At',
    cell: info => new Date(info.getValue()).toLocaleString()
  }),
  columnHelper.display({
    id: 'actions',
    header: 'Actions',
    cell: props => h('div', { class: 'actions' }, [
      h(resolveComponent('IconButton'), {
        icon: 'edit',
        variant: 'text',
        color: 'primary',
        onClick: () => navigateTo(`/clients/${props.row.original.id}`)
      }),
      h(resolveComponent('IconButton'), {
        icon: 'trash',
        variant: 'text',
        color: 'danger',
        onClick: () => handleDelete(props.row.original.id),
        disabled: deletingClientId.value === props.row.original.id
      })
    ])
  })
]

const table = useVueTable({
  get data() {
    return clients.value || []
  },
  columns,
  getCoreRowModel: getCoreRowModel()
})

const handleDelete = async (id: string) => {
  if (!confirm('Are you sure you want to delete this client? This action cannot be undone.')) {
    return
  }

  deletingClientId.value = id
  try {
    const result = await deleteClient(id)
    if (result.error.value) {
      throw result.error.value
    }
    await refresh()
  } catch (err) {
    console.error('Failed to delete client:', err)
    alert('Failed to delete client')
  } finally {
    deletingClientId.value = null
  }
}

const handleClientCreated = async (clientSecret: string) => {
  showCreateModal.value = false
  newClientSecret.value = clientSecret
  showSecretModal.value = true
  await refresh()
}

const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text)
  alert('Copied to clipboard!')
}
</script>

<template>
  <div>
    <AppNav />
    <main class="container">
      <div class="header">
        <h1>OAuth Clients</h1>
        <button class="btn-primary" @click="showCreateModal = true">Create New Client</button>
      </div>

      <div v-if="loading" class="loading">Loading clients...</div>
      <div v-else-if="error" class="error">{{ error.message || error }}</div>
      <div v-else-if="!clients || clients.length === 0" class="empty">
        No OAuth clients found. Create one to get started.
      </div>

      <table v-else class="clients-table">
        <thead>
          <tr v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
            <th v-for="header in headerGroup.headers" :key="header.id">
              <FlexRender
                v-if="!header.isPlaceholder"
                :render="header.column.columnDef.header"
                :props="header.getContext()"
              />
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in table.getRowModel().rows" :key="row.id">
            <td v-for="cell in row.getVisibleCells()" :key="cell.id">
              <FlexRender
                :render="cell.column.columnDef.cell"
                :props="cell.getContext()"
              />
            </td>
          </tr>
        </tbody>
      </table>
    </main>

    <!-- Create Client Modal -->
    <ClientForm
      v-if="showCreateModal"
      @close="showCreateModal = false"
      @created="handleClientCreated"
    />

    <!-- Client Secret Modal -->
    <div v-if="showSecretModal" class="modal-overlay" @click="showSecretModal = false">
      <div class="modal" @click.stop>
        <h2>Client Created Successfully</h2>
        <p class="warning">
          Save this client secret now. You won't be able to see it again!
        </p>
        <div class="secret-container">
          <code>{{ newClientSecret }}</code>
          <button class="btn-copy" @click="copyToClipboard(newClientSecret)">Copy</button>
        </div>
        <button class="btn-primary" @click="showSecretModal = false">Done</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

h1 {
  color: #333;
  margin: 0;
}

.btn-primary {
  background: #007bff;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  transition: background 0.2s;
}

.btn-primary:hover {
  background: #0056b3;
}

.loading,
.error,
.empty {
  padding: 2rem;
  text-align: center;
  background: #f5f5f5;
  border-radius: 8px;
}

.error {
  background: #ffebee;
  color: #d32f2f;
}

.clients-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.clients-table thead {
  background: #f5f5f5;
}

.clients-table th,
.clients-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #e0e0e0;
}

.clients-table th {
  font-weight: 600;
  color: #555;
}

.clients-table tbody tr:hover {
  background: #fafafa;
}

.clients-table tbody tr:last-child td {
  border-bottom: none;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
}

.modal h2 {
  margin-top: 0;
  color: #333;
}

.warning {
  color: #dc3545;
  font-weight: 500;
  margin-bottom: 1rem;
}

.secret-container {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
  padding: 1rem;
  background: #f5f5f5;
  border-radius: 4px;
}

.secret-container code {
  flex: 1;
  word-break: break-all;
  font-family: monospace;
  font-size: 0.875rem;
}

.btn-copy {
  padding: 0.5rem 1rem;
  background: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  white-space: nowrap;
}

.btn-copy:hover {
  background: #5a6268;
}
</style>
