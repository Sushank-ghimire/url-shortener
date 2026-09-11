"use client";

import { useMutation } from "@tanstack/react-query";
import { verification } from "../api";

export function useMailVerification() {
  return useMutation({
    mutationKey: ["email-verification"],
    mutationFn: verification,
  });
}
