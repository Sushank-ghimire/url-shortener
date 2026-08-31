export type ApiErrorBody = {
  message?: string;
  code?: string;
  details?: unknown;
};

export class ApiError extends Error {
  readonly status: number;
  readonly code?: string;
  readonly details?: unknown;

  constructor(status: number, body?: ApiErrorBody) {
    super(body?.message ?? `API request failed with status ${status}`);

    this.name = "ApiError";
    this.status = status;
    this.code = body?.code;
    this.details = body?.details;

    Object.setPrototypeOf(this, ApiError.prototype);
  }
}
