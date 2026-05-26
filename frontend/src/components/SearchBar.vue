<script setup lang="ts">
import { useWeaponStore } from '../stores/weapons'

const store = useWeaponStore()
</script>

<template>
  <div class="search-section">
    <input
      v-model="store.searchTerm"
      placeholder="🔍 Search weapons..."
      class="search-input"
    />
  </div>
</template>

<style scoped>
.search-section {
  margin-bottom: 20px;
}

.search-input {
  width: 100%;
  padding: 12px;
  border: 1px solid #333;
  border-radius: 8px;
  background-color: #1a1a1a;
  color: white;
  font-size: 1rem;
  outline: none;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.search-input:focus {
  border-color: #42b883;
}

.search-input::placeholder {
  color: #666;
}
</style>
