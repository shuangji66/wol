// 自动检测 Base Path
const BASE_PATH = (() => {
  const path = window.location.pathname;
  let base = path.replace(/\/index\.html$/, '').replace(/\/$/, '');
  return base === '' ? '' : base;
})();

export function apiUrl(endpoint) {
  return BASE_PATH + endpoint;
}