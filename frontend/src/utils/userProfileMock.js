const PROFILE_STORAGE_KEY = "bill_user_profile_detail";
const PASSWORD_STORAGE_KEY = "bill_user_password_mock";
const SESSION_NOTICE_KEY = "bill_user_session_notice";

const DEFAULT_PASSWORD = "123456";

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

function buildProfileFromSession() {
  const sessionProfile = safeParse(localStorage.getItem("bill_user_profile")) || {};
  const account = sessionProfile.account || sessionProfile.username || "demo_user";

  return {
    account,
    username: account,
    nickname: sessionProfile.nickname || "演示用户",
    phone: sessionProfile.phone || "13800000000",
    email: sessionProfile.email || "demo@example.com",
    avatar:
      sessionProfile.avatar ||
      "https://images.unsplash.com/photo-1557862921-37829c790f19?auto=format&fit=crop&w=240&q=80",
    avatarCompressed: sessionProfile.avatarCompressed || "",
    avatarUpdatedAt: sessionProfile.avatarUpdatedAt || ""
  };
}

function writeSessionProfile(profile) {
  const sessionProfile = {
    account: profile.account,
    username: profile.username,
    nickname: profile.nickname,
    phone: profile.phone,
    email: profile.email,
    avatar: profile.avatar,
    avatarCompressed: profile.avatarCompressed,
    avatarUpdatedAt: profile.avatarUpdatedAt
  };

  localStorage.setItem("bill_user_profile", JSON.stringify(sessionProfile));
}

export function getUserProfileMock() {
  const storedProfile = safeParse(localStorage.getItem(PROFILE_STORAGE_KEY));

  if (storedProfile) {
    return storedProfile;
  }

  const defaultProfile = buildProfileFromSession();
  localStorage.setItem(PROFILE_STORAGE_KEY, JSON.stringify(defaultProfile));
  writeSessionProfile(defaultProfile);
  return defaultProfile;
}

export function saveUserProfileMock(nextProfile) {
  const savedProfile = Object.assign({}, getUserProfileMock(), nextProfile);

  localStorage.setItem(PROFILE_STORAGE_KEY, JSON.stringify(savedProfile));
  writeSessionProfile(savedProfile);
  return savedProfile;
}

export function getUserPasswordMock() {
  const password = localStorage.getItem(PASSWORD_STORAGE_KEY);

  if (password) {
    return password;
  }

  localStorage.setItem(PASSWORD_STORAGE_KEY, DEFAULT_PASSWORD);
  return DEFAULT_PASSWORD;
}

export function updateUserPasswordMock(form) {
  const currentPassword = (form.currentPassword || "").trim();
  const nextPassword = (form.nextPassword || "").trim();
  const confirmPassword = (form.confirmPassword || "").trim();
  const storedPassword = getUserPasswordMock();

  if (!currentPassword || !nextPassword || !confirmPassword) {
    return {
      ok: false,
      message: "请完整填写原密码、新密码和确认密码"
    };
  }

  if (currentPassword !== storedPassword) {
    return {
      ok: false,
      message: "原密码校验失败，请重新输入"
    };
  }

  if (nextPassword.length < 6) {
    return {
      ok: false,
      message: "新密码长度至少为 6 位"
    };
  }

  if (nextPassword !== confirmPassword) {
    return {
      ok: false,
      message: "两次输入的新密码不一致"
    };
  }

  localStorage.setItem(PASSWORD_STORAGE_KEY, nextPassword);
  return {
    ok: true,
    message: "密码修改成功，请重新登录"
  };
}

export function setSessionNotice(notice) {
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

export function consumeSessionNotice() {
  const notice = safeParse(localStorage.getItem(SESSION_NOTICE_KEY));
  localStorage.removeItem(SESSION_NOTICE_KEY);
  return notice;
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

  return setSessionNotice(noticeMap[reason] || noticeMap.logout);
}

export function createAvatarPlaceholder(fileName) {
  if (!fileName) {
    return "未命名头像";
  }

  return fileName.replace(/\.[^.]+$/, "");
}
