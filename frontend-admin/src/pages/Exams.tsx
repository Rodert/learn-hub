import { useState, useEffect } from 'react'
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Spin,
  InputNumber,
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '@/services/api'

interface Exam {
  id: number
  title: string
  description: string
  total_score: number
  pass_score: number
  time_limit: number
  status: string
  created_at: string
}

export default function Exams() {
  const [exams, setExams] = useState<Exam[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [pagination, setPagination] = useState({ page: 1, limit: 10, total: 0 })

  // 获取试卷列表
  const fetchExams = async (page = 1) => {
    setLoading(true)
    try {
      const response = await api.get('/exams', {
        params: { page, limit: pagination.limit },
      })
      const { items, total } = response.data.data
      setExams(items || [])
      setPagination({ ...pagination, page, total })
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取试卷列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchExams()
  }, [])

  // 打开新增/编辑对话框
  const handleOpenModal = (exam?: Exam) => {
    if (exam) {
      setEditingId(exam.id)
      form.setFieldsValue(exam)
    } else {
      setEditingId(null)
      form.resetFields()
    }
    setModalVisible(true)
  }

  // 保存试卷
  const handleSave = async (values: any) => {
    try {
      if (editingId) {
        await api.put(`/exams/${editingId}`, values)
        message.success('更新成功')
      } else {
        await api.post('/exams', values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchExams(pagination.page)
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 删除试卷
  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/exams/${id}`)
      message.success('删除成功')
      fetchExams(pagination.page)
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败')
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '标题', dataIndex: 'title', key: 'title', ellipsis: true },
    { title: '总分', dataIndex: 'total_score', key: 'total_score', width: 80 },
    { title: '及格分', dataIndex: 'pass_score', key: 'pass_score', width: 80 },
    {
      title: '时间限制',
      dataIndex: 'time_limit',
      key: 'time_limit',
      width: 100,
      render: (time: number) => `${time}分钟`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => {
        const statusMap: Record<string, string> = {
          draft: '草稿',
          published: '已发布',
          archived: '已归档',
        }
        return statusMap[status] || status
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: Exam) => (
        <Space>
          <Button
            type="primary"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleOpenModal(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="删除试卷"
            description="确定要删除这个试卷吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button danger size="small" icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Row justify="space-between" align="middle">
          <Col>
            <h2>考试管理</h2>
          </Col>
          <Col>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => handleOpenModal()}
            >
              新增试卷
            </Button>
          </Col>
        </Row>
      </Card>

      <Card>
        <Spin spinning={loading}>
          <Table
            columns={columns}
            dataSource={exams}
            rowKey="id"
            pagination={{
              current: pagination.page,
              pageSize: pagination.limit,
              total: pagination.total,
              onChange: (page: number) => fetchExams(page),
            }}
          />
        </Spin>
      </Card>

      <Modal
        title={editingId ? '编辑试卷' : '新增试卷'}
        open={modalVisible}
        onOk={() => form.submit()}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSave}
          style={{ marginTop: 20 }}
        >
          <Form.Item
            name="title"
            label="试卷标题"
            rules={[{ required: true, message: '请输入试卷标题' }]}
          >
            <Input placeholder="请输入试卷标题" />
          </Form.Item>

          <Form.Item
            name="description"
            label="试卷描述"
            rules={[{ required: true, message: '请输入试卷描述' }]}
          >
            <Input.TextArea rows={3} placeholder="请输入试卷描述" />
          </Form.Item>

          <Form.Item
            name="total_score"
            label="总分"
            rules={[{ required: true, message: '请输入总分' }]}
          >
            <InputNumber min={0} max={1000} placeholder="请输入总分" />
          </Form.Item>

          <Form.Item
            name="pass_score"
            label="及格分"
            rules={[{ required: true, message: '请输入及格分' }]}
          >
            <InputNumber min={0} max={1000} placeholder="请输入及格分" />
          </Form.Item>

          <Form.Item
            name="time_limit"
            label="时间限制（分钟）"
            rules={[{ required: true, message: '请输入时间限制' }]}
          >
            <InputNumber min={0} max={1000} placeholder="请输入时间限制" />
          </Form.Item>

          <Form.Item
            name="status"
            label="状态"
            rules={[{ required: true, message: '请选择状态' }]}
          >
            <Select placeholder="请选择状态">
              <Select.Option value="draft">草稿</Select.Option>
              <Select.Option value="published">已发布</Select.Option>
              <Select.Option value="archived">已归档</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
