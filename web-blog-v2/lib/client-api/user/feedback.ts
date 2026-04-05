/**
 * 用户反馈接口（需要认证）
 */

import { authRequest } from '../http';
import type { Feedback } from '../types';

export interface CreateFeedbackRequest {
  content: string;
}

export async function createFeedback(data: CreateFeedbackRequest): Promise<void> {
  return authRequest<void>('/feedback/create', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getUserFeedbacks(): Promise<Feedback[]> {
  return authRequest<Feedback[]>('/feedback/info');
}
