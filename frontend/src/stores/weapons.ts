import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export interface WeaponItem {
  id: number
  name: string
  type: string
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
