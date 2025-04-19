-- View 1: Active missions with crew count and destination details
CREATE OR REPLACE VIEW active_mission_summary AS
SELECT 
  m.id AS mission_id,
  m.name AS mission_name,
  m.status,
  p.name AS destination,
  p.type AS planet_type,
  m.launch_date,
  m.return_date,
  m.budget,
  COUNT(mc.astronaut_id) AS crew_count,
  STRING_AGG(a.last_name, ', ') AS crew_names
FROM mission m
JOIN planet p ON m.destination_planet_id = p.id
LEFT JOIN mission_crew mc ON m.id = mc.mission_id
LEFT JOIN astronaut a ON mc.astronaut_id = a.id
WHERE m.status IN ('planned', 'in_progress')
GROUP BY m.id, m.name, m.status, p.name, p.type, m.launch_date, m.return_date, m.budget
ORDER BY launch_date;

-- View 2: Resource availability by planet with rarity metrics
CREATE OR REPLACE VIEW planet_resource_summary AS
SELECT
  p.id AS planet_id,
  p.name AS planet_name,
  p.type AS planet_type,
  COUNT(r.id) AS resource_count,
  SUM(CASE WHEN r.rarity = 'common' THEN 1 ELSE 0 END) AS common_resources,
  SUM(CASE WHEN r.rarity = 'uncommon' THEN 1 ELSE 0 END) AS uncommon_resources,
  SUM(CASE WHEN r.rarity = 'rare' THEN 1 ELSE 0 END) AS rare_resources,
  SUM(CASE WHEN r.rarity = 'ultra_rare' THEN 1 ELSE 0 END) AS ultra_rare_resources,
  STRING_AGG(r.name, ', ') AS resource_list
FROM planet p
LEFT JOIN resource r ON p.id = r.planet_id
GROUP BY p.id, p.name, p.type
ORDER BY p.id;