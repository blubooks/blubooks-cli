// Composables
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: "/",
    meta: { 
      requiresAuth: true
    },
    component: () => import("../views/Home.vue"),
  },
  {
    name: "BadNotFound",
    path: "/:pathMatch(.*)",
    component: () => import("../views/Home.vue"),
  }    

]

const router = createRouter({
  history: createWebHistory("/public/"),
  //mode: 'history', 
  routes,
})

export default router
