import { PlusOutlined } from '@ant-design/icons';
import {
  type ActionType,
  ModalForm,
  ProFormRadio,
  ProFormText,
  ProFormTextArea,
  ProFormDigit,
} from '@ant-design/pro-components';
import { useIntl, useRequest } from '@umijs/max';
import { Button, message } from 'antd';
import type { FC } from 'react';
import { createCourse } from '@/services/ant-design-pro/course';

interface CreateFormProps {
  reload?: ActionType['reload'];
}

const CreateForm: FC<CreateFormProps> = (props) => {
  const { reload } = props;
  const intl = useIntl();
  const [messageApi, contextHolder] = message.useMessage();

  const { run, loading } = useRequest(createCourse, {
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
        title="新建课程"
        trigger={
          <Button type="primary" icon={<PlusOutlined />}>
            新建
          </Button>
        }
        width="800px"
        modalProps={{ okButtonProps: { loading } }}
        onFinish={async (values) => {
          await run(values as API.CourseListItem);
          return true;
        }}
      >
        <ProFormText
          name="title"
          label="课程标题"
          rules={[{ required: true, message: '请输入课程标题!' }]}
          width="md"
        />
        <ProFormTextArea
          name="description"
          label="课程描述"
          width="md"
          fieldProps={{ rows: 4 }}
        />
        <ProFormText
          name="coverImage"
          label="封面图片URL"
          width="md"
          placeholder="https://example.com/image.jpg"
        />
        <ProFormRadio.Group
          name="contentType"
          label="内容类型"
          rules={[{ required: true, message: '请选择内容类型!' }]}
          options={[
            { label: '视频', value: 1 },
            { label: '文本', value: 2 },
            { label: '混合', value: 3 },
          ]}
        />
        <ProFormText
          name="videoUrl"
          label="视频URL"
          width="md"
          placeholder="https://example.com/video.mp4"
          dependencies={['contentType']}
          rules={[
            ({ getFieldValue }) => ({
              validator: (_, value) => {
                const contentType = getFieldValue('contentType');
                if (contentType === 1 || contentType === 3) {
                  if (!value) {
                    return Promise.reject(new Error('视频类型必须提供视频URL'));
                  }
                }
                return Promise.resolve();
              },
            }),
          ]}
        />
        <ProFormTextArea
          name="textContent"
          label="文本内容"
          width="md"
          fieldProps={{ rows: 8 }}
          dependencies={['contentType']}
          rules={[
            ({ getFieldValue }) => ({
              validator: (_, value) => {
                const contentType = getFieldValue('contentType');
                if (contentType === 2 || contentType === 3) {
                  if (!value) {
                    return Promise.reject(new Error('文本类型必须提供文本内容'));
                  }
                }
                return Promise.resolve();
              },
            }),
          ]}
        />
        <ProFormDigit
          name="duration"
          label="视频时长（秒）"
          width="md"
          min={0}
          placeholder="仅视频类型需要填写"
        />
        <ProFormRadio.Group
          name="status"
          label="状态"
          initialValue={0}
          options={[
            { label: '草稿', value: 0 },
            { label: '已发布', value: 1 },
            { label: '已下架', value: 2 },
          ]}
        />
        <ProFormDigit
          name="sortOrder"
          label="排序"
          width="md"
          initialValue={0}
          min={0}
        />
      </ModalForm>
    </>
  );
};

export default CreateForm;

