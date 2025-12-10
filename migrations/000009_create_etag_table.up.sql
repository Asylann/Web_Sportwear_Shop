CREATE TABLE IF NOT EXISTS public.etag_versions
(
    id serial NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    version integer NOT NULL,
    CONSTRAINT etag_versions_pkey PRIMARY KEY (id)
);