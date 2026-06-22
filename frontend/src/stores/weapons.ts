import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'

export interface Ingredient {
  id: number;
  name: string;
  amount?: number;
  icon?: string;
}

export interface Ammo {
  id: number;
  name: string;
  icon?: string;
}

export interface WeaponMod {
  id: number;
  name: string;
  icon?: string;
}

export interface AmmoDetail extends Ammo {
  compatible_weapons: WeaponItem[]
}

export interface ModDetail extends WeaponMod {
  compatible_weapons: WeaponItem[]
}

export interface Suggestion {
  id: number
  user_id: number
  username?: string
  payload: any
  status: string
  created_at: string
  reviewed_at?: string
  reviewed_by?: number
}

export interface WeaponItem {
  id: number;
  name: string;
  type: string;
  description?: string;
  shortname?: string; 
  icon?: string;
  stacksize: number;
  firemode: string;
  capacity?: number;  
  craftable: boolean;
  timeToCraft?: number; 
  
  ingredients?: Ingredient[];
  ammo?: Ammo[];
  mods?: WeaponMod[];
}

const CACHE_TTL = 2 * 60 * 1000 // 2 минуты

export const useWeaponStore = defineStore('weapons', () => {
  const weapons = ref<WeaponItem[]>([])
  const searchTerm = ref('')

  const ammoList = ref<Ammo[]>([])
  const modList = ref<WeaponMod[]>([])
  const ingredientList = ref<Ingredient[]>([])

  const weaponCache = ref<Record<number, WeaponItem>>({})
  const weaponTimestamps = ref<Record<number, number>>({})

  const ammoCache = ref<Record<number, AmmoDetail>>({})
  const ammoTimestamps = ref<Record<number, number>>({})

  const modCache = ref<Record<number, ModDetail>>({})
  const modTimestamps = ref<Record<number, number>>({})
  
  const API_BASE = ''

  async function fetchWeapons() {
    const res = await fetch(`${API_BASE}/api/weapons`)
    weapons.value = await res.json()
  }

  async function fetchAllAmmo() {
    const res = await fetch(`${API_BASE}/api/ammo`)
    ammoList.value = await res.json()
  }

  async function fetchAllMods() {
    const res = await fetch(`${API_BASE}/api/mods`)
    modList.value = await res.json()
  }

  async function fetchAllIngredients() {
    const res = await fetch(`${API_BASE}/api/ingredients`)
    ingredientList.value = await res.json()
  }

  async function fetchWeapon(id: number): Promise<WeaponItem> {
    if (weaponCache.value[id] && Date.now() - weaponTimestamps.value[id] < CACHE_TTL) {
      return weaponCache.value[id]
    }
    const res = await fetch(`${API_BASE}/api/weapons/${id}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    weaponCache.value[id] = data
    weaponTimestamps.value[id] = Date.now()
    return data
  }

  async function fetchAmmo(id: number): Promise<AmmoDetail> {
    if (ammoCache.value[id] && Date.now() - ammoTimestamps.value[id] < CACHE_TTL) {
      return ammoCache.value[id]
    }
    const res = await fetch(`${API_BASE}/api/ammo/${id}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    ammoCache.value[id] = data
    ammoTimestamps.value[id] = Date.now()
    return data
  }

  async function fetchMod(id: number): Promise<ModDetail> {
    if (modCache.value[id] && Date.now() - modTimestamps.value[id] < CACHE_TTL) {
      return modCache.value[id]
    }
    const res = await fetch(`${API_BASE}/api/mods/${id}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    modCache.value[id] = data
    modTimestamps.value[id] = Date.now()
    return data
  }

  const filteredWeapons = computed(() => {
    const s = searchTerm.value.toLowerCase()
    return weapons.value.filter(item =>
      item.name.toLowerCase().includes(s) ||
      item.type.toLowerCase().includes(s)
    )
  })

  function authHeaders(): Record<string, string> {
    const auth = useAuthStore()
    if (auth.token) {
      return { Authorization: `Bearer ${auth.token}` }
    }
    return {}
  }

  async function uploadIcon(file: File): Promise<string> {
    const formData = new FormData()
    formData.append('icon', file)
    const headers = authHeaders()
    const res = await fetch(`${API_BASE}/api/upload`, {
      method: 'POST',
      headers,
      body: formData,
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.error || `HTTP ${res.status}`)
    }
    const result = await res.json()
    return result.icon
  }

  async function createWeapon(data: any): Promise<number> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...authHeaders(),
    }
    const res = await fetch(`${API_BASE}/api/weapons`, {
      method: 'POST',
      headers,
      body: JSON.stringify(data),
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.error || `HTTP ${res.status}`)
    }
    const result = await res.json()
    return result.id
  }

  async function deleteWeapon(id: number): Promise<void> {
    const res = await fetch(`${API_BASE}/api/weapons/${id}`, {
      method: 'DELETE',
      headers: authHeaders(),
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
  }

  const suggestionCache = ref<Suggestion[]>([])

  async function createSuggestion(data: any): Promise<Suggestion> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...authHeaders(),
    }
    const res = await fetch(`${API_BASE}/api/suggestions`, {
      method: 'POST',
      headers,
      body: JSON.stringify(data),
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.error || `HTTP ${res.status}`)
    }
    return await res.json()
  }

  async function fetchSuggestions(): Promise<Suggestion[]> {
    const res = await fetch(`${API_BASE}/api/suggestions`, {
      headers: authHeaders(),
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    suggestionCache.value = data
    return data
  }

  async function fetchSuggestion(id: number): Promise<Suggestion> {
    const res = await fetch(`${API_BASE}/api/suggestions/${id}`, {
      headers: authHeaders(),
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    return await res.json()
  }

  async function approveSuggestion(id: number) {
    const res = await fetch(`${API_BASE}/api/suggestions/${id}/approve`, {
      method: 'PUT',
      headers: authHeaders(),
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.error || `HTTP ${res.status}`)
    }
    return await res.json()
  }

  async function rejectSuggestion(id: number) {
    const res = await fetch(`${API_BASE}/api/suggestions/${id}/reject`, {
      method: 'PUT',
      headers: authHeaders(),
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.error || `HTTP ${res.status}`)
    }
    return await res.json()
  }

  async function deleteSuggestion(id: number) {
    const res = await fetch(`${API_BASE}/api/suggestions/${id}`, {
      method: 'DELETE',
      headers: authHeaders(),
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.error || `HTTP ${res.status}`)
    }
  }

  return { weapons, searchTerm, ammoList, modList, ingredientList, filteredWeapons, fetchWeapons, fetchAllAmmo, fetchAllMods, fetchAllIngredients, fetchWeapon, fetchAmmo, fetchMod, createWeapon, uploadIcon, deleteWeapon, suggestionCache, createSuggestion, fetchSuggestions, fetchSuggestion, approveSuggestion, rejectSuggestion, deleteSuggestion }
})
