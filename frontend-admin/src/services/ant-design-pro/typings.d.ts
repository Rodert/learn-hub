// @ts-ignore
/* eslint-disable */

declare namespace API {
  type CurrentUser = {
    name?: string;
    avatar?: string;
    userid?: string;
    email?: string;
    signature?: string;
    title?: string;
    group?: string;
    tags?: { key?: string; label?: string }[];
    notifyCount?: number;
    unreadCount?: number;
    country?: string;
    access?: string;
    geographic?: {
      province?: { label?: string; key?: string };
      city?: { label?: string; key?: string };
    };
    address?: string;
    phone?: string;
  };

  type LoginResult = {
    status?: string;
    type?: string;
    currentAuthority?: string;
  };

  type PageParams = {
    current?: number;
    pageSize?: number;
  };

  type RuleListItem = {
    key?: number;
    disabled?: boolean;
    href?: string;
    avatar?: string;
    name?: string;
    owner?: string;
    desc?: string;
    callNo?: number;
    status?: number;
    updatedAt?: string;
    createdAt?: string;
    progress?: number;
  };

  type RuleList = {
    data?: RuleListItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type FakeCaptcha = {
    code?: number;
    status?: string;
  };

  type LoginParams = {
    username?: string;
    password?: string;
    autoLogin?: boolean;
    type?: string;
  };

  type ErrorResponse = {
    /** 业务约定的错误码 */
    errorCode: string;
    /** 业务上的错误信息 */
    errorMessage?: string;
    /** 业务上的请求是否成功 */
    success?: boolean;
  };

  type NoticeIconList = {
    data?: NoticeIconItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type NoticeIconItemType = 'notification' | 'message' | 'event';

  type NoticeIconItem = {
    id?: string;
    extra?: string;
    key?: string;
    read?: boolean;
    avatar?: string;
    title?: string;
    status?: string;
    datetime?: string;
    description?: string;
    type?: NoticeIconItemType;
  };

  type UserListItem = {
    id?: number;
    username?: string;
    name?: string;
    email?: string;
    phone?: string;
    avatar?: string;
    userid?: string;
    access?: string;
    status?: number;
    roles?: string[];
    createdAt?: string;
    updatedAt?: string;
  };

  type RoleListItem = {
    id?: number;
    code?: string;
    name?: string;
    description?: string;
    status?: number;
    userCount?: number;
    menuCount?: number;
    createdAt?: string;
    updatedAt?: string;
  };

  type RoleOption = {
    id?: number;
    code?: string;
    name?: string;
  };

  type MenuOption = {
    id?: number;
    parentId?: number;
    name?: string;
    path?: string;
    component?: string;
    icon?: string;
    access?: string;
  };

  type CourseListItem = {
    id?: number;
    title?: string;
    description?: string;
    coverImage?: string;
    contentType?: number; // 1-视频，2-文本，3-混合
    videoUrl?: string;
    textContent?: string;
    duration?: number; // 视频时长（秒）
    status?: number; // 0-草稿，1-已发布，2-已下架
    sortOrder?: number;
    createdAt?: string;
    updatedAt?: string;
  };

  type ProgressListItem = {
    userId?: number;
    username?: string;
    name?: string;
    progress?: number;
    duration?: number;
    isCompleted?: boolean;
    completedAt?: string;
    lastStudyAt?: string;
  };

  type UserProgressItem = {
    courseId?: number;
    courseTitle?: string;
    progress?: number;
    duration?: number;
    isCompleted?: boolean;
    completedAt?: string;
    lastStudyAt?: string;
  };
}
