export type HttpMethod = "GET" | "POST" | "PATCH";

type MockRequestOptions = {
  delayMs?: number;
  signal?: AbortSignal;
};

export async function requestJson<T>(url: string, method: HttpMethod = "GET", body?: unknown): Promise<T> {
  const response = await fetch(url, {
    method,
    headers: { "Content-Type": "application/json" },
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!response.ok) {
    throw new Error(`Request failed: ${response.status}`);
  }

  return (await response.json()) as T;
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
