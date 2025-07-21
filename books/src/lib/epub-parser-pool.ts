// src/lib/epub-parser-pool.ts
import { browser } from '$app/environment';

export class EpubParserPool {
    private workers: Worker[] = [];
    private queue: Array<{
        file: File;
        resolve: (value: any) => void;
        reject: (error: any) => void;
    }> = [];
    private busyWorkers = new Set<Worker>();
    private workerCount: number;
    private initialized = false;

    constructor(workerCount = browser ? (navigator.hardwareConcurrency || 4) : 1) {
        this.workerCount = Math.min(workerCount, 8);
        if (browser) {
            this.initializeWorkers();
        }
    }

    private initializeWorkers() {
        if (!browser || this.initialized) return;

        for (let i = 0; i < this.workerCount; i++) {
            try {
                const worker = new Worker(
                    new URL('./epub-parser.worker.ts', import.meta.url),
                    { type: 'module' }
                );

                worker.onmessage = (event) => {
                    const { success, payload, error } = event.data;
                    const task = this.queue.shift();

                    if (task) {
                        if (success) {
                            task.resolve(payload);
                        } else {
                            task.reject(new Error(error));
                        }
                    }

                    // Mark worker as available
                    this.busyWorkers.delete(worker);
                    this.processQueue();
                };

                worker.onerror = (error) => {
                    console.error('Worker error:', error);
                    this.busyWorkers.delete(worker);
                    this.fallbackParse();
                };

                this.workers.push(worker);
            } catch (error) {
                console.warn('Failed to create worker, falling back to main thread');
                break;
            }
        }
        this.initialized = true;
    }

    private processQueue() {
        if (this.queue.length === 0) return;

        // Find available worker
        const availableWorker = this.workers.find(worker => !this.busyWorkers.has(worker));

        if (availableWorker && this.queue.length > 0) {
            const task = this.queue[0];
            this.busyWorkers.add(availableWorker);
            availableWorker.postMessage({ file: task.file });
        }
    }

    private async fallbackParse() {
        if (this.queue.length === 0) return;

        const task = this.queue.shift();
        if (!task) return;

        try {
            const result = await this.parseInMainThread(task.file);
            task.resolve(result);
        } catch (error) {
            task.reject(error);
        }

        // Continue processing queue
        if (this.queue.length > 0) {
            setTimeout(() => this.fallbackParse(), 0);
        }
    }

    private async parseInMainThread(file: File): Promise<{
        id: string;
        title: string;
        author: string;
        coverUrl: string | null;
    }> {
        const ePub = (await import('epubjs')).default;
        const arrayBuffer = await file.arrayBuffer();
        const book = ePub(arrayBuffer);

        await book.ready;

        const metadata = book.packaging?.metadata || {};

        let coverUrl = null;
        try {
            const coverPromise = book.coverUrl();
            const timeoutPromise = new Promise<never>((_, reject) =>
                setTimeout(() => reject(new Error('timeout')), 2000)
            );
            coverUrl = await Promise.race([coverPromise, timeoutPromise]) as string;
        } catch {
            // No cover available
        }

        return {
            id: crypto.randomUUID(),
            title: metadata.title || file.name.replace('.epub', '') || 'Unknown Title',
            author: metadata.creator || 'Unknown Author',
            coverUrl
        };
    }

    parseEpub(file: File): Promise<{
        id: string;
        title: string;
        author: string;
        coverUrl: string | null;
    }> {
        if (!browser) {
            return Promise.resolve({
                id: crypto.randomUUID(),
                title: file.name.replace('.epub', '') || 'Unknown Title',
                author: 'Unknown Author',
                coverUrl: null
            });
        }

        if (!this.initialized) {
            this.initializeWorkers();
        }

        return new Promise((resolve, reject) => {
            this.queue.push({ file, resolve, reject });

            if (this.workers.length > 0) {
                this.processQueue();
            } else {
                this.fallbackParse();
            }
        });
    }

    parseMultiple(files: File[]): Promise<Array<{
        id: string;
        title: string;
        author: string;
        coverUrl: string | null;
    }>> {
        return Promise.all(files.map(file => this.parseEpub(file)));
    }

    terminate() {
        this.workers.forEach(worker => worker.terminate());
        this.workers = [];
        this.queue = [];
        this.busyWorkers.clear();
        this.initialized = false;
    }
}

// Singleton instance
export const epubParser = new EpubParserPool();