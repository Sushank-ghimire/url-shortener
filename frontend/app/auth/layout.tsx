import React from "react";

interface IAuthLayoutProps {
  children: React.ReactNode;
}

export default function AuthLayout({ children }: IAuthLayoutProps) {
  return (
    <main className="flex min-h-screen w-screen overflow-x-hidden justify-center items-center">
      {children}
    </main>
  );
}
