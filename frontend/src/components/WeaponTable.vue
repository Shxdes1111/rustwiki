<script setup lang="ts">
import { useWeaponStore } from '../stores/weapons'
import { useToast } from 'vue-toastification'
import { useRouter } from 'vue-router'
import SearchBar from './SearchBar.vue'

const store = useWeaponStore()
const toast = useToast()
const router = useRouter()

function showNotImplemented() {
  toast.info('Страница не заполнена')
}
function goToDetails(id: number) {
  router.push({ name: 'weapon-details', params: { id } })
}
</script>

<template>
  <SearchBar />
  <table class="wiki-table">
    <thead>
      <tr>
        <th>ID</th>
        <th>Name</th>
        <th>Type</th>
        <th>Action</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="item in store.filteredWeapons" :key="item.id">
        <td>{{ item.id }}</td>
        <td class="weapon-name">{{ item.name }}</td>
        <td><span class="badge">{{ item.type }}</span></td>
        <td>
          <!-- Меняем @click -->
          <button class="view-btn" @click="goToDetails(item.id)">View Details</button>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<style scoped>
.wiki-table {
  width: 100%;
  border-collapse: collapse;
  background-color: #242424;
  border-radius: 8px;
  overflow: hidden;
}

.wiki-table th {
  background-color: #42b883;
  color: #1a1a1a;
  text-align: left;
  padding: 12px;
}

.wiki-table td {
  padding: 12px;
  border-bottom: 1px solid #333;
}

.wiki-table tr:hover {
  background-color: #2f2f2f;
}

.weapon-name {
  font-weight: bold;
}

.badge {
  background: #444;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.8rem;
}

.view-btn {
  background: #42b883;
  border: none;
  color: white;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: opacity 0.2s;
}

.view-btn:hover {
  opacity: 0.8;
}
</style>
