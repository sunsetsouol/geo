import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue')
    },
    {
      path: '/prompts',
      name: 'prompts',
      component: () => import('../views/PromptsView.vue')
    },
    {
      path: '/articles',
      name: 'articles',
      component: () => import('../views/ArticlesView.vue')
    },
    {
      path: '/articles/generate',
      name: 'article-generate',
      component: () => import('../views/ArticleGenerateView.vue')
    },
    {
      path: '/articles/edit/:id',
      name: 'article-edit',
      component: () => import('../views/ArticleEditView.vue')
    }
  ]
})

export default router
