import { ApiError, type ApiErrorBody } from "./errors";
import secureLocalStorage from "react-secure-storage";

export const API_URL = process.env.NEXT_PUBLIC_API_URL;

let refreshPromise: Promise<void> | null = null;

async function refreshAccessToken(): Promise<void> {
  const refreshToken = secureLocalStorage.getItem("refresh_token");

  if (!refreshToken || typeof refreshToken !== "string") {
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

        const body = await response.json();

        console.log("REFRESH RESPONSE:", body);

        if (!response.ok) {
          throw new ApiError(response.status, body);
        }

        const accessToken = body.data?.access_token;
        const newRefreshToken = body.data?.refresh_token;

        if (!accessToken) {
          throw new ApiError(401, {
            message: "Refresh response does not contain access token",
          });
        }

        secureLocalStorage.setItem("access_token", accessToken);

        if (newRefreshToken) {
          secureLocalStorage.setItem("refresh_token", newRefreshToken);
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
      ...(accessToken
        ? {
            Authorization: `Bearer ${accessToken}`,
          }
        : {}),
      ...options.headers,
    },
  });

  if (response.ok) {
    return response.json();
  }

  const errorBody: ApiErrorBody = await response.json();

  // Access token expired
  if (response.status === 401 && retry) {
    try {
      await refreshAccessToken();

      return api<T>(path, options, false);
    } catch (error) {
      secureLocalStorage.removeItem("access_token");
      secureLocalStorage.removeItem("refresh_token");

      throw error instanceof ApiError
        ? error
        : new ApiError(401, { message: "Authentication failed" });
    }
  }

  throw new ApiError(response.status, errorBody);
}
