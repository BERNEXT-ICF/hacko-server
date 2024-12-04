-- Membuat tipe ENUM untuk enrollment_status
-- Create enum type for enrollment status
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'enrollment_status') THEN
        CREATE TYPE enrollment_status AS ENUM ('active', 'completed', 'dropped', 'removed');
    END IF;
END
$$;

-- Membuat tabel users_classes
CREATE TABLE IF NOT EXISTS users_classes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    class_id INT NOT NULL,
    enrollment_status enrollment_status DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (class_id) REFERENCES classes (id) ON DELETE CASCADE
);
