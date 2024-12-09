DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'question_type') THEN
        CREATE TYPE question_type AS ENUM ('basics', 'sorting');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS questions_quiz (
    id SERIAL PRIMARY KEY,
    quiz_id INT NOT NULL,
    creator_question_quiz_id UUID NOT NULL,
    type question_type NOT NULL, 
    question VARCHAR(255) NOT NULL,
    answers JSONB NOT NULL,  
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (quiz_id) REFERENCES quiz(id) ON DELETE CASCADE
);
