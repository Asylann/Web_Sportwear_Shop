CREATE TABLE IF NOT EXISTS public.transactions
(
    id serial NOT NULL,
    from_wallet_id integer,
    to_wallet_id integer,
    amount numeric(10, 2) NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT transactions_pkey PRIMARY KEY (id)
);