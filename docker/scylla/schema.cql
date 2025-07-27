CREATE TABLE series_by_domain (
  domain TEXT,
  series_id UUID,
  series_name TEXT,
  url TEXT,
  avatar TEXT,
  genre TEXT,
  PRIMARY KEY (domain, series_id)
);


CREATE TABLE chapters_by_series (
  series_id UUID,
  chapter_id UUID,
  chapter_url TEXT,
  title TEXT,
  created_at TIMESTAMP,
  PRIMARY KEY (series_id, chapter_id)
);


CREATE TABLE images_by_chapter (
  chapter_id UUID,
  image_index INT,
  image_url TEXT,
  PRIMARY KEY (chapter_id, image_index)
);
