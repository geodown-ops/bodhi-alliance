-- 系統會員（2026-10-06 Geodown）：覺行小組與世界佛教教育協會共用同一份會員資料（volunteer），
-- 會員自己選擇加入哪一邊、隨時可以切換。加入覺行小組才有菩提幣錢包與共修活動；
-- 加入協會才看得到協會會刊、會員行事曆與協會通知。
ALTER TABLE volunteer ADD COLUMN in_groups boolean NOT NULL DEFAULT true;
ALTER TABLE volunteer ADD COLUMN in_association boolean NOT NULL DEFAULT false;
ALTER TABLE volunteer ADD COLUMN association_joined_at timestamptz;

-- 協會會刊：一期一筆，連到線上版或 PDF
CREATE TABLE association_issue (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  title      text NOT NULL,
  issued_on  date NOT NULL,
  summary    text NOT NULL DEFAULT '',
  url        text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX association_issue_date ON association_issue(issued_on DESC);

-- 會員行事曆
CREATE TABLE association_event (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  title       text NOT NULL,
  starts_at   timestamptz NOT NULL,
  ends_at     timestamptz,
  location    text NOT NULL DEFAULT '',
  description text NOT NULL DEFAULT '',
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  CHECK (ends_at IS NULL OR ends_at > starts_at)
);
CREATE INDEX association_event_starts ON association_event(starts_at);

-- 協會通知：會員登入後在個人頁看到
CREATE TABLE association_notice (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  title      text NOT NULL,
  body       text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX association_notice_created ON association_notice(created_at DESC);
