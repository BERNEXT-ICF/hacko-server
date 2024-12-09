DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'quiz_status') THEN
        CREATE TYPE quiz_status AS ENUM ('public', 'draft');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS quiz (
    id SERIAL PRIMARY KEY,
    class_id INT NOT NULL,
    creator_quiz_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    status quiz_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (class_id) REFERENCES class(id) ON DELETE CASCADE
);
