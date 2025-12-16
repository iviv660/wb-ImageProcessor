CREATE TABLE images (
                        id BIGSERIAL PRIMARY KEY,
                        file_name TEXT NOT NULL,
                        format TEXT NOT NULL,
                        original_url TEXT NOT NULL,
                        changed_url TEXT,
                        created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


