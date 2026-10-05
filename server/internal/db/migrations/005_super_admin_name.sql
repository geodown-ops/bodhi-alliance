-- 後台角色「聯盟管理員」改名為「超級管理員」（角色代碼仍是 alliance_admin）。
-- 第一個管理員帳號是用預設名稱建立的，名稱一起改；自己改過名稱的帳號不動。

UPDATE app_user SET display_name = '超級管理員' WHERE display_name = '聯盟管理員';
