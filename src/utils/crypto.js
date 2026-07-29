const XOR_KEY = "wol-fnos-seed";

export function xorStr(str) {
  let res = '';
  for (let i = 0; i < str.length; i++) {
    res += String.fromCharCode(str.charCodeAt(i) ^ XOR_KEY.charCodeAt(i % XOR_KEY.length));
  }
  return res;
}

export function encryptCreds(user, pass) {
  const payload = JSON.stringify({ user, pass });
  return btoa(xorStr(payload));
}

export function decryptCreds(encrypted) {
  try {
    const dec = JSON.parse(xorStr(atob(encrypted)));
    return { user: dec.user || '', pass: dec.pass || '' };
  } catch {
    return null;
  }
}