import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // Data is considered fresh for 1 minute.
      staleTime: 60 * 1000,
      // Keep unused/inactive queries in the cache for 5 minutes.
      gcTime: 5 * 60 * 1000,

      // Retry failed requests twice.
      retry: 2,

      // gets focus.
      refetchOnWindowFocus: false,

      // Refetch when the component mounts only if the data is stale.
      refetchOnMount: true,

      // Refetch when the browser reconnects.
      refetchOnReconnect: true,
    },

    mutations: {
      retry: 0,
    },
  },
});
