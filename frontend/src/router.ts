import { createRouter, createWebHistory } from 'vue-router';
import TranslationWorkspace from './components/TranslationWorkspace.vue';
import Settings from './components/Settings.vue';

const routes = [
  {
    path: '/',
    name: 'Translation',
    component: TranslationWorkspace,
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;