-- Insert sample users for testing
-- Password for all users is "password123" (bcrypt hashed)
INSERT INTO users (id, email, password, first_name, last_name, roles, is_active, created_at, updated_at) 
VALUES 
('admin-user-1', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Admin', 'User', '{"admin", "user"}', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('regular-user-1', 'user@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Regular', 'User', '{"user"}', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('test-user-1', 'test@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Test', 'User', '{"user"}', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;