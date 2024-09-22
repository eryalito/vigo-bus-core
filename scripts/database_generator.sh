#!/bin/bash

# Run the Python file
echo "Deleting the old database..."
rm -f "stops.db"
echo "Generating the database..."
echo "Downloading the data..."
curl -s -o stops.json https://datos.vigo.org/data/transporte/paradas.json
echo "Data downloaded"
echo "Running the Python script..."
python "scripts/database_generator_stops_lines.py"
echo "Database generated"
