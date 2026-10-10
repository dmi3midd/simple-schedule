import type { ApiUserError } from '../types/schedule';

export class ApiError extends Error {
  code: number;
  timestamp?: string;

  constructor(message: string, code: number, timestamp?: string) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.timestamp = timestamp;
  }
}

const BASE_URL = '';

export async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${BASE_URL}${endpoint}`;
  const headers = new Headers(options.headers || {});
  
  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (!response.ok) {
    let errorMessage = `Помилка запиту (${response.status})`;
    try {
      const errData: ApiUserError = await response.json();
      if (errData && errData.message) {
        errorMessage = errData.message;
      }
      throw new ApiError(errorMessage, response.status, errData?.timestamp);
    } catch (e) {
      if (e instanceof ApiError) throw e;
      throw new ApiError(errorMessage, response.status);
    }
  }

  // Handle 204 or empty responses
  if (response.status === 204) {
    return {} as T;
  }

  const contentType = response.headers.get('content-type');
  if (contentType && contentType.includes('application/json')) {
    return response.json();
  }

  return {} as T;
}
