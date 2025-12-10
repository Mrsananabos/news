-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS news (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL
    );

CREATE TABLE IF NOT EXISTS news_categories (
    news_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    PRIMARY KEY (news_id, category_id),
    CONSTRAINT fk_news FOREIGN KEY (news_id) REFERENCES news(id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS news_categories;
DROP TABLE IF EXISTS news;
-- +goose StatementEnd
