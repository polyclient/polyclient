-- Procedure 1: Launch mission (update status and crew)
CREATE OR REPLACE PROCEDURE launch_mission(
  p_mission_id INTEGER,
  p_launch_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
LANGUAGE plpgsql
AS $$
DECLARE
  v_status VARCHAR(20);
  v_crew_count INTEGER;
BEGIN
  -- Check if mission exists and is in planned status
  SELECT status INTO v_status FROM mission WHERE id = p_mission_id;
  
  IF v_status IS NULL THEN
    RAISE EXCEPTION 'Mission with ID % does not exist', p_mission_id;
  ELSIF v_status != 'planned' THEN
    RAISE EXCEPTION 'Mission with ID % is not in planned status (current status: %)', p_mission_id, v_status;
  END IF;
  
  -- Check if mission has crew assigned
  SELECT COUNT(*) INTO v_crew_count FROM mission_crew WHERE mission_id = p_mission_id;
  
  IF v_crew_count = 0 THEN
    RAISE EXCEPTION 'Cannot launch mission with ID % because no crew is assigned', p_mission_id;
  END IF;
  
  -- Update mission status and launch date
  UPDATE mission 
  SET status = 'in_progress', 
      launch_date = p_launch_date,
      updated_at = CURRENT_TIMESTAMP
  WHERE id = p_mission_id;
  
  RAISE NOTICE 'Mission % successfully launched', p_mission_id;
END;
$$;

-- Procedure 2: Complete mission (update status, return date, and log discoveries)
CREATE OR REPLACE PROCEDURE complete_mission(
  p_mission_id INTEGER,
  p_return_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  p_discovery_notes TEXT DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
DECLARE
  v_status VARCHAR(20);
  v_destination_id INTEGER;
  v_crew_rec RECORD;
BEGIN
  -- Check if mission exists and is in_progress
  SELECT status, destination_planet_id INTO v_status, v_destination_id 
  FROM mission 
  WHERE id = p_mission_id;
  
  IF v_status IS NULL THEN
    RAISE EXCEPTION 'Mission with ID % does not exist', p_mission_id;
  ELSIF v_status != 'in_progress' THEN
    RAISE EXCEPTION 'Only in_progress missions can be completed (current status: %)', v_status;
  END IF;
  
  -- Update mission status and return date
  UPDATE mission 
  SET status = 'completed', 
      return_date = p_return_date,
      updated_at = CURRENT_TIMESTAMP
  WHERE id = p_mission_id;
  
  -- Update all crew end dates that are still NULL
  UPDATE mission_crew
  SET end_date = p_return_date::DATE
  WHERE mission_id = p_mission_id AND end_date IS NULL;
  
  -- Log mission completion
  RAISE NOTICE 'Mission % completed successfully', p_mission_id;
  
  -- Log discovery notes if provided
  IF p_discovery_notes IS NOT NULL THEN
    RAISE NOTICE 'Discovery notes recorded: %', p_discovery_notes;
  END IF;
END;
$$;