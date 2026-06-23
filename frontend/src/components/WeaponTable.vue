<script setup lang="ts">
import { ref, nextTick, onMounted, onUnmounted, inject } from 'vue'
import { useWeaponStore } from '../stores/weapons'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import SearchBar from './SearchBar.vue'

const store = useWeaponStore()
const authStore = useAuthStore()
const router = useRouter()

const openAuth = inject<() => void>('openAuth', () => {})

// Ссылки на элементы для замера высоты
const tableWrapper = ref<HTMLElement | null>(null)
const innerTable = ref<HTMLElement | null>(null)

// Функция перехода на детали
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

// Автоматически отслеживаем изменение размера таблицы
let tableObserver: ResizeObserver | null = null

onMounted(async () => {
  await nextTick()
  if (!innerTable.value) return
  const updateHeight = () => {
    if (tableWrapper.value && innerTable.value) {
      tableWrapper.value.style.height = `${innerTable.value.offsetHeight}px`
    }
  }
  updateHeight()
  tableObserver = new ResizeObserver(updateHeight)
  tableObserver.observe(innerTable.value)
})

onUnmounted(() => {
  tableObserver?.disconnect()
})
</script>

<template>
  <SearchBar />
  
  <div ref="tableWrapper" class="table-blind-container">
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
            <button class="view-btn" @click="goToDetails(item.id)">View <span class="details-text">Details</span></button>
            <button v-if="authStore.isAdmin" class="delete-btn" @click="handleDelete(item.id)">×</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <div v-if="store.searchTerm && !store.filteredWeapons.length" class="no-results">
    <template v-if="authStore.isAuthenticated">
      <span>No articles found for "{{ store.searchTerm }}". </span>
      <button class="create-btn" @click="goToCreate">Create</button>
    </template>
    <span v-else>
      To create a new article, please <a class="auth-link" @click="openAuth()">log in</a>
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

@media (max-width: 768px) {
  .details-text {
    display: none;
  }
}
</style>
