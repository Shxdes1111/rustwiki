import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

const TOKEN_KEY = 'rustwiki_token'
const API_BASE = import.meta.env.VITE_API_BASE || ''

export interface User {
  id: number
  username: string
  role: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY))
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const role = computed(() => user.value?.role ?? 'guest')
  const isAdmin = computed(() => role.value === 'admin')

  async function init() {
    const saved = localStorage.getItem(TOKEN_KEY)
    if (!saved) return

    token.value = saved
    try {
      const res = await fetch(`${API_BASE}/api/users/me`, {
        headers: { Authorization: `Bearer ${saved}` },
      })
      if (!res.ok) throw new Error('Token invalid')
      user.value = await res.json()
    } catch {
      token.value = null
      user.value = null
      localStorage.removeItem(TOKEN_KEY)
    }
  }

  async function login(username: string, password: string) {
    const res = await fetch(`${API_BASE}/api/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.error || 'Login failed')
    }
    const data = await res.json()
    token.value = data.token
    user.value = { id: data.user_id, username: data.username, role: data.role }
    localStorage.setItem(TOKEN_KEY, data.token)
  }

  async function register(username: string, password: string) {
    const res = await fetch(`${API_BASE}/api/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.error || 'Registration failed')
    }
    const data = await res.json()
    token.value = data.token
    user.value = { id: data.user_id, username: data.username, role: data.role }
    localStorage.setItem(TOKEN_KEY, data.token)
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
  }

  return { token, user, isAuthenticated, role, isAdmin, init, login, register, logout }
})
