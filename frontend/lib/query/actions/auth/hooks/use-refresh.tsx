"use client";

import { useMutation } from "@tanstack/react-query";
import { verification } from "../api";

export function useRefreshToken() {
  return useMutation({
    mutationKey: ["refresh-token"],
    mutationFn: verification,
  });
}
