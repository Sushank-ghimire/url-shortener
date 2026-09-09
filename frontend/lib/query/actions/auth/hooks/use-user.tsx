import { useQuery } from "@tanstack/react-query";
import { getProfile } from "../api";

export function useUser() {
  return useQuery({
    queryKey: ["user", "profile"],
    queryFn: getProfile,
  });
}
