CREATE TABLE
	schedules (
		schedule_id TEXT PRIMARY KEY,
		schedule_name TEXT NOT NULL,
		schedule_description TEXT,
		schedule_url TEXT,
		schedule_meta TEXT,
		default_interval_days INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

CREATE TABLE
	org_schedules (
		organization_id TEXT NOT NULL,
		schedule_id TEXT NOT NULL,
		status TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (organization_id, schedule_id),
		FOREIGN KEY (organization_id) REFERENCES organizations (id),
		FOREIGN KEY (schedule_id) REFERENCES schedules (schedule_id)
	);

-- Insert dummy data into schedules table
INSERT INTO
	schedules (
		schedule_id,
		schedule_name,
		schedule_description,
		schedule_url,
		schedule_meta,
		default_interval_days
	)
VALUES
	(
		'sch_gh_trending_daily',
		'Github Trending Daily',
		'Get a daily summary of trending repositories on Github',
		'https://github.com/trending?since=daily',
		'',
		1
	),
	(
		'sch_gh_trending_weekly',
		'Github Trending Weekly',
		'Get a weekly summary of trending repositories on Github',
		'https://github.com/trending?since=weekly',
		'',
		7
	),
	(
		'sch_gh_trending_monthly',
		'Github Trending Monthly',
		'Get a monthly summary of trending repositories on Github',
		'https://github.com/trending?since=monthly',
		'',
		30
	),
	(
		'sch_ph_leaderboard_daily',
		'Product Hunt Daily Leaderboard',
		'Get a daily summary of the leaderboard on Product Hunt',
		'https://www.producthunt.com/leaderboard/daily',
		'',
		1
	),
	(
		'sch_ph_leaderboard_weekly',
		'Product Hunt Weekly Leaderboard',
		'Get a weekly summary of trending products on Product Hunt',
		'https://www.producthunt.com/leaderboard/weekly',
		'',
		7
	),
	(
		'sch_ph_leaderboard_monthly',
		'Product Hunt Monthly Leaderboard',
		'Get a monthly summary of trending products on Product Hunt',
		'https://www.producthunt.com/leaderboard/monthly',
		'',
		30
	);