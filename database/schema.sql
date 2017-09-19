/* Application database schema
 * RDBMS: PostgreSQL 9.5
 *
 * This should be the first file that the postgres container runs. If adding
 * other files, you may need to prefix the name of this one with '0'.
 */

CREATE TABLE languages (
  id serial NOT NULL PRIMARY KEY,
  language text NOT NULL UNIQUE,     -- Language name (e.g., English, Dutch)
  tag text UNIQUE                    -- Short code (e.g., en, es, nl)
);

CREATE TABLE words (
  id serial NOT NULL PRIMARY KEY,
  word text UNIQUE NOT NULL -- A single word, i.e., one without spaces
);

-- Inflections of an infinitive for a single tense
CREATE TABLE tense_inflections (
  id serial NOT NULL PRIMARY KEY,
  -- Person
  first  int REFERENCES words(id),
  second int REFERENCES words(id),
  third  int REFERENCES words(id),
  -- Number
  plural int REFERENCES words(id)
);

-- These tables store the conjugated forms of an infinitive verb.
CREATE TABLE conjugation_tables (
  id serial NOT NULL PRIMARY KEY,
  -- infinitive_id is added later in this schema file.
  present int NOT NULL REFERENCES tense_inflections(id), -- Present tense
  past    int NOT NULL REFERENCES tense_inflections(id)  -- Past tense
);

CREATE TABLE verbs (
  id serial NOT NULL PRIMARY KEY,
  word_id int NOT NULL REFERENCES words(id) ON DELETE CASCADE,
  lang_id int NOT NULL REFERENCES languages(id) ON DELETE CASCADE,

  /* All verbs (even infinitives) are part of a conjugation table.
     A group of verbs that come from the same infinitive will share
     the same table. */
  conjugation_table int NOT NULL REFERENCES conjugation_tables(id),
  UNIQUE (word_id, lang_id)
);

ALTER TABLE conjugation_tables
ADD COLUMN infinitive_id int REFERENCES verbs(id)
ON DELETE CASCADE;

----------------------- Function Definitions --------------------------

/* addInfinitive creates a conjugation table and verb entry for a
   word that is marked as an infinitive verb.
   It returns the id of the conjugation table. */
CREATE FUNCTION add_infinitive(word_id int, lang_id int) RETURNS INTEGER AS $$
DECLARE
    present_id INTEGER;
    past_id    INTEGER;
    conj_id    INTEGER;
    verb_id    INTEGER;
BEGIN
    -- Create two tense inflections for present and past tenses.
    INSERT INTO tense_inflections DEFAULT VALUES
    RETURNING id INTO present_id;
    INSERT INTO tense_inflections DEFAULT VALUES
    RETURNING id INTO past_id;
    -- Create a conjugation table with the tense inflections.
    INSERT INTO conjugation_tables (present, past)
    VALUES (present_id, past_id)
    RETURNING id INTO conj_id;
    -- Create a verb for the infinitive.
    INSERT INTO verbs (word_id, lang_id, conjugation_table)
    VALUES (word_id, lang_id, conj_id)
    RETURNING id INTO verb_id;
    -- Add the verb's id to the conjugation table as an infinitive.
    UPDATE conjugation_tables
    SET infinitive_id = verb_id
    WHERE id = conj_id;
    -- Return the id of the conjugation table.
    RETURN conj_id;
END;
$$ LANGUAGE plpgsql;

