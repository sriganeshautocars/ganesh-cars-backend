import express from "express";
import { addCar, getCars, getCarById, updateCar, deleteCar, updateCarHoldStatus } from "../controllers/carController.js";
import { verifyToken } from "../middleware/auth.js";

const router = express.Router();

router.get("/", getCars);
router.get("/:id", getCarById);
router.post("/", verifyToken, addCar);
router.put("/:id", verifyToken, updateCar);
router.delete("/:id", verifyToken, deleteCar);
router.patch('/:id/hold', verifyToken, updateCarHoldStatus)

export default router;
