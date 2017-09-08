#!/usr/bin/env python3
import os, sys
import psycopg2
import itertools as it

def import_languages(lang_filepath, db, user, password):
  """Add the languages from a BCP47 registry file to a database.

  The registry file should follow the format displayed in
  https://www.iana.org/assignments/language-subtag-registry/language-subtag-registry

  The credentials should be for a PostgreSQL database that has
  a table called languages. The schema is defined in languages.sql.
  """
  credentials = "dbname='{}' user='{}' host='localhost' password='{}'"\
                .format(db, user, password)
  con = psycopg2.connect(credentials)
  cursor = con.cursor()

  # It's possible that the database volume already existed. In that
  # case, we don't need to add the languages; they're already there.
  cursor.execute("SELECT count(*) FROM languages")
  if cursor.fetchone()[0] != 0:
    return
  
  with open(lang_filepath, 'r') as registry:
    # The first two lines are unimportant.
    registry.readline()
    registry.readline()
    delimiter = '%%'

    # Registry files partition language definitions into sections that
    # are separated by '%%'. Read each definition and add the relevant
    # parts to the database.
    for key, group in it.groupby(registry, lambda line: line.startswith(delimiter)):
      if not key:
        language = list(group)

        tag = language[1].split(": ")[1][:-1]
        description = language[2].split(": ")[1][:-1].lower()
        
        # Some descriptions appear twice in the official registry file.
        # Because 'description' is UNIQUE, that causes problems.
        # TODO: Make the language table's primary key a tuple.
        cursor.execute(
          """
          SELECT exists(
            SELECT * FROM languages WHERE description=%s
          )
          """, (description,)
        )
              
        if cursor.fetchone()[0] == False:
          cursor.execute(
            """
            INSERT INTO languages (description, tag)
              VALUES(%s, %s)
            """, (description, tag)
          )
    con.commit()
      

if __name__ == "__main__":
  if len(sys.argv) != 2:
    print("Missing file argument")
    print("Usage: {} language-registry-file".format(sys.argv[0]))
    sys.exit()

  try:
    import_languages(
      sys.argv[1],
      os.getenv('POSTGRES_DB'),
      os.getenv('POSTGRES_USER'),
      os.getenv('POSTGRES_PASSWORD')
    )
  except Exception as e:
    print(e)
    sys.exit('Failed to establish database connection')
