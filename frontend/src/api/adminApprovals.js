import request from "@/utils/request";

function toNumberIds(userIds) {
  return (userIds || [])
    .map(function(userId) {
      return Number(userId);
    })
    .filter(function(userId) {
      return Number.isFinite(userId) && userId > 0;
    });
}

export function listPendingApprovalUsers() {
  return request.get("/admin/approvals/users");
}

export function approveApprovalUser(userId) {
  return request.post(`/admin/approvals/users/${userId}/approve`, {});
}

export function rejectApprovalUser(userId, remark) {
  return request.post(`/admin/approvals/users/${userId}/reject`, {
    remark: remark || ""
  });
}

export function batchApproveApprovalUsers(userIds) {
  return request.post("/admin/approvals/users/batch-approve", {
    user_ids: toNumberIds(userIds)
  });
}

export function batchRejectApprovalUsers(userIds, remark) {
  return request.post("/admin/approvals/users/batch-reject", {
    user_ids: toNumberIds(userIds),
    remark: remark || ""
  });
}
