import type { ActionType, ProColumns } from '@ant-design/pro-components';
import {
  FooterToolbar,
  PageContainer,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl, useRequest } from '@umijs/max';
import { Button, message, Popconfirm, Tag } from 'antd';
import React, { useRef, useState } from 'react';
import {
  deleteCourse,
  getCourseList,
  publishCourse,
} from '@/services/ant-design-pro/course';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const CourseList: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [selectedRowsState, setSelectedRows] = useState<API.CourseListItem[]>([]);
  const intl = useIntl();
  const [messageApi, contextHolder] = message.useMessage();

  const { run: delRun, loading } = useRequest(deleteCourse, {
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

  const { run: publishRun } = useRequest(
    (id: number, status: number) => publishCourse(id, { status }),
    {
      manual: true,
      onSuccess: () => {
        actionRef.current?.reloadAndRest?.();
        messageApi.success('操作成功');
      },
      onError: () => {
        messageApi.error('操作失败，请重试');
      },
    },
  );

  const columns: ProColumns<API.CourseListItem>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      hideInSearch: true,
      width: 80,
    },
    {
      title: '课程标题',
      dataIndex: 'title',
      ellipsis: true,
    },
    {
      title: '类型',
      dataIndex: 'contentType',
      hideInSearch: true,
      valueType: 'select',
      valueEnum: {
        1: { text: '视频', status: 'Processing' },
        2: { text: '文本', status: 'Success' },
        3: { text: '混合', status: 'Warning' },
      },
      render: (_, record) => {
        const typeMap: Record<number, { text: string; color: string }> = {
          1: { text: '视频', color: 'blue' },
          2: { text: '文本', color: 'green' },
          3: { text: '混合', color: 'orange' },
        };
        const type = typeMap[record.contentType || 1];
        return <Tag color={type.color}>{type.text}</Tag>;
      },
    },
    {
      title: '时长',
      dataIndex: 'duration',
      hideInSearch: true,
      width: 100,
      render: (_, record) => {
        if (record.contentType === 2) {
          return '-';
        }
        const minutes = Math.floor((record.duration || 0) / 60);
        const seconds = (record.duration || 0) % 60;
        return `${minutes}:${seconds.toString().padStart(2, '0')}`;
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      valueType: 'select',
      valueEnum: {
        0: { text: '草稿', status: 'Default' },
        1: { text: '已发布', status: 'Success' },
        2: { text: '已下架', status: 'Error' },
      },
    },
    {
      title: '排序',
      dataIndex: 'sortOrder',
      hideInSearch: true,
      width: 80,
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
      width: 200,
      render: (_, record) => [
        <UpdateForm
          key="edit"
          trigger={<a>编辑</a>}
          values={record}
          onOk={actionRef.current?.reload}
        />,
        record.status === 1 ? (
          <a
            key="publish"
            onClick={() => {
              if (record.id) {
                publishRun(record.id, 2);
              }
            }}
          >
            下架
          </a>
        ) : (
          <a
            key="publish"
            onClick={() => {
              if (record.id) {
                publishRun(record.id, 1);
              }
            }}
          >
            发布
          </a>
        ),
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
      <ProTable<API.CourseListItem, API.PageParams>
        headerTitle="课程管理"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <CreateForm key="create" reload={actionRef.current?.reload} />,
        ]}
        request={async (params) => {
          const result = await getCourseList({
            current: params.current,
            pageSize: params.pageSize,
            title: params.title,
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

export default CourseList;

