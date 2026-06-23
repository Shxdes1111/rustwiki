<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useWeaponStore, type Suggestion } from '../stores/weapons'
import { useToast } from 'vue-toastification'

const router = useRouter()
const store = useWeaponStore()
const toast = useToast()

const suggestions = ref<Suggestion[]>([])
const loading = ref(false)
const isMobile = ref(window.innerWidth < 768)

function checkWidth() {
  isMobile.value = window.innerWidth < 768
}

onMounted(async () => {
  window.addEventListener('resize', checkWidth)
  loading.value = true
  try {
    suggestions.value = await store.fetchMySuggestions()
  } catch (err) {
    toast.error(`Failed to load suggestions: ${err instanceof Error ? err.message : 'Unknown error'}`)
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', checkWidth)
})

const goToDetails = (id: number) => {
  router.push(`/my/suggestions/${id}`)
}

const goToEdit = (id: number) => {
  router.push(`/weapon/create?edit=${id}`)
}
</script>

<template>
  <div class="page">
    <button class="back-btn" @click="router.push('/')">← Back to list</button>
    <h1 class="page-title">My Suggestions</h1>

    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="!suggestions.length" class="empty">You haven't submitted any suggestions yet.</div>

    <template v-else>
      <table v-if="!isMobile" class="suggestion-table">
        <thead>
          <tr>
            <th>№</th>
            <th>Weapon</th>
            <th>Status</th>
            <th>Created</th>
            <th>Rejection reason</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(s, index) in suggestions" :key="s.id" class="suggestion-row" @click="goToDetails(s.id)">
            <td>{{ index + 1 }}</td>
            <td>{{ s.payload?.name || 'Unknown' }}</td>
            <td><span :class="['badge', `badge-${s.status}`]">{{ s.status }}</span></td>
            <td>{{ new Date(s.created_at).toLocaleDateString() }}</td>
            <td class="reason-cell">{{ s.status === 'rejected' ? (s.rejection_reason || '—') : '—' }}</td>
            <td class="actions-cell" @click.stop>
              <button v-if="s.status === 'rejected'" class="btn-edit" @click="goToEdit(s.id)">Edit</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-else class="card-list">
        <div v-for="(s, index) in suggestions" :key="s.id" class="suggestion-card" @click="goToDetails(s.id)">
          <div class="card-header">
            <span class="card-id">#{{ index + 1 }}</span>
            <span :class="['badge', `badge-${s.status}`]">{{ s.status }}</span>
          </div>
          <div class="card-body">
            <div class="card-row"><span class="card-label">Weapon</span><span class="card-value">{{ s.payload?.name || 'Unknown' }}</span></div>
            <div class="card-row"><span class="card-label">Created</span><span class="card-value">{{ new Date(s.created_at).toLocaleDateString() }}</span></div>
            <div v-if="s.status === 'rejected'" class="card-row">
              <span class="card-label">Reason</span>
              <span class="card-value reason-text">{{ s.rejection_reason || '—' }}</span>
            </div>
          </div>
          <div v-if="s.status === 'rejected'" class="card-actions" @click.stop>
            <button class="btn-edit" @click="goToEdit(s.id)">Edit</button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
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
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #333;
}

td:nth-child(2) {
  max-width: 0;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.reason-cell {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #fecaca;
  font-size: 0.85rem;
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

.btn-edit {
  background: #2563eb;
  color: white;
  border: none;
  padding: 6px 14px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: background 0.2s;
}

.btn-edit:hover { background: #1d4ed8; }

.card-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.suggestion-card {
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.suggestion-card:hover { background: #262626; }

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.card-id {
  font-weight: bold;
  color: #94a3b8;
  font-size: 0.9rem;
}

.card-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.card-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-label {
  color: #94a3b8;
  font-size: 0.85rem;
}

.card-value {
  color: #f8fafc;
  font-size: 0.9rem;
  text-align: right;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 60%;
}

.reason-text {
  color: #fecaca;
}

.card-actions {
  display: flex;
  gap: 8px;
  justify-content: stretch;
}

.card-actions .btn-edit {
  flex: 1;
  text-align: center;
  padding: 8px;
}
</style>
