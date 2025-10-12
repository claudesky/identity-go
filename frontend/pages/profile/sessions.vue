<script setup lang="ts">
import { useSelf, type Session } from '~/composables/useSelf'
import { useVueTable, getCoreRowModel, createColumnHelper, FlexRender } from '@tanstack/vue-table'

definePageMeta({
  middleware: 'auth'
})

useHead({
  title: 'Sessions'
})

const { getSessions } = useSelf()

const sessions = ref<Session[]>([])
const loading = ref(true)
const error = ref('')

const loadSessions = async () => {
  loading.value = true
  error.value = ''

  try {
    const data = (await getSessions()).data
    sessions.value = data.sort((a, b) =>
      new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    )
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to load sessions'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadSessions()
})

const columnHelper = createColumnHelper<Session>()

const columns = [
  columnHelper.accessor('created_at', {
    header: 'Created At',
    cell: info => new Date(info.getValue()).toLocaleString()
  }),
  columnHelper.accessor('last_issued_at', {
    header: 'Last Active',
    cell: info => new Date(info.getValue()).toLocaleString()
  }),
  columnHelper.accessor('expires_at', {
    header: 'Expires At',
    cell: info => new Date(info.getValue()).toLocaleString()
  }),
  columnHelper.accessor('revoked', {
    header: 'Status',
    cell: info => {
      if (info.getValue()) return 'Revoked'
      const expiresAt = new Date(info.row.original.expires_at)
      return expiresAt < new Date() ? 'Expired' : 'Active'
    }
  })
]

const table = useVueTable({
  get data() {
    return sessions.value
  },
  columns,
  getCoreRowModel: getCoreRowModel()
})
</script>

<template>
  <div>
    <AppNav />
    <main class="container">
      <h1>Active Sessions</h1>

      <div v-if="loading" class="loading">Loading sessions...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else-if="sessions.length === 0" class="empty">No active sessions found</div>

      <table v-else class="sessions-table">
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
  </div>
</template>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

h1 {
  margin-bottom: 2rem;
  color: #333;
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

.sessions-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.sessions-table thead {
  background: #f5f5f5;
}

.sessions-table th,
.sessions-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #e0e0e0;
}

.sessions-table th {
  font-weight: 600;
  color: #555;
}

.sessions-table tbody tr:hover {
  background: #fafafa;
}

.sessions-table tbody tr:last-child td {
  border-bottom: none;
}
</style>
