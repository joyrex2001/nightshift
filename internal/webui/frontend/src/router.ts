import Vue from 'vue';
import Router from 'vue-router';
import ScannersOverview from './views/ScannersOverview.vue';

Vue.use(Router);

export default new Router({
  mode: 'hash',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      redirect: '/scanners',
    },
    {
      path: '/scanners',
      name: 'scanners',
      component: ScannersOverview,
    },
    {
      path: '/objects',
      name: 'objects',
      component: () => import('./views/ObjectsOverview.vue'),
    },
    {
      path: '/triggers',
      name: 'triggers',
      component: () => import('./views/TriggersOverview.vue'),
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "about" */ './views/About.vue'),
    },
  ],
});
