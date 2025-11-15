// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取课程列表 GET /api/course/list */
export async function getCourseList(
  params: {
    current?: number;
    pageSize?: number;
    title?: string;
    status?: number;
  },
  options?: { [key: string]: any },
) {
  return request<{
    data?: API.CourseListItem[];
    total?: number;
    success?: boolean;
  }>('/api/course/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 获取课程详情 GET /api/course/:id */
export async function getCourseDetail(
  id: number,
  options?: { [key: string]: any },
) {
  return request<API.CourseListItem>(`/api/course/${id}`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 创建课程 POST /api/course */
export async function createCourse(
  body: API.CourseListItem,
  options?: { [key: string]: any },
) {
  return request<API.CourseListItem>('/api/course', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新课程 PUT /api/course/:id */
export async function updateCourse(
  id: number,
  body: API.CourseListItem,
  options?: { [key: string]: any },
) {
  return request<API.CourseListItem>(`/api/course/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除课程 DELETE /api/course/:id */
export async function deleteCourse(
  id: number,
  options?: { [key: string]: any },
) {
  return request<Record<string, any>>(`/api/course/${id}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}

/** 发布/下架课程 POST /api/course/:id/publish */
export async function publishCourse(
  id: number,
  body: { status: number },
  options?: { [key: string]: any },
) {
  return request<Record<string, any>>(`/api/course/${id}/publish`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

