"use client";

import Spinner from "@/components/ui/spinner";
import { useMailVerification } from "@/lib/query/actions/auth/hooks/use-verification";
import { useSearchParams } from "next/navigation";
import { useEffect } from "react";

export default function VerifyUserEmail() {
  const searchParams = useSearchParams();
  const token = searchParams.get("token");

  const mailVerificationMutation = useMailVerification();

  useEffect(() => {
    if (!token) return;
    mailVerificationMutation.mutate({ token });
  }, [token]);

  if (!token) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-red-500">Invalid verification link.</p>
      </div>
    );
  }

  if (mailVerificationMutation.isPending) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Spinner />
        <p className="ml-2">Verifying your email...</p>
      </div>
    );
  }

  if (mailVerificationMutation.isSuccess) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-green-600">
          Your email has been verified successfully.
        </p>
      </div>
    );
  }

  if (mailVerificationMutation.isError) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-red-500">
          Email verification failed. The link may be invalid or expired.
        </p>
      </div>
    );
  }

  return null;
}
