"use client";

import { queryClient } from "@/lib/query/client";
import { QueryClientProvider } from "@tanstack/react-query";
import React from "react";

interface ProviderProps {
  children: React.ReactNode;
}

export default function Providers({ children }: ProviderProps) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
