/*
Reference from https://cloud.google.com/spanner/docs/schema-and-data-model
*/

CREATE TABLE singers (
 singer_id   BIGINT PRIMARY KEY,
 first_name  VARCHAR(1024),
 last_name   VARCHAR(1024),
 singer_info BYTEA
 );

CREATE TABLE albums (
 singer_id     BIGINT,
 album_id      BIGINT,
 album_title   VARCHAR,
 PRIMARY KEY (singer_id, album_id)
 )
 INTERLEAVE IN PARENT singers ON DELETE CASCADE;

CREATE TABLE songs (
 singer_id     BIGINT,
 album_id      BIGINT,
 track_id      BIGINT,
 song_name     VARCHAR,
 PRIMARY KEY (singer_id, album_id, track_id)
 )
 INTERLEAVE IN PARENT albums ON DELETE CASCADE;
