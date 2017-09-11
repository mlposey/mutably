/* Application database schema
 *
 * RDBMS: PostgreSQL 9.5
 */

CREATE TABLE languages (
  id serial NOT NULL PRIMARY KEY,
  description text NOT NULL UNIQUE,  -- Language name (e.g., English, Dutch)
  tag text NOT NULL UNIQUE           -- Short code (e.g., en, es, nl)
);

CREATE TABLE words (
  id serial NOT NULL PRIMARY KEY,
  word text UNIQUE NOT NULL -- A single word, i.e., one without spaces
);

-- A verb is a word with matching set of templates.
-- Multiple rows may share the same word/lang pair but have a different
-- id. This is intended. An example is the verb lie, which has two
-- meanings in English. These two meanings will each have their own row.
CREATE TABLE verbs (
  id serial NOT NULL PRIMARY KEY,
  word_id int NOT NULL REFERENCES words(id) ON DELETE CASCADE,
  lang_id int NOT NULL REFERENCES languages(id) ON DELETE CASCADE
);

-- Templates provide rules for how a verb is conjugated.
CREATE TABLE templates (
  id serial NOT NULL PRIMARY KEY,
  lang_id int NOT NULL REFERENCES languages(id) ON DELETE CASCADE,
  template text NOT NULL UNIQUE
);

CREATE TABLE verb_templates (
  verb_id int NOT NULL REFERENCES verbs(id) ON DELETE CASCADE,
  template_id int NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
  PRIMARY KEY (verb_id, template_id)
);

