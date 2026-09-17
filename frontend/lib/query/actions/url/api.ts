import { api } from "@/lib/api/client";
import type {
  RedirectRequest,
  RedirectResponse,
  ResolveRequest,
  ResolveResponse,
  ShortURLRequest,
  ShortURLResponse,
} from "./types";

export async function redirect(
  data: RedirectRequest,
): Promise<RedirectResponse> {
  return api<RedirectResponse>(`/api/v1/url/${data.shortcode}`, {
    method: "GET",
  });
}

export async function resolveURL(
  data: ResolveRequest,
): Promise<ResolveResponse> {
  return api<ResolveResponse>(`/api/v1/url/resolve/${data.shortcode}`, {
    method: "GET",
  });
}

export async function createShortURL(
  data: ShortURLRequest,
): Promise<ShortURLResponse> {
  return api<ShortURLResponse>(`/api/v1/url`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}
