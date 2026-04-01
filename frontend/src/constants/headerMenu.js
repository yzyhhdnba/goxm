export const headerLeftMenu = [
  { text: "首页", url: "/", requiresAuth: false, sag: "" },
  { text: "分区", url: "/partition", requiresAuth: false, sag: "" },
  { text: "搜索", url: "/search", requiresAuth: false, sag: "" },
  { text: "历史", url: "/history", requiresAuth: true, sag: "" },
];

export const headerRightMenu = [
  { text: "消息", url: "/notice", requiresAuth: true, sag: "" },
  { text: "个人中心", url: "/personal", requiresAuth: true, sag: "" },
  { text: "管理后台", url: "/management", requiresAuth: true, sag: "" },
];
