<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { reactive, ref, onMounted, computed } from 'vue'
import { useWeaponStore } from '../stores/weapons'
import { useAuthStore } from '../stores/auth'
import { useToast } from 'vue-toastification'

const router = useRouter()
const route = useRoute()
const store = useWeaponStore()
const authStore = useAuthStore()
const toast = useToast()

const editParam = Array.isArray(route.query.edit) ? route.query.edit[0] : route.query.edit
const suggestionId = ref<number | null>(editParam ? Number(editParam) : null)
const isEditing = computed(() => suggestionId.value !== null)
const loadingSuggestion = ref(false)

const form = reactive({
  name: '',
  type: 'range',
  firemode: 'semi',
  craftable: true,
  stacksize: 1,
  description: '',
  shortname: '',
  icon: '',
  capacity: null as number | null,
  time_to_craft: null as number | null,
})

const iconFile = ref<File | null>(null)
const iconPreview = ref('')

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

const setIngredientAmount = (id: number, raw: string, el?: HTMLInputElement) => {
  const item = selectedIngredients.value.find(i => i.id === id)
  if (!item) return
  const cleaned = raw.replace(/\D/g, '').replace(/^0+/, '')
  item.amount = parseInt(cleaned) || 1
  if (el) el.value = cleaned
}

const preventNonDigit = (e: KeyboardEvent) => {
  if (!/[0-9]/.test(e.key) && e.key !== 'Backspace' && e.key !== 'Delete' && e.key !== 'Tab' && e.key !== 'ArrowLeft' && e.key !== 'ArrowRight') {
    e.preventDefault()
  }
}

const handleBlur = (id: number, raw: string, el?: HTMLInputElement) => {
  const trimmed = raw.trim()
  if (!trimmed || trimmed === '0') {
    toggleIngredient(id)
  } else {
    setIngredientAmount(id, raw, el)
  }
}

onMounted(async () => {
  if (!store.ammoList.length) store.fetchAllAmmo()
  if (!store.modList.length) store.fetchAllMods()
  if (!store.ingredientList.length) store.fetchAllIngredients()

  if (isEditing.value && suggestionId.value) {
    loadingSuggestion.value = true
    try {
      const s = await store.fetchMySuggestion(suggestionId.value)
      const p = s.payload
      form.name = p.name || ''
      form.type = p.type || 'range'
      form.firemode = p.firemode || 'semi'
      form.craftable = p.craftable ?? true
      form.stacksize = p.stacksize ?? 1
      form.description = p.description || ''
      form.shortname = p.shortname || ''
      form.icon = p.icon || ''
      form.capacity = p.capacity ?? null
      form.time_to_craft = p.time_to_craft ?? null
      selectedAmmo.value = p.ammo_ids || []
      selectedMods.value = p.mod_ids || []
      selectedIngredients.value = (p.ingredients || []).map((i: any) => ({ id: i.id, amount: i.amount }))
      if (p.icon) {
        iconPreview.value = p.icon
      }
    } catch (err) {
      toast.error(`Failed to load suggestion: ${err instanceof Error ? err.message : 'Unknown error'}`)
      router.push('/my/suggestions')
    } finally {
      loadingSuggestion.value = false
    }
  }
})

const WEAPON_CATEGORY_ID = 1

const handleFileSelect = (e: Event) => {
  const target = e.target
  if (!(target instanceof HTMLInputElement)) return
  const file = target.files?.[0]
  if (!file) return

  if (iconPreview.value) {
    URL.revokeObjectURL(iconPreview.value)
  }
  iconFile.value = file
  iconPreview.value = URL.createObjectURL(file)
}

function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

const handleSubmit = async () => {
  let icon = form.icon || ''
  let iconBase64 = ''

  if (iconFile.value) {
    if (authStore.isAdmin) {
      icon = await store.uploadIcon(iconFile.value)
    } else {
      iconBase64 = await fileToBase64(iconFile.value)
    }
  }

  const payload: Record<string, unknown> = {
    name: form.name,
    type: form.type,
    firemode: form.firemode,
    craftable: form.craftable,
    stacksize: form.stacksize,
    description: form.description,
    shortname: form.shortname,
    icon,
    icon_base64: iconBase64,
    capacity: form.capacity || null,
    time_to_craft: form.time_to_craft || null,
    category_id: WEAPON_CATEGORY_ID,
    ammo_ids: selectedAmmo.value,
    mod_ids: selectedMods.value,
    ingredients: selectedIngredients.value,
  }

  try {
    if (isEditing.value && suggestionId.value) {
      await store.resubmitSuggestion(suggestionId.value, payload)
      toast.success('Suggestion resubmitted for review!')
      router.push('/my/suggestions')
    } else if (authStore.isAdmin) {
      const id = await store.createWeapon(payload)
      toast.success('Weapon created successfully!')
      await store.fetchWeapons()
      router.push(`/weapon/${id}`)
    } else {
      await store.createSuggestion(payload)
      toast.success('Suggestion submitted for review!')
      router.push('/')
    }
  } catch (err) {
    toast.error(`Failed to create weapon: ${err instanceof Error ? err.message : 'Unknown error'}`)
  }
}
</script>

