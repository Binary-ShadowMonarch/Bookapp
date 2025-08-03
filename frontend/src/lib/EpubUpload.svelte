<!-- this is the EPUB upload component -->
<!-- it handles uploading new books and displaying the library -->
<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Upload } from 'lucide-svelte';

	// these are the types I use for books in my app
	export type BookStatus = 'read' | 'unread' | 'finished';
	export interface Book {
		id: string;
		title: string;
		author: string;
		image: string;
		status: BookStatus;
		completion: number;
		fileUrl: string;
	}

	// these are the types that come from my backend API
	interface FileInfo {
		id: string;
		name: string;
		url: string;
		size?: number;
		mimeType?: string;
	}
	interface LibraryResponse {
		files: FileInfo[];
		total: number;
	}

	// props that get passed from the parent component
	interface Props {
		onBookAdded: (book: Book) => void; // callback when a new book is added
	}
	let { onBookAdded }: Props = $props();

	// DOM reference and state
	let fileInput: HTMLInputElement; // reference to the file input element
	let isUploading = $state(false); // shows loading state during upload

	// this function extracts metadata from an EPUB file
	// it reads the title, author, and cover image from the EPUB
	async function parseEpubMetadata(file: File): Promise<{
		id: string;
		title: string;
		author: string;
		coverUrl: string | null;
	}> {
		console.log('DEBUG: Parsing EPUB metadata for:', file.name);
		return new Promise((resolve, reject) => {
			const reader = new FileReader();
			reader.onload = async (e) => {
				try {
					const arrayBuffer = e.target?.result as ArrayBuffer;
					const ePub = (await import('epubjs')).default;
					const book = ePub(arrayBuffer);
					await book.ready;

					// extract basic metadata
					const meta = {
						id: crypto.randomUUID(), // generate a unique ID
						title: book.packaging.metadata.title || 'Unknown Title',
						author: book.packaging.metadata.creator || 'Unknown Author',
						coverUrl: null as string | null
					};

					// try to extract the cover image
					try {
						const cover = await book.coverUrl();
						if (cover) {
							// convert blob URLs to data URLs to avoid security issues
							if (cover.startsWith('blob:')) {
								const response = await fetch(cover);
								const blob = await response.blob();
								const dataUrl = await new Promise<string>((resolve) => {
									const fileReader = new FileReader();
									fileReader.onload = () => resolve(fileReader.result as string);
									fileReader.readAsDataURL(blob);
								});
								meta.coverUrl = dataUrl;

								// clean up the blob URL to free memory
								URL.revokeObjectURL(cover);
							} else {
								meta.coverUrl = cover;
							}
						}
					} catch (error) {
						console.warn('DEBUG: Could not extract book cover:', error);
						// no cover available, will use default
					}

					console.log('DEBUG: EPUB metadata extracted:', meta);
					resolve(meta);
				} catch (err) {
					console.error('DEBUG: Error parsing EPUB metadata:', err);
					reject(err);
				}
			};
			reader.onerror = reject;
			reader.readAsArrayBuffer(file);
		});
	}

	// load existing books when component mounts
	onMount(async () => {
		console.log('DEBUG: EpubUpload component mounting');
		try {
			const res = await fetch('/api/protected/library', {
				credentials: 'include'
			});
			if (!res.ok) return;

			// 1) Destructure the JSON into its `files` array
			const { files }: LibraryResponse = await res.json();

			// 2) Iterate each FileInfo object
			for (const fileInfo of files) {
				const url = fileInfo.url;
				const blobRes = await fetch(url, { credentials: 'include' });
				if (!blobRes.ok) continue;

				const blob = await blobRes.blob();
				const filename = fileInfo.name; // you already have the name
				const file = new File([blob], filename, { type: blob.type });

				const meta = await parseEpubMetadata(file);
				const book: Book = {
					id: meta.id,
					title: meta.title,
					author: meta.author,
					image: meta.coverUrl || '/candle.webp',
					status: 'unread',
					completion: 0,
					fileUrl: url
				};

				onBookAdded(book);
			}
		} catch (e) {
			console.error('Error loading existing EPUBs:', e);
		}
	});

	function triggerFileSelect() {
		fileInput.click();
	}

	async function handleFileSelect(e: Event) {
		const target = e.target as HTMLInputElement;
		const file = target.files?.[0];
		if (!file) return;
		if (file.type !== 'application/epub+zip' && !file.name.endsWith('.epub')) {
			return alert('Please select a valid EPUB file');
		}

		isUploading = true;
		try {
			const resc = await fetch('/api/protected/profile', { method: 'GET', credentials: 'include' });
			if (!resc.ok) {
				goto('/login');
				return;
			}
			const meta = await parseEpubMetadata(file);
			const form = new FormData();
			form.append('file', file);

			const uploadRes = await fetch('/api/protected/upload', {
				method: 'POST',
				credentials: 'include',
				body: form
			});
			if (!uploadRes.ok) {
				const txt = await uploadRes.text();
				throw new Error(txt);
			}
			const { url: fileUrl } = await uploadRes.json();

			const book: Book = {
				id: meta.id,
				title: meta.title,
				author: meta.author,
				image: meta.coverUrl || '/candle.webp',
				status: 'unread',
				completion: 0,
				fileUrl
			};
			onBookAdded(book);
			fileInput.value = '';
		} catch (err) {
			console.error('Upload error:', err);
			alert((err as Error).message || 'Upload failed');
		} finally {
			isUploading = false;
		}
	}
</script>

<input
	bind:this={fileInput}
	type="file"
	accept=".epub,application/epub+zip"
	onchange={handleFileSelect}
	style="display: none"
/>
<button
	onclick={triggerFileSelect}
	disabled={isUploading}
	class="flex cursor-pointer items-center gap-2 rounded-full bg-green-500 px-3 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-green-600 disabled:opacity-50"
>
	{#if isUploading}
		<Upload /><span class="hidden sm:inline">Uploading</span>
	{:else}
		<Upload /><span class="hidden sm:inline">Upload</span>
	{/if}
</button>
