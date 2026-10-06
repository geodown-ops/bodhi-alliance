-- 共修活動可以掛在場域底下（2026-10-06 Geodown）：覺行小組頁列出活動場域，點了看該場域的活動並報名。
-- 線上活動與自訂地點的活動沒有場域。
ALTER TABLE practice_event ADD COLUMN venue_id uuid REFERENCES venue(id) ON DELETE SET NULL;
CREATE INDEX practice_event_venue ON practice_event(venue_id) WHERE venue_id IS NOT NULL;
