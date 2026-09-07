import { api } from "@/lib/api/client";

import type {
  GetProfileResponse,
  LoginInput,
  LoginResponse,
  LogoutResponse,
  MailVerificationRequest,
  MailVerificationResponse,
  RefreshTokenRequest,
  RefreshTokenResponse,
  RegisterInput,
  RegisterResponse,
} from "./types";

export async function login(data: LoginInput): Promise<LoginResponse> {
  return api<LoginResponse>("/api/v1/auth/login", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function register(data: RegisterInput): Promise<RegisterResponse> {
  return api<RegisterResponse>("/api/v1/auth/register", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function verification(
  data: MailVerificationRequest,
): Promise<MailVerificationResponse> {
  return api<MailVerificationResponse>("/api/v1/email/verify", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function refresh(
  data: RefreshTokenRequest,
): Promise<RefreshTokenResponse> {
  return api<RefreshTokenResponse>("/api/v1/auth/refresh", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function getProfile(): Promise<GetProfileResponse> {
  const apiRes = api<GetProfileResponse>("/api/v1/auth/me", {
    method: "GET",
  });
  return apiRes;
}

export async function logout(): Promise<LogoutResponse> {
  return api<LogoutResponse>("/api/v1/auth/logout", {
    method: "POST",
  });
}

export async function logoutAll(): Promise<LogoutResponse> {
  return api<LogoutResponse>("/api/v1/auth/logout-all", {
    method: "POST",
  });
}
