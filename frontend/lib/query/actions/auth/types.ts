export type LoginInput = {
  email: string;
  password: string;
};

export type User = {
  id: string;
  email: string;
  name: string;
  is_verified: boolean;
};

export type LoginResponse = {
  data: {
    access_token: string;
    refresh_token: string;
    user: User;
  };
  success: boolean;
};

export type RegisterInput = {
  email: string;
  password: string;
  name: string;
};

export type RegisterResponse = {
  user: User;
  success: boolean;
};

export type RefreshTokenRequest = {
  refresh_token: string;
};

export type RefreshTokenResponse = {
  access_token: string;
  refresh_token: string;
  success: boolean;
};

export type MailVerificationRequest = {
  token: string;
};

export type MailVerificationResponse = {
  success: boolean;
};

export type GetProfileResponse = {
  message: string;
  success: boolean;
  data: User;
};

export type LogoutResponse = {
  success: boolean;
  message: string;
};
