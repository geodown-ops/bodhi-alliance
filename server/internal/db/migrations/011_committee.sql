-- 菩提幣決策小組成員（2026-10-08 Geodown）：會員身分，後台會員名冊打勾設定。
-- 目前只是標記，不附帶後台權限（菩提幣審核仍由超級管理員操作）。
ALTER TABLE volunteer ADD COLUMN is_committee boolean NOT NULL DEFAULT false;
