-- Insert sample email logs for testing
INSERT INTO email_logs (id, to_email, from_email, subject, body, is_html, status, created_at, updated_at) 
VALUES 
('sample-email-1', 'test@example.com', 'noreply@example.com', 'Welcome to our service', 'Welcome! Thank you for joining our service.', false, 'sent', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sample-email-2', 'user@example.com', 'noreply@example.com', 'Password Reset', '<p>Click <a href="#" >here</a> to reset your password.</p>', true, 'sent', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;