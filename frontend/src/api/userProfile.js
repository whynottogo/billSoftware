import request from "@/utils/request";

const SESSION_NOTICE_KEY = "bill_user_session_notice";

function safeParse(rawValue) {
  if (!rawValue) {
    return null;
  }

  try {
    return JSON.parse(rawValue);
  } catch (error) {
    return null;
  }
}

function extractPayload(result) {
  if (result && result.data) {
    return result.data;
  }

  return result || {};
}

function asText(value, fallback) {
  if (value === undefined || value === null) {
    return fallback || "";
  }

  return String(value);
}

function buildSessionProfile(profile) {
  const current = safeParse(localStorage.getItem("bill_user_profile")) || {};

  return {
    account: profile.account || current.account || profile.username || current.username || "",
    username: profile.username || current.username || profile.account || current.account || "",
    nickname: profile.nickname || "",
    phone: profile.phone || "",
    email: profile.email || "",
    avatar: profile.avatar || "",
    avatarCompressed: profile.avatarCompressed || "",
    avatarUpdatedAt: profile.avatarUpdatedAt || ""
  };
}

export function getUserProfile() {
  return request.get("/user/profile");
}

export function updateUserProfile(payload) {
  return request.put("/user/profile", payload);
}

export function updateUserPassword(payload) {
  return request.put("/user/profile/password", payload);
}

export function normalizeUserProfilePayload(result, fallbackProfile) {
  const rawPayload = extractPayload(result);
  const payload = rawPayload && rawPayload.profile ? rawPayload.profile : rawPayload;
  const fallback = fallbackProfile || {};
  const avatarOriginal = asText(
    payload.avatar_original !== undefined ? payload.avatar_original : payload.avatarOriginal,
    fallback.avatarOriginal || ""
  );
  const avatarCompressed = asText(
    payload.avatar_compressed !== undefined ? payload.avatar_compressed : payload.avatarCompressed,
    fallback.avatarCompressed || ""
  );
  const avatar = asText(payload.avatar, avatarCompressed || avatarOriginal || fallback.avatar || "");

  return {
    id: Number(payload.id || fallback.id || 0),
    account: asText(payload.account, payload.username || fallback.account || fallback.username || ""),
    username: asText(payload.username, payload.account || fallback.username || fallback.account || ""),
    nickname: asText(payload.nickname, fallback.nickname || ""),
    phone: asText(payload.phone, fallback.phone || ""),
    email: asText(payload.email, fallback.email || ""),
    avatar: avatar,
    avatarOriginal: avatarOriginal,
    avatarCompressed: avatarCompressed,
    avatarUpdatedAt: asText(
      payload.avatar_updated_at !== undefined ? payload.avatar_updated_at : payload.updated_at,
      fallback.avatarUpdatedAt || ""
    ),
    updatedAt: asText(payload.updated_at, fallback.updatedAt || "")
  };
}

export function persistUserSessionProfile(profile) {
  const sessionProfile = buildSessionProfile(profile);
  localStorage.setItem("bill_user_profile", JSON.stringify(sessionProfile));
  return sessionProfile;
}

export function setUserSessionNotice(notice) {
  const payload = {
    type: notice.type || "warning",
    title: notice.title || "会话提示",
    text: notice.text || "",
    reason: notice.reason || "",
    createdAt: new Date().toISOString()
  };

  localStorage.setItem(SESSION_NOTICE_KEY, JSON.stringify(payload));
  return payload;
}

export function clearUserSession(reason) {
  localStorage.removeItem("bill_user_token");
  localStorage.removeItem("bill_user_profile");

  const noticeMap = {
    logout: {
      type: "success",
      title: "已退出登录",
      text: "你的账号已安全退出，如需继续使用请重新登录。",
      reason: "logout"
    },
    kicked: {
      type: "warning",
      title: "会话已失效",
      text: "账号已在其他地方登录，请重新登录。",
      reason: "kicked"
    },
    password_reset: {
      type: "warning",
      title: "请重新登录",
      text: "密码已修改，为保护账号安全请重新登录。",
      reason: "password_reset"
    }
  };

  return setUserSessionNotice(noticeMap[reason] || noticeMap.logout);
}

export function buildUserProfileError(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }

  return fallback;
}
