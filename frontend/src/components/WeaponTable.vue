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

const goToCreate = () => {
  router.push('/weapon/create')
}

const handleDelete = async (id: number) => {
  if (!confirm('Delete this weapon?')) return
  try {
    await store.deleteWeapon(id)
    await store.fetchWeapons()
  } catch {
    alert('Failed to delete weapon')
  }
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
          <th>№</th>
          <th>Name</th>
          <th>Type</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index) in store.filteredWeapons" :key="item.id" class="table-row">
          <td>{{ index + 1 }}</td>
          <td class="weapon-name">{{ item.name }}</td>
          <td><span class="badge">{{ item.type }}</span></td>
          <td class="actions-cell">
            <button class="view-btn" @click="goToDetails(item.id)">View Details</button>
            <button class="delete-btn" @click="handleDelete(item.id)">×</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <div v-if="store.searchTerm && !store.filteredWeapons.length" class="no-results">
    <span>No weapon found for "{{ store.searchTerm }}". </span>
    <button class="create-btn" @click="goToCreate">Create</button>
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

.actions-cell {
  display: flex;
  align-items: center;
  gap: 4px;
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

.delete-btn {
  background: none;
  border: none;
  color: #ef4444;
  font-size: 1.3rem;
  cursor: pointer;
  padding: 4px 8px;
  margin-left: 6px;
  transition: color 0.2s;
  line-height: 1;
}

.delete-btn:hover {
  color: #dc2626;
}

.no-results {
  text-align: center;
  padding: 40px 20px;
  color: #94a3b8;
  font-size: 1.1rem;
}

.create-btn {
  background-color: #ce422b;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  margin-left: 8px;
  transition: background-color 0.2s;
}

.create-btn:hover {
  background-color: #a8321f;
}
</style>
