-- Function to update timestamp before update
CREATE OR REPLACE FUNCTION update_modified_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_planet_timestamp
BEFORE UPDATE ON planet
FOR EACH ROW
EXECUTE FUNCTION update_modified_timestamp();

CREATE TRIGGER update_astronaut_timestamp
BEFORE UPDATE ON astronaut
FOR EACH ROW
EXECUTE FUNCTION update_modified_timestamp();

CREATE TRIGGER update_mission_timestamp
BEFORE UPDATE ON mission
FOR EACH ROW
EXECUTE FUNCTION update_modified_timestamp();

CREATE TRIGGER update_resource_timestamp
BEFORE UPDATE ON resource
FOR EACH ROW
EXECUTE FUNCTION update_modified_timestamp();

-- Create resource discovery logging trigger
CREATE OR REPLACE FUNCTION log_resource_discovery()
RETURNS TRIGGER AS $$
BEGIN
  RAISE NOTICE 'New resource discovered: % (%) on planet ID %', 
    NEW.name, 
    NEW.type, 
    NEW.planet_id;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER resource_discovery_trigger
AFTER INSERT ON resource
FOR EACH ROW
EXECUTE FUNCTION log_resource_discovery();