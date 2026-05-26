<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useWeaponStore, type WeaponItem } from '../stores/weapons'
import { useRoute, useRouter } from 'vue-router'

const props = defineProps<{ id: string }>()
const store = useWeaponStore()
const router = useRouter()

// Находим нужное оружие в сторе по ID
const weapon = computed<WeaponItem | undefined>(() => {
  return store.weapons.find(w => w.id === Number(props.id))
})

// На случай, если пользователь обновил страницу, и стор пустой — перезапрашиваем данные
onMounted(async () => {
  if (store.weapons.length === 0) {
    await store.fetchWeapons()
  }
})
</script>

<template>
  <div v-if="weapon" class="wiki-page">
    <!-- Кнопка Назад -->
    <button class="back-btn" @click="router.push('/')">← Назад к списку</button>

    <div class="wiki-layout">
      <!-- ЛЕВАЯ КОЛОНКА: Описание и крафт -->
      <div class="main-content">
        <h1 class="page-title">{{ weapon.name }}</h1>
        
        <p class="description-text">
          {{ weapon.description || `The ${weapon.name} is a ${weapon.type} weapon. Detailed description will be loaded from DB.` }}
        </p>

        <div v-if="weapon.craftable" class="crafting-section">
          <h2>Crafting</h2>
          <p>The {{ weapon.name }} can be crafted using:</p>
          <ul>
            <li v-for="ing in weapon.ingredients" :key="ing.id">
              <span class="bullet">•</span> {{ ing.amount }}x {{ ing.name }}
            </li>
          </ul>
        </div>
      </div>

      <!-- ПРАВАЯ КОЛОНКА: Инфобокс (Карточка характеристик) -->
      <aside class="infobox">
        <div class="infobox-header">{{ weapon.name }}</div>
        
        <div class="infobox-image-box">
          <!-- Вместо заглушки тут будет динамическая картинка по имени: `/images/${weapon.shortname}.png` -->
          <div class="placeholder-img">🎯</div>
          <span class="infobox-subtitle">A professional rust item viewer.</span>
        </div>

        <!-- Раздел: General -->
        <div class="infobox-section-title">General</div>
        <div class="infobox-row"><span>Shortname</span><span class="value">{{ weapon.shortname || 'N/A' }}</span></div>
        <div class="infobox-row"><span>Type</span><span class="value">{{ weapon.type }}</span></div>
        <div class="infobox-row"><span>Stacksize</span><span class="value">{{ weapon.stacksize }}</span></div>

        <!-- Раздел: Weapon Stats -->
        <div class="infobox-section-title">Weapon Stats</div>
        <div class="infobox-row"><span>Fire Mode</span><span class="value">{{ weapon.firemode }}</span></div>
        <div class="infobox-row"><span>Capacity</span><span class="value">{{ weapon.capacity || '—' }}</span></div>

        <!-- Раздел: Ammunition -->
        <div class="infobox-section-title">Ammunition</div>
        <div class="infobox-grid">
          <div v-for="ammo in weapon.ammo" :key="ammo.id" class="grid-item" :title="ammo.name">
            📦 <span class="grid-item-text">{{ ammo.name }}</span>
          </div>
          <div v-if="!weapon.ammo?.length" class="empty-text">No custom ammo slots</div>
        </div>

        <!-- Раздел: Weapon Mods -->
        <div class="infobox-section-title">Weapon Mods</div>
        <div class="infobox-grid">
          <div v-for="mod in weapon.mods" :key="mod.id" class="grid-item" :title="mod.name">
            🔧 <span class="grid-item-text">{{ mod.name }}</span>
          </div>
          <div v-if="!weapon.mods?.length" class="empty-text">No mods supported</div>
        </div>

        <!-- Раздел: Crafting Info -->
        <div class="infobox-section-title">Crafting</div>
        <div class="infobox-row"><span>Craftable</span><span class="value">{{ weapon.craftable ? 'Yes' : 'No' }}</span></div>
        <div class="infobox-row" v-if="weapon.timeToCraft"><span>Time To Craft</span><span class="value">{{ weapon.timeToCraft }} s</span></div>

        <!-- Раздел: Ingredients Icons -->
        <div class="infobox-section-title">Ingredients</div>
        <div class="infobox-grid bg-darker">
          <div v-for="ing in weapon.ingredients" :key="ing.id" class="grid-item">
            💎 <span class="count">x{{ ing.amount || 1 }}</span>
          </div>
        </div>
      </aside>
    </div>
  </div>
  <div v-else class="loading">
    Загрузка данных об оружии...
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
  border-bottom: 1px solid #334155;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.description-text {
  line-height: 1.6;
  font-size: 1.1rem;
  color: #cbd5e1;
}

.crafting-section {
  margin-top: 40px;
}

.crafting-section h2 {
  font-size: 1.8rem;
  border-bottom: 1px solid #334155;
  padding-bottom: 5px;
}

.crafting-section ul {
  list-style: none;
  padding-left: 0;
  margin-top: 15px;
}

.crafting-section li {
  margin-bottom: 8px;
  font-size: 1.1rem;
}

.bullet {
  color: #ef4444;
  margin-right: 8px;
}

/* СТИЛИ ИНФОБОКСА (Правая панель) */
.infobox {
  width: 320px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  overflow: hidden;
  align-self: flex-start;
}

.infobox-header {
  background-color: #ef4444;
  color: white;
  text-align: center;
  padding: 10px;
  font-size: 1.3rem;
  font-weight: bold;
}

.infobox-image-box {
  padding: 20px;
  text-align: center;
  background-color: #0f172a;
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
  background-color: #ef4444;
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
  border-bottom: 1px solid #334155;
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
  background-color: #0f172a;
}

.grid-item {
  background: #1e293b;
  border: 1px solid #334155;
  padding: 8px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.8rem;
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
</style>
