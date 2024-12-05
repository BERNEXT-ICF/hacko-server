CREATE TABLE IF NOT EXISTS modules (
    id SERIAL PRIMARY KEY, -- Kolom ID dengan auto increment
    creator_materials_id UUID NOT NULL,
    materials_id INT NOT NULL, --Kolom 
    title VARCHAR(255) NOT NULL, -- Kolom untuk judul material
    content TEXT, -- Kolom untuk judul material
    attachments TEXT[], -- Array string untuk lampiran
    videos TEXT[], -- Array string untuk video
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL, -- Kolom waktu pembuatan
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL, -- Kolom waktu pembaruan
    FOREIGN KEY (creator_materials_id) REFERENCES users(id) ON DELETE CASCADE, -- Foreign key ke tabel `users`
    FOREIGN KEY (materials_id) REFERENCES materials(id) ON DELETE CASCADE -- Foreign key ke tabel `users`
);
