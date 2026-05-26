import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export interface Ingredient {
  id: number;
  name: string;
  amount?: number; // количество для крафта (например, x125)
  icon?: string;   // путь к картинке
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

  const filteredWeapons = computed(() => {
    const s = searchTerm.value.toLowerCase()
    return weapons.value.filter(item =>
      item.name.toLowerCase().includes(s) ||
      item.type.toLowerCase().includes(s)
    )
  })

  return { weapons, searchTerm, filteredWeapons, fetchWeapons }
})
