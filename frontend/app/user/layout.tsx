"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/lib/query/actions/auth/hooks/use-user";
import { useAuthStore } from "@/store/auth/store";
import Spinner from "@/components/ui/spinner";

export interface IChildrenProps {
  children: React.ReactNode;
}

export default function DashboardAuthLayout({ children }: IChildrenProps) {
  const router = useRouter();
  const { setUser } = useAuthStore();

  const { data, isPending, isError } = useUser();
  useEffect(() => {
    if (isPending) return;

    if (data?.success || data?.data) {
      setUser({
        email: data.data.email,
        id: data.data.id,
        is_verified: data.data.is_verified,
        name: data.data.name,
      });
      router.replace("/user/dashboard");
    } else if (isError) {
      router.replace("/auth/login");
    }
  }, [data, isPending, isError, router]);

  if (isPending) {
    return (
      <div className="flex min-h-screen w-screen overflow-x-hidden justify-center items-center">
        <div className="flex min-h-screen items-center justify-center">
          <Spinner />
          <p className="ml-2">Checking authentication...</p>
        </div>
      </div>
    );
  }

  if (!data?.success || !data.data) {
    return null;
  }

  return children;
}
