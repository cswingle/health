CREATE TABLE public.health (
    id serial PRIMARY KEY
    username text,
    ts timestamp with time zone DEFAULT now(),
    variable text,
    value numeric,
    UNIQUE (ts, variable)
);

CREATE TABLE public.ref_variables (
    id serial PRIMARY KEY,
    variable text,
    description text,
    units text,
    sequence integer,
    superseded date,
    superseded_note text,
    UNIQUE (variable)
);
ALTER TABLE ONLY public.health
    ADD CONSTRAINT health_username_fkey FOREIGN KEY (username) REFERENCES public.users(username);
ALTER TABLE ONLY public.health
    ADD CONSTRAINT health_variable_fkey FOREIGN KEY (variable) REFERENCES public.ref_variables(variable) ON UPDATE CASCADE;
