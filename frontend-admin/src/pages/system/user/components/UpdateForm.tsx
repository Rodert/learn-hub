import {
  type ActionType,
  ModalForm,
  ProFormSelect,
  ProFormText,
} from '@ant-design/pro-components';
import { useIntl, useRequest } from '@umijs/max';
import { message } from 'antd';
import type { FC } from 'react';
import { updateUser, getAllRoles } from '@/services/ant-design-pro/user';

export type UpdateFormProps = {
  trigger?: React.ReactElement<any>;
  values: Partial<API.UserListItem>;
  onOk?: () => void;
};

const UpdateForm: FC<UpdateFormProps> = (props) => {
  const { trigger, values, onOk } = props;
  const intl = useIntl();
  const [messageApi, contextHolder] = message.useMessage();

  const { data: roles = [] } = useRequest(getAllRoles);

  const { run, loading } = useRequest(
    (data) => updateUser(values.id!, data),
    {
      manual: true,
      onSuccess: () => {
        messageApi.success('更新成功');
        onOk?.();
      },
      onError: () => {
        messageApi.error('更新失败，请重试！');
      },
    },
  );

  return (
    <>
      {contextHolder}
      <ModalForm
        title="编辑用户"
        trigger={trigger}
        width="600px"
        modalProps={{ okButtonProps: { loading } }}
        initialValues={values}
        onFinish={async (formValues) => {
          await run({
            ...formValues,
            roleIds: formValues.roleIds || [],
          } as API.UserListItem);
          return true;
        }}
      >
        <ProFormText name="username" label="用户名" disabled width="md" />
        <ProFormText.Password
          name="password"
          label="密码"
          tooltip="留空则不修改密码"
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

export default UpdateForm;

