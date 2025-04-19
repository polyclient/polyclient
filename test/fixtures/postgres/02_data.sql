-- Insert planets
INSERT INTO planet (name, type, diameter_km, mass_kg, coordinates, has_atmosphere, has_moons, has_rings, has_oceans, habitability_status, description)
VALUES 
  ('Earth', 'terrestrial', 12742.00, 5.97e24, '{"x": 0, "y": 0, "z": 0}', true, true, false, true, 'habitable', 'Our home planet'),
  ('Mars', 'terrestrial', 6779.00, 6.42e23, '{"x": 1.5, "y": 0.2, "z": 0.1}', true, true, false, false, 'non_habitable', 'The red planet'),
  ('Jupiter', 'gaseous', 139820.00, 1.90e27, '{"x": 5.2, "y": 0.3, "z": 0.1}', true, true, true, false, 'non_habitable', 'Largest planet in our solar system'),
  ('Saturn', 'gaseous', 116460.00, 5.68e26, '{"x": 8.3, "y": 0.4, "z": 0.2}', true, true, true, false, 'non_habitable', 'Planet with beautiful rings'),
  ('Uranus', 'icy', 50724.00, 8.68e25, '{"x": 19.1, "y": 0.5, "z": 0.3}', true, true, true, false, 'non_habitable', 'Ice giant'),
  ('Proxima Centauri b', 'terrestrial', 12742.00, 5.97e24, '{"x": 40000, "y": 5000, "z": 2000}', true, false, false, false, 'potentially_habitable', 'Exoplanet orbiting Proxima Centauri'),
  ('Kepler-186f', 'terrestrial', 11542.00, 4.80e24, '{"x": 558000, "y": 12000, "z": 8000}', true, false, false, false, 'potentially_habitable', 'Potentially habitable exoplanet'),
  ('HD 189733 b', 'gaseous', 17600.00, 1.13e27, '{"x": 600000, "y": 10000, "z": 5000}', true, false, false, false, 'non_habitable', 'Hot Jupiter exoplanet'),
  ('55 Cancri e', 'terrestrial', 16500.00, 8.63e24, '{"x": 700000, "y": 11000, "z": 6000}', true, false, false, false, 'non_habitable', 'Super-Earth exoplanet'),
  ('TRAPPIST-1 e', 'terrestrial', 11300.00, 5.30e24, '{"x": 850000, "y": 13000, "z": 7000}', true, false, false, false, 'potentially_habitable', 'One of the seven Earth-sized planets orbiting TRAPPIST-1');

-- Insert astronauts
INSERT INTO astronaut (home_planet_id, first_name, last_name, specialty, birth_date, join_date, active, bio)
VALUES 
  (1, 'John', 'Doe', 'scientist', '1990-01-01', '2020-01-01', true, 'I am a scientist'),
  (1, 'Jane', 'Doe', 'engineer', '1995-05-05', '2020-01-01', true, 'I am an engineer'),
  (1, 'Bob', 'Smith', 'pilot', '1980-10-10', '2020-01-01', true, 'I am a pilot'),
  (1, 'Alice', 'Johnson', 'medic', '1975-12-12', '2020-01-01', true, 'I am a medic'),
  (1, 'Charlie', 'Brown', 'navigator', '1960-03-03', '2020-01-01', true, 'I am a navigator'),
  (1, 'Eve', 'Green', 'navigator', '1955-07-07', '2020-01-01', true, 'I am a navigator'),
  (1, 'Frank', 'White', 'squad_leader', '1940-09-09', '2020-01-01', true, 'I am a squad leader'),
  (1, 'Grace', 'Black', 'scientist', '1935-11-11', '2020-01-01', true, 'I am a scientist'),
  (1, 'Hank', 'Davis', 'engineer', '1920-02-02', '2020-01-01', false, 'I am an engineer'),
  (1, 'Ivy', 'Wilson', 'pilot', '1915-04-04', '2020-01-01', false, 'I am a pilot');    