<template>
  <div class="wiki-page">
    <button class="back-btn" @click="router.push('/')">← Back to list</button>
    <h1 class="page-title">{{ isEditing ? 'Edit Weapon' : 'Create Weapon' }}</h1>

    <div v-if="loadingSuggestion" class="loading">Loading suggestion...</div>
    <form v-else @submit.prevent="handleSubmit" class="create-form">
      <div class="form-group">
        <label for="name">Name</label>
        <input id="name" v-model="form.name" type="text" class="form-input" required />
      </div>

      <div class="form-group">
        <label for="description">Description</label>
        <textarea id="description" v-model="form.description" class="form-input" rows="4" maxlength="500"></textarea>
      </div>

      <div class="form-group">
        <label for="shortname">Shortname</label>
        <input id="shortname" v-model="form.shortname" type="text" class="form-input" required />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="type">Type</label>
          <select id="type" v-model="form.type" class="form-input">
            <option value="range">Range</option>
            <option value="melee">Melee</option>
          </select>
        </div>

        <div class="form-group">
          <label for="firemode">Firemode</label>
          <select id="firemode" v-model="form.firemode" class="form-input">
            <option value="semi">Semi</option>
            <option value="automatic">Automatic</option>
            <option value="double">Double</option>
            <option value="none">None</option>
          </select>
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="craftable">
            <input id="craftable" type="checkbox" v-model="form.craftable" class="inline-checkbox" />
            Craftable
          </label>
        </div>

        <div class="form-group">
          <label for="stacksize">Stacksize</label>
          <input id="stacksize" v-model.number="form.stacksize" type="number" min="1" class="form-input" />
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="capacity">Capacity</label>
          <input id="capacity" v-model.number="form.capacity" type="number" min="1" class="form-input" placeholder="e.g. 30" />
        </div>

        <div class="form-group">
          <label for="time_to_craft">Time to Craft (seconds)</label>
          <input id="time_to_craft" v-model.number="form.time_to_craft" type="number" min="1" class="form-input" placeholder="e.g. 15" />
        </div>
      </div>

      <div class="form-group">
        <label for="icon">Icon</label>
        <div class="icon-upload-area">
          <label class="file-input-label">
            <input type="file" accept="image/avif,image/jpeg,image/png,image/webp" @change="handleFileSelect" class="file-input" />
            <span>Choose image</span>
          </label>
          <img v-if="iconPreview" :src="iconPreview" class="icon-preview" alt="Icon preview" />
          <span v-else class="icon-placeholder">No icon selected</span>
        </div>
      </div>

      <div class="form-group">
        <label>Ammunition</label>
        <div class="checkbox-grid">
          <label v-for="a in store.ammoList" :key="a.id" class="checkbox-item">
            <input type="checkbox" :value="a.id" v-model="selectedAmmo" />
            <span class="item-name">{{ a.name }}</span>
            <img v-if="a.icon" :src="a.icon" :alt="a.name" class="grid-icon" />
          </label>
          <div v-if="!store.ammoList.length" class="empty-text">Loading...</div>
        </div>
      </div>

      <div class="form-group">
        <label>Weapon Mods</label>
        <div class="checkbox-grid">
          <label v-for="m in store.modList" :key="m.id" class="checkbox-item">
            <input type="checkbox" :value="m.id" v-model="selectedMods" />
            <span class="item-name">{{ m.name }}</span>
            <img v-if="m.icon" :src="m.icon" :alt="m.name" class="grid-icon" />
          </label>
          <div v-if="!store.modList.length" class="empty-text">Loading...</div>
        </div>
      </div>

      <div class="form-group">
        <label>Ingredients</label>
        <div class="checkbox-grid">
          <div v-for="ing in store.ingredientList" :key="ing.id" class="checkbox-item" :data-title="ing.name" @click="toggleIngredient(ing.id)">
            <label class="checkbox-content">
              <input
                type="checkbox"
                :checked="isIngredientSelected(ing.id)"
                @click.stop="toggleIngredient(ing.id)"
              />
            </label>
            <input
              type="tel"
              min="1"
              maxlength="5"
              class="amount-input"
              :value="getIngredientAmount(ing.id)"
              @input="setIngredientAmount(ing.id, ($event.target as HTMLInputElement).value, $event.target as HTMLInputElement)"
              @blur="handleBlur(ing.id, ($event.target as HTMLInputElement).value, $event.target as HTMLInputElement)"
              @keydown="preventNonDigit"
              @click.stop
              :disabled="!isIngredientSelected(ing.id)"
            />
            <img v-if="ing.icon" :src="ing.icon" :alt="ing.name" class="grid-icon" />
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
  --clr-accent: #ce422b;
  --clr-accent-hover: #a8321f;
  --clr-bg-card: #464646;
  --clr-bg-input: #1a1a1a;
  --clr-bg-grid: #2a2a2a;
  --clr-border: #444;
  --clr-border-card: #5d5d5d;
  --clr-text: #e2e8f0;
  --clr-text-secondary: #94a3b8;
  --clr-text-muted: #64748b;
  --clr-bg-btn: #333;
  --clr-tooltip-bg: #222;
  --clr-text-btn: #fff;
  --radius: 4px;

  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  color: var(--clr-text);
}

