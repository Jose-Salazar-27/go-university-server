-- Migration: insert core users, faculties, departments, professors, and degrees
-- Order:
-- 1) users
-- 2) faculties (dean_id NULL)
-- 3) departments (head_id NULL)
-- 4) professors (reference departments)
-- 5) set departments.head_id and faculties.dean_id via updates
-- 6) degrees

-- 1) Users (professors + admin + sample students)
INSERT INTO users (email, password_hash, first_name, last_name, user_type, avatar_url, is_active)
VALUES
('john.doe@university.edu', '$2b$10$YourHashedPasswordHere', 'John', 'Doe', 'professor', 'https://example.com/avatars/john.jpg', true),
('jane.smith@university.edu', '$2b$10$YourHashedPasswordHere', 'Jane', 'Smith', 'professor', 'https://example.com/avatars/jane.jpg', true),
('robert.chen@university.edu', '$2b$10$YourHashedPasswordHere', 'Robert', 'Chen', 'professor', NULL, true),
('michael.johnson@university.edu', '$2b$10$YourHashedPasswordHere', 'Michael', 'Johnson', 'professor', NULL, true),
('sarah.wilson@university.edu', '$2b$10$YourHashedPasswordHere', 'Sarah', 'Wilson', 'professor', NULL, true),
('david.brown@university.edu', '$2b$10$YourHashedPasswordHere', 'David', 'Brown', 'professor', NULL, true),
('admin@university.edu', '$2b$10$YourHashedPasswordHere', 'System', 'Administrator', 'admin', NULL, true),
('student1@university.edu', '$2b$10$YourHashedPasswordHere', 'Alice', 'Johnson', 'student', NULL, true),
('student2@university.edu', '$2b$10$YourHashedPasswordHere', 'Bob', 'Williams', 'student', NULL, true)
ON CONFLICT (email) DO UPDATE SET
	first_name = EXCLUDED.first_name,
	last_name = EXCLUDED.last_name,
	user_type = EXCLUDED.user_type,
	avatar_url = EXCLUDED.avatar_url,
	is_active = EXCLUDED.is_active
;

-- 2) Faculties (dean_id left NULL; will set after professors exist)
INSERT INTO faculties (code, name, description)
VALUES
('FI', 'Faculty of Engineering', 'Engineering and technology education'),
('FM', 'Faculty of Medicine', 'Health sciences and human medicine'),
('FC', 'Faculty of Sciences', 'Basic and applied sciences'),
('FH', 'Faculty of Humanities', 'Social sciences and humanities'),
('FE', 'Faculty of Economics', 'Economic and administrative sciences')
ON CONFLICT (code) DO UPDATE SET
	name = EXCLUDED.name,
	description = EXCLUDED.description
;

-- 3) Departments (head_id NULL for now)
INSERT INTO departments (faculty_id, code, name)
VALUES
((SELECT id FROM faculties WHERE code = 'FI'), 'CS', 'Computer Science'),
((SELECT id FROM faculties WHERE code = 'FI'), 'CE', 'Civil Engineering'),
((SELECT id FROM faculties WHERE code = 'FI'), 'IE', 'Industrial Engineering'),
((SELECT id FROM faculties WHERE code = 'FI'), 'EE', 'Electrical Engineering'),

((SELECT id FROM faculties WHERE code = 'FM'), 'ANAT', 'Anatomy'),
((SELECT id FROM faculties WHERE code = 'FM'), 'PEDI', 'Pediatrics'),
((SELECT id FROM faculties WHERE code = 'FM'), 'SURG', 'Surgery'),

((SELECT id FROM faculties WHERE code = 'FC'), 'MATH', 'Mathematics'),
((SELECT id FROM faculties WHERE code = 'FC'), 'PHYS', 'Physics'),
((SELECT id FROM faculties WHERE code = 'FC'), 'CHEM', 'Chemistry'),

((SELECT id FROM faculties WHERE code = 'FH'), 'PHIL', 'Philosophy'),
((SELECT id FROM faculties WHERE code = 'FH'), 'HIST', 'History'),
((SELECT id FROM faculties WHERE code = 'FH'), 'LIT', 'Literature'),

((SELECT id FROM faculties WHERE code = 'FE'), 'ECON', 'Economics'),
((SELECT id FROM faculties WHERE code = 'FE'), 'ADM', 'Administration')
ON CONFLICT (faculty_id, code) DO UPDATE SET
	name = EXCLUDED.name
;

