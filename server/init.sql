CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,                    -- Автоинкрементный ID
    title TEXT NOT NULL,             -- Заголовок заметки
    content TEXT,                             -- Содержимое заметки
    created_at TIMESTAMP DEFAULT NOW()       -- Дата создания
);