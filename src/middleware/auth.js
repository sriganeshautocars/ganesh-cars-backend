import jwt from "jsonwebtoken";
import dotenv from "dotenv";
dotenv.config();

export const verifyToken = (req, res, next) => {

    if (req.method === "OPTIONS") {
        return next();
    }

    const token = req.cookies.auth_token;
    if (!token) return res.status(401).json({ message: "Access denied: No token" });

    try {
        req.user = jwt.verify(token, process.env.JWT_SECRET);
        next();
    } catch {
        res.status(401).json({ message: "Invalid or expired token" });
    }
};
