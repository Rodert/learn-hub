import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { useIntl, useRequest } from '@umijs/max';
import { Select, Tag } from 'antd';
import React, { useRef, useState } from 'react';
import { getCourseList } from '@/services/ant-design-pro/course';
import { getCourseProgress } from '@/services/ant-design-pro/progress';

const CourseProgress: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [selectedCourseId, setSelectedCourseId] = useState<number | undefined>();
  const intl = useIntl();

  const { data: courses = [] } = useRequest(getCourseList, {
    defaultParams: [{ current: 1, pageSize: 100 }],
  });

  const columns: ProColumns<API.ProgressListItem>[] = [
    {
      title: '用户ID',
      dataIndex: 'userId',
      hideInSearch: true,
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      ellipsis: true,
    },
    {
      title: '姓名',
      dataIndex: 'name',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '学习进度',
      dataIndex: 'progress',
      hideInSearch: true,
      width: 150,
      render: (_, record) => (
        <div>
          <div style={{ marginBottom: 4 }}>
            {record.progress}%
          </div>
          <div
            style={{
              width: '100%',
              height: 8,
              backgroundColor: '#f0f0f0',
              borderRadius: 4,
              overflow: 'hidden',
            }}
          >
            <div
              style={{
                width: `${record.progress || 0}%`,
                height: '100%',
                backgroundColor: record.progress === 100 ? '#52c41a' : '#1890ff',
                transition: 'width 0.3s',
              }}
            />
          </div>
        </div>
      ),
    },
    {
      title: '已学习时长',
      dataIndex: 'duration',
      hideInSearch: true,
      width: 120,
      render: (_, record) => {
        const minutes = Math.floor((record.duration || 0) / 60);
        const seconds = (record.duration || 0) % 60;
        return `${minutes}:${seconds.toString().padStart(2, '0')}`;
      },
    },
    {
      title: '完成状态',
      dataIndex: 'isCompleted',
      hideInSearch: true,
      width: 100,
      render: (_, record) => (
        <Tag color={record.isCompleted ? 'green' : 'blue'}>
          {record.isCompleted ? '已完成' : '进行中'}
        </Tag>
      ),
    },
    {
      title: '完成时间',
      dataIndex: 'completedAt',
      valueType: 'dateTime',
      hideInSearch: true,
      render: (_, record) => record.completedAt || '-',
    },
    {
      title: '最后学习时间',
      dataIndex: 'lastStudyAt',
      valueType: 'dateTime',
      hideInSearch: true,
      sorter: true,
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.ProgressListItem, API.PageParams>
        headerTitle="学习进度查看"
        actionRef={actionRef}
        rowKey="userId"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Select
            key="course"
            style={{ width: 200 }}
            placeholder="请选择课程"
            value={selectedCourseId}
            onChange={(value) => {
              setSelectedCourseId(value);
              actionRef.current?.reload();
            }}
            allowClear
          >
            {courses.data?.map((course) => (
              <Select.Option key={course.id} value={course.id}>
                {course.title}
              </Select.Option>
            ))}
          </Select>,
        ]}
        request={async (params) => {
          if (!selectedCourseId) {
            return {
              data: [],
              success: true,
              total: 0,
            };
          }
          const result = await getCourseProgress(
            selectedCourseId,
            {
              current: params.current,
              pageSize: params.pageSize,
              username: params.username,
            },
          );
          return {
            data: result.data || [],
            success: result.success,
            total: result.total || 0,
          };
        }}
        columns={columns}
        pagination={{
          defaultPageSize: 20,
        }}
      />
    </PageContainer>
  );
};

export default CourseProgress;

