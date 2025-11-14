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
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '@/services/api'

interface Material {
  id: number
  title: string
  description: string
  content_type: string
  status: string
  created_at: string
}

export default function Materials() {
  const [materials, setMaterials] = useState<Material[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [pagination, setPagination] = useState({ page: 1, limit: 10, total: 0 })

  // 获取资料列表
  const fetchMaterials = async (page = 1) => {
    setLoading(true)
    try {
      const response = await api.get('/materials', {
        params: { page, limit: pagination.limit, status: 'published' },
      })
      const { items, total } = response.data.data
      setMaterials(items || [])
      setPagination({ ...pagination, page, total })
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取资料列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchMaterials()
  }, [])

  // 打开新增/编辑对话框
  const handleOpenModal = (material?: Material) => {
    if (material) {
      setEditingId(material.id)
      form.setFieldsValue(material)
    } else {
      setEditingId(null)
      form.resetFields()
    }
    setModalVisible(true)
  }

  // 保存资料
  const handleSave = async (values: any) => {
    try {
      if (editingId) {
        await api.put(`/materials/${editingId}`, values)
        message.success('更新成功')
      } else {
        await api.post('/materials', values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchMaterials(pagination.page)
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 删除资料
  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/materials/${id}`)
      message.success('删除成功')
      fetchMaterials(pagination.page)
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败')
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '标题', dataIndex: 'title', key: 'title', ellipsis: true },
    {
      title: '类型',
      dataIndex: 'content_type',
      key: 'content_type',
      width: 100,
      render: (type: string) => {
        const typeMap: Record<string, string> = {
          text: '文本',
          video: '视频',
          file: '文件',
          mixed: '混合',
        }
        return typeMap[type] || type
      },
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
      render: (_: any, record: Material) => (
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
            title="删除资料"
            description="确定要删除这个资料吗？"
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
            <h2>资料管理</h2>
          </Col>
          <Col>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => handleOpenModal()}
            >
              新增资料
            </Button>
          </Col>
        </Row>
      </Card>

      <Card>
        <Spin spinning={loading}>
          <Table
            columns={columns}
            dataSource={materials}
            rowKey="id"
            pagination={{
              current: pagination.page,
              pageSize: pagination.limit,
              total: pagination.total,
              onChange: (page) => fetchMaterials(page),
            }}
          />
        </Spin>
      </Card>

      <Modal
        title={editingId ? '编辑资料' : '新增资料'}
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
            label="标题"
            rules={[{ required: true, message: '请输入标题' }]}
          >
            <Input placeholder="请输入资料标题" />
          </Form.Item>

          <Form.Item
            name="description"
            label="描述"
            rules={[{ required: true, message: '请输入描述' }]}
          >
            <Input.TextArea rows={3} placeholder="请输入资料描述" />
          </Form.Item>

          <Form.Item
            name="content_type"
            label="类型"
            rules={[{ required: true, message: '请选择类型' }]}
          >
            <Select placeholder="请选择资料类型">
              <Select.Option value="text">文本</Select.Option>
              <Select.Option value="video">视频</Select.Option>
              <Select.Option value="file">文件</Select.Option>
              <Select.Option value="mixed">混合</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="content"
            label="内容"
            rules={[{ required: true, message: '请输入内容' }]}
          >
            <Input.TextArea rows={5} placeholder="请输入资料内容" />
          </Form.Item>

          <Form.Item name="file_url" label="文件 URL">
            <Input placeholder="请输入文件 URL（可选）" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
