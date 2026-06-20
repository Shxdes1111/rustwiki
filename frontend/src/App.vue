<script setup lang="ts">
import { ref, onMounted, onUnmounted, provide } from 'vue'
import { useWeaponStore } from './stores/weapons'
import { useAuthStore } from './stores/auth'
import AuthForm from './components/AuthForm.vue'

const store = useWeaponStore()
const authStore = useAuthStore()

const showAuth = ref(false)
const dropdownOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

provide('openAuth', () => { showAuth.value = true })

function toggleDropdown() {
  dropdownOpen.value = !dropdownOpen.value
}

function handleClickOutside(e: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target as Node)) {
    dropdownOpen.value = false
  }
}

function handleLogout() {
  dropdownOpen.value = false
  authStore.logout()
  window.location.reload()
}

onMounted(async () => {
  await authStore.init()
  store.fetchWeapons()
  store.fetchAllAmmo()
  store.fetchAllMods()
  store.fetchAllIngredients()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="wiki-container">
    <div class="header">
      <h1>RustWiki Weapons</h1>
      <div class="header-right">
        <template v-if="authStore.isAuthenticated">
          <div ref="dropdownRef" class="user-menu">
            <button class="user-avatar" @click="toggleDropdown">
              {{ authStore.role.charAt(0).toUpperCase() + authStore.role.slice(1) }}
              <span class="arrow">{{ dropdownOpen ? '▲' : '▼' }}</span>
            </button>
            <div v-if="dropdownOpen" class="dropdown">
              <div class="dropdown-item username">{{ authStore.user?.username }}</div>
              <div class="dropdown-item role">Role: {{ authStore.role }}</div>
              <div class="dropdown-divider"></div>
              <button class="dropdown-item logout-btn" @click="handleLogout">Logout</button>
            </div>
          </div>
        </template>
        <button v-else class="login-btn" @click="showAuth = true">Login</button>
      </div>
    </div>
    <RouterView />
    <AuthForm v-if="showAuth" @close="showAuth = false" />
  </div>
</template>

<style scoped>
.wiki-container {
  padding: 40px;
  max-width: 1000px;
  margin: 0 auto;
  color: white;
}

.header {
  margin-top: 40px;
  margin-bottom: 40px;
}

.header h1 {
  margin: 0;
  text-align: center;
  font-size: 3rem;
  font-weight: 800;
  text-transform: uppercase;
  color: #ef4444;
  text-shadow: 0 2px 10px rgba(239, 68, 68, 0.3);
}

.header-right {
  width: 100%;
  text-align: right;
  margin-top: 12px;
}

.login-btn {
  height: 35px;
  width: 90px;
  margin: 0px 20px;
  text-align: center;
  align-items: center;
  background: #ce422b;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  transition: background 0.2s;
}

.login-btn:hover {
  background: #a8321f;
}

.user-menu {
  position: relative;
  display: inline-block;
}

.user-avatar {
  height: 35px;
  width: 90px;
  margin: 0px 20px;
  border-radius: 8px;
  background: #ce422b;
  color: white;
  border: none;
  font-size: 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2px;
  transition: background 0.2s;
}

.user-avatar:hover {
  background: #a8321f;
}

.arrow {
  font-size: 0.6rem;
}

.dropdown {
  position: absolute;
  right: 0;
  top: 48px;
  text-align: left;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  min-width: 180px;
  z-index: 100;
  overflow: hidden;
}

.dropdown-item {
  padding: 10px 16px;
  font-size: 0.9rem;
}

.dropdown-item.username {
  font-weight: 600;
  color: white;
}

.dropdown-item.role {
  color: #94a3b8;
  font-size: 0.8rem;
}

.dropdown-divider {
  height: 1px;
  background: #333;
  margin: 4px 0;
}

.logout-btn {
  width: 100%;
  text-align: left;
  background: none;
  border: none;
  color: #ef4444;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background 0.2s;
}

.logout-btn:hover {
  background: #222;
}
</style>
