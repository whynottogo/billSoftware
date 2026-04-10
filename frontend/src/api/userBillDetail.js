import request from "@/utils/request";

export function getUserBillMonthDetail(month) {
  return request.get(`/user/bills/month/${month}`);
}
