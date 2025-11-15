// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取角色列表 GET /api/role/list */
export async function getRoleList(
  params: {
    current?: number;
    pageSize?: number;
    code?: string;
    status?: number;
  },
  options?: { [key: string]: any },
) {
  return request<{
    data?: API.RoleListItem[];
    total?: number;
    success?: boolean;
  }>('/api/role/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建角色 POST /api/role */
export async function createRole(body: API.RoleListItem, options?: { [key: string]: any }) {
  return request<API.RoleListItem>('/api/role', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新角色 PUT /api/role/:id */
export async function updateRole(
  id: number,
  body: API.RoleListItem,
  options?: { [key: string]: any },
) {
  return request<API.RoleListItem>(`/api/role/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除角色 DELETE /api/role/:id */
export async function deleteRole(id: number, options?: { [key: string]: any }) {
  return request<Record<string, any>>(`/api/role/${id}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}

/** 获取所有菜单 GET /api/role/menus */
export async function getAllMenus(options?: { [key: string]: any }) {
  return request<API.MenuOption[]>('/api/role/menus', {
    method: 'GET',
    ...(options || {}),
  });
}

