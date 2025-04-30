DROP TABLE IF EXISTS task_categories;

ALTER TABLE tasks
    ADD COLUMN category_id UUID,
    ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT now();

ALTER TABLE tasks
    ADD CONSTRAINT fk_tasks_category
    FOREIGN KEY (category_id)
    REFERENCES categories(id)
    ON DELETE SET NULL;
