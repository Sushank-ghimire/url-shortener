import { useQuery } from "@tanstack/react-query";
import { getURLs } from "../api";

export function useURLs() {
  return useQuery({
    queryKey: ["urls"],
    queryFn: getURLs,
  });
}
