const express = require('express');
const { Pool } = require('pg');
const cors = require('cors');

const app = express();
app.use(express.json());
app.use(cors());

// Database configuration
const pool = new Pool({
  host: process.env.DB_HOST || 'postgres',
  port: process.env.DB_PORT || 5432,
  user: process.env.DB_USER || 'postgres',
  password: process.env.DB_PASS || 'password',
  database: process.env.DB_NAME || 'microservices_db',
});

// Wait for database and connect
const connectDb = async () => {
    let retries = 5;
    while(retries) {
        try {
            await pool.query('SELECT 1');
            console.log("Connected to PostgreSQL DB.");
            break;
        } catch (err) {
            console.log(err);
            retries -= 1;
            console.log(`Unable to connect to DB, retrying... (${retries} retries left)`);
            await new Promise(res => setTimeout(res, 5000));
        }
    }
}
connectDb();


app.get('/health', (req, res) => {
    res.status(200).json({ status: 'UP', service: 'item-service' });
});

// GET /items
app.get('/items', async (req, res) => {
    try {
        const result = await pool.query('SELECT * FROM item_schema.items ORDER BY id DESC');
        res.json(result.rows);
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: 'Internal Server Error' });
    }
});

// GET /items/:id
app.get('/items/:id', async (req, res) => {
    try {
        const { id } = req.params;
        const result = await pool.query('SELECT * FROM item_schema.items WHERE id = $1', [id]);
        
        if (result.rows.length === 0) {
            return res.status(404).json({ error: 'Item not found' });
        }
        res.json(result.rows[0]);
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: 'Internal Server Error' });
    }
});

// POST /items
app.post('/items', async (req, res) => {
    try {
        const { name } = req.body;
        if (!name) {
            return res.status(400).json({ error: 'Name is required' });
        }
        
        const result = await pool.query(
            'INSERT INTO item_schema.items (name) VALUES ($1) RETURNING *',
            [name]
        );
        res.status(201).json(result.rows[0]);
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: 'Internal Server Error' });
    }
});

const PORT = process.env.PORT || 8081;
app.listen(PORT, () => {
    console.log(`Item server running on port ${PORT}`);
});
