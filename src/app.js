import express from "express";
import cors from "cors";
import cookieParser from "cookie-parser";
import authRoutes from "./routes/authRoutes.js";
import carRoutes from "./routes/carRoutes.js";
import uploadRoute from "./routes/uploadRoute.js";

const allowedOrigins = ['http://localhost:5173', 'https://ganesh-cars-frontend.vercel.app'];

const corsOptions = {
    origin: function (origin, callback) {
        // allow requests with no origin (like Postman & mobile apps)
        if (!origin) return callback(null, true);

        if (allowedOrigins.includes(origin)) {
            callback(null, true);
        } else {
            callback(new Error("Not allowed by CORS"));
        }
    },
    credentials: true,         // allow cookies & authorization headers
    methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"],
};


const app = express();
app.use(cors(corsOptions));
app.use(cookieParser());
app.use(express.json({ limit: "10mb" }));

app.use("/api/auth", authRoutes);
app.use("/api/cars", carRoutes);
app.use("/api/upload", uploadRoute);

export default app;
