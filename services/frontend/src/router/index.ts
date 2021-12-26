import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Home",
    component: import("../views/Home.vue")
  },
  {
    path: "/users",
    name: "Users",
    component: import("../views/Users.vue")
  },
  {
    path: "/jobs",
    name: "Jobs",
    component: import("../views/Jobs.vue")
  }
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
});

export default router;
