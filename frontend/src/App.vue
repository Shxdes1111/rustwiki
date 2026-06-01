<script setup lang="ts">
import { onMounted } from 'vue'
import { useWeaponStore } from './stores/weapons'

const store = useWeaponStore()

onMounted(() => store.fetchWeapons())
</script>

<template>
  <div class="wiki-container">
    <h1>RustWiki Weapons</h1>
    <RouterView />
  </div>
</template>

<style scoped>
.wiki-container {
  padding: 40px;
  max-width: 1000px;
  margin: 0 auto;
  color: white;
}

.wiki-container h1 {
  /* 1. Положение */
  text-align: center; 
  margin-top: 40px;     /* center — по центру, left — влево, right — вправо */
  margin-bottom: 80px;     /* Отступ снизу (до таблицы или карточки) */

  /* 2. Размер */
  font-size: 3rem;         /* Размер шрифта (можно в rem, em или px, например: 42px) */
  font-weight: 800;        /* Жирность шрифта (от 100 до 900) */
  text-transform: uppercase; /* Сделать все буквы заглавными (по желанию) */

  /* 3. Цвет */
  color: #ef4444;          /* Красный цвет под стиль вашей карточки (цвета Rust) */
  
  /* Бонус: красивое свечение текста */
  text-shadow: 0 2px 10px rgba(239, 68, 68, 0.3); 
}
</style>
