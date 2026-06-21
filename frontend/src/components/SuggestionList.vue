<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useWeaponStore } from '../stores/weapons'
import { useToast } from 'vue-toastification'

const router = useRouter()
const store = useWeaponStore()
const toast = useToast()

const suggestions = ref<any[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    suggestions.value = await store.fetchSuggestions()
  } catch (err) {
    toast.error(`Failed to load suggestions: ${err instanceof Error ? err.message : 'Unknown error'}`)
  } finally {
    loading.value = false
  }
})

const goToDetails = (id: number) => {
  router.push(`/admin/suggestions/${id}`)
}

const handleApprove = async (id: number) => {
  if (!confirm('Approve this suggestion? The weapon will be created.')) return
  try {
    await store.approveSuggestion(id)
    toast.success('Suggestion approved!')
    suggestions.value = await store.fetchSuggestions()
    await store.fetchWeapons()
  } catch (err) {
    toast.error(`Failed to approve: ${err instanceof Error ? err.message : 'Unknown error'}`)
  }
}

const handleReject = async (id: number) => {
  if (!confirm('Reject this suggestion?')) return
  try {
    await store.rejectSuggestion(id)
    toast.success('Suggestion rejected')
    suggestions.value = await store.fetchSuggestions()
  } catch (err) {
    toast.error(`Failed to reject: ${err instanceof Error ? err.message : 'Unknown error'}`)
  }
}
</script>

<template>
  <div class="page">
    <button class="back-btn" @click="router.push('/')">← Back to list</button>
    <h1 class="page-title">Suggestion Requests</h1>

    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="!suggestions.length" class="empty">No suggestions yet.</div>
    <table v-else class="suggestion-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Author</th>
          <th>Weapon</th>
          <th>Status</th>
          <th>Created</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in suggestions" :key="s.id" class="suggestion-row" @click="goToDetails(s.id)">
          <td>{{ s.id }}</td>
          <td>{{ s.username || `User #${s.user_id}` }}</td>
          <td>{{ s.payload?.name || 'Unknown' }}</td>
          <td><span :class="['badge', `badge-${s.status}`]">{{ s.status }}</span></td>
          <td>{{ new Date(s.created_at).toLocaleDateString() }}</td>
          <td class="actions-cell" @click.stop>
            <button v-if="s.status === 'pending'" class="btn-approve" @click="handleApprove(s.id)">Approve</button>
            <button v-if="s.status === 'pending'" class="btn-reject" @click="handleReject(s.id)">Reject</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
  font-family: sans-serif;
  color: #e2e8f0;
}

.back-btn {
  background: #333;
  color: #fff;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  margin-bottom: 20px;
}

.back-btn:hover { background: #444; }

.page-title {
  font-size: 2.5rem;
  border-bottom: 1px solid #5d5d5d;
  padding-bottom: 10px;
  margin-bottom: 20px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.loading, .empty {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
  font-size: 1.1rem;
}

.suggestion-table {
  width: 100%;
  border-collapse: collapse;
  background-color: #1a1a1a;
}

@media (max-width: 768px) {
  .page-title { white-space: normal; font-size: 1.8rem; }
  table { display: block; overflow-x: auto; }
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #333;
}

td:nth-child(2),
td:nth-child(3) {
  max-width: 0;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.suggestion-row {
  cursor: pointer;
  transition: background 0.2s;
}

.suggestion-row:hover {
  background-color: #262626;
}

.badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.85rem;
  text-transform: capitalize;
}

.badge-pending {
  background: #854d0e;
  color: #fef08a;
}

.badge-approved {
  background: #166534;
  color: #bbf7d0;
}

.badge-rejected {
  background: #991b1b;
  color: #fecaca;
}

.actions-cell {
  white-space: nowrap;
}

.btn-approve {
  background: #16a34a;
  color: white;
  border: none;
  padding: 6px 14px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  margin-right: 6px;
  transition: background 0.2s;
}

.btn-approve:hover { background: #15803d; }

.btn-reject {
  background: #dc2626;
  color: white;
  border: none;
  padding: 6px 14px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: background 0.2s;
}

.btn-reject:hover { background: #b91c1c; }
</style>