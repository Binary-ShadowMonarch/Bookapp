--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5 (Debian 17.5-1.pgdg120+1)
-- Dumped by pg_dump version 17.5 (Debian 17.5-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: email_verifications; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.email_verifications (
    email text NOT NULL,
    hashed_password text NOT NULL,
    code text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.email_verifications OWNER TO admin;

--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.refresh_tokens (
    id integer NOT NULL,
    token text NOT NULL,
    user_id integer NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.refresh_tokens OWNER TO admin;

--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.refresh_tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.refresh_tokens_id_seq OWNER TO admin;

--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email text NOT NULL,
    hashed_password text,
    provider text NOT NULL,
    provider_id text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.users OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: refresh_tokens id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: email_verifications; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.email_verifications (email, hashed_password, code, expires_at, created_at) FROM stdin;
\.


--
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.refresh_tokens (id, token, user_id, expires_at, created_at) FROM stdin;
5	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InNhdXJhYnp6MzE2NkBnbWFpbC5jb20iLCJzdWIiOiJzYXVyYWJ6ejMxNjZAZ21haWwuY29tIiwiZXhwIjoxNzUzOTk2MDk2LCJpYXQiOjE3NTMzOTEyOTYsImp0aSI6IjAwYWIxZDdjLTAzNzQtNDljMC1hOTUwLWJlNDNiZWZjMzlmOSJ9.0uK_K7w_qyhPXBEGXG8DPTOlKhTI-SSXgzjflSR9OAw	2	2025-07-31 21:08:16.435654+00	2025-07-24 21:08:16.435857+00
6	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InNhdXJhYnp6MzE2NkBnbWFpbC5jb20iLCJzdWIiOiJzYXVyYWJ6ejMxNjZAZ21haWwuY29tIiwiZXhwIjoxNzUzOTk2MTAxLCJpYXQiOjE3NTMzOTEzMDEsImp0aSI6IjBkNTIzMmE2LTRjMGYtNDRlMC1hMDUwLTUxY2M0YWJlMGZjYyJ9.IyKZxnZonVXweGrbFhCRzLWB9xlueeyj5qxbl-YTCpw	2	2025-07-31 21:08:21.852412+00	2025-07-24 21:08:21.852618+00
8	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InNhdXJhYnBvdWRlbDAwMkBnbWFpbC5jb20iLCJzdWIiOiJzYXVyYWJwb3VkZWwwMDJAZ21haWwuY29tIiwiZXhwIjoxNzU0MDM3MDAwLCJpYXQiOjE3NTM0MzIyMDAsImp0aSI6IjUxOGIwOGEwLTNhN2QtNDg0OC04ZDljLWFhY2Q2ZGNlNzg2YiJ9.uEATuXgfIPNxizZ2oZ7EXUjmI06WB-M4ksghIZjkRd4	1	2025-08-01 08:30:00.959996+00	2025-07-25 08:30:00.960264+00
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.users (id, email, hashed_password, provider, provider_id, created_at, updated_at) FROM stdin;
1	saurabpoudel002@gmail.com		google	112145290310712777109	2025-07-24 20:02:47.243421+00	2025-07-24 20:02:47.243421+00
2	saurabzz3166@gmail.com		google	118006921727762154030	2025-07-24 20:28:20.433247+00	2025-07-24 20:28:20.433247+00
3	poudelsaurab04@gmail.com	$2a$10$vPRfV4kLBaT5HFcygLdG7OMVcsj5hw.4mM357mTSYEGJhZz3GwA/K	local		2025-07-26 18:52:04.687083+00	2025-07-26 18:52:04.687083+00
\.


--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.refresh_tokens_id_seq', 18, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- Name: email_verifications email_verifications_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.email_verifications
    ADD CONSTRAINT email_verifications_pkey PRIMARY KEY (email);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_key UNIQUE (token);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens_expires_at_idx; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX refresh_tokens_expires_at_idx ON public.refresh_tokens USING btree (expires_at);


--
-- Name: refresh_tokens_user_id_idx; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX refresh_tokens_user_id_idx ON public.refresh_tokens USING btree (user_id);


--
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

