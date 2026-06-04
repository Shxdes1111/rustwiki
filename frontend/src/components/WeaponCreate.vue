<script setup lang="ts">
import { useRouter } from 'vue-router'
import { reactive, ref, onMounted } from 'vue'
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

onMounted(() => {
  if (!store.ammoList.length) store.fetchAllAmmo()
  if (!store.modList.length) store.fetchAllMods()
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
  gap: 6px;
  padding: 6px 8px;
  background: #464646;
  border: 1px solid #5d5d5d;
  border-radius: 4px;
  font-size: 0.8rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.checkbox-item:hover {
  background: #5a5a5a;
}

.checkbox-item input[type="checkbox"] {
  flex-shrink: 0;
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
