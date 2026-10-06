-- 覺行共修活動與菩提幣核發（2026-10-06 Geodown）。
-- 任何人在覺行小組頁報名（建立帳號）後，可以發起三人以上的共修活動；活動結束後發起人送審，
-- 超級管理員（代表菩提幣決策小組）核准後，發起人與協辦志工的菩提幣入帳。

-- 報名只要真實姓名、電子郵件、密碼；所屬中心改成選填，另外可以留 LINE ID。
ALTER TABLE volunteer ALTER COLUMN home_center_id DROP NOT NULL;
ALTER TABLE volunteer ADD COLUMN line_id text NOT NULL DEFAULT '';

CREATE TABLE practice_event (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  organizer_id uuid NOT NULL REFERENCES volunteer(id) ON DELETE RESTRICT,
  title        text NOT NULL,
  is_online    boolean NOT NULL DEFAULT false,
  -- 線下是地址，線上是會議連結或說明
  location     text NOT NULL DEFAULT '',
  starts_at    timestamptz NOT NULL,
  ends_at      timestamptz NOT NULL,
  -- 開放人數（含發起人），至少三人
  capacity     int NOT NULL CHECK (capacity BETWEEN 3 AND 500),
  description  text NOT NULL DEFAULT '',
  status       text NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'cancelled')),
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now(),
  CHECK (ends_at > starts_at)
);
CREATE INDEX practice_event_starts ON practice_event(status, starts_at);

-- organizer：發起人；helper：協辦志工或減壓教練；participant：一般參加者
CREATE TABLE event_participant (
  event_id     uuid NOT NULL REFERENCES practice_event(id) ON DELETE CASCADE,
  volunteer_id uuid NOT NULL REFERENCES volunteer(id) ON DELETE CASCADE,
  role         text NOT NULL DEFAULT 'participant' CHECK (role IN ('organizer', 'helper', 'participant')),
  joined_at    timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (event_id, volunteer_id)
);
CREATE INDEX event_participant_volunteer ON event_participant(volunteer_id);

-- 菩提幣審核申請：一場活動一筆，退回後發起人可以修改再送
CREATE TABLE coin_claim (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id     uuid NOT NULL UNIQUE REFERENCES practice_event(id) ON DELETE RESTRICT,
  submitted_by uuid NOT NULL REFERENCES volunteer(id) ON DELETE RESTRICT,
  attendance   int NOT NULL CHECK (attendance >= 3),
  report       text NOT NULL DEFAULT '',
  status       text NOT NULL DEFAULT 'submitted' CHECK (status IN ('submitted', 'approved', 'rejected')),
  review_note  text NOT NULL DEFAULT '',
  reviewed_by  uuid REFERENCES app_user(id) ON DELETE SET NULL,
  reviewed_at  timestamptz,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX coin_claim_status ON coin_claim(status, created_at);

-- 這次申請要核發給誰、多少（送審時依梯級表建議，核准時可以調整）
CREATE TABLE claim_recipient (
  claim_id     uuid NOT NULL REFERENCES coin_claim(id) ON DELETE CASCADE,
  volunteer_id uuid NOT NULL REFERENCES volunteer(id) ON DELETE RESTRICT,
  role         text NOT NULL CHECK (role IN ('organizer', 'helper')),
  amount       bigint NOT NULL CHECK (amount >= 0),
  PRIMARY KEY (claim_id, volunteer_id)
);

-- 錢包帳本：餘額 = 加總。目前只有核發；兌換與核銷在券上線時加入。
CREATE TABLE coin_ledger (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  volunteer_id uuid NOT NULL REFERENCES volunteer(id) ON DELETE RESTRICT,
  amount       bigint NOT NULL,
  kind         text NOT NULL CHECK (kind IN ('issue')),
  claim_id     uuid REFERENCES coin_claim(id) ON DELETE RESTRICT,
  memo         text NOT NULL DEFAULT '',
  created_at   timestamptz NOT NULL DEFAULT now(),
  UNIQUE (claim_id, volunteer_id)
);
CREATE INDEX coin_ledger_volunteer ON coin_ledger(volunteer_id, created_at);
