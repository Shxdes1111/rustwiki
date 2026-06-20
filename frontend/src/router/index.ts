import { createRouter, createWebHistory } from 'vue-router'
import WeaponTable from '../components/WeaponTable.vue'
import WeaponDetails from '../components/WeaponDetails.vue'
import WeaponCreate from '../components/WeaponCreate.vue'
import AmmoDetails from '../components/AmmoDetails.vue'
import ModDetails from '../components/ModDetails.vue'
import SuggestionList from '../components/SuggestionList.vue'
import SuggestionDetails from '../components/SuggestionDetails.vue'

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
      path: '/weapon/create',
      name: 'weapon-create',
      component: WeaponCreate
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
    },
    {
      path: '/admin/suggestions',
      name: 'suggestion-list',
      component: SuggestionList
    },
    {
      path: '/admin/suggestions/:id',
      name: 'suggestion-details',
      component: SuggestionDetails,
      props: true
    }
  ]
})

export default router