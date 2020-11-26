import { createRouter, createWebHistory } from "vue-router";
import Home from './views/Home.vue';
import Login from './views/Login.vue';
import { hasPermission, isAuthenticated } from './utils';

const routes = [
  {
    path: '/:pathMatch(.*)*',
    redirect: '/home'
  },
  {
    path: "/home",
    name: "Home",
    component: Home,
    meta: { requirePermission: "ADMIN" },
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
    meta: { guest: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  if (to.matched.some((record) => record.meta.requirePermission)) {
    if (hasPermission("ADMIN")) {
      next();
      return;
    }
    next("/login");
  } else {
    next();
  }
});

router.beforeEach((to, from, next) => {
  if (to.params.error === "true") {
    next();
    return;
  }
  if (to.matched.some((record) => record.meta.guest)) {
    if (isAuthenticated()) {
      next("/home");
      return;
    }
    next();
  } else {
    next();
  }
});

export default router;
