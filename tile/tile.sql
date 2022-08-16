CREATE OR REPLACE FUNCTION
    public.tile_omt(z integer, x integer, y integer, query_params json)
    RETURNS bytea
AS
$$
DECLARE
    aerodrome_label_mvt     bytea;
    aeroway_mvt             bytea;
    boundary_mvt            bytea;
    building_mvt            bytea;
    housenumber_mvt         bytea;
    landcover_mvt           bytea;
    landuse_mvt             bytea;
    mountain_peak_mvt       bytea;
    park_mvt                bytea;
    place_mvt               bytea;
    poi_mvt                 bytea;
    transportation_mvt      bytea;
    transportation_name_mvt bytea;
    water_mvt               bytea;
    water_name_mvt          bytea;
    waterway_mvt            bytea;
BEGIN
    --- aerodrome_label
    SELECT INTO aerodrome_label_mvt ST_AsMVT(tile, 'aerodrome_label', 4096, 'geom', 'id')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 id,
                 name,
                 name_en,
                 name_de,
                 tags,
                 class,
                 iata,
                 icao,
                 ele,
                 ele_ft
          FROM layer_aerodrome_label(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    --- aeroway
    SELECT INTO aeroway_mvt ST_AsMVT(tile, 'aeroway', 4096, 'geom')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 class,
                 ref
          FROM layer_aeroway(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    -- boundary
    SELECT INTO boundary_mvt ST_AsMVT(tile, 'boundary', 4096, 'geom')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 admin_level,
                 disputed,
                 disputed_name,
                 claimed_by,
                 maritime
          FROM layer_boundary(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    -- building
    SELECT INTO building_mvt ST_AsMVT(tile, 'building', 4096, 'geom', 'osm_id')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 osm_id,
                 render_height,
                 render_min_height,
                 colour,
                 hide_3d
          FROM layer_building(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    -- housenumber
    SELECT INTO housenumber_mvt ST_AsMVT(tile, 'housenumber', 4096, 'geom', 'osm_id')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 osm_id,
                 housenumber
          FROM layer_housenumber(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    -- landcover
    SELECT INTO landcover_mvt ST_AsMVT(tile, 'landcover', 4096, 'geom')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 class,
                 subclass
          FROM layer_landcover(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    -- landuse
    SELECT INTO landuse_mvt ST_AsMVT(tile, 'landuse', 4096, 'geom', 'osm_id')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 osm_id,
                 class
          FROM layer_landuse(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    -- mountain_peak
    SELECT INTO mountain_peak_mvt ST_AsMVT(tile, 'mountain_peak', 4096, 'geom', 'osm_id')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 osm_id,
                 name,
                 name_en,
                 name_de,
                 class,
                 tags,
                 ele,
                 ele_ft,
                 rank
          FROM layer_mountain_peak(TileBBox(z, x, y, 3857), z, 256)) AS tile
    WHERE geom IS NOT NULL;

    -- park
    SELECT INTO park_mvt ST_AsMVT(tile, 'park', 4096, 'geom', 'osm_id')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 osm_id,
                 class,
                 name,
                 name_en,
                 name_de,
                 tags,
                 rank
          FROM layer_park(TileBBox(z, x, y, 3857), z, 256)) AS tile
    WHERE geom IS NOT NULL;

    -- place
    SELECT INTO place_mvt ST_AsMVT(tile, 'place', 4096, 'geom', 'osm_id')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 osm_id,
                 name,
                 name_en,
                 name_de,
                 capital,
                 class,
                 iso_a2,
                 rank
          FROM layer_place(TileBBox(z, x, y, 3857), z, 256)) AS tile
    WHERE geom IS NOT NULL;

    -- poi
    SELECT INTO poi_mvt ST_AsMVT(tile, 'poi', 4096, 'geom', 'osm_id')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 osm_id,
                 name,
                 name_en,
                 name_de,
                 class,
                 subclass,
                 rank,
                 agg_stop,
                 level,
                 layer,
                 indoor
          FROM layer_poi(TileBBox(z, x, y, 3857), z, 256)) AS tile
    WHERE geom IS NOT NULL;

    -- transportation
    SELECT INTO transportation_mvt ST_AsMVT(tile, 'transportation', 4096, 'geom', 'osm_id')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 osm_id,
                 class,
                 subclass,
                 ramp,
                 oneway,
                 brunnel,
                 service,
                 layer,
                 level,
                 indoor,
                 bicycle,
                 foot,
                 horse,
                 mtb_scale,
                 surface
          FROM layer_transportation(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    -- transportation_name
    SELECT INTO transportation_name_mvt ST_AsMVT(tile, 'transportation_name', 4096, 'geom')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 name,
                 ref,
                 ref_length,
                 network,
                 class,
                 subclass,
                 brunnel,
                 layer,
                 level,
                 indoor
          FROM layer_transportation_name(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    -- water
    SELECT INTO water_mvt ST_AsMVT(tile, 'water', 4096, 'geom')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 class,
                 brunnel,
                 intermittent
          FROM layer_water(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    -- water_name
    SELECT INTO water_name_mvt ST_AsMVT(tile, 'water_name', 4096, 'geom')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 name,
                 name_en,
                 name_de,
                 class,
                 intermittent
          FROM layer_water_name(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    -- waterway
    SELECT INTO waterway_mvt ST_AsMVT(tile, 'waterway', 4096, 'geom')
    FROM (SELECT ST_AsMVTGeom(geometry, TileBBox(z, x, y, 3857), 4096, 64, true) AS geom,
                 name,
                 name_en,
                 name_de,
                 class,
                 brunnel,
                 intermittent
          FROM layer_waterway(TileBBox(z, x, y, 3857), z)) AS tile
    WHERE geom IS NOT NULL;

    RETURN aerodrome_label_mvt ||
           aeroway_mvt ||
           boundary_mvt ||
           building_mvt ||
           housenumber_mvt ||
           landcover_mvt ||
           landuse_mvt ||
           mountain_peak_mvt ||
           park_mvt ||
           place_mvt ||
           poi_mvt ||
           transportation_mvt ||
           transportation_name_mvt ||
           water_mvt ||
           water_name_mvt ||
           waterway_mvt;
END
$$
    LANGUAGE plpgsql IMMUTABLE
                     STRICT
                     PARALLEL SAFE;