-- 4) Professors (assign to departments by code)
-- We set a stable `professor_id` value to identify them for later updates
INSERT INTO professors (id, professor_id, department_id, hire_date, academic_title, office_location)
VALUES
((SELECT id FROM users WHERE email = 'john.doe@university.edu'), 'PROF001', (SELECT id FROM departments WHERE code = 'CS'), '2020-01-15', 'PhD in Computer Science', 'Room 101, CS Building'),
((SELECT id FROM users WHERE email = 'jane.smith@university.edu'), 'PROF002', (SELECT id FROM departments WHERE code = 'CS'), '2018-03-20', 'PhD in Software Engineering', 'Room 102, CS Building'),
((SELECT id FROM users WHERE email = 'robert.chen@university.edu'), 'PROF003', (SELECT id FROM departments WHERE code = 'CS'), '2015-08-10', 'PhD in Artificial Intelligence', 'Room 103, CS Building'),
((SELECT id FROM users WHERE email = 'michael.johnson@university.edu'), 'PROF004', (SELECT id FROM departments WHERE code = 'MATH'), '2019-06-01', 'PhD in Mathematics', 'Room 201, Math Building'),
((SELECT id FROM users WHERE email = 'sarah.wilson@university.edu'), 'PROF005', (SELECT id FROM departments WHERE code = 'EE'), '2017-09-15', 'PhD in Electrical Engineering', 'Room 301, EE Building'),
((SELECT id FROM users WHERE email = 'david.brown@university.edu'), 'PROF006', (SELECT id FROM departments WHERE code = 'ADM'), '2021-01-10', 'MBA, PhD in Management', 'Room 401, Business Building')
ON CONFLICT (id) DO UPDATE SET
	professor_id = EXCLUDED.professor_id,
	department_id = EXCLUDED.department_id,
	hire_date = EXCLUDED.hire_date,
	academic_title = EXCLUDED.academic_title,
	office_location = EXCLUDED.office_location
;

-- 5) Populate department heads and faculty deans using the professors inserted above
-- Assign department heads where we have matching professors
UPDATE departments SET head_id = p.id
FROM professors p
WHERE departments.code = 'CS' AND p.professor_id = 'PROF001';

UPDATE departments SET head_id = p.id
FROM professors p
WHERE departments.code = 'MATH' AND p.professor_id = 'PROF004';

UPDATE departments SET head_id = p.id
FROM professors p
WHERE departments.code = 'EE' AND p.professor_id = 'PROF005';

UPDATE departments SET head_id = p.id
FROM professors p
WHERE departments.code = 'ADM' AND p.professor_id = 'PROF006';

-- Assign some faculty deans (example mapping)
UPDATE faculties SET dean_id = p.id
FROM professors p
WHERE faculties.code = 'FI' AND p.professor_id = 'PROF001';

UPDATE faculties SET dean_id = p.id
FROM professors p
WHERE faculties.code = 'FC' AND p.professor_id = 'PROF004';

UPDATE faculties SET dean_id = p.id
FROM professors p
WHERE faculties.code = 'FE' AND p.professor_id = 'PROF006';

-- 6) Degrees (careful to reference department ids)
INSERT INTO degrees (department_id, code, name, degree_type, total_credits, duration_years)
VALUES
((SELECT id FROM departments WHERE code = 'CS' AND faculty_id = (SELECT id FROM faculties WHERE code = 'FI')), 'CS-BS', 'Bachelor of Computer Science', 'bachelor', 240, 4),
((SELECT id FROM departments WHERE code = 'CS' AND faculty_id = (SELECT id FROM faculties WHERE code = 'FI')), 'CS-MS', 'Master of Computer Science', 'master', 120, 2),
((SELECT id FROM departments WHERE code = 'CS' AND faculty_id = (SELECT id FROM faculties WHERE code = 'FI')), 'CS-PHD', 'PhD in Computer Science', 'phd', 180, 4),

((SELECT id FROM departments WHERE code = 'CE' AND faculty_id = (SELECT id FROM faculties WHERE code = 'FI')), 'CE-BS', 'Bachelor of Civil Engineering', 'bachelor', 260, 5),

((SELECT id FROM departments WHERE code = 'ANAT' AND faculty_id = (SELECT id FROM faculties WHERE code = 'FM')), 'MED-BS', 'Bachelor of Medicine', 'bachelor', 350, 6),

((SELECT id FROM departments WHERE code = 'MATH' AND faculty_id = (SELECT id FROM faculties WHERE code = 'FC')), 'MATH-BS', 'Bachelor of Mathematics', 'bachelor', 220, 4),

((SELECT id FROM departments WHERE code = 'ECON' AND faculty_id = (SELECT id FROM faculties WHERE code = 'FE')), 'ECON-BS', 'Bachelor of Economics', 'bachelor', 230, 4),

((SELECT id FROM departments WHERE code = 'PHIL' AND faculty_id = (SELECT id FROM faculties WHERE code = 'FH')), 'PHIL-BS', 'Bachelor of Philosophy', 'bachelor', 210, 4),

((SELECT id FROM departments WHERE code = 'ADM' AND faculty_id = (SELECT id FROM faculties WHERE code = 'FE')), 'ADM-BS', 'Bachelor of Business Administration', 'bachelor', 225, 4),

((SELECT id FROM departments WHERE code = 'EE' AND faculty_id = (SELECT id FROM faculties WHERE code = 'FI')), 'EE-BS', 'Bachelor of Electrical Engineering', 'bachelor', 255, 4)
ON CONFLICT (department_id, code) DO UPDATE SET
	name = EXCLUDED.name,
	degree_type = EXCLUDED.degree_type,
	total_credits = EXCLUDED.total_credits,
	duration_years = EXCLUDED.duration_years
;

-- End of data migration

