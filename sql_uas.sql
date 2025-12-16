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
('achievement:detail', 'achievement', 'read', 'Melihat detail [prestasi]');

CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

select * from role_permissions;
select * from roles;
select * from achievement_references;
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
    (SELECT id FROM roles WHERE name = 'Admin'),
    p.id
FROM permissions p
WHERE name IN ('achievement:list');

-- acddc6df-1376-46b5-96bf-593ec3248e04

select * from permissions;
select * from role_permissions;

select * from achievement_references;

SELECT r.name AS role_name, p.name AS permission_name, p.id AS permission_id
FROM roles r
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE p.name = 'achievement:list';

-- insert permissions
INSERT INTO permissions (name, resource, action, description) VALUES
('achievement:list:self', 'achievement', 'list', 'Melihat daftar achievement student');

-- assign permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN (
    'achievement:list:self'
)
WHERE r.name = 'Mahasiswa';

-- check hasil permissions
SELECT r.name AS role_name, p.name AS permission_name
FROM roles r
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON p.id = rp.permission_id
ORDER BY r.name, p.name;

ALTER TYPE status ADD VALUE 'deleted';

SELECT id, student_id, status
FROM achievement_references
WHERE id = 'e6846390-0f33-4e5d-89fc-c221b345fda0';

select* from users;
select * from students;


SELECT id, student_id, status
FROM achievement_references
WHERE student_id::text = '8015efbe-6a87-4335-830d-a02214721ccc';

SELECT DISTINCT student_id FROM achievement_references;

-- check kesediaan data
SELECT id, student_id, status
FROM achievement_references
WHERE id = 'e6846390-0f33-4e5d-89fc-c221b345fda0';

select * from students;
select * from lecturers;

SELECT id, student_id, status
FROM achievement_references
WHERE id = '5f3c3af6-b9c1-40c4-9cd8-a4a7768df196';

SELECT id FROM students WHERE user_id = '8015efbe-6a87-4335-830d-a02214721ccc';

SELECT id, student_id, status
FROM achievement_references
WHERE id = '5f3c3af6-b9c1-40c4-9cd8-a4a7768df196';

SELECT
  id,
  student_id,
  status,
  length(status),
  status = 'draft' AS is_exact_draft
FROM achievement_references
WHERE id = '48425e70-70a9-439d-aad3-c7c7b31c6f25';

SELECT
    column_name,
    data_type,
    udt_name
FROM information_schema.columns
WHERE table_name = 'achievement_references';

SELECT DISTINCT status FROM achievement_references;

SELECT id, student_id, status
FROM achievement_references
WHERE id = '48425e70-70a9-439d-aad3-c7c7b31c6f25';

SELECT id, student_id, status
FROM achievement_references
WHERE id = '48425e70-70a9-439d-aad3-c7c7b31c6f25'
  AND student_id = '06060f6a-882d-4e6b-837f-ed9a09da7213';

UPDATE achievement_references
SET status = 'deleted',
    updated_at = NOW()
WHERE id = '48425e70-70a9-439d-aad3-c7c7b31c6f25'
  AND student_id = '06060f6a-882d-4e6b-837f-ed9a09da7213'
  AND status = 'draft'
RETURNING id, student_id, status;

ALTER TABLE achievement_references
DROP CONSTRAINT achievement_references_status_check;

ALTER TABLE achievement_references
ADD CONSTRAINT achievement_references_status_check
CHECK (status IN ('draft', 'submitted', 'verified', 'rejected', 'deleted'));

select * from users;
