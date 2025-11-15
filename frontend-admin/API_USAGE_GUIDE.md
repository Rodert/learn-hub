# Ant Design Pro 项目用法与接口数据格式指南

## 📋 目录
1. [接口数据格式规范](#接口数据格式规范)
2. [请求方法使用](#请求方法使用)
3. [组件使用模式](#组件使用模式)
4. [状态管理](#状态管理)
5. [权限控制](#权限控制)
6. [国际化使用](#国际化使用)

---

## 🔌 接口数据格式规范

### 1. 标准响应格式

所有接口统一使用以下响应格式：

```typescript
interface ResponseStructure {
  success: boolean;        // 请求是否成功
  data: any;              // 响应数据
  errorCode?: number;     // 错误码（失败时）
  errorMessage?: string;  // 错误信息（失败时）
  showType?: ErrorShowType; // 错误展示类型
}
```

### 2. 列表接口响应格式

列表查询接口返回格式：

```typescript
type RuleList = {
  data?: RuleListItem[];  // 列表数据
  total?: number;         // 总条数
  success?: boolean;      // 是否成功
  pageSize?: number;      // 每页条数（可选）
  current?: number;       // 当前页码（可选）
};
```

**示例：**
```json
{
  "success": true,
  "data": [
    {
      "key": 1,
      "name": "规则名称",
      "desc": "描述",
      "status": 1,
      "updatedAt": "2024-01-01"
    }
  ],
  "total": 100,
  "pageSize": 10,
  "current": 1
}
```

### 3. 登录接口响应格式

```typescript
type LoginResult = {
  status?: 'ok' | 'error';  // 登录状态
  type?: string;            // 登录类型（account/mobile）
  currentAuthority?: string; // 用户权限（admin/user/guest）
};
```

**成功示例：**
```json
{
  "status": "ok",
  "type": "account",
  "currentAuthority": "admin"
}
```

**失败示例：**
```json
{
  "status": "error",
  "type": "account",
  "currentAuthority": "guest"
}
```

### 4. 用户信息接口响应格式

```typescript
type CurrentUser = {
  name?: string;
  avatar?: string;
  userid?: string;
  email?: string;
  signature?: string;
  title?: string;
  group?: string;
  tags?: { key?: string; label?: string }[];
  notifyCount?: number;
  unreadCount?: number;
  country?: string;
  access?: string;  // 权限标识：'admin' | 'user' | 'guest'
  geographic?: {
    province?: { label?: string; key?: string };
    city?: { label?: string; key?: string };
  };
  address?: string;
  phone?: string;
};
```

**响应示例：**
```json
{
  "success": true,
  "data": {
    "name": "Serati Ma",
    "avatar": "https://...",
    "userid": "00000001",
    "email": "antdesign@alipay.com",
    "access": "admin",
    "tags": [
      { "key": "0", "label": "很有想法的" }
    ]
  }
}
```

### 5. 分页参数格式

```typescript
type PageParams = {
  current?: number;   // 当前页码，从1开始
  pageSize?: number;  // 每页条数
};
```

---

## 🌐 请求方法使用

### 1. 基础请求方法

使用 `@umijs/max` 的 `request` 方法：

```typescript
import { request } from '@umijs/max';

// GET 请求
export async function getData(params: API.PageParams) {
  return request<API.RuleList>('/api/rule', {
    method: 'GET',
    params: {
      ...params,
    },
  });
}

// POST 请求
export async function createData(body: API.RuleListItem) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'POST',
    data: body,
  });
}

// PUT 请求
export async function updateData(body: API.RuleListItem) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'PUT',
    data: body,
  });
}

// DELETE 请求
export async function deleteData(id: number) {
  return request(`/api/rule/${id}`, {
    method: 'DELETE',
  });
}
```

### 2. 使用 useRequest Hook

推荐使用 `useRequest` 进行请求管理：

```typescript
import { useRequest } from '@umijs/max';
import { message } from 'antd';
import { addRule } from '@/services/ant-design-pro/api';

const { run, loading } = useRequest(addRule, {
  manual: true,  // 手动触发
  onSuccess: () => {
    message.success('操作成功');
    // 刷新列表
    actionRef.current?.reload();
  },
  onError: () => {
    message.error('操作失败，请重试');
  },
});

// 调用
await run({ data: formValues });
```

### 3. 请求拦截器配置

在 `src/requestErrorConfig.ts` 中配置：

```typescript
// 请求拦截器 - 添加 token
requestInterceptors: [
  (config: RequestOptions) => {
    const url = config?.url?.concat('?token=123');
    return { ...config, url };
  },
],

// 响应拦截器 - 统一处理响应
responseInterceptors: [
  (response) => {
    const { data } = response as unknown as ResponseStructure;
    if (data?.success === false) {
      message.error('请求失败！');
    }
    return response;
  },
],
```

### 4. 错误处理

错误类型枚举：

```typescript
enum ErrorShowType {
  SILENT = 0,        // 静默处理
  WARN_MESSAGE = 1,  // 警告提示
  ERROR_MESSAGE = 2, // 错误提示
  NOTIFICATION = 3,  // 通知提示
  REDIRECT = 9,      // 重定向
}
```

后端返回错误时，会自动根据 `showType` 进行相应处理。

---

## 🧩 组件使用模式

### 1. ProTable 使用

```typescript
import { ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';

const columns: ProColumns<API.RuleListItem>[] = [
  {
    title: '规则名称',
    dataIndex: 'name',
    valueType: 'text',
  },
  {
    title: '状态',
    dataIndex: 'status',
    valueType: 'select',
    valueEnum: {
      0: { text: '关闭', status: 'Default' },
      1: { text: '运行中', status: 'Processing' },
      2: { text: '已上线', status: 'Success' },
      3: { text: '异常', status: 'Error' },
    },
  },
  {
    title: '更新时间',
    dataIndex: 'updatedAt',
    valueType: 'dateTime',
    sorter: true,
  },
];

<ProTable<API.RuleListItem, API.PageParams>
  headerTitle="查询表格"
  actionRef={actionRef}
  rowKey="key"
  search={{
    labelWidth: 120,
  }}
  request={rule}  // 请求方法
  columns={columns}
  rowSelection={{
    onChange: (_, selectedRows) => {
      setSelectedRows(selectedRows);
    },
  }}
/>
```

**关键点：**
- `request` 属性接收一个函数，自动处理分页、排序、筛选
- `valueType` 支持多种类型：`text`、`select`、`dateTime`、`textarea` 等
- `valueEnum` 用于枚举值映射
- `sorter: true` 启用排序

### 2. ModalForm 使用

```typescript
import { ModalForm, ProFormText, ProFormTextArea } from '@ant-design/pro-components';

<ModalForm
  title="新建规则"
  trigger={<Button type="primary">新建</Button>}
  width="400px"
  onFinish={async (values) => {
    await run({ data: values as API.RuleListItem });
    return true;  // 返回 true 会关闭弹窗
  }}
>
  <ProFormText
    name="name"
    label="规则名称"
    rules={[{ required: true, message: '请输入规则名称' }]}
  />
  <ProFormTextArea name="desc" label="描述" />
</ModalForm>
```

### 3. StepsForm 使用

```typescript
import { StepsForm, ProFormText, ProFormSelect } from '@ant-design/pro-components';

<StepsForm
  onFinish={async (values) => {
    await run({ data: values });
  }}
>
  <StepsForm.StepForm
    title="基本信息"
    initialValues={values}
  >
    <ProFormText name="name" label="名称" />
  </StepsForm.StepForm>
  
  <StepsForm.StepForm
    title="配置属性"
  >
    <ProFormSelect
      name="type"
      label="类型"
      valueEnum={{
        0: '类型一',
        1: '类型二',
      }}
    />
  </StepsForm.StepForm>
</StepsForm>
```

### 4. ProDescriptions 使用

```typescript
import { ProDescriptions } from '@ant-design/pro-components';

<ProDescriptions<API.RuleListItem>
  column={2}
  title="详情"
  request={async () => ({
    data: currentRow || {},
  })}
  columns={columns as ProDescriptionsItemProps<API.RuleListItem>[]}
/>
```

---

## 🔄 状态管理

### 1. 全局初始状态

使用 `@@initialState` 管理全局状态：

```typescript
// 在 app.tsx 中定义
export async function getInitialState() {
  const fetchUserInfo = async () => {
    try {
      const msg = await queryCurrentUser();
      return msg.data;
    } catch (error) {
      history.push('/user/login');
    }
    return undefined;
  };

  const currentUser = await fetchUserInfo();
  return {
    fetchUserInfo,
    currentUser,
    settings: defaultSettings,
  };
}
```

### 2. 使用全局状态

```typescript
import { useModel } from '@umijs/max';

const { initialState, setInitialState } = useModel('@@initialState');
const { currentUser } = initialState || {};

// 更新状态
setInitialState((s) => ({
  ...s,
  currentUser: userInfo,
}));
```

### 3. 刷新用户信息

```typescript
const fetchUserInfo = async () => {
  const userInfo = await initialState?.fetchUserInfo?.();
  if (userInfo) {
    flushSync(() => {
      setInitialState((s) => ({
        ...s,
        currentUser: userInfo,
      }));
    });
  }
};
```

---

## 🔐 权限控制

### 1. 权限定义

在 `src/access.ts` 中定义权限：

```typescript
export default function access(
  initialState: { currentUser?: API.CurrentUser } | undefined,
) {
  const { currentUser } = initialState ?? {};
  return {
    canAdmin: currentUser && currentUser.access === 'admin',
    canUser: currentUser && currentUser.access === 'user',
  };
}
```

### 2. 路由权限控制

在 `config/routes.ts` 中使用：

```typescript
{
  path: '/admin',
  name: 'admin',
  access: 'canAdmin',  // 需要 canAdmin 权限
  component: './Admin',
}
```

### 3. 组件内权限控制

```typescript
import { useAccess } from '@umijs/max';

const access = useAccess();

{access.canAdmin && (
  <Button>管理员操作</Button>
)}
```

---

## 🌍 国际化使用

### 1. 使用 FormattedMessage

```typescript
import { FormattedMessage } from '@umijs/max';

<FormattedMessage
  id="pages.searchTable.title"
  defaultMessage="查询表格"
/>
```

### 2. 使用 useIntl

```typescript
import { useIntl } from '@umijs/max';

const intl = useIntl();

intl.formatMessage({
  id: 'pages.searchTable.title',
  defaultMessage: '查询表格',
})
```

### 3. 在组件属性中使用

```typescript
title={intl.formatMessage({
  id: 'pages.searchTable.title',
  defaultMessage: '查询表格',
})}
```

---

## 📝 常用类型定义

### API 命名空间

所有 API 类型定义在 `API` 命名空间下：

```typescript
// 在 typings.d.ts 中定义
declare namespace API {
  type CurrentUser = { ... };
  type LoginResult = { ... };
  type RuleListItem = { ... };
  type PageParams = { ... };
}

// 使用
const user: API.CurrentUser = { ... };
```

---

## 🎯 最佳实践

### 1. 接口调用

```typescript
// ✅ 推荐：使用 useRequest
const { run, loading } = useRequest(addRule, {
  manual: true,
  onSuccess: () => {
    message.success('成功');
    reload();
  },
});

// ❌ 不推荐：直接调用
const handleSubmit = async () => {
  try {
    await addRule(data);
    message.success('成功');
  } catch (error) {
    message.error('失败');
  }
};
```

### 2. 表单提交

```typescript
// ✅ 推荐：在 onFinish 中处理
<ModalForm
  onFinish={async (values) => {
    await run({ data: values });
    return true;  // 返回 true 关闭弹窗
  }}
>

// ❌ 不推荐：在外部处理
const handleSubmit = () => {
  form.validateFields().then(values => {
    run({ data: values });
  });
};
```

### 3. 列表刷新

```typescript
// ✅ 推荐：使用 actionRef
const actionRef = useRef<ActionType>();

<ProTable
  actionRef={actionRef}
  request={rule}
/>

// 刷新
actionRef.current?.reload();
actionRef.current?.reloadAndRest();  // 重置并刷新
```

---

## 🔗 相关文档

- [Umi Max 文档](https://umijs.org/docs/max/introduce)
- [Ant Design Pro Components](https://procomponents.ant.design/)
- [Ant Design](https://ant.design/)





