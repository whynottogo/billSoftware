import request from "@/utils/request";

export function getAdminUserBillsOverview(userId) {
  return request.get(`/admin/users/${userId}/bills/overview`);
}
