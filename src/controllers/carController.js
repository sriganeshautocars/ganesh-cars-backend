import { pool } from "../db.js";

export const addCar = async (req, res) => {
    try {
        const car = req.body;

        const query = `
      INSERT INTO cars (
        thumbnail, brand, name, variant, km_driven, fuel_type, body_type,
        transmission_type, price, location, insurance, no_of_seats, reg_number,
        ownership, engine_displacement, highway_mileage, make_year, reg_year,
        features, specifications, images, created_at, updated_at
      )
      VALUES (
        $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,
        $19,$20,$21,NOW(),NOW()
      ) RETURNING *;
    `;

        const values = [
            car.thumbnail, car.brand, car.name, car.variant, car.km_driven,
            car.fuel_type, car.body_type, car.transmission_type, car.price,
            car.location, car.insurance, car.no_of_seats, car.reg_number,
            car.ownership, car.engine_displacement, car.highway_mileage,
            car.make_year, car.reg_year,
            JSON.stringify(car.features), JSON.stringify(car.specifications),
            JSON.stringify(car.images)
        ];

        const result = await pool.query(query, values);
        res.status(201).json(result.rows[0]);
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};

export const getCars = async (req, res) => {
    try {
        const result = await pool.query(`
            SELECT
                id,
                thumbnail,
                brand,
                name,
                variant,
                km_driven,
                fuel_type,
                body_type,
                transmission_type,
                price,
                location,
                insurance,
                no_of_seats,
                reg_number,
                ownership,
                engine_displacement,
                highway_mileage,
                make_year,
                reg_year,
                created_at,
                updated_at
            FROM cars
            ORDER BY created_at DESC
        `);
        res.json(result.rows);
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};

export const getCarById = async (req, res) => {
    try {
        const { id } = req.params;
        const result = await pool.query("SELECT * FROM cars WHERE id = $1", [id]);
        if (result.rowCount === 0) return res.status(404).json({ message: "Car not found" });
        res.json(result.rows[0]);
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};

export const updateCar = async (req, res) => {
    try {
        const { id } = req.params;
        const fields = Object.keys(req.body);
        const values = Object.values(req.body);

        if (fields.length === 0) return res.status(400).json({ message: "No fields to update" });

        const setClause = fields.map((f, i) => `${f} = $${i + 1}`).join(", ");
        const query = `UPDATE cars SET ${setClause}, updated_at = NOW() WHERE id = $${fields.length + 1} RETURNING *`;

        const result = await pool.query(query, [...values, id]);
        res.json(result.rows[0]);
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};

// Delete car by ID
export const deleteCar = async (req, res) => {
    const { id } = req.params;

    try {
        const result = await pool.query("DELETE FROM cars WHERE id = $1 RETURNING *", [id]);

        if (result.rowCount === 0) {
            return res.status(404).json({ message: "Car not found" });
        }

        res.json({ message: "Car deleted successfully", deletedCar: result.rows[0] });
    } catch (error) {
        console.error("Error deleting car:", error);
        res.status(500).json({ message: "Internal server error" });
    }
};

