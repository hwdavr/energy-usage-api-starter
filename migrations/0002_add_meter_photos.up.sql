-- Add table for storing meter photos (faulty/fixed)
CREATE TABLE IF NOT EXISTS meter_photos (
    id SERIAL PRIMARY KEY,
    meter_id UUID NOT NULL REFERENCES meters(id),
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(16) NOT NULL, -- 'faulty' or 'fixed'
    path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);