CREATE TABLE recaps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slug VARCHAR(64) NOT NULL UNIQUE,
    week_start DATE NOT NULL,
    week_end DATE NOT NULL,
    stats JSONB NOT NULL,
    narrative TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, week_start)
);

CREATE INDEX idx_recaps_slug ON recaps(slug);
