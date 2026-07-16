--
-- PostgreSQL database dump
--

\restrict DC7EABV2XUDlLmo8HN3OQ2chvflK5PNfujhhRiyidWr9ba6L67yJENkJcs2vE9Y

-- Dumped from database version 17.7
-- Dumped by pg_dump version 17.7

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
-- Name: items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.items (
    id integer NOT NULL,
    note text,
    payload text
);


--
-- Name: TABLE items; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.items IS 'demo; with semicolon';


--
-- Name: v_items; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.v_items AS
 SELECT id,
    note
   FROM public.items
  WHERE (id > 0);


--
-- Data for Name: items; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.items (id, note, payload) FROM stdin;
1	semi;colon	tab\there
2	\N	x;y;z
\.


--
-- Name: items items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict DC7EABV2XUDlLmo8HN3OQ2chvflK5PNfujhhRiyidWr9ba6L67yJENkJcs2vE9Y

