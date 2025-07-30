# Books - Your Personal Digital Library
#### Video Demo: <URL HERE>
#### Description:

Alright, so I built this thing called "Books" - it's basically my attempt to create a digital library that doesn't suck. You know how most e-reader apps feel like they were designed in 2010? Yeah, I wanted to fix that.

## The Origin Story (Because Every Good Project Has One)

This whole thing started because I was frustrated with Microsoft's Aquile Reader. Don't get me wrong, it was actually pretty good - great UI, customizations, looked awesome. But here's the thing: I use Linux for most tasks, and while it worked great on Windows and Android, it was just another Microsoft whim. They offered free premium during development, then introduced paid plans and ads. Classic Microsoft move, right?

So I decided to recreate that great UI experience for EPUB reading, but make it cross-platform and actually free. I've been wanting to build this for a while, and CS50x gave me the perfect excuse to finally do it. It's still not complete (I have a lot of features to implement before I'd call it a decent EPUB reader), but it's getting there.

## What the Heck is This?

Think of it as your personal digital bookshelf, but actually good. You can upload your EPUB files, organize them, read them in a beautiful interface, and track your reading progress (yet to save on backend and load on page load, only works client-side for now and is lost on reload). It's like having a library in your pocket that includes everything, from cooking recipes to documentation or fun stories.

## The Tech Stack (Because Apparently That Matters)

I went full-stack on this one, which means I probably overcomplicated things but hey, that's how we learn, right?

**Frontend**: SvelteKit with TypeScript - I was originally planning to use Electron or Tauri for cross-platform, but SvelteKit looked promising and I had some web development experience and didn't want to overcomplicate and abandon this project. Plus, I can always wrap it in Tauri later if I want a native app (knew this from showcase on SvelteKit Discord - very welcoming and friendly).

**Backend**: Go (Golang) - Started with a simple Go auth tutorial from YouTube (Consulting Ninja), but it evolved into something much more complex. I appreciate Go's simplicity and the fact that it compiles to a single binary and is concurrent and fast (for the CS50 finance problem set I wanted to make concurrent requests and caching, but failed to do so).

**Database**: PostgreSQL - Reliable, well-documented, and supports all the features I needed (I didn't even know I could add custom functions and other SQL features like in C and Python - I never used it but good to know). Plus, it's free and I was introduced to SQLite3 in the course, so PostgreSQL it is.

**File Storage**: MinIO - I wanted S3-compatible storage (for reliable backups and industry standard). ChatGPT suggested this, and it's perfect for self-hosting. Companies use this, so if my app grows, I won't regret the choice.

**Authentication**: JWT tokens with refresh tokens, plus Google OAuth. The authentication flow was probably the hardest part to get right.

## The Architecture (AKA How I Made This Mess)

### Frontend Structure
- **`frontend/src/lib/`** - This is where all the good stuff lives:
  - `NavBar.svelte` - The navigation bar that stays at the top. Pretty straightforward, but it took me way too long to get the responsive design right, and I have slots here. I can import wherever I need and add different buttons. On the homepage there's a logo on the left, nothing in the center, and login/signup on the right. On the library page, it's the same navbar but with search in the center and upload/store on the right with a user icon on the far right. I can toggle the additional view that has filter purposes on the library page but it's hidden on the homepage.
  - `Login.svelte` & `Signup.svelte` - The authentication pages. I added some nice animations and validation because users deserve better than "invalid input" errors.
  - `BookCard.svelte` - Each book gets its own card with a cover image, title, author, and reading progress. I'm pretty proud of the hover effects on this one.
  - `BookReader.svelte` - The actual e-reader component. This was the hardest part - integrating EPUB.js to render books properly. It supports dark/light mode, chapter navigation, and progress tracking (only client-side for now as I haven't implemented the backend logic to store the progress, so it's lost on page reset for now).
  - `EpubUpload.svelte` - Handles file uploads and metadata extraction. I even added a worker pool for parsing EPUBs in the background so the UI doesn't freeze (ChatGPT helped me here).
  - `SearchBar.svelte` - A search component that looks good and actually works. Revolutionary, I know (got it from uiverse.io).

### Backend Structure
- **`backend/internal/handlers/`** - All the HTTP endpoints:
  - `router.go` - The traffic controller that routes requests to the right handlers.
  - `login.go` & `google.go` - Handle authentication. The Google OAuth flow was surprisingly straightforward once I figured out the redirect URLs.
  - `library.go` - Manages the book library, file operations, and user permissions. Users can only access their own books.
  - `upload.go` - Handles file uploads to MinIO. I added proper validation and error handling because users will upload anything.
  - `refresh.go` & `logout.go` - Token management. The refresh token system prevents users from getting logged out every 15 minutes.

- **`backend/internal/auth/`** - The security stuff:
  - `service.go` - The main authentication service. It's like the brain of the security system, where all the service functions and their helper functions live.
  - `jwt.go` - JWT token creation and validation. I used the v4 library because it's the latest and greatest.
  - `bcrypt.go` - Password hashing. Because storing plain text passwords is not secure.

- **`backend/internal/store/`** - Database operations (PostgreSQL)
- **`backend/internal/models/`** - Data structures
- **`backend/internal/middleware/`** - CORS basically outright rejects calls from domains other than the configured ones (localhost:4353 for dev and books.saurabpoudel.com.np for production) and JWT middleware validates the token.

### Infrastructure
- **`docker-compose.yml`** - Orchestrates all the services. I have:
  - Nginx reverse proxy (port 4353)
  - SvelteKit frontend
  - Go backend
  - PostgreSQL database
  - MinIO file storage
  - Automatic database backups (because losing data is not cool)

## Key Features (The Stuff That Actually Works)

### Authentication System
I built a proper authentication system with JWT tokens, refresh tokens, and Google OAuth. The JWT tokens expire every 15 minutes for security, but the refresh tokens last 7 days so users don't get annoyed. I also added email verification for new accounts because spam accounts are the worst, and I have Cloudflare DDoS protection enabled on the Cloudflare dashboard that adds some layer of protection as I don't think my app is foolproof secure.

### EPUB Reader
The reader component is probably the most complex part. It uses EPUB.js to render books, supports continuous scrolling (because pagination was annoying - I'll be adding an option but for now it's the only option), has dark/light mode, and tracks reading progress. I even added a table of contents and chapter navigation.

### File Management
Users can upload EPUB files, and the system automatically extracts metadata (title, author, cover image) using a web worker pool. This prevents the UI from freezing when parsing large books. Files are stored in MinIO with proper access controls - users can only access their own books.

### Responsive Design
The whole thing works on desktop, tablet, and mobile. I used Tailwind CSS because writing custom CSS is like writing poetry - beautiful but time-consuming (It's not true, CSS is good and I have some, and it might seem crazy to use with Tailwind but it works and I wanted to give it a shot. I'm scared to use CSS because I don't have good experience with it and frequently broke stuff, and it's my least favorite to work with).

## Design Decisions (And Why I Made Them)

### Why SvelteKit?
I was originally planning to use Electron or Tauri for cross-platform development, but I discovered SvelteKit through Fireship's videos and it looked promising. I had some web development experience, so it made sense to start there. Plus, I can always wrap it in Tauri later if I want a native app, and the community on Discord was friendly and I saw their projects with Tauri and SvelteKit.

### Why Go for the Backend?
Started with a simple Go auth tutorial from YouTube, but it evolved into something much more complex. I appreciate Go's simplicity and the fact that it compiles to a single binary. The standard library is comprehensive, and there's official support for things like MinIO (which to be fair I didn't know when I chose this - the main reason was file uploads and concurrency; I tried to use parallel processes in Python for the CS50 finance problem set to make HTTP requests concurrent but at that time I couldn't get it to work).

### Why PostgreSQL?
It's reliable, well-documented, and supports all the features I needed (that is SQLite3 features). Plus, it's free and doesn't have the weird limitations of some NoSQL databases (it's just on paper online when researching before app development).

### Why MinIO?
I wanted S3-compatible storage (from my research this is what companies use and it's reliable) for reliable backup for the future just in case. ChatGPT suggested this, and it's perfect for self-hosting. Companies use this, so if my app grows, I won't regret the choice.

## The Struggles (Because Nothing Works on the First Try)

### Authentication Flow - The Biggest Headache
This was probably the hardest part. I struggled the most with understanding how SvelteKit server-side, frontend (client browser), and backend communication worked with tokens. The issue was with refresh tokens - SvelteKit server would request a refresh from the backend, but the client browser had no idea the refresh token was used (and I delete used refresh tokens). I almost removed the refresh token implementation entirely, but decided to integrate everything to the frontend instead. Not sure how secure this is, but it works.

### EPUB Parsing
Getting EPUB.js to work properly was challenging. The documentation is sparse, and the API changes between versions. I ended up creating a worker pool to handle parsing in the background, which was actually a fun challenge.

### Theme Switching Issues
When implementing smooth scrolling, I broke the dark mode on EPUB.js. The issue was that it would only change theme once loaded and on scroll to next page (that was EPUB rendering and I was applying theme on rendition after registering theme). This took me a while to figure out. Here is my commit message
 `󰣇 ~/Documents/Final project   master  !? ❯ git checkout backend/                                                     23:01 
backend/       FETCH_HEAD     frontend/      HEAD           master         ORIG_HEAD      origin/master
f4b3611  -- [HEAD^^]  introduced & fixed bug on previous commit theme change fix (22 hours ago)
aeed0fd  -- [HEAD~3]  minor svelte.config csp fixes for styling as blobs for eupb (23 hours ago)`

### File Uploads
Handling large file uploads with proper progress tracking (only client-side progress tracking - I don't save it to the database and refresh wipes it out for now) and error handling was trickier than expected. I added file size limits, type validation, and proper cleanup for failed uploads.

### Docker Setup
Getting all the services to work together in Docker was like herding cats. The networking between containers, environment variables, and health checks took way longer than I'd like to admit, and ChatGPT came to help here.

## What I Learned (The Good Stuff)

This project taught me a lot about full-stack development, containerization, and building a real application that people might actually use. I learned about:

- Web Workers and background processing
- JWT authentication and security best practices like session-based and JSON web tokens (like how I get automatically logged in on sites I visit like Notion - I get it now)
- Docker orchestration and microservices (I still don't get it fully and it barely works on my current setup)
- EPUB file format and parsing
- Responsive design and user experience (with no users 😆)
- Database design and optimization

### The Backend Developer Meme
At one point, I looked at my app and realized the majority of my efforts were in the backend, while the frontend was still pretty basic (just as I left it weeks ago). I finally understood the meme about backend developers staring at frontend developers who just changed a button to yellow. I couldn't see as a client (when visiting the web app) what I did on the backend myself. But I don't regret it - I now have a solid understanding of cookies, authentication, and I can reuse this exact Go backend for future apps' authentication flow.

### Testing and Development Tools
I mostly used curl and scripts for testing (like the `checker.sh` script in the backend folder), but I discovered Postman and it's actually pretty useful for testing API endpoints. I also learned to appreciate Git more after losing 6 hours of work to VSCode crashing. Now I stage changes regularly and commit frequently.

## Future Improvements (If I Ever Get Around to It)

- Add support for other book formats (PDF, MOBI) - My dad wants to add his books (Nepali, Hindi, and Sanskrit) as he likes to read stories and poems (Ramayana, Puranas). I might need to add conversion or accept more formats.
- Implement book recommendations
- Add social features (sharing, reviews)
- Add reading statistics and analytics
- Implement offline reading support
- Add language options (as requested by my first client - my dad)
- Page flipping animations
- Global EPUB store (though I have no idea how to implement this legally - I don't want to illegally distribute authors' work, only books that authors themselves made available for free)
- Forgot password functionality
- Delete books feature

## How to Run This Thing

2. Open .env and fill in your configuration (as per the names)
3. Run `docker-compose --env-file ./.env up --build -d`
4. Visit `http://localhost:4353`

The setup is pretty straightforward thanks to Docker, but you'll need to configure Google OAuth and SendGrid for email verification if you want all the features to work.

## Deployment

I deployed this on Oracle Cloud Free Tier with Ubuntu server and CloudPanel which I installed and set up myself. I used my domain `books.saurabpoudel.com.np` (provided by Mercantile Communications free of charge for Nepali students/businesses to encourage digital adoption). The app runs on port 4353 with Nginx as a reverse proxy, forwarding `/api` requests to the Go backend and everything else to SvelteKit on port 3000, and behind another reverse proxy of CloudPanel that handles SSL certificates for me.

## Final Thoughts

Building this was a lot of fun, even though it took way longer than I expected. There's something satisfying about creating a tool that you'd actually want to use yourself. The code is probably not perfect (what code ever is?), but it works, it's secure (best of my knowledge), and it doesn't make me want to throw my computer out the window (because I couldn't hit X on Microsoft ads).

CS50x gave me the confidence to tackle this project, especially the data structures and SQL parts. I never worked with databases before. Data structures made me realize that “data” is just bits and bytes I could do anything with data - the formats are just different ways to address certain scenarios.

I've found YouTube creators incredibly helpful. Huge thanks to the creators who kept me sane during development - Christian Lempa (Nextcloud/nginx), NetworkChuck (Docker/Portainer), Fireship (for 5 years of experience in 100 seconds), Primeagen, Low Level Learning, Pirate Software, Consulting Ninja (first Go auth backend implementation), Bro Code (initial days of CS50x concepts on C), and many others. The open-source community, especially SvelteKit's Discord, has been amazing.

If you're reading this, thanks for checking out my project! Feel free to steal any ideas, but maybe give me credit if you do something cool with them. 😄

 