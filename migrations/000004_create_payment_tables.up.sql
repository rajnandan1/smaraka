CREATE TABLE
	subscription (
		"id" TEXT PRIMARY KEY,
		"user_id" TEXT UNIQUE,
		"success_order_id" TEXT,
		"current_state" TEXT,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES Users (id)
	);

CREATE TABLE
	orders (
		"id" TEXT PRIMARY KEY,
		"user_id" TEXT,
		"customer_name" TEXT,
		"customer_email" TEXT,
		"subscription_id" TEXT,
		"variant_id" INTEGER,
		"store_id" INTEGER,
		"current_state" TEXT,
		"order_object" TEXT,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES Users (id),
		FOREIGN KEY (subscription_id) REFERENCES Subscription (id)
	);