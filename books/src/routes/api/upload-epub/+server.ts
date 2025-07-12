// src/routes/api/upload-epub/+server.ts
import { json } from '@sveltejs/kit';
import { Client } from 'minio';
import type { RequestHandler } from './$types';
import { Buffer } from 'node:buffer'; // <-- ensure Node Buffer available

const minioClient = new Client({
  endPoint: 'localhost',
  port: 9000,
  useSSL: false,
  accessKey: 'admin',     // <-- match your Docker credentials
  secretKey: 'secretpassword'
});

const BUCKET_NAME = 'epub-books';

async function ensureBucket() {
  if (!(await minioClient.bucketExists(BUCKET_NAME))) {
    await minioClient.makeBucket(BUCKET_NAME);
  }
}

export const POST: RequestHandler = async ({ request }) => {
  await ensureBucket();

  const form = await request.formData();
  const file = form.get('file') as File;
  const bookId = String(form.get('bookId') || '');

  if (!file || !bookId) {
    return json({ error: 'File and bookId are required' }, { status: 400 });
  }

  // Convert browser File to Node.js Buffer
  const arrayBuffer = await file.arrayBuffer();
  const buffer = Buffer.from(arrayBuffer);
  const filename = `${bookId}.epub`;

  try {
    // omit size parameter—let MinIO infer if using Buffer
    await minioClient.putObject(
      BUCKET_NAME,
      filename,
      buffer,
      buffer.length,
      { 'Content-Type': file.type }
    );
  } catch (err: any) {
    console.error('Upload error code=', err.code, 'message=', err.message);
    return json(
      { error: `Upload failed with ${err.code}: ${err.message}` },
      { status: 500 }
    );
  }

  const fileUrl = `http://localhost:9000/${BUCKET_NAME}/${filename}`;
  return json({ success: true, fileUrl, filename });
};
