--up
ALTER TABLE jobs ADD COLUMN salary TEXT DEFAULT 'Not disclosed';
ALTER TABLE jobs ADD COLUMN location TEXT DEFAULT 'Remote';

