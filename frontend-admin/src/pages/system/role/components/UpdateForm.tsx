import {
  type ActionType,
  ModalForm,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { useIntl, useRequest } from '@umijs/max';
import { message } from 'antd';
import type { FC } from 'react';
import { updateRole, getAllMenus } from '@/services/ant-design-pro/role';

export type UpdateFormProps = {
  trigger?: React.ReactElement<any>;
  values: Partial<API.RoleListItem>;
  onOk?: () => void;
};

const UpdateForm: FC<UpdateFormProps> = (props) => {
  const { trigger, values, onOk } = props;
  const intl = useIntl();
  const [messageApi, contextHolder] = message.useMessage();

  const { data: menus = [] } = useRequest(getAllMenus);

  const { run, loading } = useRequest(
    (data) => updateRole(values.id!, data),
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
        title="编辑角色"
        trigger={trigger}
        width="600px"
        modalProps={{ okButtonProps: { loading } }}
        initialValues={values}
        onFinish={async (formValues) => {
          await run({
            ...formValues,
            menuIds: formValues.menuIds || [],
          } as API.RoleListItem);
          return true;
        }}
      >
        <ProFormText name="code" label="角色代码" disabled width="md" />
        <ProFormText
          name="name"
          label="角色名称"
          rules={[{ required: true, message: '请输入角色名称!' }]}
          width="md"
        />
        <ProFormTextArea name="description" label="描述" width="md" />
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
          name="menuIds"
          label="分配菜单"
          width="md"
          mode="multiple"
          options={menus.map((menu: API.MenuOption) => ({
            label: menu.name || menu.path,
            value: menu.id,
          }))}
        />
      </ModalForm>
    </>
  );
};

export default UpdateForm;

