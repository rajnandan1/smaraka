DROP TABLE IF EXISTS job_queue;

DROP TABLE IF EXISTS url_organizations;

DROP TABLE IF EXISTS url_store;

DROP TABLE IF EXISTS user_organizations;

DROP TABLE IF EXISTS organizations;

DROP TABLE IF EXISTS users;

CALL paradedb.drop_bm25 (index_name => 'url_store_idx');