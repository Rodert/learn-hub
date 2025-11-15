import { PlusOutlined } from '@ant-design/icons';
import {
  type ActionType,
  ModalForm,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { useIntl, useRequest } from '@umijs/max';
import { Button, message } from 'antd';
import type { FC } from 'react';
import { createRole, getAllMenus } from '@/services/ant-design-pro/role';

interface CreateFormProps {
  reload?: ActionType['reload'];
}

const CreateForm: FC<CreateFormProps> = (props) => {
  const { reload } = props;
  const intl = useIntl();
  const [messageApi, contextHolder] = message.useMessage();

  const { data: menus = [] } = useRequest(getAllMenus);

  const { run, loading } = useRequest(createRole, {
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
        title="新建角色"
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
            menuIds: values.menuIds || [],
          } as API.RoleListItem);
          return true;
        }}
      >
        <ProFormText
          name="code"
          label="角色代码"
          rules={[{ required: true, message: '请输入角色代码!' }]}
          width="md"
          tooltip="唯一标识，如：admin, user"
        />
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
          initialValue={1}
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

export default CreateForm;

