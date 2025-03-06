CREATE TABLE
	secrets (
		"id" TEXT PRIMARY KEY,
		"organization_id" TEXT,
		"secret_type" TEXT,
		"secret_value" TEXT UNIQUE,
		"current_state" TEXT,
		"creator_id" TEXT,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"last_used_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"secret_name" TEXT,
		FOREIGN KEY (organization_id) REFERENCES Organizations (id),
		FOREIGN KEY (creator_id) REFERENCES Users (id)
	);