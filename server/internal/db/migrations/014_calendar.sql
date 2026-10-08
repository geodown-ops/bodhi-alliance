-- 會員行事曆改用 dengo 的行程管理設計（2026-10-08 Geodown）：每筆行程有類型（決定行事曆上的顏色），
-- 可以掛一個活動標籤（像 hashtag，用來分類與篩選）。刪掉標籤不會刪行程，只是拿掉標籤。
CREATE TABLE association_event_tag (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name       text NOT NULL UNIQUE,
  created_at timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE association_event
  ADD COLUMN kind text NOT NULL DEFAULT 'general' CHECK (kind IN ('general', 'training', 'meeting', 'activity', 'ceremony')),
  ADD COLUMN tag_id uuid REFERENCES association_event_tag(id) ON DELETE SET NULL;
