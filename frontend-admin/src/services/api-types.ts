/**
 * API 类型定义
 */

// ============ 通用类型 ============

export interface Response<T = any> {
  code: number
  message: string
  data: T
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  limit: number
}

// ============ 认证相关 ============

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: UserInfo
}

export interface UserInfo {
  id: number
  username: string
  nickname: string
  status: 'active' | 'inactive' | 'banned'
}

// ============ 用户相关 ============

export interface User {
  id: number
  username: string
  nickname: string
  status: 'active' | 'inactive' | 'banned'
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  username: string
  password: string
  nickname: string
}

export interface UpdateUserRequest {
  nickname?: string
  status?: 'active' | 'inactive' | 'banned'
}

// ============ 资料相关 ============

export interface Material {
  id: number
  title: string
  description: string
  content_type: 'text' | 'video' | 'file' | 'mixed'
  content: string
  file_url?: string
  file_size?: number
  cover_url?: string
  status: 'draft' | 'published' | 'archived'
  created_at: string
  updated_at: string
}

export interface CreateMaterialRequest {
  title: string
  description: string
  content_type: 'text' | 'video' | 'file' | 'mixed'
  content: string
  file_url?: string
  cover_url?: string
}

export interface UpdateMaterialRequest {
  title?: string
  description?: string
  status?: 'draft' | 'published' | 'archived'
}

// ============ 题库相关 ============

export interface Question {
  id: number
  exam_id?: number
  question_type: 'single_choice' | 'multiple_choice' | 'fill_blank'
  content: string
  options?: string
  answer: string
  explanation?: string
  score: number
  created_at: string
  updated_at: string
}

export interface CreateQuestionRequest {
  exam_id?: number
  question_type: 'single_choice' | 'multiple_choice' | 'fill_blank'
  content: string
  answer: string
  explanation?: string
  score: number
}

export interface UpdateQuestionRequest {
  content?: string
  answer?: string
  explanation?: string
  score?: number
}

// ============ 考试相关 ============

export interface Exam {
  id: number
  title: string
  description: string
  total_score: number
  pass_score: number
  time_limit: number
  status: 'draft' | 'published' | 'archived'
  created_at: string
  updated_at: string
  questions?: Question[]
}

export interface CreateExamRequest {
  title: string
  description: string
  total_score: number
  pass_score: number
  time_limit: number
}

export interface UpdateExamRequest {
  title?: string
  description?: string
  total_score?: number
  pass_score?: number
  time_limit?: number
  status?: 'draft' | 'published' | 'archived'
}

export interface ExamRecord {
  id: number
  user_id: number
  exam_id: number
  score?: number
  status: 'in_progress' | 'submitted' | 'graded'
  answers: string
  start_time?: string
  submit_time?: string
  created_at: string
  updated_at: string
}

export interface StartExamResponse {
  exam_record_id: number
  exam: Exam
  questions: Question[]
}

export interface SubmitExamRequest {
  exam_record_id: number
  answers: Array<{
    question_id: number
    answer: string
  }>
}

// ============ 学习记录相关 ============

export interface CourseRecord {
  id: number
  user_id: number
  material_id: number
  status: 'not_started' | 'in_progress' | 'completed'
  progress_percent: number
  view_duration: number
  completed_at?: string
  created_at: string
  updated_at: string
}

export interface UpdateCourseRecordRequest {
  progress_percent?: number
  view_duration?: number
}

// ============ 角色相关 ============

export interface Role {
  id: number
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface CreateRoleRequest {
  name: string
  description: string
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
}

// ============ 权限相关 ============

export interface Permission {
  id: number
  name: string
  description: string
  resource: string
  action: string
  created_at: string
  updated_at: string
}

export interface CreatePermissionRequest {
  name: string
  description: string
  resource: string
  action: string
}

// ============ 菜单相关 ============

export interface Menu {
  id: number
  name: string
  path: string
  icon?: string
  component?: string
  parent_id?: number
  order_num: number
  visible: boolean
  type: 'menu' | 'button'
  permission?: string
  children?: Menu[]
}

// ============ 查询参数 ============

export interface PaginationParams {
  page?: number
  limit?: number
}

export interface MaterialListParams extends PaginationParams {
  status?: 'draft' | 'published' | 'archived'
}

export interface QuestionListParams extends PaginationParams {
  type?: 'single_choice' | 'multiple_choice' | 'fill_blank'
}

export interface ExamListParams extends PaginationParams {
  status?: 'draft' | 'published' | 'archived'
}

export interface UserListParams extends PaginationParams {
  status?: 'active' | 'inactive' | 'banned'
}
