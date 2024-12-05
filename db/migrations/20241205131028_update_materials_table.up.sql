CREATE TABLE IF NOT EXISTS materials (
    id SERIAL PRIMARY KEY, -- Kolom ID dengan auto increment
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL, -- Kolom untuk judul material
    class_id INT NOT NULL, -- Kolom untuk ID kelas, foreign key ke tabel `classes`
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL, -- Kolom waktu pembuatan
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL, -- Kolom waktu pembaruan
    FOREIGN KEY (class_id) REFERENCES class(id) ON DELETE CASCADE, -- Foreign key ke tabel `classes`
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE -- Foreign key ke tabel `users`
);
