import { createRouter, createWebHashHistory   , RouteRecordRaw } from "vue-router";

import Home from "../views/Home.vue";
import Videomanage1 from "../views/management/videomanage1.vue";
import videomanage2 from '../views/management/videomanage2.vue';
import Statistics from "../views/management/statistics.vue";
import Statistics2 from "../views/management/statistics2.vue";
import Management from "../views/management/main.vue";
import managementWelcome from "../views/management/welcome.vue";
const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/ErrorMessage",
    name: "ErrorMessage",
    component: () => import("@/views/ErrorMessage.vue"),
  },
  {
    path: "/upload",
    name: "Upload",
    component: () => import("../views/upload/upload.vue"),
    children: [
      {
        path: "/upload",
        name: "mianUpload",
        component: () => import("../views/upload/mainUpload.vue"),
      },{
        path: "/upload/toBePass",
        name: "toBePass",
        component: () => import("../views/upload/toBePass.vue"),
      },
      {
        path: "/upload/notPass",
        name: "notPass",
        component: () => import("../views/upload/notPass.vue"),
      },
      {
        path: "/upload/pass",
        name: "pass",
        component: () => import("@/views/upload/pass.vue"),
      },
    ],
  },
  {
    path: "/personal",
    name: "PersonalCenter",
    component: () => import("../views/personal/index.vue"),
  },
  {
    path: "/ce",
    name: "ce",
    component: () => import("../views/personal/ceshi.vue"),
  },
  {
    path: "/history",
    name: "history",
    component: () => import("../views/history/index.vue"),
  },
  {
    path: "/search",
    name: "Search",
    component: () => import("../views/search/chaall.vue"),
  },
  {
    path: "/sevideo",
    name: "Sevideo",
    component: () => import("../views/search/chavideo.vue"),
  },
  {
    path: "/seuser",
    name: "Seuser",
    component: () => import("../views/search/chauser.vue"),
  },
  {
    path: "/notice",
    name: "notice",
    component: () => import("../views/notice/index.vue"),
  },
  {
    path: "/video",
    name: "video",
    component: () => import("../views/video/index.vue"),
  },
  {
    path: "/partition",
    name: "partition",
    component: () => import("../views/partition/index.vue"),
  },
  {
    path: "/managementLogin",
    name: "managementLogin",
    component: () => import("../views/management/login.vue"),
  },
  {
    path: "/management",
    name: "management",
    component: Management,
    children: [
      {
        path: "/management",
        name: "managementWelcome",
        component: managementWelcome,
      },
      {
        path: "/management/videomanage1",
        name: "videomanage1",
        component: Videomanage1,
      },
      {
        path: "/management/videomanage2",
        name: "videomanage2",
        component: videomanage2,
      },
      {
        path: "/management/statistics",
        name: "statistics",
        component: Statistics,
      },
      {
        path: "/management/statistics2",
        name: "statistics2",
        component:Statistics2,
      },
    ],
  },
];

const router = createRouter({
  history: createWebHashHistory(process.env.BASE_URL),
  routes,
});

export default router;
