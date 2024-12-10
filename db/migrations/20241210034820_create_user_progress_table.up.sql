DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'progress_status') THEN
        CREATE TYPE progress_type AS ENUM ('not_started', 'on_progress', 'done');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS users_progress(
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    class_id INT NOT NULL,
    users_classes_id INT NOT NULL,
    material_id INT NOT NULL,
    module_id INT NOT NULL,
    quiz_id INT,
    progress DECIMAL,
    status progress_status DEFAULT 'not_started' NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quiz(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (class_id) REFERENCES class(id) ON DELETE CASCADE,
    FOREIGN KEY (users_classes_id) REFERENCES users_classes(id) ON DELETE CASCADE,
    FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE,
    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE
);
