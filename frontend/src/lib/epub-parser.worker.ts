// this is a web worker that parses EPUB files in the background
// it runs in a separate thread so it doesn't freeze the main UI
import ePub from 'epubjs';

// the metadata structure that gets extracted from EPUB files
interface BookMetadata {
	id: string;
	title: string;
	author: string;
	coverUrl: string | null;
}

// cache parsed books to avoid re-parsing the same file multiple times
// this makes the app faster when users upload the same book again
const bookCache = new Map<string, BookMetadata>();

// this is the main message handler for the worker
// it receives files from the main thread and sends back parsed metadata
self.onmessage = async (event: MessageEvent) => {
	const { file }: { file: File } = event.data;

	try {
		// console.log('DEBUG: Worker received file for parsing:', file.name);

		// create a cache key from file size and name
		// this helps identify if we've parsed this exact file before
		const cacheKey = `${file.name}-${file.size}`;

		// check if we've already parsed this file
		if (bookCache.has(cacheKey)) {
			// console.log('DEBUG: Found cached metadata for:', file.name);
			self.postMessage({
				success: true,
				payload: bookCache.get(cacheKey)
			});
			return;
		}

		// convert the file to an array buffer so we can parse it
		const arrayBuffer = await file.arrayBuffer();
		const book = ePub(arrayBuffer);

		// wait for the book to be fully loaded before extracting metadata
		await book.ready;

		// get the metadata from the EPUB file
		const metadata = book.packaging?.metadata || {};

		// create the book metadata object
		const bookMeta: BookMetadata = {
			id: crypto.randomUUID(), // generate a unique ID
			title: metadata.title || file.name.replace('.epub', '') || 'Unknown Title',
			author: metadata.creator || 'Unknown Author',
			coverUrl: null
		};

		// try to extract the cover image with a timeout
		// sometimes cover extraction can hang, so I add a 2-second timeout
		try {
			const coverPromise = book.coverUrl();
			const timeoutPromise = new Promise<never>((_, reject) =>
				setTimeout(() => reject(new Error('timeout')), 2000)
			);
			const coverUrl = (await Promise.race([coverPromise, timeoutPromise])) as string;
			bookMeta.coverUrl = coverUrl;
			// console.log('DEBUG: Successfully extracted cover for:', file.name);
		} catch {
			// console.log('DEBUG: No cover found or timeout for:', file.name);
			// no cover available or timeout occurred - coverUrl remains null
		}

		// save the result in the cache for future use
		bookCache.set(cacheKey, bookMeta);

		// send the parsed metadata back to the main thread
		// console.log('DEBUG: Sending parsed metadata for:', file.name);
		self.postMessage({
			success: true,
			payload: bookMeta
		});
	} catch (error) {
		console.error('DEBUG: Error parsing EPUB in worker:', error);
		// send error back to main thread
		self.postMessage({
			success: false,
			error: (error as Error).message
		});
	}
};
