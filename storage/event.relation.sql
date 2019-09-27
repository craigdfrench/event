CREATE TABLE public.event
(
    "Id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "CreatedAt" timestamp with time zone NOT NULL DEFAULT now(),
    "Email" text COLLATE pg_catalog."default" NOT NULL,
    "Environment" text COLLATE pg_catalog."default" NOT NULL,
    "Component" text COLLATE pg_catalog."default" NOT NULL,
    "Message" text COLLATE pg_catalog."default",
    "Data" json, 
    CONSTRAINT event_pkey PRIMARY KEY ("Id")
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;
CREATE INDEX "Component_index"
    ON public.event USING btree
    ("Component" COLLATE pg_catalog."default" text_pattern_ops)
    TABLESPACE pg_default;
CREATE INDEX "CreatedAt"
    ON public.event USING btree
    ("CreatedAt")
    TABLESPACE pg_default;
CREATE INDEX "Email_index"
    ON public.event USING btree
    ("Email" COLLATE pg_catalog."default" text_pattern_ops)
    TABLESPACE pg_default;
CREATE INDEX "Environment_index"
    ON public.event USING btree
    ("Environment" COLLATE pg_catalog."default" text_pattern_ops)
    TABLESPACE pg_default;
    