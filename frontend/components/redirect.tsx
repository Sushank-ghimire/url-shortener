"use client";

import { useEffect, useState } from "react";
import Spinner from "./ui/spinner";
import { resolveURL } from "@/lib/query/actions/url/api";

interface IRedirectParams {
  shortcode: string;
}

export default function RedirectPage({ shortcode }: IRedirectParams) {
  const [error, setError] = useState(false);

  useEffect(() => {
    if (!shortcode) {
      setError(true);
      return;
    }

    const checkAndRedirect = async () => {
      try {
        const data = await resolveURL({ shortcode });
        if (!data.success || !data.data.exists) {
          setError(true);
          return;
        }
        window.location.href = data.data.long_url;
      } catch (error) {
        setError(true);
      }
    };

    checkAndRedirect();
  }, [shortcode]);

  if (error) {
    return (
      <main className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <h1 className="text-xl font-semibold">URL Not Found</h1>
          <p className="mt-2 text-gray-500">
            The short URL you entered doesn't exist or has expired.
          </p>
        </div>
      </main>
    );
  }

  return (
    <main className="flex min-h-screen items-center justify-center">
      <div className="flex items-center">
        <Spinner />
        <p className="ml-2">Redirecting user...</p>
      </div>
    </main>
  );
}
