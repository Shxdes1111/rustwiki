<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useWeaponStore } from '../stores/weapons'
import { useToast } from 'vue-toastification'

const router = useRouter()
const store = useWeaponStore()
const toast = useToast()

const weapons = ref<any[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    weapons.value = await store.fetchMyWeapons()
  } catch (err) {
    toast.error(`Failed to load: ${err instanceof Error ? err.message : 'Unknown error'}`)
  } finally {
    loading.value = false
  }
})

const goToDetails = (id: number) => {
  router.push(`/weapon/${id}`)
}
</script>

<template>
  <div class="page">
    <button class="back-btn" @click="router.push('/')">← Back to list</button>
    <h1 class="page-title">My Articles</h1>

    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="!weapons.length" class="empty">
      You haven't created any articles yet. Suggest a weapon and get it approved!
    </div>

    <table v-else class="weapon-table">
      <thead>
        <tr>
          <th>№</th>
          <th>Name</th>
          <th>Type</th>
          <th>Views</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index) in weapons" :key="item.id" class="weapon-row" @click="goToDetails(item.id)">
          <td>{{ index + 1 }}</td>
          <td class="weapon-name">{{ item.name }}</td>
          <td><span class="badge">{{ item.type }}</span></td>
          <td class="views-cell">{{ item.views ?? 0 }}</td>
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
}

.loading, .empty {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
  font-size: 1.1rem;
}

.weapon-table {
  width: 100%;
  border-collapse: collapse;
  background: #1a1a1a;
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #333;
}

.weapon-row {
  cursor: pointer;
  transition: background 0.2s;
}

.weapon-row:hover {
  background: #262626;
}

.weapon-name {
  font-weight: 600;
}

.badge {
  background: #854d0e;
  color: #fef08a;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.8rem;
  text-transform: capitalize;
}

.views-cell {
  color: #94a3b8;
}
</style>
