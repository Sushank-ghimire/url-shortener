import { ApiError, type ApiErrorBody } from "./errors";
import secureLocalStorage from "react-secure-storage";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

let refreshPromise: Promise<void> | null = null;

async function refreshAccessToken(): Promise<void> {
  const refreshToken = secureLocalStorage.getItem("refresh_token");

  if (!refreshToken) {
    throw new ApiError(401, {
      message: "No refresh token available",
    });
  }

  if (!refreshPromise) {
    refreshPromise = (async () => {
      try {
        const response = await fetch(`${API_URL}/api/v1/auth/refresh`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            refresh_token: refreshToken,
          }),
        });

        if (!response.ok) {
          throw new ApiError(response.status, await response.json());
        }

        const data = await response.json();
        secureLocalStorage.setItem("access_token", data.access_token);

        if (data.refresh_token) {
          secureLocalStorage.setItem("refresh_token", data.refresh_token);
        }
      } finally {
        refreshPromise = null;
      }
    })();
  }

  return refreshPromise;
}

export async function api<T>(
  path: string,
  options: RequestInit = {},
  retry = true,
): Promise<T> {
  const accessToken = secureLocalStorage.getItem("access_token");

  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
      ...options.headers,
    },
  });

  if (response.ok) {
    return response.json();
  }

  const errorBody: ApiErrorBody = await response.json();

  if (response.status === 401 && retry && path !== "/api/v1/auth/refresh") {
    try {
      await refreshAccessToken();
      return api<T>(path, options, false);
    } catch {
      secureLocalStorage.removeItem("access_token");
      secureLocalStorage.removeItem("refresh_token");

      throw new ApiError(401, errorBody);
    }
  }

  throw new ApiError(response.status, errorBody);
}
