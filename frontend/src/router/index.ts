import { createRouter, createWebHistory } from 'vue-router'
import WeaponTable from '../components/WeaponTable.vue'
import WeaponDetails from '../components/WeaponDetails.vue'
import AmmoDetails from '../components/AmmoDetails.vue'
import ModDetails from '../components/ModDetails.vue'

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
      props: true
    },
    {
      path: '/ammo/:id',
      name: 'ammo-details',
      component: AmmoDetails,
      props: true
    },
    {
      path: '/mods/:id',
      name: 'mod-details',
      component: ModDetails,
      props: true
    }
  ]
})

export default router