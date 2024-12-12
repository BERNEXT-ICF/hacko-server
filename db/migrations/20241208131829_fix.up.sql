DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'submission_status') THEN
        CREATE TYPE submission_status AS ENUM ('submitted', 'rated');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS assignments (
    id SERIAL PRIMARY KEY,
    class_id INT NOT NULL,
    creator_assignment_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (class_id) REFERENCES class(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS submissions (
    id SERIAL PRIMARY KEY,
    assignment_id INT NOT NULL,
    student_id UUID NOT NULL,
    link TEXT,
    grade INT,
    feedback TEXT,
    status submission_status DEFAULT 'submitted' NOT NULL, 
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    graded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE
);
