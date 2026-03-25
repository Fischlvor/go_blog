/**
 * 用户接口（需要认证）
 */

import { authRequest } from '../http';
import type { UserProfile } from '../types';

export async function getUserInfo(): Promise<UserProfile> {
  return authRequest<UserProfile>('/user/info');
}

export interface UpdateUserInfoRequest {
  nickname: string;
  avatar: string;
  address?: string;
  signature?: string;
}

export async function updateUserInfo(data: UpdateUserInfoRequest): Promise<void> {
  return authRequest<void>('/user/changeInfo', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function resetPassword(password: string, newPassword: string): Promise<void> {
  return authRequest<void>('/user/resetPassword', {
    method: 'PUT',
    body: JSON.stringify({ password, new_password: newPassword }),
  });
}

export async function logout(): Promise<void> {
  return authRequest<void>('/user/logout', { method: 'POST' });
}

export interface ForgotPasswordRequest {
  email: string;
  verification_code: string;
  new_password: string;
}

export async function forgotPassword(data: ForgotPasswordRequest): Promise<void> {
  return authRequest<void>('/user/forgotPassword', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}
