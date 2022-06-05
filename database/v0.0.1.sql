-- name: create-drill-table
CREATE TABLE IF NOT EXISTS  drill (
    name TEXT PRIMARY KEY,
    description TEXT
);

-- name: create-tag-table
CREATE TABLE IF NOT EXISTS tag (
    name TEXT PRIMARY KEY
);

--  name: create-drill-tags
CREATE TABLE IF NOT EXISTS drill_tags (
    drill TEXT NOT NULL,
    tag TEXT NOT NULL,
    FOREIGN KEY (drill) REFERENCES drill(name),
    FOREIGN KEY (tag) REFERENCES tag(name),
    PRIMARY KEY (drill, tag)
);
