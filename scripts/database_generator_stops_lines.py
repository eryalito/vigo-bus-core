import sqlite3
import json

def create_database(db_path):
    # Connect to the SQLite database (or create it if it doesn't exist)
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    # SQL statements to create the tables
    create_lines_table = """
    CREATE TABLE IF NOT EXISTS lines (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    );"""

    create_stops_table = """
    CREATE TABLE IF NOT EXISTS stops (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        stop_number INTEGER NOT NULL,
        stop_id INTEGER NOT NULL,
        name TEXT NOT NULL,
        lat REAL,
        lon REAL
    );"""

    create_line_stops_table = """
    CREATE TABLE IF NOT EXISTS line_stops (
        line_id INTEGER,
        stop_id INTEGER,
        PRIMARY KEY (line_id, stop_id),
        FOREIGN KEY (line_id) REFERENCES lines(id),
        FOREIGN KEY (stop_id) REFERENCES stops(id)
    );"""

    # Execute the SQL statements to create the tables
    cursor.execute(create_lines_table)
    cursor.execute(create_stops_table)
    cursor.execute(create_line_stops_table)

    # Commit the changes and close the connection
    conn.commit()
    conn.close()

def insert_data(db_path, json_path):
    # Connect to the SQLite database
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    # Read the JSON file
    with open(json_path, 'r') as f:
        stops_data = json.load(f)

    # Insert stops and lines into the database
    for stop in stops_data:

        # Insert stop
        cursor.execute("""
            INSERT INTO stops (stop_number, stop_id, name, lat, lon)
            VALUES (?, ?, ?, ?, ?)
        """, (stop['id'], stop['stop_id'], stop['nombre'], stop['lat'], stop['lon']))
        
        stop_db_id = cursor.lastrowid

        # Insert lines and line_stops
        lines = stop['lineas'].split(', ')
        for line in lines:
            # Check if the line already exists
            cursor.execute("SELECT id FROM lines WHERE name = ?", (line,))
            line_row = cursor.fetchone()
            if line_row:
                line_db_id = line_row[0]
            else:
                cursor.execute("INSERT INTO lines (name) VALUES (?)", (line,))
                line_db_id = cursor.lastrowid

            # Insert into line_stops
            cursor.execute("""
                INSERT INTO line_stops (line_id, stop_id)
                VALUES (?, ?)
            """, (line_db_id, stop_db_id))

    # Commit the changes and close the connection
    conn.commit()
    conn.close()

if __name__ == "__main__":
    # Path to the SQLite database file
    db_path = "stops.db"
    # Path to the JSON file
    json_path = "stops.json"

    # Create the database and tables
    create_database(db_path)
    # Insert data from the JSON file
    insert_data(db_path, json_path)
    print(f"Database created and data inserted at {db_path}")