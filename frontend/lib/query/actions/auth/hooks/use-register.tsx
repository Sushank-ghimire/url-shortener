"use client";

import { useMutation } from "@tanstack/react-query";
import { register } from "../api";

export function useRegister() {
  return useMutation({
    mutationKey: ["auth", "register"],
    mutationFn: register,
  });
}
