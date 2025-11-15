// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 查看课程学习进度 GET /api/admin/course/:id/progress */
export async function getCourseProgress(
  id: number,
  params: {
    current?: number;
    pageSize?: number;
    username?: string;
  },
  options?: { [key: string]: any },
) {
  return request<{
    data?: API.ProgressListItem[];
    total?: number;
    success?: boolean;
  }>(`/api/admin/course/${id}/progress`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 查看用户学习进度 GET /api/admin/user/:id/progress */
export async function getUserProgress(
  id: number,
  params: {
    current?: number;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<{
    data?: API.UserProgressItem[];
    total?: number;
    success?: boolean;
  }>(`/api/admin/user/${id}/progress`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

