<script setup lang="ts">
import { ref } from 'vue'
import { useWeaponStore } from '../stores/weapons'

const store = useWeaponStore()

const placeholderText = ref('🔍 Search weapons...')

const inputRef = ref<HTMLInputElement | null>(null)

const handleFocus = () => {
  placeholderText.value = ''
}

const handleBlur = () => {
  placeholderText.value = '🔍 Search weapons...'
}

const clearSearch = () => {
  store.searchTerm = ''
  inputRef.value?.focus()
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    clearSearch()
  }
}
</script>

<template>
  <div class="search-section">
    <div class="search-wrapper">
      <input
        ref="inputRef"
        v-model="store.searchTerm"
        :placeholder="placeholderText"
        @focus="handleFocus"
        @blur="handleBlur"
        @keydown="handleKeydown"
        class="search-input"
      />
      <button
        v-if="store.searchTerm"
        class="clear-btn"
        @click="clearSearch"
        aria-label="Clear search"
      >
        ×
      </button>
    </div>
  </div>
</template>

<style scoped>
.search-section {
  margin-bottom: 20px;
}

.search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-input {
  width: 100%;
  padding: 12px 36px 12px 12px;
  border: 1px solid #333;
  border-radius: 8px;
  background-color: #1a1a1a;
  color: white;
  font-size: 1rem;
  outline: none;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.search-input:focus {
  border-color: #ef4444;
}

.search-input::placeholder {
  color: #666;
}

.clear-btn {
  position: absolute;
  right: 8px;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  border-radius: 50%;
  color: #666;
  font-size: 1.2rem;
  cursor: pointer;
  transition: color 0.15s, background 0.15s;
  line-height: 1;
}

.clear-btn:hover {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
}
</style>
