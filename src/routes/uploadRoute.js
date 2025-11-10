import express from "express";
import multer from "multer";
import { uploadToR2, uploadMultipleToR2 } from "../services/r2Service.js";
import { verifyToken } from "../middleware/auth.js";

const router = express.Router();
const upload = multer({ storage: multer.memoryStorage() });

// Single file upload
router.post("/single", verifyToken, upload.single("image"), async (req, res) => {
    try {
        if (!req.file) return res.status(400).json({ message: "No file uploaded" });
        const imageUrl = await uploadToR2(req.file);
        res.json({ imageUrl });
    } catch (error) {
        console.error("Upload error:", error);
        res.status(500).json({ message: "Failed to upload image" });
    }
});

// Multiple file upload
router.post("/multiple", verifyToken, upload.array("images", 10), async (req, res) => {
    try {
        if (!req.files || req.files.length === 0) {
            return res.status(400).json({ message: "No files uploaded" });
        }

        const imageUrls = await uploadMultipleToR2(req.files);
        res.json({ imageUrls });
    } catch (error) {
        console.error("Multiple upload error:", error);
        res.status(500).json({ message: "Failed to upload images" });
    }
});

export default router;
