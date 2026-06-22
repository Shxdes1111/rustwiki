import { createRouter, createWebHistory } from 'vue-router'
import WeaponTable from '../components/WeaponTable.vue'
import WeaponDetails from '../components/WeaponDetails.vue'
import WeaponCreate from '../components/WeaponCreate.vue'
import AmmoDetails from '../components/AmmoDetails.vue'
import ModDetails from '../components/ModDetails.vue'
import SuggestionList from '../components/SuggestionList.vue'
import SuggestionDetails from '../components/SuggestionDetails.vue'
import MySuggestions from '../components/MySuggestions.vue'
import MySuggestionDetails from '../components/MySuggestionDetails.vue'
import MyWeapons from '../components/MyWeapons.vue'

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
    },
    {
      path: '/my/suggestions',
      name: 'my-suggestions',
      component: MySuggestions
    },
    {
      path: '/my/suggestions/:id',
      name: 'my-suggestion-details',
      component: MySuggestionDetails,
      props: true
    },
    {
      path: '/my/weapons',
      name: 'my-weapons',
      component: MyWeapons
    }
  ]
})

export default router