import request from "@/utils/request";

export function listAdminUsers() {
  return request.get("/admin/users");
}

export function changeAdminUserStatus(userId, status) {
  return request.put(`/admin/users/${userId}/status`, {
    status: status
  });
}
