/* Application database schema
 * RDBMS: PostgreSQL 9.5
 *
 * This should be the first file that the postgres container runs. If adding
 * other files, you may need to prefix the name of this one with '0'.
 */

CREATE TABLE languages (
    id serial PRIMARY KEY,
    name text UNIQUE NOT NULL,
    tag text UNIQUE -- Short code (e.g., en, es, nl)
);

CREATE TABLE words (
    id serial PRIMARY KEY,
    word text UNIQUE NOT NULL
);

-- Grammatical tense (e.g., present, past)
CREATE TABLE tenses (
    id serial PRIMARY KEY,
    tense text UNIQUE NOT NULL
);
INSERT INTO tenses (tense)
VALUES ('present'), ('past');

CREATE TABLE verb_forms (
    id serial PRIMARY KEY,
    lang_id  int NOT NULL REFERENCES languages(id),
    word_id  int NOT NULL REFERENCES words(id),
    inf_id   int NOT NULL REFERENCES words(id),
    tense_id int NOT NULL REFERENCES tenses(id),
    person   int, -- Plural verbs won't have a person.
    num      int NOT NULL -- 1 is singular; not 1 is plural
);
