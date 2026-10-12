-- 使用問題：Sunny 答不了的網站使用問題、或訪客說遇到的狀況，由 Sunny 記下來，
-- 超級管理員在後台回答；登入的會員在個人頁看得到回答，沒登入的訪客可以留 email。
CREATE TABLE guide.site_issue (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  conversation  text NOT NULL DEFAULT '',          -- 前台一次對話的代號，用來補上 email
  user_id       uuid REFERENCES app_user(id) ON DELETE SET NULL,
  contact_email text NOT NULL DEFAULT '',
  summary       text NOT NULL,                     -- Sunny 整理的問題摘要
  transcript    jsonb NOT NULL DEFAULT '[]',       -- 記錄當下最近幾句對話
  status        text NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'answered')),
  answer        text NOT NULL DEFAULT '',
  answered_by   uuid REFERENCES app_user(id) ON DELETE SET NULL,
  answered_at   timestamptz,
  created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX site_issue_open ON guide.site_issue(status, created_at DESC);
CREATE INDEX site_issue_user ON guide.site_issue(user_id, created_at DESC);
CREATE INDEX site_issue_conversation ON guide.site_issue(conversation) WHERE conversation <> '';
