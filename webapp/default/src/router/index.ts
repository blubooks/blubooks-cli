// Composables
import { createRouter/*, createWebHistory */, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: "/",
    name: "home",
    meta: { 
      requiresAuth: true
    },
    component: () => import("../views/Home.vue"),
  }, 
  {
    name: "print",
    path: "/_doc/:pathMatch(.*)",
    component: () => import("../views/Print.vue"),
  },     
  {
    name: "home",
    path: "/:pathMatch(.*)",
    component: () => import("../views/Home.vue"),
  },  
]

const router = createRouter({
  //history: createWebHistory(),
  history: createWebHashHistory(),
  //mode: 'history', 
  routes,
})

export default router
