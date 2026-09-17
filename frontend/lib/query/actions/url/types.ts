// Redirect Types
export type RedirectResponse = {
  data: {
    long_url: string;
  };
  message: string;
};

export type RedirectRequest = {
  shortcode: string;
};

// Resolve URL Types
export type ResolveRequest = RedirectRequest;

// URL types
export type URL = {
  id: number;
  short_code: string;
  long_url: string;
  title: string;
  custom_alias: string;
  is_active: boolean;
  expires_at: string | null;
  created_at: string;
  updated_at: string;
};

export type ShortURLResponse = {
  success: boolean;
  message: string;
  data: {
    url: URL;
  };
};

export type ShortURLRequest = {
  title: string;
  custom_alias: string;
  long_url: string;
  expires_at?: string | null;
};

export interface ResolveResponse {
  success: boolean;
  data: {
    exists: boolean;
    long_url: string;
  };
  message: string;
}
