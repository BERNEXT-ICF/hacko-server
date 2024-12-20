-- Create table 'kelas' if not exists with a foreign key to 'users' table
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'visibility') THEN
        CREATE TYPE visibility AS ENUM ('public', 'draf');
    END IF;
END
$$;

-- Create table 'kelas' if not exists with a foreign key to 'users' table
CREATE TABLE IF NOT EXISTS class (
    id SERIAL PRIMARY KEY,
    creator_class_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image TEXT,
    video TEXT,
    status visibility NOT NULL DEFAULT 'draf',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Foreign key constraint to 'users' table
    CONSTRAINT fk_creatorclass FOREIGN KEY (creator_class_Id) REFERENCES users(id) ON DELETE CASCADE
);
