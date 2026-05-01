<script setup>
import { ref, onMounted, computed } from 'vue'

const weapons = ref([])
const searchTerm = ref('') 

onMounted(async () => {
  const res = await fetch('http://localhost:8080/api/weapons')
  weapons.value = await res.json()
})

// Logic: Filters the list based on the name matching the search box
const filteredWeapons = computed(() => {
  const s = searchTerm.value.toLowerCase()
  return weapons.value.filter(item =>{
    return item.name.toLowerCase().includes(s) ||
  item.type.toLowerCase().includes(s)})

})
</script>

<template>
  <div class="wiki-container">
    <h1>RustWiki Weapons</h1>
     <div class="search-section">
      <input 
        v-model="searchTerm" 
        placeholder="🔍 Search weapons..." 
        class="search-input"
      />
    </div>
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
        <tr v-for="item in filteredWeapons" :key="item.id">
          <td>{{ item.id }}</td>
          <td class="weapon-name">{{ item.name }}</td>
          <td><span class="badge">{{ item.type }}</span></td>
          <td><button class="view-btn">View Details</button></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
/* This makes the page dark and centered */
.wiki-container {
  padding: 40px;
  max-width: 1000px;
  margin: 0 auto;
  color: white;
}

/* This fixes the table width and spacing */
.wiki-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
  background-color: #242424;
  border-radius: 8px;
  overflow: hidden; /* Keeps the rounded corners */
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