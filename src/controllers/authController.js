import { pool } from "../db.js";
import jwt from "jsonwebtoken";
import dotenv from "dotenv";
import { comparePassword } from "../utils/hash.js";

dotenv.config();

export const login = async (req, res) => {
    const { username, password } = req.body;

    try {
        const result = await pool.query("SELECT * FROM users WHERE username = $1", [username]);
        if (result.rowCount === 0) return res.status(401).json({ message: "Invalid credentials" });

        const user = result.rows[0];
        const match = await comparePassword(password, user.password);
        if (!match) return res.status(401).json({ message: "Invalid credentials" });

        const token = jwt.sign({ id: user.id, username: user.username }, process.env.JWT_SECRET, {
            expiresIn: "1d",
        });

        res.json({ token });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};
