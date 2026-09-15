"use client";

import Spinner from "@/components/ui/spinner";

export default function AuthLoading() {
  return (
    <div className="flex min-h-screen w-screen overflow-x-hidden justify-center items-center">
      <div className="flex min-h-screen items-center justify-center">
        <Spinner />
        <p className="ml-2">Checking authentication...</p>
      </div>
    </div>
  );
}
