-- Insert core stat types
INSERT INTO stat_types (name, display_name, min_value, max_value, default_value) VALUES
('strength', 'Strength', 1, 10, 1),
('perception', 'Perception', 1, 10, 1),
('endurance', 'Endurance', 1, 10, 1),
('charisma', 'Charisma', 1, 10, 1),
('intelligence', 'Intelligence', 1, 10, 1),
('agility', 'Agility', 1, 10, 1),
('luck', 'Luck', 1, 10, 1),
('free-points', 'Free Points', 0, 99, 0),
('total-free-points', 'Total FP', 0, 99, 0);