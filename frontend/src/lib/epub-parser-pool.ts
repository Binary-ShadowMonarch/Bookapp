// this is a worker pool for parsing EPUB files
// it uses web workers to parse EPUBs in the background so the UI doesn't freeze
import { browser } from '$app/environment';

export class EpubParserPool {
	// array of web workers that can parse EPUBs
	private workers: Worker[] = [];
	// queue of files waiting to be parsed
	private queue: Array<{
		file: File;
		resolve: (value: any) => void;
		reject: (error: any) => void;
	}> = [];
	// track which workers are currently busy
	private busyWorkers = new Set<Worker>();
	// how many workers to create
	private workerCount: number;
	// prevent double initialization
	private initialized = false;

	// create the worker pool with a reasonable number of workers
	constructor(workerCount = browser ? navigator.hardwareConcurrency || 4 : 1) {
		console.log('DEBUG: Creating EPUB parser pool with', workerCount, 'workers');
		this.workerCount = Math.min(workerCount, 8); // cap at 8 workers max
		if (browser) {
			this.initializeWorkers();
		}
	}

	// create the web workers that will do the actual EPUB parsing
	private initializeWorkers() {
		if (!browser || this.initialized) return;

		console.log('DEBUG: Initializing', this.workerCount, 'EPUB parser workers');

		for (let i = 0; i < this.workerCount; i++) {
			try {
				// create a new web worker from the epub-parser.worker.ts file
				const worker = new Worker(new URL('./epub-parser.worker.ts', import.meta.url), {
					type: 'module'
				});

				// handle messages from the worker
				worker.onmessage = (event) => {
					const { success, payload, error } = event.data;
					const task = this.queue.shift(); // get the next task from the queue

					if (task) {
						if (success) {
							console.log('DEBUG: EPUB parsed successfully by worker');
							task.resolve(payload);
						} else {
							console.error('DEBUG: EPUB parsing failed in worker:', error);
							task.reject(new Error(error));
						}
					}

					// mark this worker as available again
					this.busyWorkers.delete(worker);
					this.processQueue(); // try to process more files
				};

				// handle worker errors
				worker.onerror = (error) => {
					console.error('DEBUG: Worker error:', error);
					this.busyWorkers.delete(worker);
					this.fallbackParse(); // fall back to main thread parsing
				};

				this.workers.push(worker);
			} catch (error) {
				console.warn('DEBUG: Failed to create worker, falling back to main thread');
				break;
			}
		}
		this.initialized = true;
	}

	// process the queue of files waiting to be parsed
	private processQueue() {
		if (this.queue.length === 0) return;

		// find a worker that's not busy
		const availableWorker = this.workers.find((worker) => !this.busyWorkers.has(worker));

		if (availableWorker && this.queue.length > 0) {
			const task = this.queue[0];
			this.busyWorkers.add(availableWorker);
			console.log('DEBUG: Sending file to worker for parsing');
			availableWorker.postMessage({ file: task.file });
		}
	}

	// fallback to parsing in the main thread if workers fail
	private async fallbackParse() {
		if (this.queue.length === 0) return;

		const task = this.queue.shift();
		if (!task) return;

		console.log('DEBUG: Falling back to main thread parsing');
		try {
			const result = await this.parseInMainThread(task.file);
			task.resolve(result);
		} catch (error) {
			task.reject(error);
		}

		// continue processing the queue
		if (this.queue.length > 0) {
			setTimeout(() => this.fallbackParse(), 0);
		}
	}

	// parse an EPUB file in the main thread (fallback method)
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
			coverUrl = (await Promise.race([coverPromise, timeoutPromise])) as string;
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

	parseMultiple(files: File[]): Promise<
		Array<{
			id: string;
			title: string;
			author: string;
			coverUrl: string | null;
		}>
	> {
		return Promise.all(files.map((file) => this.parseEpub(file)));
	}

	terminate() {
		this.workers.forEach((worker) => worker.terminate());
		this.workers = [];
		this.queue = [];
		this.busyWorkers.clear();
		this.initialized = false;
	}
}

// Singleton instance
export const epubParser = new EpubParserPool();
