-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url
(
    alias      VARCHAR(10) PRIMARY KEY,
    original   VARCHAR(255) NOT NULL,
    expires_at DATE         NOT NULL,
    visits     INTEGER
);

CREATE INDEX IF NOT EXISTS idx_original ON url(original);

CREATE OR REPLACE FUNCTION update_expiration()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.visits > 100 THEN
        NEW.expires_at := CURRENT_DATE + INTERVAL '1 month';
        NEW.visits := 0;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_expiration_trigger ON public.url;
CREATE TRIGGER update_expiration_trigger
    BEFORE UPDATE OF visits ON url
    FOR EACH ROW
    WHEN (OLD.visits IS DISTINCT FROM NEW.visits)
EXECUTE FUNCTION update_expiration();
-- +goose StatementEnd
