-- 覺行小組長（2026-10-08 Geodown）：和減壓教練一樣是會員身分，後台會員名冊打勾設定。
-- 已經在某個小組當組長的會員，一開始就算覺行小組長。
ALTER TABLE volunteer ADD COLUMN is_group_leader boolean NOT NULL DEFAULT false;
UPDATE volunteer v SET is_group_leader = true
WHERE EXISTS (SELECT 1 FROM group_member m WHERE m.volunteer_id = v.id AND m.role = 'leader');
