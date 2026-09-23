// A 401 rejects the stored admin session only when the request actually
// presented that session. Login checks pass a candidate token and must not
// wipe a different token that is already stored.
export function shouldClearAdminSession(input: {
  status: number;
  auth: boolean;
  requestToken?: string;
  storedToken: string;
}): boolean {
  if (input.status !== 401 || !input.auth) {
    return false;
  }
  const stored = input.storedToken.trim();
  if (!stored) {
    return false;
  }
  const explicit = input.requestToken?.trim() ?? "";
  return explicit === "" || explicit === stored;
}
