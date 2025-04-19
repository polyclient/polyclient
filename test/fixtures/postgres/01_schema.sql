SET client_encoding = 'UTF8';

-- Domain for coordinates (x,y,z)
CREATE DOMAIN coordinates AS JSONB CONSTRAINT valid_coordinates CHECK (
    jsonb_typeof(value) = 'object'
    AND jsonb_typeof(value -> 'x') = 'number'
    AND jsonb_typeof(value -> 'y') = 'number'
    AND jsonb_typeof(value -> 'z') = 'number'
);

-- Enums 
CREATE TYPE planet_type AS ENUM ('terrestrial', 'gaseous', 'icy', 'dwarf');

CREATE TYPE astronaut_specialty AS ENUM ('engineer', 'scientist', 'pilot', 'medic', 'navigator', 'squad_leader');

CREATE TYPE mission_status AS ENUM ('planned', 'in_progress', 'completed', 'failed', 'cancelled');

CREATE TYPE habitability_status AS ENUM ('habitable', 'non_habitable', 'potentially_habitable', 'unknown');

CREATE TYPE resource_rarity AS ENUM ('common', 'uncommon', 'rare', 'ultra_rare');

-- Sequences
CREATE SEQUENCE mission_id_seq
    START WITH 10000
    INCREMENT BY 42
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Tables
CREATE TABLE planet (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type planet_type NOT NULL,
    coordinates coordinates NOT NULL,
    diameter_km DOUBLE PRECISION CHECK (diameter_km > 0) NOT NULL,
    mass_kg DOUBLE PRECISION CHECK (mass_kg > 0) NOT NULL,
    has_atmosphere BOOLEAN DEFAULT FALSE NOT NULL,
    has_moons BOOLEAN DEFAULT FALSE NOT NULL,
    has_rings BOOLEAN DEFAULT FALSE NOT NULL,
    has_oceans BOOLEAN DEFAULT FALSE NOT NULL,
    habitability_status habitability_status DEFAULT 'unknown' NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE astronaut (
    id SERIAL PRIMARY KEY,
    home_planet_id INTEGER REFERENCES planet(id),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    specialty astronaut_specialty NOT NULL,
    birth_date DATE NOT NULL,
    join_date DATE NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    bio TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE mission (
    id INTEGER PRIMARY KEY DEFAULT nextval('mission_id_seq'),
    code_name VARCHAR(50) GENERATED ALWAYS AS ('M-' || id) STORED, 
    destination_planet_id INTEGER REFERENCES planet(id),
    name VARCHAR(100) UNIQUE NOT NULL,
    status mission_status DEFAULT 'planned' NOT NULL,
    launch_date TIMESTAMP WITH TIME ZONE,
    return_date TIMESTAMP WITH TIME ZONE,
    budget DOUBLE PRECISION CHECK (budget > 0) NOT NULL,
    objective TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_dates CHECK (
        return_date IS NULL
        OR launch_date IS NULL
        OR return_date > launch_date
    )
);

CREATE TABLE mission_crew (
    mission_id INTEGER REFERENCES mission(id) ON DELETE CASCADE,
    astronaut_id INTEGER REFERENCES astronaut(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    join_date DATE NOT NULL,
    end_date DATE,
    notes TEXT,
    PRIMARY KEY (mission_id, astronaut_id),
    CONSTRAINT valid_crew_dates CHECK (
        end_date IS NULL
        OR end_date >= join_date
    )
);

CREATE TABLE resource (
    id SERIAL PRIMARY KEY,
    planet_id INTEGER REFERENCES planet(id),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    rarity resource_rarity NOT NULL,
    quantity DOUBLE PRECISION CHECK (quantity >= 0) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    description TEXT,
    discovered_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_planet_type ON planet(type);

CREATE INDEX idx_astronaut_specialty ON astronaut(specialty);

CREATE INDEX idx_mission_status ON mission(status);

CREATE INDEX idx_resource_type_rarity ON resource(type, rarity);

CREATE INDEX idx_mission_dates ON mission(launch_date, return_date);