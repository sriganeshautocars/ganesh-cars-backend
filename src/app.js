import express from "express";
import cors from "cors";
import authRoutes from "./routes/authRoutes.js";
import carRoutes from "./routes/carRoutes.js";
import uploadRoute from "./routes/uploadRoute.js";

const app = express();
app.use(cors());
app.use(express.json({ limit: "10mb" }));

app.use("/api/auth", authRoutes);
app.use("/api/cars", carRoutes);
app.use("/api/upload", uploadRoute);

export default app;
