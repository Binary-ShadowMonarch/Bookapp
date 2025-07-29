<!-- src/lib/EpubUpload.svelte -->
<script lang="ts">
	import { onMount } from 'svelte';
	import { Upload } from 'lucide-svelte';
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

	interface Props {
		onBookAdded: (book: Book) => void;
	}
	let { onBookAdded }: Props = $props();

	let fileInput: HTMLInputElement;
	let isUploading = $state(false);

	// Parse metadata from an EPUB File object
	// Improved parseEpubMetadata function that ensures blob URLs are converted
	async function parseEpubMetadata(file: File): Promise<{
		id: string;
		title: string;
		author: string;
		coverUrl: string | null;
	}> {
		return new Promise((resolve, reject) => {
			const reader = new FileReader();
			reader.onload = async (e) => {
				try {
					const arrayBuffer = e.target?.result as ArrayBuffer;
					const ePub = (await import('epubjs')).default;
					const book = ePub(arrayBuffer);
					await book.ready;

					const meta = {
						id: crypto.randomUUID(),
						title: book.packaging.metadata.title || 'Unknown Title',
						author: book.packaging.metadata.creator || 'Unknown Author',
						coverUrl: null as string | null
					};

					try {
						const cover = await book.coverUrl();
						if (cover) {
							// Always convert blob URLs to data URLs to avoid CSP issues
							if (cover.startsWith('blob:')) {
								const response = await fetch(cover);
								const blob = await response.blob();
								const dataUrl = await new Promise<string>((resolve) => {
									const fileReader = new FileReader();
									fileReader.onload = () => resolve(fileReader.result as string);
									fileReader.readAsDataURL(blob);
								});
								meta.coverUrl = dataUrl;

								// Clean up the blob URL
								URL.revokeObjectURL(cover);
							} else {
								meta.coverUrl = cover;
							}
						}
					} catch (error) {
						console.warn('Could not extract book cover:', error);
						// No cover available, use default
					}

					resolve(meta);
				} catch (err) {
					console.error('Error parsing EPUB metadata:', err);
					reject(err);
				}
			};
			reader.onerror = reject;
			reader.readAsArrayBuffer(file);
		});
	}
	// On mount: load existing EPUBs
	onMount(async () => {
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
					image: meta.coverUrl || '/default.webp',
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
				image: meta.coverUrl || '/default.webp',
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
