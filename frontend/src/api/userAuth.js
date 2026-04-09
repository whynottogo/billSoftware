import request from "@/utils/request";

export function userLogin(payload) {
  return request.post("/user/auth/login", payload);
}

export function userRegister(payload) {
  return request.post("/user/auth/register", payload);
}

export function adminLogin(payload) {
  return request.post("/admin/auth/login", payload);
}

