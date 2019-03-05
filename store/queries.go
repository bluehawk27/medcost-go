package store

// PROVIDER Queries
const getProviderByID = `SELECT * FROM provider WHERE id = ? LIMIT 1`

const getProviderByProviderID = `SELECT * FROM provider WHERE provider_id = ? LIMIT 1`

const insertProvider = `INSERT INTO provider (provider_id, name, street, city, zip_code, state, hrr_description, created_at, latitude, longitude) VALUES (:provider_id, :name, :street, :city, :zip_code, :state, :hrr_description, :created_at, :latitude, :longitude)`

// DRG Queries
const getDrgByDRGID = `SELECT * FROM diagnostic_related_group WHERE code = ? LIMIT 1`

const getDrgByID = `SELECT * FROM diagnostic_related_group WHERE id = ? LIMIT 1`

const insertDrg = `INSERT INTO diagnostic_related_group (name, code) VALUES (:name, :code)`

// Inpatient Queries
const getInpatientServiceByProviderID = `SELECT * FROM inpatient WHERE provider_id = ? LIMIT 1`

const getInpatientServiceByID = `SELECT * FROM inpatient WHERE id = ? LIMIT 1`

const getInpatientServiceByProvIDDrgIDYear = `SELECT * FROM inpatient WHERE provider_id = ? AND drg_id = ? AND year = ?`

const insertInpatientService = `INSERT INTO inpatient (provider_id, drg_id, total_discharges, avg_covered_charges, avg_total_payments, avg_medicare_payments, year, created_at) VALUES (:provider_id, :drg_id, :total_discharges, :avg_covered_charges, :avg_total_payments, :avg_medicare_payments, :year, :created_at)`

// APC Queries
const getApcByAPCID = `SELECT * FROM ambulatory_payment_classification WHERE code = ? LIMIT 1`

const getApcByID = `SELECT * FROM ambulatory_payment_classification WHERE id = ? LIMIT 1`

const insertApc = `INSERT INTO ambulatory_payment_classification (name, code) VALUES (:name, :code)`

// Outpatient Queries
const getOutpatientServiceByID = `SELECT * FROM outpatient WHERE id = ? LIMIT 1`

const getOutpatientServiceByProviderID = `SELECT * FROM outpatient WHERE provider_id = ? LIMIT 1`

const getOutpatientServiceByProvIDApcIDYear = `SELECT * FROM outpatient WHERE provider_id = ? AND apc_id = ? AND year = ?`

const insertOutpatientService = `INSERT INTO outpatient (provider_id, apc_id, services_count, avg_est_submitted_charges, avg_total_payments, year, created_at) VALUES (:provider_id, :apc_id, :services_count, :avg_est_submitted_charges, :avg_total_payments, :year, :created_at)`
