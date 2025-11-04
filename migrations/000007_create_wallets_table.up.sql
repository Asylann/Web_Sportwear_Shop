CREATE TABLE IF NOT EXISTS public.wallets
(
    id serial NOT NULL,
    user_id integer NOT NULL,
    balance numeric(10, 2) NOT NULL DEFAULT 0.00,
    currency character varying(8) COLLATE pg_catalog."default" NOT NULL DEFAULT 'KZT'::character varying,
    CONSTRAINT wallets_pkey PRIMARY KEY (id),
    CONSTRAINT wallets_user_id_key UNIQUE (user_id)
);
