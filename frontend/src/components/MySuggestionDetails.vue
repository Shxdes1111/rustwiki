<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useWeaponStore } from '../stores/weapons'
import { useToast } from 'vue-toastification'

const props = defineProps<{ id: string }>()

const router = useRouter()
const store = useWeaponStore()
const toast = useToast()

const suggestion = ref<any>(null)
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    suggestion.value = await store.fetchMySuggestion(Number(props.id))
    if (!store.ammoList.length) await store.fetchAllAmmo()
    if (!store.modList.length) await store.fetchAllMods()
    if (!store.ingredientList.length) await store.fetchAllIngredients()
  } catch (err) {
    toast.error(`Failed to load suggestion: ${err instanceof Error ? err.message : 'Unknown error'}`)
  } finally {
    loading.value = false
  }
})

const p = computed(() => suggestion.value?.payload || {})

const ammoNames = computed(() =>
  (p.value?.ammo_ids || []).map((id: number) =>
    store.ammoList.find(a => a.id === id)?.name || `Ammo #${id}`
  )
)

const modNames = computed(() =>
  (p.value?.mod_ids || []).map((id: number) =>
    store.modList.find(m => m.id === id)?.name || `Mod #${id}`
  )
)

const ingredientItems = computed(() =>
  (p.value?.ingredients || []).map((item: any) => {
    const ing = store.ingredientList.find(i => i.id === item.id)
    return { name: ing?.name || `Ingredient #${item.id}`, amount: item.amount }
  })
)

const goToEdit = () => {
  router.push(`/weapon/create?edit=${props.id}`)
}
</script>

<template>
  <div class="page">
    <button class="back-btn" @click="router.push('/my/suggestions')">← Back to my suggestions</button>

    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="!suggestion" class="loading">Suggestion not found.</div>
    <template v-else>
      <h1 class="page-title">{{ p.name || 'Suggestion' }}</h1>

      <div class="meta-bar">
        <span :class="['badge', `badge-${suggestion.status}`]">{{ suggestion.status }}</span>
        <span class="meta">{{ new Date(suggestion.created_at).toLocaleString() }}</span>
        <span v-if="suggestion.reviewed_at" class="meta">Reviewed: {{ new Date(suggestion.reviewed_at).toLocaleString() }}</span>
      </div>

      <div v-if="suggestion.status === 'rejected' && suggestion.rejection_reason" class="rejection-reason">
        <strong>Rejection reason:</strong> {{ suggestion.rejection_reason }}
      </div>

      <div class="details-grid">
        <div class="detail-card">
          <h3>Basic Info</h3>
          <div class="info-row"><span class="label">Name</span><span>{{ p.name }}</span></div>
          <div class="info-row"><span class="label">Type</span><span>{{ p.type }}</span></div>
          <div class="info-row"><span class="label">Firemode</span><span>{{ p.firemode }}</span></div>
          <div class="info-row"><span class="label">Shortname</span><span>{{ p.shortname }}</span></div>
          <div class="info-row"><span class="label">Craftable</span><span>{{ p.craftable ? 'Yes' : 'No' }}</span></div>
          <div class="info-row"><span class="label">Stacksize</span><span>{{ p.stacksize }}</span></div>
          <div v-if="p.capacity" class="info-row"><span class="label">Capacity</span><span>{{ p.capacity }}</span></div>
          <div v-if="p.time_to_craft" class="info-row"><span class="label">Time to Craft</span><span>{{ p.time_to_craft }}s</span></div>
        </div>

        <div v-if="p.description" class="detail-card description-card">
          <h3>Description</h3>
          <p>{{ p.description }}</p>
        </div>

        <div v-if="ammoNames.length" class="detail-card">
          <h3>Ammunition</h3>
          <ul>
            <li v-for="name in ammoNames" :key="name">{{ name }}</li>
          </ul>
        </div>

        <div v-if="modNames.length" class="detail-card">
          <h3>Mods</h3>
          <ul>
            <li v-for="name in modNames" :key="name">{{ name }}</li>
          </ul>
        </div>

        <div v-if="ingredientItems.length" class="detail-card">
          <h3>Ingredients</h3>
          <ul>
            <li v-for="item in ingredientItems" :key="item.name">{{ item.name }} × {{ item.amount }}</li>
          </ul>
        </div>

        <div v-if="p.icon || p.icon_base64" class="detail-card">
          <h3>Icon</h3>
          <img :src="p.icon || p.icon_base64" alt="Weapon icon" class="icon-preview" />
        </div>
      </div>

      <div v-if="suggestion.status === 'rejected'" class="actions-bar">
        <button class="btn-edit" @click="goToEdit">Edit and resubmit</button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
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
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
  font-size: 1.1rem;
}

.meta-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.badge {
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 0.9rem;
  text-transform: capitalize;
}

.badge-pending { background: #854d0e; color: #fef08a; }
.badge-approved { background: #166534; color: #bbf7d0; }
.badge-rejected { background: #991b1b; color: #fecaca; }

.meta {
  color: #94a3b8;
  font-size: 0.9rem;
}

.rejection-reason {
  width: 100%;
  margin-bottom: 20px;
  padding: 12px 16px;
  background: #3b0f0f;
  border: 1px solid #991b1b;
  border-radius: 6px;
  color: #fecaca;
  font-size: 0.95rem;
  line-height: 1.4;
}

.details-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.detail-card {
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 16px;
  min-width: 0;
}

.detail-card h3 {
  margin: 0 0 12px;
  font-size: 1rem;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  border-bottom: 1px solid #2a2a2a;
  font-size: 0.9rem;
  overflow: hidden;
}

.info-row span:last-child {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  min-width: 0;
}

.info-row .label {
  color: #94a3b8;
}

.description-card p {
  margin: 0;
  font-size: 0.9rem;
  line-height: 1.5;
  color: #cbd5e1;
  overflow-wrap: break-word;
}

ul {
  margin: 0;
  padding-left: 20px;
  font-size: 0.9rem;
}

li {
  padding: 2px 0;
}

.icon-preview {
  max-width: 96px;
  max-height: 96px;
  object-fit: contain;
  border-radius: 4px;
}

.actions-bar {
  display: flex;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #333;
}

.btn-edit {
  background: #2563eb;
  color: white;
  border: none;
  padding: 10px 24px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  transition: background 0.2s;
}

.btn-edit:hover { background: #1d4ed8; }
</style>
