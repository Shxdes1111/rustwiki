<script setup lang="ts">
import { useRouter } from 'vue-router'
import { reactive, ref, onMounted, computed } from 'vue'
import { useWeaponStore } from '../stores/weapons'

const router = useRouter()
const store = useWeaponStore()

const form = reactive({
  name: '',
  description: '',
  shortname: ''
})

const selectedAmmo = ref<number[]>([])
const selectedMods = ref<number[]>([])
const selectedIngredients = ref<{ id: number; amount: number }[]>([])

const isIngredientSelected = (id: number) => selectedIngredients.value.some(i => i.id === id)

const getIngredientAmount = (id: number) => {
  return selectedIngredients.value.find(i => i.id === id)?.amount ?? 1
}

const toggleIngredient = (id: number) => {
  const idx = selectedIngredients.value.findIndex(i => i.id === id)
  if (idx >= 0) {
    selectedIngredients.value.splice(idx, 1)
  } else {
    selectedIngredients.value.push({ id, amount: 1 })
  }
}

const setIngredientAmount = (id: number, amount: number) => {
  const item = selectedIngredients.value.find(i => i.id === id)
  if (item) item.amount = Math.max(1, amount)
}

onMounted(() => {
  if (!store.ammoList.length) store.fetchAllAmmo()
  if (!store.modList.length) store.fetchAllMods()
  if (!store.ingredientList.length) store.fetchAllIngredients()
})

const handleSubmit = () => {
  // TODO: POST /api/weapons
}
</script>

<template>
  <div class="wiki-page">
    <button class="back-btn" @click="router.push('/')">← Back to list</button>
    <h1 class="page-title">Create Weapon</h1>

    <form @submit.prevent="handleSubmit" class="create-form">
      <div class="form-group">
        <label for="name">Name</label>
        <input id="name" v-model="form.name" type="text" class="form-input" required />
      </div>

      <div class="form-group">
        <label for="description">Description</label>
        <textarea id="description" v-model="form.description" class="form-input" rows="4"></textarea>
      </div>

      <div class="form-group">
        <label for="shortname">Shortname</label>
        <input id="shortname" v-model="form.shortname" type="text" class="form-input" required />
      </div>

      <div class="form-group">
        <label>Ammunition</label>
        <div class="checkbox-grid">
          <label v-for="a in store.ammoList" :key="a.id" class="checkbox-item">
            <input type="checkbox" :value="a.id" v-model="selectedAmmo" />
            <img v-if="a.icon" :src="a.icon" :alt="a.name" class="grid-icon" />
            <span>{{ a.name }}</span>
          </label>
          <div v-if="!store.ammoList.length" class="empty-text">Loading...</div>
        </div>
      </div>

      <div class="form-group">
        <label>Weapon Mods</label>
        <div class="checkbox-grid">
          <label v-for="m in store.modList" :key="m.id" class="checkbox-item">
            <input type="checkbox" :value="m.id" v-model="selectedMods" />
            <img v-if="m.icon" :src="m.icon" :alt="m.name" class="grid-icon" />
            <span>{{ m.name }}</span>
          </label>
          <div v-if="!store.modList.length" class="empty-text">Loading...</div>
        </div>
      </div>

      <div class="form-group">
        <label>Ingredients</label>
        <div class="checkbox-grid">
          <div v-for="ing in store.ingredientList" :key="ing.id" class="checkbox-item" :title="ing.name">
            <label class="checkbox-content">
              <input
                type="checkbox"
                :checked="isIngredientSelected(ing.id)"
                @change="toggleIngredient(ing.id)"
              />
              <img v-if="ing.icon" :src="ing.icon" :alt="ing.name" class="grid-icon" />
            </label>
            <input
              type="number"
              min="1"
              class="amount-input"
              :value="getIngredientAmount(ing.id)"
              @input="setIngredientAmount(ing.id, Number(($event.target as HTMLInputElement).value))"
              :disabled="!isIngredientSelected(ing.id)"
            />
          </div>
          <div v-if="!store.ingredientList.length" class="empty-text">Loading...</div>
        </div>
      </div>

      <button type="submit" class="submit-btn">Save</button>
    </form>
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

.page-title {
  font-size: 2.5rem;
  border-bottom: 1px solid #5d5d5d;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.create-form {
  max-width: 500px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 0.95rem;
  color: #94a3b8;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  background: #1a1a1a;
  border: 1px solid #444;
  border-radius: 4px;
  color: #e2e8f0;
  font-size: 1rem;
  font-family: sans-serif;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #ce422b;
}

textarea.form-input {
  resize: vertical;
}

.submit-btn {
  background-color: #ce422b;
  color: white;
  border: none;
  padding: 10px 24px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.2s;
}

.submit-btn:hover {
  background-color: #a8321f;
}

.checkbox-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(125px, 1fr));
  gap: 6px;
  padding: 10px;
  background-color: #2a2a2a;
  border-radius: 4px;
}

.checkbox-item {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 6px 8px;
  background: #464646;
  border: 1px solid #5d5d5d;
  border-radius: 4px;
  font-size: 0.8rem;
  transition: background-color 0.2s;
}

.checkbox-item:hover {
  background: #5a5a5a;
}

.checkbox-content {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
}

.checkbox-item input[type="checkbox"] {
  flex-shrink: 0;
}

.amount-input {
  width: 70px;
  padding: 10px 8px;
  background: #1a1a1a;
  border: 1px solid #444;
  border-radius: 3px;
  color: #e2e8f0;
  font-size: 0.85rem;
  text-align: center;
  flex-shrink: 0;
}

.amount-input::-webkit-outer-spin-button,
.amount-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.amount-input[type="number"] {
  -moz-appearance: textfield;
}

.amount-input:disabled {
  opacity: 0.3;
}

.amount-input:focus {
  outline: none;
  border-color: #ce422b;
}

.grid-icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
  flex-shrink: 0;
}

.empty-text {
  font-size: 0.85rem;
  color: #64748b;
  padding: 4px;
}
</style>
