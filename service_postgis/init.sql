
-- SEQUENCE: public.igcfile_id_seq

-- DROP SEQUENCE public.igcfile_id_seq;

CREATE SEQUENCE public.igcfile_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

ALTER SEQUENCE public.igcfile_id_seq
    OWNER TO crawler;

-- Table: public.igcfile

-- DROP TABLE public.igcfile;

CREATE TABLE public.igcfile
(
    id integer NOT NULL DEFAULT nextval('igcfile_id_seq'::regclass),
    hash character(64) COLLATE pg_catalog."default" NOT NULL,
    d_date bigint,
    date bigint,
    filename character varying COLLATE pg_catalog."default",
    d_url character varying COLLATE pg_catalog."default",
    CONSTRAINT igcfile_pkey PRIMARY KEY (id, hash)
)

TABLESPACE pg_default;

ALTER TABLE public.igcfile
    OWNER to crawler;



