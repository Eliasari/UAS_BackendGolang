CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO roles (name, description)
VALUES 
    ('Admin', 'Administrator sistem'),
    ('Mahasiswa', 'Mahasiswa pelapor prestasi'),
    ('Dosen Wali', 'Dosen pembimbing akademik');
    

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role_id UUID REFERENCES roles(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
select * from roles;
select * from users;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO users 
(id, username, email, password_hash, full_name, role_id, is_active)
VALUES
(
    uuid_generate_v4(),
    'admin1',
    'admin1@gmail.com',
    crypt('Password', gen_salt('bf')),
    'Administrator',
    (SELECT id FROM roles WHERE name = 'Admin'),
    true
),
(
    uuid_generate_v4(),
    'mahasiswa123',
    'mahasiswa123@gmail.com',
    crypt('Password', gen_salt('bf')),
    'John Doe',
    (SELECT id FROM roles WHERE name = 'Mahasiswa'),
    true
),
(
    uuid_generate_v4(),
    'dosenwali1',
    'dosenwali1@gmail.com',
    crypt('Password', gen_salt('bf')),
    'Budi Santoso',
    (SELECT id FROM roles WHERE name = 'Dosen Wali'),
    true
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT
);

select* from permissions;

INSERT INTO permissions (name, resource, action, description) VALUES
('achievement:create', 'achievement', 'create', 'Mahasiswa membuat prestasi'),
('achievement:read', 'achievement', 'read', 'Melihat prestasi'),
('achievement:update', 'achievement', 'update', 'Update prestasi'),
('achievement:delete', 'achievement', 'delete', 'Hapus prestasi'),
('achievement:verify', 'achievement', 'verify', 'Dosen wali memverifikasi'),
('user:manage', 'user', 'manage', 'Admin mengelola user');

CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

select * from role_permissions;

-- role_permissions admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Admin';

-- role_permissions mahasiswa
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'Mahasiswa'),
    p.id
FROM permissions p
WHERE name IN ('achievement:create','achievement:read','achievement:update','achievement:delete');

INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'Dosen Wali'),
    p.id
FROM permissions p
WHERE name IN ('achievement:read','achievement:verify');


select * from role_permissions;


