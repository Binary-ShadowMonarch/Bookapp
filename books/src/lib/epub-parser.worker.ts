// src/lib/epub-parser.worker.ts
import ePub from 'epubjs';

interface BookMetadata {
    id: string;
    title: string;
    author: string;
    coverUrl: string | null;
}

// Cache parsed books to avoid re-parsing
const bookCache = new Map<string, BookMetadata>();

self.onmessage = async (event: MessageEvent) => {
    const { file }: { file: File } = event.data;

    try {
        // Create cache key from file size + name
        const cacheKey = `${file.name}-${file.size}`;

        if (bookCache.has(cacheKey)) {
            self.postMessage({
                success: true,
                payload: bookCache.get(cacheKey)
            });
            return;
        }

        const arrayBuffer = await file.arrayBuffer();
        const book = ePub(arrayBuffer);

        // Wait for book to be ready before parsing
        await book.ready;

        const metadata = book.packaging?.metadata || {};

        const bookMeta: BookMetadata = {
            id: crypto.randomUUID(),
            title: metadata.title || file.name.replace('.epub', '') || 'Unknown Title',
            author: metadata.creator || 'Unknown Author',
            coverUrl: null
        };

        // Try to get cover (with timeout)
        try {
            const coverPromise = book.coverUrl();
            const timeoutPromise = new Promise<never>((_, reject) =>
                setTimeout(() => reject(new Error('timeout')), 2000)
            );
            const coverUrl = await Promise.race([coverPromise, timeoutPromise]) as string;
            bookMeta.coverUrl = coverUrl;
        } catch {
            // No cover or timeout - coverUrl remains null
        }

        // Cache the result
        bookCache.set(cacheKey, bookMeta);

        // Send single final result
        self.postMessage({
            success: true,
            payload: bookMeta
        });

    } catch (error) {
        self.postMessage({
            success: false,
            error: (error as Error).message
        });
    }
};