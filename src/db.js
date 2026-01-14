import pkg from "pg";
import dotenv from "dotenv";
dotenv.config();

const { Pool } = pkg;

const pool = globalThis.__pgPool || new Pool({
    connectionString: process.env.DATABASE_URL,
    ssl: { rejectUnauthorized: false }, // Neon requires SSL
});

globalThis.__pgPool = pool;

export { pool };
