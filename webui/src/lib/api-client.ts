import { getStoredAuthToken, useAuthStore } from "../features/auth/auth-store";
import { shouldClearAdminSession } from "./admin-session";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL?.trim() ?? "";

type Primitive = string | number | boolean | null;
type JsonValue = Primitive | JsonValue[] | { [key: string]: JsonValue };

export type ApiErrorBody = {
  error?: {
    code?: string;
    message?: string;
  };
};

export class ApiError extends Error {
  readonly status: number;
  readonly code: string;
  readonly body: ApiErrorBody | null;

  constructor(status: number, code: string, message: string, body: ApiErrorBody | null) {
    super(message);
    this.status = status;
    this.code = code;
    this.body = body;
  }
}

type RequestOptions = {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  body?: JsonValue;
  auth?: boolean;
  token?: string;
  signal?: AbortSignal;
};

// Browsers reject fetch() when the resolved URL keeps username or password
// copied from the page address, for example http://user:token@host/ui/.
export function apiURL(path: string): string {
  const raw = path.startsWith("http://") || path.startsWith("https://") ? path : `${API_BASE_URL}${path}`;
  const base = typeof window === "undefined" ? "http://127.0.0.1" : window.location.origin;
  const url = new URL(raw, base);
  url.username = "";
  url.password = "";
  return url.toString();
}

async function parseErrorBody(response: Response): Promise<ApiErrorBody | null> {
  const contentType = response.headers.get("content-type") ?? "";
  if (!contentType.includes("application/json")) {
    return null;
  }

  try {
    return (await response.json()) as ApiErrorBody;
  } catch {
    return null;
  }
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = "GET", body, auth = true, token, signal } = options;
  const headers = new Headers();

  if (body !== undefined) {
    headers.set("Content-Type", "application/json; charset=utf-8");
  }

  if (auth) {
    const resolvedToken = token?.trim() || getStoredAuthToken();
    if (resolvedToken) {
      headers.set("Authorization", `Bearer ${resolvedToken}`);
    }
  }

  const response = await fetch(apiURL(path), {
    method,
    headers,
    body: body === undefined ? undefined : JSON.stringify(body),
    signal,
  });

  if (!response.ok) {
    const parsed = await parseErrorBody(response);
    const code = parsed?.error?.code ?? "HTTP_ERROR";
    const message = parsed?.error?.message ?? response.statusText;
    if (
      shouldClearAdminSession({
        status: response.status,
        auth,
        requestToken: token,
        storedToken: getStoredAuthToken(),
      })
    ) {
      useAuthStore.getState().clearToken();
    }
    throw new ApiError(response.status, code, message, parsed);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  const contentType = response.headers.get("content-type") ?? "";
  if (!contentType.includes("application/json")) {
    return undefined as T;
  }

  return (await response.json()) as T;
}
