CREATE TABLE
	users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		name TEXT,
		password_hash TEXT,
		created_at TIMESTAMP,
		seen_at TIMESTAMP,
		updated_at TIMESTAMP
	);

CREATE TABLE
	organizations (
		id TEXT PRIMARY KEY,
		name TEXT,
		creator_id TEXT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);

CREATE TABLE
	user_organizations (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		organization_id TEXT,
		role TEXT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		UNIQUE (user_id, organization_id),
		FOREIGN KEY (organization_id) REFERENCES organizations (id),
		FOREIGN KEY (user_id) REFERENCES users (id)
	);

CREATE TABLE
	url_store (
		id TEXT PRIMARY KEY,
		url TEXT UNIQUE,
		domain TEXT,
		title TEXT,
		image_sm TEXT,
		image_lg TEXT,
		excerpt TEXT,
		color TEXT,
		status TEXT,
		full_content TEXT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);

CREATE TABLE
	url_organizations (
		id TEXT PRIMARY KEY,
		url_id TEXT,
		organization_id TEXT,
		status TEXT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		FOREIGN KEY (url_id) REFERENCES url_store (id),
		FOREIGN KEY (organization_id) REFERENCES organizations (id),
		UNIQUE (url_id, organization_id)
	);

CREATE TABLE
	job_queue (
		id TEXT PRIMARY KEY,
		org_id TEXT,
		job_id TEXT,
		job_data TEXT,
		status TEXT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		FOREIGN KEY (org_id) REFERENCES organizations (id)
	);

CREATE INDEX url_store_idx ON url_store USING bm25 (
	id,
	url,
	domain,
	title,
	excerpt,
	full_content,
	created_at
)
WITH
	(
		key_field = 'id',
		datetime_fields = '{
      "created_at": {"fast": true}
		}',
		text_fields = '{
        "title": {
          "tokenizer": {"type": "whitespace"}
        },
				"excerpt": {
          "tokenizer": {"type": "whitespace"}
        },
				"full_content": {
          "tokenizer": {"type": "whitespace"}
        },
				"domain": {
          "tokenizer": {"type": "raw"}
        }
    }'
	);