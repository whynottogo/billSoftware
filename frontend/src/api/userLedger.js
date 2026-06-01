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

export function uploadUserLedgerImage(file) {
  const formData = new FormData();
  formData.append("file", file);

  return request.post("/user/files", formData, {
    headers: {
      "Content-Type": "multipart/form-data"
    }
  });
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
