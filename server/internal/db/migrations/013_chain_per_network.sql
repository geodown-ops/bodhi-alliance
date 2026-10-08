-- 菩提幣正式上鏈（2026-10-08 Geodown）：同一位會員在測試鏈與 Polygon 主網各有一個地址，
-- 測試鏈的紀錄保留作為歷史，個人頁只顯示目前這條鏈。
ALTER TABLE member_chain_account DROP CONSTRAINT member_chain_account_pkey;
ALTER TABLE member_chain_account ADD PRIMARY KEY (volunteer_id, chain_id);
