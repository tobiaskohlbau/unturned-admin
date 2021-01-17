import { createRouter, createWebHistory } from "vue-router";
import Home from './views/Home.vue';
import Login from './views/Login.vue';
import { hasPermission, isActivated, isAuthenticated } from './utils';

const routes = [
  {
    path: '/:pathMatch(.*)*',
    redirect: '/home'
  },
  {
    path: "/home",
    name: "Home",
    component: Home,
    meta: { requireAuthenticated: true },
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
    if (!isActivated()) {
      next("/login");
      return;
    }
    if (!hasPermission(to.meta.requirePermission)) {
      next("/login");
      return;
    }
  } else if (to.matched.some((record) => record.meta.requireAuthenticated)) {
    if (!isActivated()) {
      next("/login");
      return;
    }
  }
  next();
});

router.beforeEach((to, from, next) => {
  if (to.params.error === "true") {
    next();
    return;
  }
  if (to.matched.some((record) => record.meta.guest)) {
    if (isAuthenticated() && isActivated()) {
      next("/home");
      return;
    }
    next();
  } else {
    next();
  }
});

export default router;
