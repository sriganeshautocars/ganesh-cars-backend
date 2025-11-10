import { S3Client, PutObjectCommand } from "@aws-sdk/client-s3";
import { v4 as uuidv4 } from "uuid";
import dotenv from "dotenv";

dotenv.config();

const s3 = new S3Client({
    region: "auto",
    endpoint: process.env.CLOUDFLARE_R2_ENDPOINT,
    credentials: {
        accessKeyId: process.env.R2_ACCESS_KEY_ID,
        secretAccessKey: process.env.R2_SECRET_ACCESS_KEY,
    },
});

export const uploadToR2 = async (file) => {
    const key = `cars/${uuidv4()}-${file?.originalname}`;
    const command = new PutObjectCommand({
        Bucket: process.env.R2_BUCKET_NAME,
        Key: key,
        Body: file?.buffer,
        ContentType: file?.mimetype,
    });

    await s3.send(command);
    return `${process.env.R2_PUBLIC_URL}/${key}`;
};

export const uploadMultipleToR2 = async (files) => {
    const uploadPromises = files.map((file) => uploadToR2(file));
    return Promise.all(uploadPromises); // returns array of image URLs
};

