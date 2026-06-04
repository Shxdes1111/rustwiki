<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useWeaponStore, type WeaponItem } from '../stores/weapons'
import { useRouter } from 'vue-router'

const props = defineProps<{ id: string }>()
const store = useWeaponStore()
const router = useRouter()

const weapon = ref<WeaponItem | null>(null)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    weapon.value = await store.fetchWeapon(Number(props.id))
  } catch {
    error.value = 'Не удалось загрузить данные. Проверьте, запущен ли бэкенд.'
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
          <img v-if="weapon.icon" :src="weapon.icon" :alt="weapon.name" class="weapon-icon" />
          <div v-else class="placeholder-img">🎯</div>
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

        <!-- Раздел: Ammunition (Переведен на умную компактную сетку) -->
        <div class="infobox-section-title">Ammunition</div>
        <div v-if="weapon.ammo?.length" class="infobox-grid">
          <router-link 
            v-for="ammo in weapon.ammo" 
            :key="ammo.id" 
            :to="`/ammo/${ammo.id}`" 
            class="grid-item" 
            :title="ammo.name"
          >
            <span class="item-icon">📦</span> 
            <span class="grid-item-text">{{ ammo.name }}</span>
          </router-link>
        </div>
        <div v-else class="empty-box">No custom ammo slots</div>

        <!-- Раздел: Weapon Mods (Переведен на умную компактную сетку) -->
        <div class="infobox-section-title">Weapon Mods</div>
        <div v-if="weapon.mods?.length" class="infobox-grid">
          <router-link 
            v-for="mod in weapon.mods" 
            :key="mod.id" 
            :to="`/mods/${mod.id}`" 
            class="grid-item" 
            :title="mod.name"
          >
            <span class="item-icon">🔧</span> 
            <span class="grid-item-text">{{ mod.name }}</span>
          </router-link>
        </div>
        <div v-else class="empty-box">No mods supported</div>

        <!-- Раздел: Crafting Info -->
        <div class="infobox-section-title">Crafting</div>
        <div class="infobox-row"><span>Craftable</span><span class="value">{{ weapon.craftable ? 'Yes' : 'No' }}</span></div>
        <div class="infobox-row" v-if="weapon.timeToCraft"><span>Time To Craft</span><span class="value">{{ weapon.timeToCraft }} s</span></div>

        <!-- Раздел: Ingredients Icons -->
        <div class="infobox-section-title">Ingredients</div>
        <div class="infobox-grid bg-darker">
          <div v-for="ing in weapon.ingredients" :key="ing.id" class="grid-item content-center">
            💎 <span class="count">x{{ ing.amount || 1 }}</span>
          </div>
        </div>
      </aside>
    </div>
  </div>
  <div v-else-if="error" class="loading error">
    {{ error }}
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
  border-bottom: 1px solid #5d5d5d;
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
  border-bottom: 1px solid #5d5d5d;
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
  background-color: #4b4b4c;
  border: 1px solid #5d5d5d;
  border-radius: 6px;
  overflow: hidden;
  align-self: flex-start;
  flex-shrink: 0;
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

.weapon-icon {
  max-width: 200px;
  max-height: 200px;
  display: block;
  margin: 0 auto 10px;
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

/* Оптимизированная CSS Grid сетка для модов и патронов */
.infobox-grid {
  display: grid;
  /* В инфобоксе шириной 320px две колонки по ~130px встанут идеально с учетом padding */
  grid-template-columns: repeat(auto-fill, minmax(125px, 1fr));
  gap: 6px;
  padding: 10px;
  background-color: #2a2a2a;
}

.bg-darker {
  background-color: #1f1f1f;
}

/* Элемент сетки (Мод / Патрон) */
.grid-item {
  background: #464646;
  border: 1px solid #5d5d5d;
  padding: 6px 8px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8rem;
  text-decoration: none;
  color: inherit;
  transition: all 0.2s ease;
  
  /* Защита от деформации и переносов */
  overflow: hidden;
  white-space: nowrap;
}

.grid-item:hover {
  background: #5a5a5a;
  border-color: #ef4444; /* Красная рамка при наведении в тон шапки */
}

.item-icon {
  flex-shrink: 0; /* Иконка эмодзи/картинки не сожмется */
}

.grid-item-text {
  overflow: hidden;
  text-overflow: ellipsis; /* Длинное название плавно скроется троеточием */
}

.content-center {
  justify-content: center;
}

.count {
  font-weight: bold;
  color: #cbd5e1;
}

/* Пустой контейнер для красивого отображения отсутствия элементов */
.empty-box {
  font-size: 0.85rem;
  color: #8392a5;
  font-style: italic;
  padding: 12px;
  background-color: #2a2a2a;
  text-align: center;
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
