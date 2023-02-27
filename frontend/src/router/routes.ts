import { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: 'index', component: () => import('pages/IndexPage.vue') },
      { path: 'index/art-room', component: () => import('pages/art-live-room.vue') },
      { path: 'index/room', component: () => import('pages/live-room.vue') },
      { path: 'index/focus', component: () => import('pages/focus-list.vue') },
    ],
  },
  { path: '/room', component: () => import('pages/live-room.vue') },

  { path: '/art-room', component: () => import('pages/art-live-room.vue') },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;
