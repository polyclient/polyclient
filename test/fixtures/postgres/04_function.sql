-- Function 1: Calculate mission duration in days
CREATE OR REPLACE FUNCTION calculate_mission_duration(mission_id INTEGER)
RETURNS INTEGER AS $$
DECLARE
  start_date TIMESTAMP WITH TIME ZONE;
  end_date TIMESTAMP WITH TIME ZONE;
  duration INTEGER;
BEGIN
  SELECT launch_date, return_date INTO start_date, end_date
  FROM mission
  WHERE id = mission_id;
  
  IF start_date IS NULL THEN
    RETURN NULL;
  ELSIF end_date IS NULL THEN
    -- For in_progress missions, calculate against current date
    RETURN EXTRACT(DAY FROM (CURRENT_TIMESTAMP - start_date));
  ELSE
    RETURN EXTRACT(DAY FROM (end_date - start_date));
  END IF;
END;
$$ LANGUAGE plpgsql;

-- Function 2: Get astronaut mission history
CREATE OR REPLACE FUNCTION get_astronaut_missions(astronaut_id INTEGER)
RETURNS TABLE (
  mission_name VARCHAR(100),
  role VARCHAR(50),
  mission_status VARCHAR(20),
  destination VARCHAR(50),
  start_date DATE,
  end_date DATE,
  days_served INTEGER
) AS $$
BEGIN
  RETURN QUERY
  SELECT 
    m.name AS mission_name,
    mc.role,
    m.status AS mission_status,
    p.name AS destination,
    mc.join_date AS start_date,
    mc.end_date,
    CASE 
      WHEN mc.end_date IS NULL AND m.status = 'completed' THEN 
        EXTRACT(DAY FROM (m.return_date - mc.join_date::timestamp))::INTEGER
      WHEN mc.end_date IS NULL THEN 
        EXTRACT(DAY FROM (CURRENT_DATE - mc.join_date))::INTEGER
      ELSE 
        EXTRACT(DAY FROM (mc.end_date - mc.join_date))::INTEGER
    END AS days_served
  FROM mission_crew mc
  JOIN mission m ON mc.mission_id = m.id
  LEFT JOIN planet p ON m.destination_planet_id = p.id
  WHERE mc.astronaut_id = astronaut_id
  ORDER BY mc.join_date DESC;
END;
$$ LANGUAGE plpgsql;

-- Function 3: Generate mission report
CREATE OR REPLACE FUNCTION generate_mission_report(mission_id INTEGER)
RETURNS TEXT AS $$
DECLARE
  report TEXT;
  mission_rec mission%ROWTYPE;
  planet_name VARCHAR(50);
  crew_count INTEGER;
  commander_name TEXT;
  resources_found TEXT;
  mission_duration INTEGER;
BEGIN
  -- Get mission data
  SELECT * INTO mission_rec FROM mission WHERE id = mission_id;
  
  IF mission_rec IS NULL THEN
    RETURN 'Mission not found';
  END IF;
  
  -- Get planet name
  SELECT name INTO planet_name FROM planet WHERE id = mission_rec.destination_planet_id;
  
  -- Get crew count
  SELECT COUNT(*) INTO crew_count FROM mission_crew as mc WHERE mc.mission_id = mission_rec.id;
  
  -- Get commander
  SELECT a.first_name || ' ' || a.last_name INTO commander_name
  FROM mission_crew mc
  JOIN astronaut a ON mc.astronaut_id = a.id
  WHERE mc.mission_id = mission_rec.id AND mc.role LIKE '%Commander%'
  LIMIT 1;
  
  -- Get resources (if mission completed)
  IF mission_rec.status = 'completed' THEN
    SELECT STRING_AGG(name || ' (' || type || ')', ', ') INTO resources_found
    FROM resource
    WHERE planet_id = mission_rec.destination_planet_id
    AND (discovered_at BETWEEN mission_rec.launch_date AND COALESCE(mission_rec.return_date, CURRENT_TIMESTAMP));
  END IF;
  
  -- Calculate duration
  mission_duration := calculate_mission_duration(mission_id);
  
  -- Build report
  report := 'MISSION REPORT: ' || mission_rec.name || E'\n';
  report := report || '----------------------------------------' || E'\n';
  report := report || 'Status: ' || mission_rec.status || E'\n';
  report := report || 'Destination: ' || planet_name || E'\n';
  report := report || 'Launch Date: ' || COALESCE(mission_rec.launch_date::TEXT, 'Not launched') || E'\n';
  report := report || 'Return Date: ' || COALESCE(mission_rec.return_date::TEXT, 'Not returned') || E'\n';
  report := report || 'Duration: ' || COALESCE(mission_duration::TEXT || ' days', 'N/A') || E'\n';
  report := report || 'Crew Size: ' || crew_count || E'\n';
  report := report || 'Commander: ' || COALESCE(commander_name, 'Not assigned') || E'\n';
  report := report || 'Budget: $' || mission_rec.budget || E'\n';
  report := report || 'Objective: ' || mission_rec.objective || E'\n';
  
  IF mission_rec.status = 'completed' AND resources_found IS NOT NULL THEN
    report := report || E'\nResources Discovered: ' || resources_found || E'\n';
  END IF;
  
  RETURN report;
END;
$$ LANGUAGE plpgsql;