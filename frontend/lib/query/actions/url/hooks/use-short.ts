import { useMutation } from "@tanstack/react-query";
import { createShortURL } from "../api";

export function useShortURL() {
  return useMutation({
    mutationKey: ["short-url"],
    mutationFn: createShortURL,
  });
}
