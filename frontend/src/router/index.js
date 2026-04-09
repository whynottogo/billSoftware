import { createRouter, createWebHistory } from "vue-router";

import UserLayout from "@/layouts/UserLayout.vue";
import AdminLayout from "@/layouts/AdminLayout.vue";
import UserLogin from "@/pages/user/UserLogin.vue";
import UserRegister from "@/pages/user/UserRegister.vue";
import UserHome from "@/pages/user/UserHome.vue";
import UserBillsMonth from "@/pages/user/UserBillsMonth.vue";
import UserBillDetail from "@/pages/user/UserBillDetail.vue";
import UserBillsYear from "@/pages/user/UserBillsYear.vue";
import UserBudgetMonth from "@/pages/user/UserBudgetMonth.vue";
import UserBudgetYear from "@/pages/user/UserBudgetYear.vue";
import UserAssets from "@/pages/user/UserAssets.vue";
import UserAssetDetail from "@/pages/user/UserAssetDetail.vue";
import UserCharts from "@/pages/user/UserCharts.vue";
import UserFamilies from "@/pages/user/UserFamilies.vue";
import UserFamilyDetail from "@/pages/user/UserFamilyDetail.vue";
import UserProfile from "@/pages/user/UserProfile.vue";
import UserProfilePassword from "@/pages/user/UserProfilePassword.vue";
import UserSessionKickout from "@/pages/user/UserSessionKickout.vue";
import AdminLogin from "@/pages/admin/AdminLogin.vue";
import AdminUsers from "@/pages/admin/AdminUsers.vue";
import AdminUserDetail from "@/pages/admin/AdminUserDetail.vue";
import AdminDashboard from "@/pages/admin/AdminDashboard.vue";
import AdminApprovals from "@/pages/admin/AdminApprovals.vue";
import AdminFamilies from "@/pages/admin/AdminFamilies.vue";

const routes = [
  {
    path: "/",
    redirect: "/user/login"
  },
  {
    path: "/user/login",
    component: UserLogin,
    meta: {
      public: true,
      side: "user"
    }
  },
  {
    path: "/user/register",
    component: UserRegister,
    meta: {
      public: true,
      side: "user"
    }
  },
  {
    path: "/user/session-kickout",
    component: UserSessionKickout,
    meta: {
      public: true,
      side: "user"
    }
  },
  {
    path: "/user",
    component: UserLayout,
    meta: {
      requiresAuth: true,
      side: "user"
    },
    children: [
      {
        path: "",
        redirect: "/user/ledger"
      },
      {
        path: "ledger",
        component: UserHome
      },
      {
        path: "bills/month",
        component: UserBillsMonth
      },
      {
        path: "bills/month/:month",
        component: UserBillDetail
      },
      {
        path: "bills/year",
        component: UserBillsYear
      },
      {
        path: "budget/month",
        component: UserBudgetMonth
      },
      {
        path: "budget/year",
        component: UserBudgetYear
      },
      {
        path: "assets",
        component: UserAssets
      },
      {
        path: "assets/:accountId",
        component: UserAssetDetail
      },
      {
        path: "charts/expense",
        component: UserCharts
      },
      {
        path: "charts/income",
        component: UserCharts
      },
      {
        path: "families",
        component: UserFamilies
      },
      {
        path: "families/:familyId",
        component: UserFamilyDetail
      },
      {
        path: "profile",
        component: UserProfile
      },
      {
        path: "profile/password",
        component: UserProfilePassword
      }
    ]
  },
  {
    path: "/admin/login",
    component: AdminLogin,
    meta: {
      public: true,
      side: "admin"
    }
  },
  {
    path: "/admin",
    component: AdminLayout,
    meta: {
      requiresAuth: true,
      side: "admin"
    },
    children: [
      {
        path: "",
        redirect: "/admin/dashboard"
      },
      {
        path: "users",
        component: AdminUsers
      },
      {
        path: "users/:userId",
        component: AdminUserDetail
      },
      {
        path: "dashboard",
        component: AdminDashboard
      },
      {
        path: "approvals",
        component: AdminApprovals
      },
      {
        path: "families",
        component: AdminFamilies
      }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const isAdminSide = to.path.startsWith("/admin");
  const tokenKey = isAdminSide ? "bill_admin_token" : "bill_user_token";
  const token = localStorage.getItem(tokenKey);

  if (to.meta.public && token) {
    next(isAdminSide ? "/admin/dashboard" : "/user/ledger");
    return;
  }

  if (to.meta.requiresAuth && !token) {
    next(isAdminSide ? "/admin/login" : "/user/login");
    return;
  }

  next();
});

export default router;
