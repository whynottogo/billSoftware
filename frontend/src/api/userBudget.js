import request from "@/utils/request";

export function getCurrentMonthBudget() {
  return request.get("/user/budgets/month/current");
}

export function updateCurrentMonthBudget(payload) {
  return request.put("/user/budgets/month/current", payload);
}

export function getYearBudgetOptions() {
  return request.get("/user/budgets/year/options");
}

export function getYearBudget(year) {
  return request.get(`/user/budgets/year/${year}`);
}

export function updateCurrentYearBudget(payload) {
  return request.put("/user/budgets/year/current", payload);
}
