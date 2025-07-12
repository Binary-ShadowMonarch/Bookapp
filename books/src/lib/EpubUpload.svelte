<script lang="ts">
  
  type BookStatus = 'read' | 'unread' | 'finished';
  
  interface Book {
    id: string;
    title: string;
    author: string;
    image: string;
    status: BookStatus;
    completion: number;
  }

  interface Props {
    onBookAdded: (book: Book) => void;
  }
  
  let { onBookAdded }: Props = $props();
  
  let fileInput: HTMLInputElement;
  let isUploading = $state(false);

  async function handleFileSelect(event: Event) {
    const target = event.target as HTMLInputElement;
    const file = target.files?.[0];
    if (!file) return;
    
    if (file.type !== 'application/epub+zip' && !file.name.endsWith('.epub')) {
      alert('Please select a valid EPUB file');
      return;
    }
    
    await processEpubFile(file);
  }
  
  async function processEpubFile(file: File) {
    isUploading = true;
    
    try {
      // Parse EPUB metadata
      const metadata = await parseEpubMetadata(file);
      
      // Upload to MinIO
      await uploadToMinio(file, metadata.id);
      
      // Create book object matching your interface
      const book: Book = {
        id: metadata.id,
        title: metadata.title || 'Unknown Title',
        author: metadata.author || 'Unknown Author',
        image: metadata.coverUrl || '/default.webp',
        status: 'unread',
        completion: 0
      };
      
      // Call the parent callback
      onBookAdded(book);
      
      // Reset file input
      fileInput.value = '';
      
    } catch (error) {
      console.error('Upload failed:', error);
      alert('Failed to upload EPUB file');
    } finally {
      isUploading = false;
    }
  }
  
  async function parseEpubMetadata(file: File): Promise<{
    id: string;
    title: string;
    author: string;
    coverUrl: string | null;
  }> {
    return new Promise(async (resolve, reject) => {
      const reader = new FileReader();
      reader.onload = async (e) => {
        try {
          const arrayBuffer = e.target?.result as ArrayBuffer;
          // Dynamic import to avoid build issues
          const ePub = await import('epubjs');
          
          const book = ePub.default(arrayBuffer);
          await book.ready;
          
          const metadata = {
            id: crypto.randomUUID(),
            title: book.packaging.metadata.title || 'Unknown Title',
            author: book.packaging.metadata.creator || 'Unknown Author',
            coverUrl: null as string | null
          };
          
          // Try to get cover image
          try {
            const coverUrl = await book.coverUrl();
            if (coverUrl) {
              metadata.coverUrl = coverUrl;
            }
          } catch (coverError) {
            console.log('No cover found, using default');
          }
          
          resolve(metadata);
        } catch (error) {
          reject(error);
        }
      };
      reader.onerror = reject;
      reader.readAsArrayBuffer(file);
    });
  }
  
  async function uploadToMinio(file: File, bookId: string) {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('bookId', bookId);
    
    const response = await fetch('/api/upload-epub', {
      method: 'POST',
      body: formData
    });
    
    if (!response.ok) {
      throw new Error('Upload failed');
    }
    
    return await response.json();
  }
  
  function triggerFileSelect() {
    fileInput.click();
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
  aria-label="Upload EPUB"
  class="flex cursor-pointer items-center gap-2 rounded-full bg-green-500 px-3 py-1.5 text-sm font-semibold text-white shadow-sm transition hover:bg-green-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600 disabled:opacity-50 disabled:cursor-not-allowed"
>
  {#if isUploading}
    <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
  {:else}
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
      <polyline points="17 8 12 3 7 8" />
      <line x1="12" x2="12" y1="3" y2="15" />
    </svg>
  {/if}
  <span class="hidden sm:inline">
    {isUploading ? 'Uploading...' : 'Upload'}
  </span>
</button>