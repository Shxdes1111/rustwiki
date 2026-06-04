<script setup lang="ts">
import { ref, watch, nextTick, onMounted } from 'vue'
import { useWeaponStore } from '../stores/weapons'
import { useRouter } from 'vue-router'
import SearchBar from './SearchBar.vue'

const store = useWeaponStore()
const router = useRouter()

// Ссылки на элементы для замера высоты
const tableWrapper = ref<HTMLElement | null>(null)
const innerTable = ref<HTMLElement | null>(null)

// Функция перехода на детали
const goToDetails = (id: number) => {
  router.push(`/weapon/${id}`)
}

// Главная магия: следим за изменением отфильтрованного списка
watch(() => store.filteredWeapons, async () => {
  // Ждем, пока Vue обновит DOM (строки внутри таблицы изменятся)
  await nextTick()
  
  if (tableWrapper.value && innerTable.value) {
    // 1. Измеряем, сколько высоты ТЕПЕРЬ нужно таблице с новыми строками
    const newHeight = innerTable.value.offsetHeight
    
    // 2. Плавно задаем эту высоту внешнему контейнеру-жалюзи
    tableWrapper.value.style.height = `${newHeight}px`
  }
}, { deep: true })

// Задаем начальную высоту при первой загрузке страницы
onMounted(async () => {
  await nextTick()
  if (tableWrapper.value && innerTable.value) {
    tableWrapper.value.style.height = `${innerTable.value.offsetHeight}px`
  }
})
</script>

<template>
  <SearchBar />
  
  <!-- Внешний контейнер-"жалюзи", который будет плавно менять высоту -->
  <div ref="tableWrapper" class="table-blind-container">
    
    <!-- Сама таблица (внутренний контент для замера высоты) -->
    <table ref="innerTable" class="wiki-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Name</th>
          <th>Type</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in store.filteredWeapons" :key="item.id" class="table-row">
          <td>{{ item.id }}</td>
          <td class="weapon-name">{{ item.name }}</td>
          <td><span class="badge">{{ item.type }}</span></td>
          <td>
            <button class="view-btn" @click="goToDetails(item.id)">View Details</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
/* --- НАСТРОЙКА ЭФФЕКТА ЖАЛЮЗИ --- */
.table-blind-container {
  width: 100%;
  overflow: hidden; /* Прячет нижнюю часть таблицы, которая "не влезает" */
  
  /* Плавное изменение высоты (скорость сжатия/расширения) */
  transition: height 0.8s cubic-bezier(0.25, 1, 0.5, 1); 
}

.wiki-table {
  width: 100%;
  border-collapse: collapse;
  background-color: #1a1a1a;
  color: white;
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #333;
}

.table-row {
  /* Плавное исчезновение/появление контента внутри строк, чтобы не было мерцания */
  transition: opacity 1s ease;
}

.table-row:hover {
  background-color: #262626;
}

.badge {
  background-color: #333;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.85rem;
}

.view-btn {
  background-color: #ce422b;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.view-btn:hover {
  background-color: #a8321f;
}
</style>
