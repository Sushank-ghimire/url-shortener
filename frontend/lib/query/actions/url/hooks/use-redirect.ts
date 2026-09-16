import { useQuery } from "@tanstack/react-query";
import { redirect } from "../api";

export function useRedirect(shortcode: string) {
  return useQuery({
    queryKey: ["redirect", "url", shortcode],
    queryFn: () => redirect({ shortcode }),
    enabled: Boolean(shortcode),
    retry: false,
  });
}
