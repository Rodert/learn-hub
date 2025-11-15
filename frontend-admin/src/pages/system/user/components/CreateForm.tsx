import { PlusOutlined } from '@ant-design/icons';
import {
  type ActionType,
  ModalForm,
  ProFormSelect,
  ProFormText,
} from '@ant-design/pro-components';
import { useIntl, useRequest } from '@umijs/max';
import { Button, message } from 'antd';
import type { FC } from 'react';
import { createUser, getAllRoles } from '@/services/ant-design-pro/user';

interface CreateFormProps {
  reload?: ActionType['reload'];
}

const CreateForm: FC<CreateFormProps> = (props) => {
  const { reload } = props;
  const intl = useIntl();
  const [messageApi, contextHolder] = message.useMessage();

  const { data: roles = [] } = useRequest(getAllRoles);

  const { run, loading } = useRequest(createUser, {
    manual: true,
    onSuccess: () => {
      messageApi.success('创建成功');
      reload?.();
    },
    onError: () => {
      messageApi.error('创建失败，请重试！');
    },
  });

  return (
    <>
      {contextHolder}
      <ModalForm
        title="新建用户"
        trigger={
          <Button type="primary" icon={<PlusOutlined />}>
            新建
          </Button>
        }
        width="600px"
        modalProps={{ okButtonProps: { loading } }}
        onFinish={async (values) => {
          await run({
            ...values,
            roleIds: values.roleIds || [],
          } as API.UserListItem);
          return true;
        }}
      >
        <ProFormText
          name="username"
          label="用户名"
          rules={[{ required: true, message: '请输入用户名!' }]}
          width="md"
        />
        <ProFormText.Password
          name="password"
          label="密码"
          rules={[{ required: true, message: '请输入密码!' }]}
          width="md"
        />
        <ProFormText name="name" label="姓名" width="md" />
        <ProFormText name="email" label="邮箱" width="md" />
        <ProFormText name="phone" label="手机号" width="md" />
        <ProFormSelect
          name="access"
          label="权限"
          width="md"
          options={[
            { label: '管理员', value: 'admin' },
            { label: '普通用户', value: 'user' },
          ]}
        />
        <ProFormSelect
          name="status"
          label="状态"
          width="md"
          initialValue={1}
          options={[
            { label: '启用', value: 1 },
            { label: '禁用', value: 0 },
          ]}
        />
        <ProFormSelect
          name="roleIds"
          label="角色"
          width="md"
          mode="multiple"
          options={roles.map((role: API.RoleOption) => ({
            label: role.name,
            value: role.id,
          }))}
        />
      </ModalForm>
    </>
  );
};

export default CreateForm;

