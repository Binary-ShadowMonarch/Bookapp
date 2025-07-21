// src/lib/epub-parser.worker.ts
import ePub from 'epubjs';

self.onmessage = async (event) => {
    const { file } = event.data;

    try {
        const arrayBuffer = await file.arrayBuffer();
        const book = ePub(arrayBuffer);
        await book.ready;

        const metadata = book.packaging.metadata;
        const coverUrl = await book.coverUrl();

        self.postMessage({
            success: true,
            payload: {
                id: crypto.randomUUID(),
                title: metadata.title || 'Unknown Title',
                author: metadata.creator || 'Unknown Author',
                coverUrl: coverUrl || null,
            },
        });
    } catch (error) {
        self.postMessage({
            success: false,
            error: (error as Error).message,
        });
    }
};