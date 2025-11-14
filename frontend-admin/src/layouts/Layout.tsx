import { useState } from 'react'
import { Layout, Menu, Dropdown, Avatar, Space } from 'antd'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  FileTextOutlined,
  QuestionOutlined,
  FileExcelOutlined,
  UserOutlined,
  SafetyOutlined,
  LogoutOutlined,
} from '@ant-design/icons'
import { useNavigate, Routes, Route } from 'react-router-dom'
import Dashboard from '@/pages/Dashboard'
import Materials from '@/pages/Materials'
import Questions from '@/pages/Questions'
import Exams from '@/pages/Exams'
import Users from '@/pages/Users'
import Roles from '@/pages/Roles'
import './Layout.css'

const { Header, Sider, Content } = Layout

export default function LayoutComponent() {
  const [collapsed, setCollapsed] = useState(false)
  const navigate = useNavigate()
  const user = JSON.parse(localStorage.getItem('user') || '{}')

  const menuItems = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: '仪表盘',
    },
    {
      key: '/materials',
      icon: <FileTextOutlined />,
      label: '资料管理',
    },
    {
      key: '/questions',
      icon: <QuestionOutlined />,
      label: '题库管理',
    },
    {
      key: '/exams',
      icon: <FileExcelOutlined />,
      label: '考试管理',
    },
    {
      key: '/users',
      icon: <UserOutlined />,
      label: '用户管理',
    },
    {
      key: '/roles',
      icon: <SafetyOutlined />,
      label: '角色权限',
    },
  ]

  const userMenu = [
    {
      key: 'profile',
      label: '个人中心',
    },
    {
      key: 'logout',
      label: '退出登录',
      icon: <LogoutOutlined />,
    },
  ]

  const handleUserMenuClick = (e: any) => {
    if (e.key === 'logout') {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      navigate('/login')
    }
  }

  return (
    <Layout style={{ height: '100vh' }}>
      <Sider trigger={null} collapsible collapsed={collapsed}>
        <div className="logo">
          <h2>{collapsed ? 'LH' : 'Learn Hub'}</h2>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          items={menuItems}
          onClick={(e) => navigate(e.key)}
        />
      </Sider>

      <Layout>
        <Header className="header">
          <div className="header-left">
            {collapsed ? (
              <MenuUnfoldOutlined
                className="trigger"
                onClick={() => setCollapsed(false)}
              />
            ) : (
              <MenuFoldOutlined
                className="trigger"
                onClick={() => setCollapsed(true)}
              />
            )}
          </div>

          <div className="header-right">
            <Dropdown menu={{ items: userMenu, onClick: handleUserMenuClick }}>
              <Space>
                <Avatar icon={<UserOutlined />} />
                <span>{user.nickname || user.username}</span>
              </Space>
            </Dropdown>
          </div>
        </Header>

        <Content className="content">
          <Routes>
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/materials" element={<Materials />} />
            <Route path="/questions" element={<Questions />} />
            <Route path="/exams" element={<Exams />} />
            <Route path="/users" element={<Users />} />
            <Route path="/roles" element={<Roles />} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  )
}
