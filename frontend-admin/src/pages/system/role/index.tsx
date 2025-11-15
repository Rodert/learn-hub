import type { ActionType, ProColumns } from '@ant-design/pro-components';
import {
  FooterToolbar,
  PageContainer,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl, useRequest } from '@umijs/max';
import { Button, message, Popconfirm, Tag } from 'antd';
import React, { useRef, useState } from 'react';
import { deleteRole, getRoleList } from '@/services/ant-design-pro/role';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const RoleList: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [selectedRowsState, setSelectedRows] = useState<API.RoleListItem[]>([]);
  const intl = useIntl();
  const [messageApi, contextHolder] = message.useMessage();

  const { run: delRun, loading } = useRequest(deleteRole, {
    manual: true,
    onSuccess: () => {
      setSelectedRows([]);
      actionRef.current?.reloadAndRest?.();
      messageApi.success('删除成功');
    },
    onError: () => {
      messageApi.error('删除失败，请重试');
    },
  });

  const columns: ProColumns<API.RoleListItem>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      hideInSearch: true,
      width: 80,
    },
    {
      title: '角色代码',
      dataIndex: 'code',
      ellipsis: true,
    },
    {
      title: '角色名称',
      dataIndex: 'name',
      ellipsis: true,
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '用户数',
      dataIndex: 'userCount',
      hideInSearch: true,
      width: 100,
    },
    {
      title: '菜单数',
      dataIndex: 'menuCount',
      hideInSearch: true,
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'status',
      valueType: 'select',
      valueEnum: {
        0: { text: '禁用', status: 'Default' },
        1: { text: '启用', status: 'Success' },
      },
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      valueType: 'dateTime',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => [
        <UpdateForm
          key="edit"
          trigger={<a>编辑</a>}
          values={record}
          onOk={actionRef.current?.reload}
        />,
        <Popconfirm
          key="delete"
          title="确定要删除吗？"
          onConfirm={() => {
            if (record.id) {
              delRun(record.id);
            }
          }}
        >
          <a style={{ color: 'red' }}>删除</a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer>
      {contextHolder}
      <ProTable<API.RoleListItem, API.PageParams>
        headerTitle="角色管理"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <CreateForm key="create" reload={actionRef.current?.reload} />,
        ]}
        request={async (params) => {
          const result = await getRoleList({
            current: params.current,
            pageSize: params.pageSize,
            code: params.code,
            status: params.status,
          });
          return {
            data: result.data || [],
            success: result.success,
            total: result.total || 0,
          };
        }}
        columns={columns}
        rowSelection={{
          onChange: (_, selectedRows) => {
            setSelectedRows(selectedRows);
          },
        }}
      />
      {selectedRowsState?.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              已选择 <a style={{ fontWeight: 600 }}>{selectedRowsState.length}</a> 项
            </div>
          }
        >
          <Button
            loading={loading}
            onClick={async () => {
              await Promise.all(
                selectedRowsState.map((row) => row.id && delRun(row.id)),
              );
              setSelectedRows([]);
              actionRef.current?.reloadAndRest?.();
            }}
          >
            批量删除
          </Button>
        </FooterToolbar>
      )}
    </PageContainer>
  );
};

export default RoleList;

