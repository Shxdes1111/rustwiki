import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

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

export interface WeaponItem {
  id: number;
  name: string;
  type: string;
  description?: string;
  shortname?: string; 
  stacksize: number;
  firemode: string;
  capacity?: number;  
  craftable: boolean;
  timeToCraft?: number; 
  
  ingredients?: Ingredient[];
  ammo?: Ammo[];
  mods?: WeaponMod[];
}

export const useWeaponStore = defineStore('weapons', () => {
  const weapons = ref<WeaponItem[]>([])
  const searchTerm = ref('')

  async function fetchWeapons() {
    const res = await fetch('http://localhost:8080/api/weapons')
    weapons.value = await res.json()
  }

  async function fetchWeapon(id: number): Promise<WeaponItem> {
    const res = await fetch(`http://localhost:8080/api/weapons/${id}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    return await res.json()
  }

  async function fetchAmmo(id: number): Promise<AmmoDetail> {
    const res = await fetch(`http://localhost:8080/api/ammo/${id}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    return await res.json()
  }

  async function fetchMod(id: number): Promise<ModDetail> {
    const res = await fetch(`http://localhost:8080/api/mods/${id}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    return await res.json()
  }

  const filteredWeapons = computed(() => {
    const s = searchTerm.value.toLowerCase()
    return weapons.value.filter(item =>
      item.name.toLowerCase().includes(s) ||
      item.type.toLowerCase().includes(s)
    )
  })

  return { weapons, searchTerm, filteredWeapons, fetchWeapons, fetchWeapon, fetchAmmo, fetchMod }
})
