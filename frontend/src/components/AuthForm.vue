<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const emit = defineEmits<{ close: [] }>()

const authStore = useAuthStore()

const tab = ref<'login' | 'register'>('login')

const username = ref('')
const password = ref('')
const showPassword = ref(false)
const error = ref('')
const loading = ref(false)

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    if (tab.value === 'login') {
      await authStore.login(username.value, password.value)
    } else {
      await authStore.register(username.value, password.value)
    }
    emit('close')
  } catch (e: any) {
    error.value = e.message || 'Something went wrong'
  } finally {
    loading.value = false
  }
}

function switchTab(t: 'login' | 'register') {
  tab.value = t
  error.value = ''
}
</script>

<template>
  <div class="auth-overlay" @click.self="emit('close')">
    <div class="auth-modal">
      <button class="close-btn" @click="emit('close')">&times;</button>

      <div class="tabs">
        <button
          :class="['tab', { active: tab === 'login' }]"
          @click="switchTab('login')"
        >
          Login
        </button>
        <button
          :class="['tab', { active: tab === 'register' }]"
          @click="switchTab('register')"
        >
          Register
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="auth-form">
        <div class="field">
          <label for="username">Username</label>
          <input
            id="username"
            v-model="username"
            type="text"
            placeholder="Username"
            required
            minlength="3"
            maxlength="50"
          />
        </div>

        <div class="field">
          <label for="password">Password</label>
          <div class="password-wrapper">
            <input
              id="password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Password"
              required
              minlength="6"
              maxlength="50"
            />
            <button type="button" class="toggle-password" @click="showPassword = !showPassword" :title="showPassword ? 'Hide password' : 'Show password'">
              {{ showPassword ? '🙈' : '👁' }}
            </button>
          </div>
        </div>

        <p v-if="error" class="error">{{ error }}</p>

        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? 'Please wait...' : tab === 'login' ? 'Sign in' : 'Sign up' }}
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.auth-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.auth-modal {
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 12px;
  padding: 32px;
  width: 400px;
  max-width: 90vw;
  position: relative;
}

.close-btn {
  position: absolute;
  top: 12px;
  right: 16px;
  background: none;
  border: none;
  color: #666;
  font-size: 1.5rem;
  cursor: pointer;
}

.close-btn:hover {
  color: white;
}

.tabs {
  display: flex;
  gap: 0;
  margin-bottom: 24px;
  border-bottom: 1px solid #333;
}

.tab {
  flex: 1;
  padding: 10px;
  background: none;
  border: none;
  color: #666;
  font-size: 1rem;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: color 0.2s, border-color 0.2s;
}

.tab.active {
  color: #ef4444;
  border-bottom-color: #ef4444;
}

.tab:hover {
  color: #ccc;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field label {
  color: #94a3b8;
  font-size: 0.9rem;
}

.field input {
  padding: 10px 12px;
  border: 1px solid #333;
  border-radius: 6px;
  background: #121212;
  color: white;
  font-size: 1rem;
  outline: none;
}

.field input:focus {
  border-color: #ef4444;
}

.password-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.password-wrapper input {
  width: 100%;
  padding-right: 40px;
}

.toggle-password {
  position: absolute;
  right: 8px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 4px;
  line-height: 1;
  color: #666;
  transition: color 0.2s;
}

.toggle-password:hover {
  color: white;
}

.error {
  color: #ef4444;
  font-size: 0.9rem;
  margin: 0;
}

.submit-btn {
  padding: 12px;
  background: #ce422b;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.submit-btn:hover {
  background: #a8321f;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
