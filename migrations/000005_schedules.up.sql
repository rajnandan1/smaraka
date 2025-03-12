CREATE TABLE
	schedules (
		schedule_id TEXT PRIMARY KEY,
		schedule_name TEXT NOT NULL,
		schedule_description TEXT,
		schedule_type TEXT,
		schedule_status TEXT NOT NULL DEFAULT 'ACTIVE',
		schedule_url TEXT NOT NULL,
		organization_id TEXT NOT NULL,
		schedule_meta TEXT,
		interval_days INTEGER NOT NULL DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (schedule_url, organization_id),
		FOREIGN KEY (organization_id) REFERENCES organizations (id)
	);