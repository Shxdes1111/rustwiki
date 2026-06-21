<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onUnmounted, inject } from 'vue'
import { useWeaponStore } from '../stores/weapons'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import SearchBar from './SearchBar.vue'

const store = useWeaponStore()
const authStore = useAuthStore()
const router = useRouter()

const openAuth = inject<() => void>('openAuth', () => {})

const isMobile = ref(window.innerWidth < 768)

function onResize() {
  isMobile.value = window.innerWidth < 768
}

onMounted(() => {
  window.addEventListener('resize', onResize)
  if (!isMobile.value) initTableHeight()
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
})

// Ссылки на элементы для замера высоты
const tableWrapper = ref<HTMLElement | null>(null)
const innerTable = ref<HTMLElement | null>(null)

async function initTableHeight() {
  await nextTick()
  if (tableWrapper.value && innerTable.value) {
    tableWrapper.value.style.height = `${innerTable.value.offsetHeight}px`
  }
}

const goToDetails = (id: number) => {
  router.push(`/weapon/${id}`)
}

const goToCreate = () => {
  router.push('/weapon/create')
}

const handleDelete = async (id: number) => {
  if (!confirm('Delete this weapon?')) return
  try {
    await store.deleteWeapon(id)
    await store.fetchWeapons()
  } catch {
    alert('Failed to delete weapon')
  }
}

watch(isMobile, (mobile) => {
  if (mobile) {
    if (tableWrapper.value) {
      tableWrapper.value.style.height = ''
    }
  } else {
    initTableHeight()
  }
})

watch(() => store.filteredWeapons, async () => {
  if (isMobile.value) return
  await nextTick()
  
  if (tableWrapper.value && innerTable.value) {
    const newHeight = innerTable.value.offsetHeight
    tableWrapper.value.style.height = `${newHeight}px`
  }
}, { deep: true })
</script>

<template>
  <SearchBar />
  
  <div ref="tableWrapper" v-show="!isMobile" class="table-blind-container">
    <table ref="innerTable" class="wiki-table">
      <thead>
        <tr>
          <th>№</th>
          <th>Name</th>
          <th>Type</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index) in store.filteredWeapons" :key="item.id" class="table-row">
          <td>{{ index + 1 }}</td>
          <td class="weapon-name">{{ item.name }}</td>
          <td><span class="badge">{{ item.type }}</span></td>
          <td class="actions-cell">
            <button class="view-btn" @click="goToDetails(item.id)">View Details</button>
            <button v-if="authStore.isAdmin" class="delete-btn" @click="handleDelete(item.id)">×</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <div v-show="isMobile" class="card-list">
    <div v-for="(item, index) in store.filteredWeapons" :key="item.id" class="weapon-card" @click="goToDetails(item.id)">
      <div class="card-header">
        <span class="card-num">#{{ index + 1 }}</span>
        <span class="badge">{{ item.type }}</span>
      </div>
      <div class="card-body">{{ item.name }}</div>
      <div class="card-actions" @click.stop>
        <button class="view-btn" @click="goToDetails(item.id)">View Details</button>
        <button v-if="authStore.isAdmin" class="delete-btn" @click="handleDelete(item.id)">×</button>
      </div>
    </div>
  </div>

  <div v-if="store.searchTerm && !store.filteredWeapons.length" class="no-results">
    <template v-if="authStore.isAuthenticated">
      <span>No weapon found for "{{ store.searchTerm }}". </span>
      <button class="create-btn" @click="goToCreate">Create</button>
    </template>
    <span v-else>
      To create a new weapon, please <a class="auth-link" @click="openAuth()">log in</a>
    </span>
  </div>
</template>

<style scoped>
.table-blind-container {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  transition: height 0.8s cubic-bezier(0.25, 1, 0.5, 1); 
}

.wiki-table {
  width: 100%;
  border-collapse: collapse;
  background-color: #1a1a1a;
  color: white;
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #333;
  vertical-align: middle;
}

.table-row {
  transition: opacity 1s ease;
}

.table-row:hover {
  background-color: #262626;
}

@media (max-width: 768px) {
  .table-blind-container {
    transition: none !important;
  }
}

.actions-cell {
  white-space: nowrap;
}

.badge {
  background-color: #333;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.85rem;
}

.view-btn {
  background-color: #ce422b;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.view-btn:hover {
  background-color: #a8321f;
}

.delete-btn {
  background: none;
  border: none;
  color: #ef4444;
  font-size: 1.3rem;
  cursor: pointer;
  padding: 4px 8px;
  margin-left: 6px;
  transition: color 0.2s;
  line-height: 1;
}

.delete-btn:hover {
  color: #dc2626;
}

.card-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.weapon-card {
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 14px;
  cursor: pointer;
  transition: background 0.2s;
}

.weapon-card:hover {
  background: #262626;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.card-num {
  color: #64748b;
  font-size: 0.85rem;
  font-weight: 600;
}

.card-body {
  font-size: 1.1rem;
  font-weight: 600;
  color: #f8fafc;
  margin-bottom: 12px;
}

.card-actions {
  display: flex;
  gap: 8px;
}

.card-actions .view-btn {
  flex: 1;
  text-align: center;
  padding: 8px;
}

.card-actions .delete-btn {
  margin-left: 0;
  padding: 8px 14px;
  background: #333;
  border-radius: 4px;
  font-size: 1.2rem;
}

.card-actions .delete-btn:hover {
  background: #444;
  color: #ef4444;
}

.no-results {
  text-align: center;
  padding: 40px 20px;
  color: #94a3b8;
  font-size: 1.1rem;
}

.create-btn {
  background-color: #ce422b;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  margin-left: 8px;
  transition: background-color 0.2s;
}

.create-btn:hover {
  background-color: #a8321f;
}

.auth-link {
  color: #ce422b;
  cursor: pointer;
  text-decoration: underline;
}

.auth-link:hover {
  color: #ef4444;
}
</style>
