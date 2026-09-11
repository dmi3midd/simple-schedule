-- +goose Up
CREATE TYPE day_of_week AS ENUM (
	'monday',
	'tuesday',
	'wednesday',
	'thursday',
	'friday',
	'saturday',
	'sunday'
);

CREATE TABLE slots (
	id UUID PRIMARY KEY,
	week_id UUID NOT NULL REFERENCES weeks(id) ON DELETE CASCADE,
	tag_id UUID REFERENCES tags(id) ON DELETE SET NULL,
	activity VARCHAR(255) NOT NULL,
	day_of_week day_of_week NOT NULL,
	start_time INTEGER NOT NULL,
	end_time INTEGER NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS slots;
DROP TYPE IF EXISTS day_of_week;
