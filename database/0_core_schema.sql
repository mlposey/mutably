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

/* add_infinitive creates a new infinitive verb.
 *
 * The verb is assigned a new conjugation table that is populated with the ids
 * of new (but empty) tense inflection rows. However, if the verb already
 * existed, nothing is created.
 *
 * params:
 *      _word the text of a word which may or may not be in the words table
 *      _lang_id the id of a language in the languages table
 *
 * returns: The id of the verb's conjugation table
 */
CREATE OR REPLACE FUNCTION add_infinitive(_word TEXT, _lang_id INTEGER)
RETURNS INTEGER AS $$
DECLARE
    _present_id INTEGER;
    _past_id    INTEGER;
    _word_id    INTEGER;
    _conj_id    INTEGER;
    _verb_id    INTEGER;
BEGIN
    SELECT id FROM words
    WHERE word = _word
    INTO _word_id;

    IF NOT FOUND THEN
        INSERT INTO words (word)
        VALUES (_word)
        RETURNING id INTO _word_id;
    ELSE
        -- Check for existence before adding.
        SELECT conjugation_table FROM verbs
        WHERE  word_id = _word_id
        AND    lang_id = _lang_id
        INTO   _conj_id;

        IF FOUND THEN
            RETURN _conj_id;
        END IF;
    END IF;

    -- Create two tense inflections for present and past tenses.
    INSERT INTO tense_inflections DEFAULT VALUES
    RETURNING id INTO _present_id;
    INSERT INTO tense_inflections DEFAULT VALUES
    RETURNING id INTO _past_id;
    -- Create a conjugation table with the tense inflections.
    INSERT INTO conjugation_tables (present, past)
    VALUES (_present_id, _past_id)
    RETURNING id INTO _conj_id;
    -- Create a verb for the infinitive.
    INSERT INTO verbs (word_id, lang_id, conjugation_table)
    VALUES (_word_id, _lang_id, _conj_id)
    RETURNING id INTO _verb_id;
    -- Add the verb's id to the conjugation table as an infinitive.
    UPDATE conjugation_tables
    SET infinitive_id = _verb_id
    WHERE id = _conj_id;
    -- Return the id of the conjugation table.
    RETURN _conj_id;
END;
$$ LANGUAGE plpgsql;

