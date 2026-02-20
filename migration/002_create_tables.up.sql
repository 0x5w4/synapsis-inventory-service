START TRANSACTION;

CREATE TABLE IF NOT EXISTS "products" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "stock" INT NOT NULL,
    "price" DECIMAL(10, 2) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "reservations" (
    "id" SERIAL PRIMARY KEY,
    "product_id" INT NOT NULL,
    "order_id" INT NOT NULL,
    "quantity" INT NOT NULL,
    "status" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "fk_reservations_product_id_products" FOREIGN KEY ("product_id") REFERENCES "products"("id") ON DELETE RESTRICT
);

COMMIT;
