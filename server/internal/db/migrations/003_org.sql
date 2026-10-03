-- 場域管理與共好企業管理（docs/data-model.md 的 center、venue、merchant）。
-- 鏈上地址、決議外鍵與額度在 M1–M3 加入；這裡先收組織資料。

CREATE TABLE center (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name         text NOT NULL UNIQUE,
  region       text NOT NULL DEFAULT '',
  address      text NOT NULL DEFAULT '',
  contact_name text NOT NULL DEFAULT '',
  email        text NOT NULL DEFAULT '',
  phone        text NOT NULL DEFAULT '',
  status       text NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended')),
  note         text NOT NULL DEFAULT '',
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);

-- 共好企業：聯盟單位（屬於某中心的餐廳、住宿…）或外部贊助商家
CREATE TABLE merchant (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  kind           text NOT NULL CHECK (kind IN ('alliance_unit', 'sponsor')),
  center_id      uuid REFERENCES center(id) ON DELETE RESTRICT,
  name           text NOT NULL,
  contact_name   text NOT NULL DEFAULT '',
  email          text NOT NULL DEFAULT '',
  phone          text NOT NULL DEFAULT '',
  region         text NOT NULL DEFAULT '',
  address        text NOT NULL DEFAULT '',
  offerings      text NOT NULL DEFAULT '',
  -- 企業自己提出的每月贊助規模（幣）；正式額度要主辦審核小組決議（M3）
  proposed_monthly_cap bigint CHECK (proposed_monthly_cap IS NULL OR proposed_monthly_cap >= 0),
  status         text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'suspended')),
  application_id uuid UNIQUE REFERENCES merchant_application(id) ON DELETE SET NULL,
  note           text NOT NULL DEFAULT '',
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now(),
  -- 聯盟單位必須屬於一個中心；贊助商家不屬於任何中心（data-model D49）
  CHECK ((kind = 'alliance_unit') = (center_id IS NOT NULL))
);

-- 場域：中心底下舉辦活動的地點（禪堂、營地…），活動與核發以場域為單位
CREATE TABLE venue (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  center_id      uuid NOT NULL REFERENCES center(id) ON DELETE RESTRICT,
  name           text NOT NULL,
  address        text NOT NULL DEFAULT '',
  description    text NOT NULL DEFAULT '',
  charges_public boolean NOT NULL DEFAULT false,
  merchant_id    uuid REFERENCES merchant(id) ON DELETE SET NULL,
  status         text NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'closed')),
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now(),
  UNIQUE (center_id, name)
);

ALTER TABLE practice_group ADD COLUMN center_id uuid REFERENCES center(id) ON DELETE SET NULL;
