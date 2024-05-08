package database

// "fmt"
// "log"

func (db *appdbimpl) GetTables() error {

	// List of table names
	/*
		tables := []string{"users"} // Add more tables as needed
		// "photos","follow", "likes", "Bans", "comments"
		//
		// Iterate over each table and print its content
		for _, table := range tables {
			// fmt.Printf("Content of table '%s':\n", table)

			// Prepare SQL statement to select all rows from the table
			// query := fmt.Sprintf("SELECT * FROM %s", table)
			// query := fmt.Sprintf("SELECT photo_id, author, upload_datetime, location, caption FROM %s", table)

			// Execute the query
			rows, err := db.c.Query(query)
			if err != nil {
				log.Printf("Error querying table %s: %v", table, err)
				continue
			}
			defer rows.Close()

			// Iterate over the rows and print each row
			for rows.Next() {
				var values []interface{}
				columns, err := rows.Columns()
				if err != nil {
					log.Printf("Error getting column names for table %s: %v", table, err)
					break
				}

				// Create a slice to store the values of each row
				for range columns {
					var value interface{}
					values = append(values, &value)
				}

				// Scan the values into the slice
				if err := rows.Scan(values...); err != nil {
					log.Printf("Error scanning row values for table %s: %v", table, err)
					continue
				}

				// Print the values of the row
				for i, value := range values {
					// fmt.Printf("%s: %v\t", columns[i], *value.(*interface{}))
				}
				// fmt.Println()
			}
			if err := rows.Err(); err != nil {
				log.Printf("Error iterating rows for table %s: %v", table, err)
			}

			// fmt.Println()
			// Drop each table

			for _, table := range tables {
				// query := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
				_, err := db.c.Exec(query)
				if err != nil {
					// return fmt.Errorf("error dropping table %s: %w", table, err)
				}
				// fmt.Printf("Table %s dropped successfully\n", table)
			}
			return nil
		}
	*/
	return nil
}
