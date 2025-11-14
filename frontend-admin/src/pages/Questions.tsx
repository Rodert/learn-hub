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

interface Question {
  id: number
  question_type: string
  content: string
  answer: string
  score: number
  created_at: string
}

export default function Questions() {
  const [questions, setQuestions] = useState<Question[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [pagination, setPagination] = useState({ page: 1, limit: 10, total: 0 })

  // 获取题目列表
  const fetchQuestions = async (page = 1) => {
    setLoading(true)
    try {
      const response = await api.get('/questions', {
        params: { page, limit: pagination.limit },
      })
      const { items, total } = response.data.data
      setQuestions(items || [])
      setPagination({ ...pagination, page, total })
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取题目列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchQuestions()
  }, [])

  // 打开新增/编辑对话框
  const handleOpenModal = (question?: Question) => {
    if (question) {
      setEditingId(question.id)
      form.setFieldsValue(question)
    } else {
      setEditingId(null)
      form.resetFields()
    }
    setModalVisible(true)
  }

  // 保存题目
  const handleSave = async (values: any) => {
    try {
      if (editingId) {
        await api.put(`/questions/${editingId}`, values)
        message.success('更新成功')
      } else {
        await api.post('/questions', values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchQuestions(pagination.page)
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 删除题目
  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/questions/${id}`)
      message.success('删除成功')
      fetchQuestions(pagination.page)
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败')
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    {
      title: '题型',
      dataIndex: 'question_type',
      key: 'question_type',
      width: 100,
      render: (type: string) => {
        const typeMap: Record<string, string> = {
          single_choice: '单选',
          multiple_choice: '多选',
          fill_blank: '填空',
        }
        return typeMap[type] || type
      },
    },
    { title: '题目内容', dataIndex: 'content', key: 'content', ellipsis: true },
    { title: '分值', dataIndex: 'score', key: 'score', width: 80 },
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
      render: (_: any, record: Question) => (
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
            title="删除题目"
            description="确定要删除这个题目吗？"
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
            <h2>题库管理</h2>
          </Col>
          <Col>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => handleOpenModal()}
            >
              新增题目
            </Button>
          </Col>
        </Row>
      </Card>

      <Card>
        <Spin spinning={loading}>
          <Table
            columns={columns}
            dataSource={questions}
            rowKey="id"
            pagination={{
              current: pagination.page,
              pageSize: pagination.limit,
              total: pagination.total,
              onChange: (page: number) => fetchQuestions(page),
            }}
          />
        </Spin>
      </Card>

      <Modal
        title={editingId ? '编辑题目' : '新增题目'}
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
            name="question_type"
            label="题型"
            rules={[{ required: true, message: '请选择题型' }]}
          >
            <Select placeholder="请选择题型">
              <Select.Option value="single_choice">单选题</Select.Option>
              <Select.Option value="multiple_choice">多选题</Select.Option>
              <Select.Option value="fill_blank">填空题</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="content"
            label="题目内容"
            rules={[{ required: true, message: '请输入题目内容' }]}
          >
            <Input.TextArea rows={4} placeholder="请输入题目内容" />
          </Form.Item>

          <Form.Item
            name="answer"
            label="标准答案"
            rules={[{ required: true, message: '请输入标准答案' }]}
          >
            <Input placeholder="请输入标准答案" />
          </Form.Item>

          <Form.Item
            name="score"
            label="分值"
            rules={[{ required: true, message: '请输入分值' }]}
          >
            <InputNumber min={0} max={100} placeholder="请输入分值" />
          </Form.Item>

          <Form.Item name="explanation" label="解释">
            <Input.TextArea rows={3} placeholder="请输入答案解释（可选）" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
