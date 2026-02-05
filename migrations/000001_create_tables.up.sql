-- Ensure `gen_random_uuid()` is available
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    user_type VARCHAR(20) NOT NULL CHECK (user_type IN ('student', 'professor', 'admin')),
    avatar_url VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    student_id VARCHAR(50) UNIQUE,
    enrollment_date DATE NOT NULL,
    graduation_date DATE,
    current_status VARCHAR(20) DEFAULT 'active'
        CHECK (current_status IN ('active', 'graduated', 'suspended', 'dropout'))
);

CREATE TABLE IF NOT EXISTS faculties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    dean_id UUID, -- kept nullable to avoid cyclic dependency during initial migration
    -- Additional denormalized/decorative columns for data imports
    dean_name VARCHAR(200),
    location VARCHAR(200),
    contact_email VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    faculty_id UUID REFERENCES faculties(id) ON DELETE CASCADE,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(200) NOT NULL,
    head_id UUID,
    -- Optional descriptive fields to support legacy/import data
    description TEXT,
    head_name VARCHAR(200),
    office_location VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(faculty_id, code)
);

CREATE TABLE IF NOT EXISTS professors (
    id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    professor_id VARCHAR(50) UNIQUE,
    department_id UUID REFERENCES departments(id),
    hire_date DATE NOT NULL,
    academic_title VARCHAR(100), -- PhD, MSc, etc.
    office_location VARCHAR(100),
    office_hours TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Note: FK from faculties.dean_id -> professors(id) is commented out to avoid
-- migration ordering/cycle problems. Add it after inserting professors, e.g.:
ALTER TABLE faculties ADD CONSTRAINT fk_faculties_dean FOREIGN KEY (dean_id) REFERENCES professors(id);

CREATE TABLE IF NOT EXISTS degrees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    department_id UUID REFERENCES departments(id),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(200) NOT NULL,
    degree_type VARCHAR(50) CHECK (degree_type IN ('bachelor', 'master', 'phd', 'associate')),
    total_credits INTEGER NOT NULL,
    short_name VARCHAR(100),
    duration_semesters INTEGER,
    duration_years INTEGER,
    modality VARCHAR(50),
    coordinator_name VARCHAR(200),
    is_active BOOLEAN DEFAULT true,
    UNIQUE(department_id, code),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    credits INTEGER NOT NULL CHECK (credits > 0),
    department_id UUID REFERENCES departments(id),
    prerequisites TEXT, 
    theoretical_hours INTEGER,
    practical_hours INTEGER,
    level INTEGER,
    course_type VARCHAR(50),
    is_core BOOLEAN DEFAULT false,
    is_elective BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(code, department_id) -- the code could be repetead between departments
);

-- Table to store explicit course prerequisites (relational)
CREATE TABLE IF NOT EXISTS course_prerequisites (
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    prerequisite_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    is_mandatory BOOLEAN DEFAULT true,
    minimum_grade DECIMAL(4,2),
    PRIMARY KEY (course_id, prerequisite_id)
);

-- Table to represent elective groups for degrees
CREATE TABLE IF NOT EXISTS degree_electives (
    degree_id UUID REFERENCES degrees(id) ON DELETE CASCADE,
    elective_group VARCHAR(100) NOT NULL,
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    credits_awarded INTEGER,
    PRIMARY KEY (degree_id, elective_group, course_id)
);

CREATE TABLE IF NOT EXISTS degree_courses (
    degree_id UUID REFERENCES degrees(id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    semester INTEGER CHECK (semester BETWEEN 1 AND 12),
    is_required BOOLEAN DEFAULT true,
    credits_required INTEGER,
    PRIMARY KEY (degree_id, course_id)
);

CREATE TABLE IF NOT EXISTS academic_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    year INTEGER NOT NULL,
    period VARCHAR(20) NOT NULL CHECK (period IN ('fall', 'spring', 'summer', 'winter')),
    name VARCHAR(100) NOT NULL, -- Example: "Fall 2024"
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT false,
    registration_start DATE,
    registration_end DATE,
    UNIQUE(year, period)
);

CREATE TABLE IF NOT EXISTS course_sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    academic_period_id UUID REFERENCES academic_periods(id),
    section_code VARCHAR(10) NOT NULL, -- Ej: "01", "A", "VIRTUAL"
    professor_id UUID REFERENCES professors(id),
    schedule JSONB, -- {days: ["Mon", "Wed"], start_time: "14:00", end_time: "15:30", room: "A-101"}
    capacity INTEGER DEFAULT 30,
    modality VARCHAR(20) DEFAULT 'in_person' 
        CHECK (modality IN ('in_person', 'virtual', 'hybrid')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(course_id, academic_period_id, section_code)
);

CREATE TABLE IF NOT EXISTS enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    course_section_id UUID REFERENCES course_sections(id) ON DELETE CASCADE,
    enrollment_date DATE DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'enrolled' 
        CHECK (status IN ('enrolled', 'dropped', 'completed', 'failed')),
    final_grade DECIMAL(4,2) CHECK (final_grade BETWEEN 0 AND 100),
    letter_grade VARCHAR(2), -- A, B+, C, etc.
    credits_earned INTEGER,
    UNIQUE(student_id, course_section_id)
);

CREATE INDEX idx_enrollments_student ON enrollments(student_id);
CREATE INDEX idx_enrollments_section ON enrollments(course_section_id);

CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_section_id UUID REFERENCES course_sections(id) ON DELETE CASCADE,
    professor_id UUID REFERENCES professors(id),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    assignment_type VARCHAR(50) CHECK (assignment_type IN ('homework', 'project', 'exam', 'quiz', 'essay')),
    max_points DECIMAL(5,2) NOT NULL,
    due_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    is_published BOOLEAN DEFAULT false,
    attachments JSONB -- URLs de archivos: [{name: "instrucciones.pdf", url: "...", type: "pdf"}]
);

CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id UUID REFERENCES assignments(id) ON DELETE CASCADE,
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    submitted_at TIMESTAMPTZ DEFAULT NOW(),
    content TEXT, -- Texto de la entrega
    attachments JSONB, -- Archivos adjuntos del estudiante
    status VARCHAR(20) DEFAULT 'submitted' 
        CHECK (status IN ('submitted', 'late', 'graded', 'missing')),
    points_earned DECIMAL(5,2),
    grade_percentage DECIMAL(5,2),
    UNIQUE(assignment_id, student_id)
);

CREATE TABLE IF NOT EXISTS submission_feedbacks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID REFERENCES submissions(id) ON DELETE CASCADE,
    professor_id UUID REFERENCES professors(id),
    comment TEXT NOT NULL,
    rubric_evaluation JSONB, -- Evaluación por rúbrica si aplica
    created_at TIMESTAMPTZ DEFAULT NOW(),
    is_draft BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS final_grades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_id UUID REFERENCES enrollments(id) ON DELETE CASCADE,
    professor_id UUID REFERENCES professors(id),
    final_score DECIMAL(5,2) CHECK (final_score BETWEEN 0 AND 100),
    letter_grade VARCHAR(2),
    comments TEXT,
    published_at TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT false -- Si requiere aprobación del departamento
);

-- Additional helpful indexes for joins/filters
CREATE INDEX IF NOT EXISTS idx_departments_faculty_id ON departments(faculty_id);
CREATE INDEX IF NOT EXISTS idx_courses_department_id ON courses(department_id);
CREATE INDEX IF NOT EXISTS idx_degrees_department_id ON degrees(department_id);
CREATE INDEX IF NOT EXISTS idx_professors_department_id ON professors(department_id);
