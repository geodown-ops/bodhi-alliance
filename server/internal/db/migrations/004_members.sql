-- M2 會員與組織：志工帳號、覺行小組成員、中心管理員（docs/data-model.md 的 volunteer）。
-- 志工自己註冊後是「待核可」；所屬中心的管理員核可實名後才能列入核發名單（M4）。

ALTER TABLE role_assignment DROP CONSTRAINT role_assignment_role_check;
ALTER TABLE role_assignment ADD CONSTRAINT role_assignment_role_check
  CHECK (role IN ('alliance_admin', 'knowledge_manager', 'center_admin'));
-- 中心管理員的 scope 是 center:{中心 id}
ALTER TABLE role_assignment ADD CONSTRAINT role_assignment_center_scope
  CHECK (role <> 'center_admin' OR scope ~ '^center:[0-9a-f-]{36}$');

CREATE TABLE volunteer (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         uuid NOT NULL UNIQUE REFERENCES app_user(id) ON DELETE CASCADE,
  home_center_id  uuid NOT NULL REFERENCES center(id) ON DELETE RESTRICT,
  legal_name      text NOT NULL,
  phone           text NOT NULL DEFAULT '',
  photo_url       text,
  -- 自己申請當禪修教練；is_coach 由中心管理員確認
  wants_coach     boolean NOT NULL DEFAULT false,
  is_coach        boolean NOT NULL DEFAULT false,
  status          text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'verified', 'rejected')),
  review_note     text NOT NULL DEFAULT '',
  verified_by     uuid REFERENCES app_user(id) ON DELETE SET NULL,
  verified_at     timestamptz,
  qr_frozen_at    timestamptz,
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now(),
  CHECK ((status = 'verified') = (verified_at IS NOT NULL))
);
CREATE INDEX volunteer_center ON volunteer(home_center_id, status);

CREATE TABLE group_member (
  group_id     uuid NOT NULL REFERENCES practice_group(id) ON DELETE CASCADE,
  volunteer_id uuid NOT NULL REFERENCES volunteer(id) ON DELETE CASCADE,
  role         text NOT NULL DEFAULT 'member' CHECK (role IN ('member', 'leader')),
  joined_at    timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (group_id, volunteer_id)
);
CREATE INDEX group_member_volunteer ON group_member(volunteer_id);
