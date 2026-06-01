<script setup lang="ts">
import { ref } from 'vue'
import { useWeaponStore } from '../stores/weapons'

const store = useWeaponStore()

// 1. Создаем переменную для динамического placeholder
const placeholderText = ref('🔍 Search weapons...')

// 2. Функция, когда пользователь кликнул на инпут (Фокус)
const handleFocus = () => {
  placeholderText.value = ''
}

// 3. Функция, когда пользователь убрал клик с инпута (Блур)
const handleBlur = () => {
  placeholderText.value = '🔍 Search weapons...'
}
</script>

<template>
  <div class="search-section">
    <!-- Связываем placeholder с переменной и вешаем слушатели событий -->
    <input
      v-model="store.searchTerm"
      :placeholder="placeholderText"
      @focus="handleFocus"
      @blur="handleBlur"
      class="search-input"
    />
  </div>
</template>

<style scoped>
.search-section {
  margin-bottom: 20px;
}

.search-input {
  width: 100%;
  padding: 12px;
  border: 1px solid #333;
  border-radius: 8px;
  background-color: #1a1a1a;
  color: white;
  font-size: 1rem;
  outline: none;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

/* Изменяем цвет рамки при фокусе (если нужно) */
.search-input:focus {
  border-color: #ef4444; /* Цвет Vue (зеленый) */
}

.search-input::placeholder {
  color: #666;
}
</style>
