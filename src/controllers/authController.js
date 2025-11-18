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

        res.cookie("auth_token", token, { httpOnly: true, secure: false, sameSite: 'lax', maxAge: 24 * 60 * 60 * 1000 });
        res.json({ message: "Login successful", token });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};

export const logout = (req, res) => {
    res.clearCookie("auth_token", { httpOnly: true, secure: process.env.NODE_ENV === "production" });
    res.json({ message: "Logout successful" });
};

export const checkAuth = (req, res) => {
    const token = req.cookies.auth_token;

    // No cookie → not logged in
    if (!token) {
        return res.json({ loggedIn: false });
    }

    try {
        const decoded = jwt.verify(token, process.env.JWT_SECRET);

        return res.json({
            loggedIn: true,
            user: {
                id: decoded.id,
                email: decoded.email,
            }
        });
    } catch (err) {
        return res.json({ loggedIn: false });
    }
};

