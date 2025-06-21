-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  events(
                                      id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                                      event_type TEXT NOT NULL,
                                      payload TEXT NOT NULL,
                                      status TEXT NOT NULL  DEFAULT  'new' CHECK ( status IN('new','processed','done') ),
                                      reserved_to TIMESTAMP NOT NULL ,
                                      created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
