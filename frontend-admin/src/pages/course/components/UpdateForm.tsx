import {
  type ActionType,
  ModalForm,
  ProFormRadio,
  ProFormText,
  ProFormTextArea,
  ProFormDigit,
} from '@ant-design/pro-components';
import { useIntl, useRequest } from '@umijs/max';
import { message } from 'antd';
import type { FC } from 'react';
import { updateCourse } from '@/services/ant-design-pro/course';

export type UpdateFormProps = {
  trigger?: React.ReactElement<any>;
  values: Partial<API.CourseListItem>;
  onOk?: () => void;
};

const UpdateForm: FC<UpdateFormProps> = (props) => {
  const { trigger, values, onOk } = props;
  const intl = useIntl();
  const [messageApi, contextHolder] = message.useMessage();

  const { run, loading } = useRequest(
    (data) => updateCourse(values.id!, data),
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
        title="编辑课程"
        trigger={trigger}
        width="800px"
        modalProps={{ okButtonProps: { loading } }}
        initialValues={values}
        onFinish={async (formValues) => {
          await run(formValues as API.CourseListItem);
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
        />
        <ProFormTextArea
          name="textContent"
          label="文本内容"
          width="md"
          fieldProps={{ rows: 8 }}
        />
        <ProFormDigit
          name="duration"
          label="视频时长（秒）"
          width="md"
          min={0}
        />
        <ProFormRadio.Group
          name="status"
          label="状态"
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
          min={0}
        />
      </ModalForm>
    </>
  );
};

export default UpdateForm;

