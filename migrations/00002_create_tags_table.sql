-- +goose Up
CREATE TABLE tags (
	id UUID PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	hex_color VARCHAR(7) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE tags;
