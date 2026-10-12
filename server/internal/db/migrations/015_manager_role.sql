-- 管理員（manager）：會員名冊打勾授予，可用後台大部分功能，但不能管理帳號、審核菩提幣或指派管理員。
ALTER TABLE role_assignment DROP CONSTRAINT role_assignment_role_check;
ALTER TABLE role_assignment ADD CONSTRAINT role_assignment_role_check
  CHECK (role IN ('alliance_admin', 'knowledge_manager', 'center_admin', 'manager'));
