import request from "@/utils/request";

export function getUserLedger(month) {
  return request.get("/user/ledger", {
    params: {
      month: month
    }
  });
}

export function createUserLedger(payload) {
  return request.post("/user/ledger", payload);
}

export function listUserCategories(type) {
  return request.get("/user/categories", {
    params: {
      type: type
    }
  });
}

export function createUserCategory(payload) {
  return request.post("/user/categories", payload);
}

export function deleteUserCategory(categoryId) {
  return request.delete(`/user/categories/${categoryId}`);
}
