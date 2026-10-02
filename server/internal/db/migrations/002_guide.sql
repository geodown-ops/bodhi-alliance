-- AI 組長自己的 schema；核心帳本表不在這裡，guide 服務也不讀寫它們。
CREATE SCHEMA guide;

-- kind：knowledge 一般知識／script 共修腳本（帶領共修時逐步使用）
CREATE TABLE guide.document (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  title        text NOT NULL,
  category     text NOT NULL CHECK (category IN ('practice', 'script', 'coin', 'association', 'faq')),
  status       text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
  current_version int NOT NULL DEFAULT 1,
  uploaded_by  uuid REFERENCES app_user(id) ON DELETE SET NULL,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);

-- 每次上傳或修改文字都新增一版，舊版保留供稽核
CREATE TABLE guide.document_version (
  document_id uuid NOT NULL REFERENCES guide.document(id) ON DELETE CASCADE,
  version     int NOT NULL,
  source_name text NOT NULL DEFAULT '',
  body        text NOT NULL,
  edited_by   uuid REFERENCES app_user(id) ON DELETE SET NULL,
  created_at  timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (document_id, version)
);

-- 只有上架文件的現行版本會切段，檢索只讀這張表
CREATE TABLE guide.chunk (
  document_id uuid NOT NULL REFERENCES guide.document(id) ON DELETE CASCADE,
  seq         int NOT NULL,
  heading     text NOT NULL DEFAULT '',
  body        text NOT NULL,
  PRIMARY KEY (document_id, seq)
);

-- 角色設定：名字、自我介紹、語氣與界線，每次儲存都是新版本，最新一版生效
CREATE TABLE guide.persona_version (
  version    serial PRIMARY KEY,
  name       text NOT NULL,
  prompt     text NOT NULL,
  edited_by  uuid REFERENCES app_user(id) ON DELETE SET NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

-- 用量與每月費用上限
CREATE TABLE guide.usage (
  id            bigserial PRIMARY KEY,
  model         text NOT NULL,
  input_tokens  bigint NOT NULL,
  output_tokens bigint NOT NULL,
  cache_read_tokens bigint NOT NULL DEFAULT 0,
  cache_write_tokens bigint NOT NULL DEFAULT 0,
  created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX guide_usage_month ON guide.usage(created_at);

CREATE TABLE guide.setting (
  key   text PRIMARY KEY,
  value text NOT NULL
);
INSERT INTO guide.setting (key, value) VALUES ('monthly_budget_usd', '50'), ('chat_enabled', 'true');
