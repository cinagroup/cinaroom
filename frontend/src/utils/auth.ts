/**
 * 权限工具函数 - 参照 CinaToken 权限系统
 * user.role >= 10  → 管理员
 * user.role >= 100 → root
 */

export function isAdmin(user: any): boolean {
  return user && typeof user.role === 'number' && user.role >= 10
}

export function isRoot(user: any): boolean {
  return user && typeof user.role === 'number' && user.role >= 100
}
