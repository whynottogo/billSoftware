import axios from "axios";
import { stripAppBase, withAppBase } from "@/utils/appBase";

const request = axios.create({
  baseURL: "/api",
  timeout: 10000
});

request.interceptors.request.use((config) => {
  const url = config.url || "";
  const isAdminRequest = url.startsWith("/admin");
  const tokenKey = isAdminRequest ? "bill_admin_token" : "bill_user_token";
  const token = localStorage.getItem(tokenKey);

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error && error.response && error.response.status === 401) {
      const currentPath = stripAppBase(window.location.pathname);

      if (currentPath.startsWith("/admin")) {
        localStorage.removeItem("bill_admin_token");
        window.location.href = withAppBase("/admin/login");
      } else {
        localStorage.removeItem("bill_user_token");
        localStorage.removeItem("bill_user_profile");
        window.alert("账号已在其他地方登录，请重新登录。");
        window.location.href = withAppBase("/user/login");
      }
    }

    return Promise.reject(error);
  }
);

export default request;
