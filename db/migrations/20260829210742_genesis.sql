-- +goose Up

-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS ltree;
-- +goose StatementEnd

CREATE TYPE sku_size AS ENUM ('XS', 'S', 'M', 'L', 'XL');

CREATE TABLE stock_keeping_units (
  id bigint
    GENERATED ALWAYS AS IDENTITY (START WITH 10000)
    PRIMARY KEY,
  name text NOT NULL,
  size sku_size NOT NULL,
  created_at timestaptz DEFAULT NOW()
);

CREATE INDEX "idx_products_name_trgm" ON stock_keeping_units USING gin (name gin_trgm_ops);

CREATE TABLE locations (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  parent_id bigint REFERENCES locations(id)  ON DELETE CASCADE,
  
  -- ltree path that is generated with trigger function
  path ltree NOT NULL,
  
  name text NOT NULL,
  barcode text COLLATE "C" UNIQUE,
  is_leaf boolean NOT NULL DEFAULT true
);

-- Fast tree lookups (ancestor, descendant, sibling matches)
CREATE INDEX idx_locations_path_gist ON locations USING gist (path);
CREATE INDEX idx_locations_path_btree ON locations (path);
CREATE INDEX idx_locations_parent_id ON locations (parent_id);


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_location_id_path()
RETURNS TRIGGER AS 
$$
  DECLARE
    parent_path ltree;
  BEGIN
    IF NEW.parent_id IS NULL THEN
      IF NEW.id IS NULL THEN
        NEW.id := nextval(pg_get_serial_sequence('locations', 'id'));
      END IF;
      NEW.path := text2ltree(NEW.id::text);
      RETURN NEW;
    END IF;

    SELECT path INTO parent_path 
    FROM locations 
    WHERE id = NEW.parent_id;

    IF parent_path IS NULL THEN
      RAISE EXCEPTION 'Parent location ID % does not exist or has no path', NEW.parent_id;
    END IF;

    IF NEW.id IS NULL THEN
      NEW.id := nextval(pg_get_serial_sequence('locations', 'id'));
    END IF;

    NEW.path := parent_path || text2ltree(NEW.id::text);
    
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_set_location_id_path
BEFORE INSERT ON locations
FOR EACH ROW
EXECUTE FUNCTION set_location_id_path();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_location_id_path()
RETURNS TRIGGER AS $$
DECLARE
  parent_path ltree;
BEGIN
  IF NEW.parent_id IS NULL THEN
    IF NEW.id IS NULL THEN
      NEW.id := nextval(pg_get_serial_sequence('locations', 'id'));
    END IF;
    NEW.path := text2ltree(NEW.id::text);
    RETURN NEW;
  END IF;

  SELECT path INTO parent_path 
  FROM locations 
  WHERE id = NEW.parent_id;

  IF parent_path IS NULL THEN
    RAISE EXCEPTION 'Parent location ID % does not exist or has no path', NEW.parent_id;
  END IF;

  IF NEW.id IS NULL THEN
    NEW.id := nextval(pg_get_serial_sequence('locations', 'id'));
  END IF;

  NEW.path := parent_path || text2ltree(NEW.id::text);
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_set_location_id_path
BEFORE INSERT ON locations
FOR EACH ROW
EXECUTE FUNCTION set_location_id_path();

-- +goose Down
DROP TRIGGER IF EXISTS trg_update_location_subtree_path ON locations;
DROP FUNCTION IF EXISTS update_location_subtree_path();

DROP TRIGGER IF EXISTS trg_set_location_id_path ON locations;
DROP FUNCTION IF EXISTS set_location_id_path();

DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS stock_keeping_units;
DROP TYPE IF EXISTS sku_size;

