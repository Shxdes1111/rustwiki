<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useWeaponStore, type ModDetail } from '../stores/weapons'
import { useRouter } from 'vue-router'

const props = defineProps<{ id: string }>()
const store = useWeaponStore()
const router = useRouter()

const mod = ref<ModDetail | null>(null)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    mod.value = await store.fetchMod(Number(props.id))
  } catch {
    error.value = 'Не удалось загрузить данные. Проверьте, запущен ли бэкенд.'
  }
})
</script>

<template>
  <div v-if="mod" class="wiki-page">
    <button class="back-btn" @click="router.push('/')">← Назад к списку</button>

    <div class="wiki-layout">
      <div class="main-content">
        <h1 class="page-title">{{ mod.name }}</h1>
        <p class="description-text">{{ mod.name }} — модуль для оружия в Rust.</p>
      </div>

      <aside class="infobox">
        <div class="infobox-header">{{ mod.name }}</div>

        <div class="infobox-image-box">
          <div class="placeholder-img">🔧</div>
          <span class="infobox-subtitle">Weapon Mod</span>
        </div>

        <div class="infobox-section-title">General</div>
        <div class="infobox-row"><span>ID</span><span class="value">{{ mod.id }}</span></div>
        <div class="infobox-row"><span>Name</span><span class="value">{{ mod.name }}</span></div>

        <div class="infobox-section-title">Compatible Weapons</div>
        <div class="infobox-grid">
          <router-link v-for="weapon in mod.compatible_weapons" :key="weapon.id" :to="`/weapon/${weapon.id}`" class="grid-item" :title="weapon.name">
            🔫 <span class="grid-item-text">{{ weapon.name }}</span>
          </router-link>
          <div v-if="!mod.compatible_weapons?.length" class="empty-text">No compatible weapons</div>
        </div>
      </aside>
    </div>
  </div>
  <div v-else-if="error" class="loading error">
    {{ error }}
  </div>
  <div v-else class="loading">
    Загрузка данных о модулях...
  </div>
</template>

<style scoped>
.wiki-page {
  max-width: 1200px;
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

.wiki-layout {
  display: flex;
  gap: 40px;
}

.main-content {
  flex: 1;
}

.page-title {
  font-size: 2.5rem;
  border-bottom: 1px solid #5d5d5d;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.description-text {
  line-height: 1.6;
  font-size: 1.1rem;
  color: #cbd5e1;
}

.infobox {
  width: 320px;
  background-color: #4b4b4c;
  border: 1px solid #5d5d5d;
  border-radius: 6px;
  overflow: hidden;
  align-self: flex-start;
}

.infobox-header {
  background-color: #9f2f2f;
  color: white;
  text-align: center;
  padding: 10px;
  font-size: 1.3rem;
  font-weight: bold;
}

.infobox-image-box {
  padding: 20px;
  text-align: center;
  background-color: #2a2a2a;
}

.placeholder-img {
  font-size: 4rem;
  margin-bottom: 10px;
}

.infobox-subtitle {
  font-size: 0.85rem;
  color: #94a3b8;
  font-style: italic;
}

.infobox-section-title {
  background-color: #9f2f2f;
  color: white;
  padding: 6px 12px;
  font-size: 0.95rem;
  font-weight: bold;
  text-transform: uppercase;
  margin-top: 1px;
}

.infobox-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 12px;
  border-bottom: 1px solid #676767;
  font-size: 0.9rem;
}

.infobox-row .value {
  font-weight: bold;
  color: #f8fafc;
}

.infobox-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px;
  background-color: #2a2a2a;
}

.grid-item {
  background: #464646;
  border: 1px solid #5d5d5d;
  padding: 8px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.8rem;
  text-decoration: none;
  color: inherit;
  transition: background-color 0.2s;
}

.grid-item:hover {
  background: #5a5a5a;
}

.empty-text {
  font-size: 0.85rem;
  color: #64748b;
  padding: 4px;
}

.loading {
  text-align: center;
  padding: 40px;
  font-size: 1.2rem;
}

.loading.error {
  color: #9f2f2f;
}
</style>
