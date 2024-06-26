-- Stores user records with a minimum of information
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    phoneNumber VARCHAR(15) NOT NULL
);
CREATE TABLE IF NOT EXISTS user_references (
    id VARCHAR(14) PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL,
    CONSTRAINT FK_REFERENCES FOREIGN KEY (user_id) REFERENCES user (id) NOT DEFERRABLE INITIALLY IMMEDIATE
);
-- Stores user notes, linking to their attachments
CREATE TABLE IF NOT EXISTS note (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    user_id INTEGER NOT NULL,
    details CLOB NOT NULL,
    CONSTRAINT FK_NOTE_USER FOREIGN KEY (user_id) REFERENCES user (id) NOT DEFERRABLE INITIALLY IMMEDIATE
);
-- Stores the note attachments
CREATE TABLE IF NOT EXISTS attachment (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    note_id INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    filename TEXT NOT NULL,
    file BLOB NOT NULL,
    CONSTRAINT FK_NOTE FOREIGN KEY (note_id) REFERENCES note (id) NOT DEFERRABLE INITIALLY IMMEDIATE
);
-- Add indexes to the three tables
CREATE UNIQUE INDEX IF NOT EXISTS UNIQ_USER_NAME ON user (name);
CREATE UNIQUE INDEX IF NOT EXISTS UNIQ_USER_EMAIL ON user (email);
CREATE UNIQUE INDEX IF NOT EXISTS UNIQ_USER_PHONE ON user (phoneNumber);
CREATE INDEX IDX_ATTACHMENT_NOTEID ON attachment (note_id);
CREATE INDEX IDX_NOTE_DETAILS ON note (details);
CREATE INDEX IDX_NOTE_USERID ON note (user_id);