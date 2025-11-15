// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取用户列表 GET /api/user/list */
export async function getUserList(
  params: {
    current?: number;
    pageSize?: number;
    username?: string;
    status?: number;
  },
  options?: { [key: string]: any },
) {
  return request<{
    data?: API.UserListItem[];
    total?: number;
    success?: boolean;
  }>('/api/user/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建用户 POST /api/user */
export async function createUser(body: API.UserListItem, options?: { [key: string]: any }) {
  return request<API.UserListItem>('/api/user', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新用户 PUT /api/user/:id */
export async function updateUser(
  id: number,
  body: API.UserListItem,
  options?: { [key: string]: any },
) {
  return request<API.UserListItem>(`/api/user/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除用户 DELETE /api/user/:id */
export async function deleteUser(id: number, options?: { [key: string]: any }) {
  return request<Record<string, any>>(`/api/user/${id}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}

/** 获取所有角色 GET /api/user/roles */
export async function getAllRoles(options?: { [key: string]: any }) {
  return request<API.RoleOption[]>('/api/user/roles', {
    method: 'GET',
    ...(options || {}),
  });
}

