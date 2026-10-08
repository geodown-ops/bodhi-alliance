-- 菩提幣上鏈（2026-10-08 Geodown）：每位系統會員都有一個平台保管的鏈上地址，
-- 入會時金庫轉 1 枚菩提幣給他，交易可在公開的區塊瀏覽器查到。
-- 這修改了 SPEC D33「志工沒有個人鏈上地址」；先在測試鏈（Polygon Amoy）運作。

-- 鏈上設定：各鏈的合約地址、部署中的交易
CREATE TABLE chain_setting (
  key        text PRIMARY KEY,
  value      text NOT NULL,
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- 會員的鏈上地址；私鑰不存，由平台的種子金鑰依會員 id 推導
CREATE TABLE member_chain_account (
  volunteer_id uuid PRIMARY KEY REFERENCES volunteer(id) ON DELETE CASCADE,
  chain_id     bigint NOT NULL,
  address      text NOT NULL,
  created_at   timestamptz NOT NULL DEFAULT now(),
  UNIQUE (chain_id, address)
);

-- 鏈上撥幣：kind = join 是入會贈幣，每位會員在每條鏈上只有一筆
CREATE TABLE chain_grant (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  volunteer_id uuid NOT NULL REFERENCES volunteer(id) ON DELETE CASCADE,
  chain_id     bigint NOT NULL,
  kind         text NOT NULL CHECK (kind IN ('join')),
  amount       bigint NOT NULL CHECK (amount > 0), -- 最小單位，100 = 1 枚
  to_address   text NOT NULL,
  status       text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'confirmed', 'failed')),
  tx_hash      text,
  block_number bigint,
  last_error   text NOT NULL DEFAULT '',
  created_at   timestamptz NOT NULL DEFAULT now(),
  sent_at      timestamptz,
  confirmed_at timestamptz,
  UNIQUE (volunteer_id, chain_id, kind)
);
CREATE INDEX chain_grant_open ON chain_grant(status, created_at) WHERE status IN ('pending', 'sent');
