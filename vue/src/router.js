import { createWebHistory, createRouter } from "vue-router";

const routes = [
  {
    path: "/",
    name: "Home",
    component: () => import("./views/Home.vue"),
  },
  {
    path: "/tilmeld",
    name: "tilmeld",
    component: () => import("./views/Deltager.vue"),
  },
  {
    path: "/deltager/:id",
    name: "deltager",
    component: () => import("./views/Deltager.vue"),
  },
  {
    path: "/:catchAll(.*)",
    component: () => import("./views/NotFound.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
