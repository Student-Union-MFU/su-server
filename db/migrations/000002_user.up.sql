CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    user_type     VARCHAR(20) NOT NULL,

    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    avatar_url    VARCHAR(255),

    student_id    VARCHAR(30),
    major         VARCHAR(255),
    school        VARCHAR(255),
    
    oauth_subject VARCHAR(255) NOT NULL UNIQUE,
    
    is_flagged    BOOLEAN DEFAULT FALSE,
  
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW()
);
