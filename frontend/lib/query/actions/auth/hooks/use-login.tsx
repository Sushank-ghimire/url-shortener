"use client";

import { useMutation } from "@tanstack/react-query";
import { login } from "../api";

export function useLogin() {
  return useMutation({
    mutationKey: ["auth", "login"],
    mutationFn: login,
  });
}
