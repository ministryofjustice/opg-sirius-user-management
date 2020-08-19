import { createRouter, createWebHistory } from "vue-router";
import Blank from "../views/Blank.vue";

const routes = [
  {
    path: "/",
    redirect: "/change-password"
  },
  {
    path: "/home",
    name: "Users",
    component: Blank
  },
  {
    path: "/teams",
    name: "Teams",
    component: Blank
  },
  {
    path: "/my-details",
    name: "My details",
    component: Blank
  },
  {
    path: "/change-password",
    name: "Change password",
    component: () =>
      import(
        /* webpackChunkName: "change-password" */ "../views/ChangePassword.vue"
      )
  }
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
});

export default router;
