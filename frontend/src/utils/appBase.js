const rawAppBase = typeof __APP_BASE_PATH__ === "string" ? __APP_BASE_PATH__ : "/";

function normalizeBase(base) {
  if (!base || base === "/") {
    return "/";
  }

  const trimmed = base.replace(/\/+$/, "");
  return trimmed.startsWith("/") ? `${trimmed}/` : `/${trimmed}/`;
}

const appBase = normalizeBase(rawAppBase);
const appBasePrefix = appBase === "/" ? "" : appBase.replace(/\/$/, "");

export function getAppBase() {
  return appBase;
}

export function withAppBase(path) {
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;

  if (!appBasePrefix) {
    return normalizedPath;
  }

  return `${appBasePrefix}${normalizedPath}`;
}

export function stripAppBase(pathname) {
  if (!pathname) {
    return "/";
  }

  if (!appBasePrefix) {
    return pathname;
  }

  if (pathname === appBasePrefix) {
    return "/";
  }

  if (pathname.startsWith(`${appBasePrefix}/`)) {
    return pathname.slice(appBasePrefix.length);
  }

  return pathname;
}
