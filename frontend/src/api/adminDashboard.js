import request from "@/utils/request";

export function getAdminDashboardOverview() {
  return request.get("/admin/dashboard/overview");
}