.back-btn {
  background: var(--clr-bg-btn);
  color: var(--clr-text-btn);
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  margin-bottom: 20px;
}

.back-btn:hover { background: var(--clr-border); }

.page-title {
  font-size: 2.5rem;
  border-bottom: 1px solid var(--clr-border-card);
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.create-form {
  display: flex;
  flex-direction: column;
  max-width: 800px;
  margin: 0 auto;
}

.form-row {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

@media (max-width: 768px) {
  .form-row { flex-direction: column; }
}

@media (max-width: 480px) {
  .checkbox-grid {
    display: flex;
    flex-direction: column;
  }
  .checkbox-item.checkbox-item {
    padding: 10px 12px;
    min-height: 56px;
  }
  .icon-upload-area {
    flex-direction: column;
    align-items: flex-start;
  }
}

.form-row .form-group {
  flex: 1;
}

.inline-checkbox {
  width: 18px;
  height: 18px;
  accent-color: var(--clr-accent);
  margin-right: 8px;
  vertical-align: middle;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 0.95rem;
  color: var(--clr-text-secondary);
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  background: var(--clr-bg-input);
  border: 1px solid var(--clr-border);
  border-radius: var(--radius);
  color: var(--clr-text);
  font-size: 1rem;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: var(--clr-accent);
}

textarea.form-input {
  resize: vertical;
}

.submit-btn {
  background-color: var(--clr-accent);
  color: white;
  border: none;
  padding: 10px 24px;
  border-radius: var(--radius);
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.2s;
}

.submit-btn:hover {
  background-color: var(--clr-accent-hover);
}

.checkbox-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 4px;
  padding: 10px;
  background-color: var(--clr-bg-grid);
  border-radius: var(--radius);
}

.checkbox-item.checkbox-item {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 10px;
  padding: 6px 10px;
  min-height: 50px;
  background: var(--clr-bg-card);
  border: 1px solid var(--clr-border-card);
  border-radius: var(--radius);
  font-size: 0.85rem;
  transition: background-color 0.2s;
  overflow: hidden;
  cursor: pointer;
}

.checkbox-item[data-title]:hover::after {
  content: attr(data-title);
  position: absolute;
  bottom: calc(100% + 4px);
  left: 50%;
  transform: translateX(-50%);
  background: var(--clr-tooltip-bg);
  color: var(--clr-text);
  padding: 4px 8px;
  border-radius: var(--radius);
  font-size: 0.8rem;
  white-space: nowrap;
  pointer-events: none;
  z-index: 10;
}

.checkbox-content {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
}

.checkbox-item input[type="checkbox"] {
  width: 18px;
  height: 18px;
  margin: 0;
  accent-color: var(--clr-accent);
  cursor: pointer;
  flex-shrink: 0;
}

.amount-input {
  width: 70px;
  padding: 10px 8px;
  background: var(--clr-bg-input);
  border: 1px solid var(--clr-border);
  border-radius: var(--radius);
  color: var(--clr-text);
  font-size: 0.85rem;
  text-align: center;
  justify-self: center;
  flex-shrink: 0;
}

.amount-input::-webkit-outer-spin-button,
.amount-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.amount-input {
  -moz-appearance: textfield; 
  appearance: textfield;    
}

.amount-input:disabled {
  opacity: 0.3;
}

.amount-input:focus {
  outline: none;
  border-color: var(--clr-accent);
}

.grid-icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
  flex-shrink: 0;
}

.empty-text {
  font-size: 0.85rem;
  color: var(--clr-text-muted);
  padding: 4px;
}

.icon-upload-area {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  background-color: var(--clr-bg-grid);
  border-radius: var(--radius);
}

.file-input-label {
  display: inline-block;
  padding: 8px 16px;
  background-color: var(--clr-accent);
  color: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.2s;
  white-space: nowrap;
}

.file-input-label:hover {
  background-color: var(--clr-accent-hover);
}

.file-input {
  display: none;
}

.icon-preview {
  width: 64px;
  height: 64px;
  object-fit: contain;
  border-radius: var(--radius);
  border: 1px solid var(--clr-border);
}

.icon-placeholder {
  font-size: 0.85rem;
  color: var(--clr-text-muted);
}

.loading {
  text-align: center;
  padding: 40px;
  color: var(--clr-text-secondary);
  font-size: 1.1rem;
}
</style>
