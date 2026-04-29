export type HttpMethod = "GET" | "POST" | "PATCH";

type RequestOptions = {
  signal?: AbortSignal;
  token?: string;
};

type MockRequestOptions = {
  delayMs?: number;
  signal?: AbortSignal;
};

type ApiErrorBody = {
  error?: string;
};

const DEFAULT_NODE_BASE_URL = "http://localhost:8080";
const isBrowser = typeof window !== "undefined";
const configuredBase = import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, "") ?? "";
const API_BASE_URL = configuredBase || (isBrowser ? "" : DEFAULT_NODE_BASE_URL);

export class ApiError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.status = status;
    this.name = "ApiError";
  }
}

function buildUrl(path: string): string {
  if (!path.startsWith("/")) {
    throw new Error(`API path must start with '/': ${path}`);
  }
  return `${API_BASE_URL}${path}`;
}

async function parseErrorMessage(response: Response): Promise<string> {
  try {
    const body = (await response.json()) as ApiErrorBody;
    if (body.error) return body.error;
  } catch {
    // ignore parse errors and use fallback
  }
  return `Request failed: ${response.status}`;
}

async function request(path: string, method: HttpMethod, body?: unknown, options: RequestOptions = {}): Promise<Response> {
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (options.token) {
    headers.Authorization = `Bearer ${options.token}`;
  }

  const response = await fetch(buildUrl(path), {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
    signal: options.signal,
  });

  if (!response.ok) {
    throw new ApiError(response.status, await parseErrorMessage(response));
  }

  return response;
}

export async function requestJson<T>(path: string, method: HttpMethod = "GET", body?: unknown, options: RequestOptions = {}): Promise<T> {
  const response = await request(path, method, body, options);
  if (response.status === 204) {
    return undefined as T;
  }
  return (await response.json()) as T;
}

export async function requestVoid(path: string, method: HttpMethod, body?: unknown, options: RequestOptions = {}): Promise<void> {
  await request(path, method, body, options);
}

export async function mockRequestJson<T>(factory: () => T, options: MockRequestOptions = {}): Promise<T> {
  const { delayMs = 120, signal } = options;

  return new Promise<T>((resolve, reject) => {
    const timeout = window.setTimeout(() => {
      if (signal?.aborted) {
        reject(new DOMException("Request aborted", "AbortError"));
        return;
      }
      try {
        resolve(factory());
      } catch (error) {
        reject(error);
      }
    }, delayMs);

    if (signal) {
      signal.addEventListener(
        "abort",
        () => {
          window.clearTimeout(timeout);
          reject(new DOMException("Request aborted", "AbortError"));
        },
        { once: true }
      );
    }
  });
}
