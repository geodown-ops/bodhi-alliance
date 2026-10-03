-- 階段 A：帳號、角色、覺行小組、報名與企業登記。
-- 帳本、券、決議等表在 M2 之後依 docs/data-model.md 加入。

CREATE TABLE app_user (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email         text NOT NULL UNIQUE,
  display_name  text NOT NULL,
  password_hash text NOT NULL,
  disabled_at   timestamptz,
  created_at    timestamptz NOT NULL DEFAULT now()
);

-- scope：alliance／guide（階段 A），之後加 center:{id}、venue:{id}、merchant:{id}、group:{id}
CREATE TABLE role_assignment (
  user_id    uuid NOT NULL REFERENCES app_user(id) ON DELETE CASCADE,
  role       text NOT NULL CHECK (role IN ('alliance_admin', 'knowledge_manager')),
  scope      text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, role, scope)
);

CREATE TABLE user_session (
  token_hash bytea PRIMARY KEY,
  user_id    uuid NOT NULL REFERENCES app_user(id) ON DELETE CASCADE,
  expires_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX user_session_user ON user_session(user_id);

-- 覺行小組：中心底下的共修小組（docs/architecture.md〈待確認〉第 2 點）
CREATE TABLE practice_group (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name        text NOT NULL,
  region      text NOT NULL,
  center_name text NOT NULL DEFAULT '',
  schedule    text NOT NULL DEFAULT '',
  description text NOT NULL DEFAULT '',
  is_online   boolean NOT NULL DEFAULT false,
  is_listed   boolean NOT NULL DEFAULT true,
  sort_order  int NOT NULL DEFAULT 0,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE group_application (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  group_id    uuid REFERENCES practice_group(id) ON DELETE SET NULL,
  name        text NOT NULL,
  email       text NOT NULL,
  phone       text NOT NULL DEFAULT '',
  region      text NOT NULL DEFAULT '',
  wants_coach boolean NOT NULL DEFAULT false,
  message     text NOT NULL DEFAULT '',
  status      text NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'contacted', 'accepted', 'declined')),
  admin_note  text NOT NULL DEFAULT '',
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

-- 共好企業登記（招募頁的登記表移到這裡）；上架與額度要等 M3 主辦審核小組決議
CREATE TABLE merchant_application (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  kind           text NOT NULL CHECK (kind IN ('center', 'sponsor', 'other')),
  org_name       text NOT NULL,
  contact_name   text NOT NULL,
  email          text NOT NULL,
  phone          text NOT NULL DEFAULT '',
  region         text NOT NULL DEFAULT '',
  offerings      text NOT NULL DEFAULT '',
  monthly_scale  text NOT NULL DEFAULT '',
  message        text NOT NULL DEFAULT '',
  status         text NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'contacted', 'accepted', 'declined')),
  admin_note     text NOT NULL DEFAULT '',
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now()
);
