import { useState, useEffect } from 'react'
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Spin,
  Tree,
  Tabs,
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '@/services/api'

interface Role {
  id: number
  name: string
  description: string
  created_at: string
}

interface Permission {
  id: number
  name: string
  description: string
  resource: string
  action: string
}

export default function Roles() {
  const [roles, setRoles] = useState<Role[]>([])
  const [permissions, setPermissions] = useState<Permission[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [permModalVisible, setPermModalVisible] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [selectedRoleId, setSelectedRoleId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [permForm] = Form.useForm()
  const [pagination, setPagination] = useState({ page: 1, limit: 10, total: 0 })

  // 获取角色列表
  const fetchRoles = async (page = 1) => {
    setLoading(true)
    try {
      const response = await api.get('/admin/roles', {
        params: { page, limit: pagination.limit },
      })
      setRoles(response.data.data || [])
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取角色列表失败')
    } finally {
      setLoading(false)
    }
  }

  // 获取权限列表
  const fetchPermissions = async () => {
    try {
      const response = await api.get('/admin/permissions')
      setPermissions(response.data.data || [])
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取权限列表失败')
    }
  }

  useEffect(() => {
    fetchRoles()
    fetchPermissions()
  }, [])

  // 打开新增/编辑对话框
  const handleOpenModal = (role?: Role) => {
    if (role) {
      setEditingId(role.id)
      form.setFieldsValue(role)
    } else {
      setEditingId(null)
      form.resetFields()
    }
    setModalVisible(true)
  }

  // 保存角色
  const handleSave = async (values: any) => {
    try {
      if (editingId) {
        await api.put(`/admin/roles/${editingId}`, values)
        message.success('更新成功')
      } else {
        await api.post('/admin/roles', values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchRoles()
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 删除角色
  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/admin/roles/${id}`)
      message.success('删除成功')
      fetchRoles()
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败')
    }
  }

  // 打开权限分配对话框
  const handleOpenPermModal = (roleId: number) => {
    setSelectedRoleId(roleId)
    setPermModalVisible(true)
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '角色名称', dataIndex: 'name', key: 'name' },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
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
      width: 200,
      render: (_: any, record: Role) => (
        <Space>
          <Button
            type="primary"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleOpenModal(record)}
          >
            编辑
          </Button>
          <Button
            size="small"
            onClick={() => handleOpenPermModal(record.id)}
          >
            权限
          </Button>
          <Popconfirm
            title="删除角色"
            description="确定要删除这个角色吗？"
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

  const permissionTreeData = [
    {
      title: '资料管理',
      key: 'materials',
      children: [
        { title: '查看', key: 'materials:view' },
        { title: '创建', key: 'materials:create' },
        { title: '编辑', key: 'materials:edit' },
        { title: '删除', key: 'materials:delete' },
      ],
    },
    {
      title: '题库管理',
      key: 'questions',
      children: [
        { title: '查看', key: 'questions:view' },
        { title: '创建', key: 'questions:create' },
        { title: '编辑', key: 'questions:edit' },
        { title: '删除', key: 'questions:delete' },
      ],
    },
    {
      title: '考试管理',
      key: 'exams',
      children: [
        { title: '查看', key: 'exams:view' },
        { title: '创建', key: 'exams:create' },
        { title: '编辑', key: 'exams:edit' },
        { title: '删除', key: 'exams:delete' },
      ],
    },
    {
      title: '用户管理',
      key: 'users',
      children: [
        { title: '查看', key: 'users:view' },
        { title: '创建', key: 'users:create' },
        { title: '编辑', key: 'users:edit' },
        { title: '删除', key: 'users:delete' },
      ],
    },
  ]

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Row justify="space-between" align="middle">
          <Col>
            <h2>角色权限管理</h2>
          </Col>
          <Col>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => handleOpenModal()}
            >
              新增角色
            </Button>
          </Col>
        </Row>
      </Card>

      <Card>
        <Spin spinning={loading}>
          <Table
            columns={columns}
            dataSource={roles}
            rowKey="id"
            pagination={false}
          />
        </Spin>
      </Card>

      <Modal
        title={editingId ? '编辑角色' : '新增角色'}
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
            name="name"
            label="角色名称"
            rules={[{ required: true, message: '请输入角色名称' }]}
          >
            <Input placeholder="请输入角色名称" />
          </Form.Item>

          <Form.Item
            name="description"
            label="描述"
            rules={[{ required: true, message: '请输入描述' }]}
          >
            <Input.TextArea rows={3} placeholder="请输入角色描述" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="分配权限"
        open={permModalVisible}
        onOk={() => setPermModalVisible(false)}
        onCancel={() => setPermModalVisible(false)}
        width={600}
      >
        <Tree
          checkable
          defaultExpandAll
          treeData={permissionTreeData}
        />
      </Modal>
    </div>
  )
}