-- Insert missions
INSERT INTO mission (name, destination_planet_id, status, launch_date, return_date, budget, objective)
VALUES
  ('Mars Pathfinder', 2, 'completed', '2022-05-10', '2023-06-15', 1250000000.00, 'Initial survey of potential colony sites'),
  ('Jupiter Orbital', 3, 'in_progress', '2023-09-22', NULL, 3750000000.00, 'Long-term study of Jovian atmosphere'),
  ('Europa Submersible', 4, 'planned', '2024-07-15', '2026-08-30', 4200000000.00, 'Search for life in subsurface ocean'),
  ('Red Soil Sample', 2, 'completed', '2021-03-05', '2021-11-20', 890000000.00, 'Collection of Martian soil samples'),
  ('Proxima Expedition', 5, 'planned', '2025-01-10', '2035-01-10', 12500000000.00, 'First interstellar colony mission');

-- Insert mission crew assignments
INSERT INTO mission_crew (mission_id, astronaut_id, role, join_date, end_date, notes)
VALUES
  (10000, 1, 'Mission Commander', '2022-05-10', '2023-06-15', 'Led first successful Mars landing'),
  (10000, 3, 'Engineer', '2022-05-10', '2023-06-15', 'Maintained life support systems'),
  (10000, 5, 'Security Chief', '2022-05-10', '2023-06-15', NULL),
  (10042, 1, 'Chief Engineer', '2023-09-22', NULL, 'Supervising atmospheric probe deployment'),
  (10042, 2, 'Science Officer', '2023-09-22', NULL, 'Leading research team'),
  (10042, 4, 'Research Specialist', '2023-09-22', NULL, 'Specializing in gas giant composition analysis'),
  (10084, 1, 'Mission Commander', '2024-07-15', NULL, 'Pre-mission training in progress'),
  (10084, 2, 'Lead Scientist', '2024-07-15', NULL, 'Designing submersible protocols'),
  (10126, 6, 'Team Lead', '2021-03-05', '2021-11-20', 'Coordinated sample collection'),
  (10126, 7, 'Strategic Advisor', '2021-03-05', '2021-03-30', 'Initial mission setup only'),
  (10168, 1, 'Mission Commander', '2025-01-10', '2035-01-10', 'Leading interstellar exploration team'),
  (10168, 2, 'Lead Scientist', '2025-01-10', '2035-01-10', 'Designing interstellar exploration protocols');

-- Insert resources
INSERT INTO resource (name, type, quantity, unit, planet_id, rarity, description, discovered_at)
VALUES
  ('Water Ice', 'H2O', 1500000.00, 'tons', 2, 'common', 'Frozen water deposits at Mars poles', '2022-06-12 08:30:00+00'),
  ('Helium-3', 'Isotope', 325.50, 'kg', 3, 'rare', 'Potential fusion energy source', '2023-10-05 14:22:00+00'),
  ('Europan Bacteria', 'Biological', 0.02, 'grams', 4, 'ultra_rare', 'Simple life forms found in ice samples', '2024-02-18 11:45:00+00'),
  ('Iron Oxide', 'Mineral', 32000000.00, 'tons', 2, 'common', 'Gives Mars its red color', '2022-05-28 10:15:00+00'),
  ('Deuterium', 'Isotope', 560.00, 'kg', 3, 'uncommon', 'Heavy hydrogen isotope', '2023-11-12 09:40:00+00'),
  ('Silica', 'Mineral', 8500000.00, 'tons', 2, 'common', 'Silicon dioxide', '2022-07-03 16:20:00+00'),
  ('Methane Ice', 'Hydrocarbon', 42000.00, 'tons', 4, 'uncommon', 'Frozen methane deposits', '2024-03-05 13:10:00+00'),
  ('Xenon', 'Noble Gas', 120.00, 'kg', 3, 'rare', 'Used in ion propulsion systems', '2023-12-18 07:50:00+00');