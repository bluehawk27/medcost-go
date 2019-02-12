package store

const GetProviderByID = `SELECT * FROM provider WHERE id = ?`

const GetProviderByProviderID = `SELECT * FROM provider WHERE provider_id = ?`

const InsertProviderNolatlong = `INSERT INTO provider (provider_id, name, street, city, zip_code, state, hrr_description, created_at, lat, long) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
