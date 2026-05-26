<script setup lang="ts">
import { onMounted } from 'vue'
import { useWeaponStore } from './stores/weapons'
import SearchBar from './components/SearchBar.vue'
import WeaponTable from './components/WeaponTable.vue'

const store = useWeaponStore()

onMounted(() => store.fetchWeapons())
</script>

<template>
  <div class="wiki-container">
    <h1>RustWiki Weapons</h1>
    <SearchBar />
    <WeaponTable />
  </div>
</template>

<style scoped>
.wiki-container {
  padding: 40px;
  max-width: 1000px;
  margin: 0 auto;
  color: white;
}
</style>
