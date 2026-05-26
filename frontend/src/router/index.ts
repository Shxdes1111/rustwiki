import { createRouter, createWebHistory } from 'vue-router'
import WeaponTable from '../components/WeaponTable.vue' // или где у вас главная страница
import WeaponDetails from '../components/WeaponDetails.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: WeaponTable
    },
    {
      path: '/weapon/:id',
      name: 'weapon-details',
      component: WeaponDetails,
      props: true // позволяет передавать id как prop в компонент
    }
  ]
})

export default router