-- PostgreSQL database dump
-- Database: kanban-master

-- Initialize settings
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

-- Create the role/user if it does not exist
DO
$$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'vsys-kanban-user') THEN
      CREATE ROLE "vsys-kanban-user" LOGIN PASSWORD 'NewPassword123';
   END IF;
END
$$;

-- Create the database if it doesn't exist
DO
$$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'kanban-master') THEN
      CREATE DATABASE "kanban-master" WITH ENCODING = 'UTF8' OWNER "vsys-kanban-user";
   END IF;
END
$$;

-- Connect to the newly created database
\connect kanban-master;

-- Create Tables and Indexes

SET default_tablespace = '';
SET default_with_oids = false;

DROP SCHEMA public;

CREATE SCHEMA public AUTHORIZATION pg_database_owner;

COMMENT ON SCHEMA public IS 'standard public schema';

-- DROP SEQUENCE public.compounds_id_seq;

CREATE SEQUENCE public.compounds_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.compounds_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.compounds_id_seq1;

CREATE SEQUENCE public.compounds_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.compounds_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.compounds_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.financialyears_srno_seq;

CREATE SEQUENCE public.financialyears_srno_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.financialyears_srno_seq OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.financialyears_srno_seq TO "vsys-kanban-user";

-- DROP SEQUENCE public.inventory_id_seq;

CREATE SEQUENCE public.inventory_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.inventory_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.inventory_id_seq1;

CREATE SEQUENCE public.inventory_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.inventory_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.inventory_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.kb_data_id_seq;

CREATE SEQUENCE public.kb_data_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.kb_data_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.kb_data_id_seq1;

CREATE SEQUENCE public.kb_data_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.kb_data_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.kb_data_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.kb_extension_id_seq;

CREATE SEQUENCE public.kb_extension_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.kb_extension_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.kb_extension_id_seq1;

CREATE SEQUENCE public.kb_extension_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.kb_extension_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.kb_extension_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.kb_root_id_seq;

CREATE SEQUENCE public.kb_root_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.kb_root_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.kb_root_id_seq1;

CREATE SEQUENCE public.kb_root_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.kb_root_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.kb_root_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.kb_transaction_id_seq;

CREATE SEQUENCE public.kb_transaction_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.kb_transaction_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.kb_transaction_id_seq1;

CREATE SEQUENCE public.kb_transaction_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.kb_transaction_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.kb_transaction_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.ldapconfig_id_seq;

CREATE SEQUENCE public.ldapconfig_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.ldapconfig_id_seq OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.ldapconfig_id_seq TO "vsys-kanban-user";

-- DROP SEQUENCE public.monthmasters_srno_seq;

CREATE SEQUENCE public.monthmasters_srno_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.monthmasters_srno_seq OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.monthmasters_srno_seq TO "vsys-kanban-user";

-- DROP SEQUENCE public.permissions_id_seq;

CREATE SEQUENCE public.permissions_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.permissions_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.permissions_id_seq1;

CREATE SEQUENCE public.permissions_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.permissions_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.permissions_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.prod_line_id_seq;

CREATE SEQUENCE public.prod_line_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.prod_line_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.prod_line_id_seq1;

CREATE SEQUENCE public.prod_line_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.prod_line_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.prod_line_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.prod_process_id_seq;

CREATE SEQUENCE public.prod_process_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.prod_process_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.prod_process_id_seq1;

CREATE SEQUENCE public.prod_process_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.prod_process_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.prod_process_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.prod_process_line_id_seq;

CREATE SEQUENCE public.prod_process_line_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.prod_process_line_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.prod_process_line_id_seq1;

CREATE SEQUENCE public.prod_process_line_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.prod_process_line_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.prod_process_line_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.recipe_id_seq;

CREATE SEQUENCE public.recipe_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.recipe_id_seq OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.recipe_id_seq TO "vsys-kanban-user";

-- DROP SEQUENCE public.sambaconfig_id_seq;

CREATE SEQUENCE public.sambaconfig_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.sambaconfig_id_seq OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.sambaconfig_id_seq TO "vsys-kanban-user";

-- DROP SEQUENCE public.stage_id_seq;

CREATE SEQUENCE public.stage_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.stage_id_seq OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.stage_id_seq TO "vsys-kanban-user";

-- DROP SEQUENCE public.systemdefaults_id_seq;

CREATE SEQUENCE public.systemdefaults_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.systemdefaults_id_seq OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.systemdefaults_id_seq TO "vsys-kanban-user";

-- DROP SEQUENCE public.systemlogs_log_id_seq;

CREATE SEQUENCE public.systemlogs_log_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.systemlogs_log_id_seq OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.systemlogs_log_id_seq TO "vsys-kanban-user";

-- DROP SEQUENCE public.user_roles_id_seq;

CREATE SEQUENCE public.user_roles_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.user_roles_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.user_roles_id_seq1;

CREATE SEQUENCE public.user_roles_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.user_roles_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.user_roles_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.user_to_vendor_id_seq;

CREATE SEQUENCE public.user_to_vendor_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.user_to_vendor_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.user_to_vendor_id_seq1;

CREATE SEQUENCE public.user_to_vendor_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.user_to_vendor_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.user_to_vendor_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.users_id_seq;

CREATE SEQUENCE public.users_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.users_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.users_id_seq1;

CREATE SEQUENCE public.users_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.users_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.users_id_seq1 TO "vsys-kanban-user";

-- DROP SEQUENCE public.usertorole_id_seq;

CREATE SEQUENCE public.usertorole_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.usertorole_id_seq OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.usertorole_id_seq TO "vsys-kanban-user";

-- DROP SEQUENCE public.vendors_id_seq;

CREATE SEQUENCE public.vendors_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.vendors_id_seq OWNER TO "vsys-kanban-user";

-- DROP SEQUENCE public.vendors_id_seq1;

CREATE SEQUENCE public.vendors_id_seq1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.vendors_id_seq1 OWNER TO "vsys-kanban-user";
GRANT ALL ON SEQUENCE public.vendors_id_seq1 TO "vsys-kanban-user";
-- public.compounds definition

-- Drop table

-- DROP TABLE public.compounds;

CREATE TABLE public.compounds (
	id serial4 NOT NULL,
	compound_name varchar(100) NOT NULL,
	description text NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	status bool DEFAULT false NULL,
	CONSTRAINT compounds_compound_name_key UNIQUE (compound_name),
	CONSTRAINT compounds_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_compound_name ON public.compounds USING btree (compound_name);

-- Permissions

ALTER TABLE public.compounds OWNER TO "vsys-kanban-user";


-- public.financialyears definition

-- Drop table

-- DROP TABLE public.financialyears;

CREATE TABLE public.financialyears (
	srno serial4 NOT NULL,
	financialyear varchar(10) NULL,
	yearcode bpchar(1) NULL,
	createdby varchar(254) NULL,
	createdon timestamptz DEFAULT now() NULL,
	modifiedby varchar(254) NULL,
	modifiedon timestamptz DEFAULT now() NULL,
	CONSTRAINT financialyears_pkey PRIMARY KEY (srno)
);
CREATE INDEX idx_financialyears_yearcode ON public.financialyears USING btree (yearcode);

-- Permissions

ALTER TABLE public.financialyears OWNER TO "vsys-kanban-user";
GRANT ALL ON TABLE public.financialyears TO "vsys-kanban-user";


-- public.ldapconfig definition

-- Drop table

-- DROP TABLE public.ldapconfig;

CREATE TABLE public.ldapconfig (
	id serial4 NOT NULL,
	ldap_url text NOT NULL,
	bind_dn text NOT NULL,
	base_dn text NOT NULL,
	"password" text NOT NULL,
	unique_identifier text DEFAULT 'uid'::text NOT NULL,
	tls_insecure bool DEFAULT false NULL,
	is_default bool DEFAULT false NOT NULL,
	createdon timestamptz DEFAULT now() NULL,
	createdby varchar(254) NULL,
	modifiedon timestamptz DEFAULT now() NULL,
	modifiedby varchar(254) NULL,
	CONSTRAINT ldapconfig_pkey PRIMARY KEY (id)
);

-- Permissions

ALTER TABLE public.ldapconfig OWNER TO "vsys-kanban-user";


-- public.monthmasters definition

-- Drop table

-- DROP TABLE public.monthmasters;

CREATE TABLE public.monthmasters (
	srno serial4 NOT NULL,
	"month" varchar(10) NULL,
	monthcode bpchar(2) NULL,
	createdby varchar(254) NULL,
	createdon timestamptz DEFAULT now() NULL,
	modifiedby varchar(254) NULL,
	modifiedon timestamptz DEFAULT now() NULL,
	CONSTRAINT monthmasters_pkey PRIMARY KEY (srno)
);
CREATE INDEX idx_monthmasters_monthcode ON public.monthmasters USING btree (monthcode);

-- Permissions

ALTER TABLE public.monthmasters OWNER TO "vsys-kanban-user";
GRANT ALL ON TABLE public.monthmasters TO "vsys-kanban-user";


-- public.permissions definition

-- Drop table

-- DROP TABLE public.permissions;

CREATE TABLE public.permissions (
	id serial4 NOT NULL,
	permission_name varchar(100) NOT NULL,
	description text NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	CONSTRAINT permissions_permission_name_key UNIQUE (permission_name),
	CONSTRAINT permissions_pkey PRIMARY KEY (id)
);

-- Permissions

ALTER TABLE public.permissions OWNER TO "vsys-kanban-user";


-- public.prod_line definition

-- Drop table

-- DROP TABLE public.prod_line;

CREATE TABLE public.prod_line (
	id serial4 NOT NULL,
	"name" varchar(100) NOT NULL,
	icon varchar(255) NULL,
	description text NULL,
	status bool NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	CONSTRAINT prod_line_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_prod_line_name ON public.prod_line USING btree (name);

-- Permissions

ALTER TABLE public.prod_line OWNER TO "vsys-kanban-user";


-- public.prod_process definition

-- Drop table

-- DROP TABLE public.prod_process;

CREATE TABLE public.prod_process (
	id serial4 NOT NULL,
	"name" varchar(100) NOT NULL,
	link varchar(255) NULL,
	icon varchar(255) NULL,
	description text NULL,
	status varchar(50) NULL,
	line_visibility bool DEFAULT false NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	expected_mean_time varchar(15) NULL,
	CONSTRAINT prod_process_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_prod_process_name ON public.prod_process USING btree (name);

-- Permissions

ALTER TABLE public.prod_process OWNER TO "vsys-kanban-user";


-- public.sambaconfig definition

-- Drop table

-- DROP TABLE public.sambaconfig;

CREATE TABLE public.sambaconfig (
	id serial4 NOT NULL,
	workgroup text NOT NULL,
	server_string text NOT NULL,
	"security" text NOT NULL,
	is_default bool DEFAULT false NOT NULL,
	createdon timestamptz DEFAULT now() NULL,
	createdby varchar(254) NULL,
	modifiedon timestamptz DEFAULT now() NULL,
	modifiedby varchar(254) NULL,
	CONSTRAINT sambaconfig_pkey PRIMARY KEY (id)
);

-- Permissions

ALTER TABLE public.sambaconfig OWNER TO "vsys-kanban-user";


-- public.stage definition

-- Drop table

-- DROP TABLE public.stage;

CREATE TABLE public.stage (
	id serial4 NOT NULL,
	"name" varchar(255) NULL,
	headers _text NULL,
	created_on timestamptz NULL,
	created_by varchar(50) NULL,
	modified_on timestamptz NULL,
	modified_by varchar(50) NULL,
	active bool NULL,
	CONSTRAINT stage_pkey PRIMARY KEY (id)
);

-- Permissions

ALTER TABLE public.stage OWNER TO "vsys-kanban-user";
GRANT ALL ON TABLE public.stage TO "vsys-kanban-user";


-- public.systemdefaults definition

-- Drop table

-- DROP TABLE public.systemdefaults;

CREATE TABLE public.systemdefaults (
	id serial4 NOT NULL,
	"name" varchar(512) NOT NULL,
	system_code varchar(50) NOT NULL,
	home_logo_url varchar(512) NULL,
	login_logo_url varchar(512) NULL,
	registration_logo_url varchar(512) NULL,
	cold_store varchar(255) NULL,
	flow_chart_heading varchar(512) NULL,
	cold_store_board varchar(50) NULL,
	cold_store_menu varchar(50) NULL,
	CONSTRAINT systemdefaults_name_key UNIQUE (name),
	CONSTRAINT systemdefaults_pkey PRIMARY KEY (id),
	CONSTRAINT systemdefaults_system_code_key UNIQUE (system_code)
);
CREATE INDEX idx_systemdefaults_system_code ON public.systemdefaults USING btree (system_code);

-- Permissions

ALTER TABLE public.systemdefaults OWNER TO "vsys-kanban-user";
GRANT ALL ON TABLE public.systemdefaults TO "vsys-kanban-user";


-- public.systemlogs definition

-- Drop table

-- DROP TABLE public.systemlogs;

CREATE TABLE public.systemlogs (
	log_id serial4 NOT NULL,
	"timestamp" timestamptz DEFAULT now() NULL,
	message text NOT NULL,
	message_type varchar(50) NOT NULL,
	is_critical bool DEFAULT false NULL,
	icon varchar(255) NULL,
	created_by varchar(254) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(254) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	CONSTRAINT systemlogs_pkey PRIMARY KEY (log_id)
);
CREATE INDEX idx_systemlogs_is_critical ON public.systemlogs USING btree (is_critical);
CREATE INDEX idx_systemlogs_message_type ON public.systemlogs USING btree (message_type);

-- Permissions

ALTER TABLE public.systemlogs OWNER TO "vsys-kanban-user";


-- public.user_roles definition

-- Drop table

-- DROP TABLE public.user_roles;

CREATE TABLE public.user_roles (
	id serial4 NOT NULL,
	role_name varchar(50) NOT NULL,
	description text NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	deny _text NULL,
	CONSTRAINT user_roles_pkey PRIMARY KEY (id),
	CONSTRAINT user_roles_role_name_key UNIQUE (role_name)
);

-- Permissions

ALTER TABLE public.user_roles OWNER TO "vsys-kanban-user";


-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id serial4 NOT NULL,
	username varchar(50) NOT NULL,
	email varchar(100) NOT NULL,
	password_hash varchar(255) NOT NULL,
	approved_by varchar(50) NULL,
	approved_on timestamptz NULL,
	rejected_by varchar(50) NULL,
	rejected_on timestamptz NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	isactive bool DEFAULT true NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_key UNIQUE (username)
);
CREATE INDEX idx_users_email ON public.users USING btree (email);
CREATE INDEX idx_users_username ON public.users USING btree (username);

-- Permissions

ALTER TABLE public.users OWNER TO "vsys-kanban-user";


-- public.vendors definition

-- Drop table

-- DROP TABLE public.vendors;

CREATE TABLE public.vendors (
	id serial4 NOT NULL,
	vendor_code varchar(255) NOT NULL,
	vendor_name varchar(100) NOT NULL,
	contact_info varchar(255) NULL,
	address text NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	isactive bool DEFAULT true NULL,
	per_day_lot_config int4 NULL,
	per_month_lot_config int4 NULL,
	CONSTRAINT vendors_pkey PRIMARY KEY (id),
	CONSTRAINT vendors_vendor_code_key UNIQUE (vendor_code),
	CONSTRAINT vendors_vendor_name_key UNIQUE (vendor_name)
);
CREATE INDEX idx_vendor_code ON public.vendors USING btree (vendor_code);
CREATE INDEX idx_vendor_name ON public.vendors USING btree (vendor_name);

-- Permissions

ALTER TABLE public.vendors OWNER TO "vsys-kanban-user";


-- public.inventory definition

-- Drop table

-- DROP TABLE public.inventory;

CREATE TABLE public.inventory (
	id serial4 NOT NULL,
	compound_id int4 NULL,
	min_quantity int4 NOT NULL,
	max_quantity int4 NOT NULL,
	product_type varchar(50) NOT NULL,
	description text NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	available_quantity int4 NULL,
	CONSTRAINT inventory_max_quantity_check CHECK ((max_quantity >= 0)),
	CONSTRAINT inventory_min_quantity_check CHECK ((min_quantity >= 0)),
	CONSTRAINT inventory_pkey PRIMARY KEY (id),
	CONSTRAINT inventory_product_type_check CHECK (((product_type)::text = ANY (ARRAY[('Rarely Required'::character varying)::text, ('Required Always'::character varying)::text, ('No Requirement'::character varying)::text]))),
	CONSTRAINT inventory_compound_id_fkey FOREIGN KEY (compound_id) REFERENCES public.compounds(id)
);
CREATE INDEX idx_inventory_compound_id ON public.inventory USING btree (compound_id);

-- Permissions

ALTER TABLE public.inventory OWNER TO "vsys-kanban-user";


-- public.kb_extension definition

-- Drop table

-- DROP TABLE public.kb_extension;

CREATE TABLE public.kb_extension (
	id serial4 NOT NULL,
	order_id int4 NULL,
	code varchar(100) NULL,
	status varchar(50) NULL,
	vendor_id int4 NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	CONSTRAINT kb_extension_pkey PRIMARY KEY (id),
	CONSTRAINT kb_extension_vendor_id_fkey FOREIGN KEY (vendor_id) REFERENCES public.vendors(id)
);
CREATE INDEX idx_kb_extension_vendor_id ON public.kb_extension USING btree (vendor_id);

-- Permissions

ALTER TABLE public.kb_extension OWNER TO "vsys-kanban-user";


-- public.prod_process_line definition

-- Drop table

-- DROP TABLE public.prod_process_line;

CREATE TABLE public.prod_process_line (
	id serial4 NOT NULL,
	prod_process_id int4 NULL,
	prod_line_id int4 NULL,
	"order" int4 NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	isgroup bool DEFAULT false NULL,
	group_name varchar(512) NULL,
	CONSTRAINT prod_process_line_pkey PRIMARY KEY (id),
	CONSTRAINT prod_process_line_prod_line_id_fkey FOREIGN KEY (prod_line_id) REFERENCES public.prod_line(id),
	CONSTRAINT prod_process_line_prod_process_id_fkey FOREIGN KEY (prod_process_id) REFERENCES public.prod_process(id)
);

-- Permissions

ALTER TABLE public.prod_process_line OWNER TO "vsys-kanban-user";


-- public.recipe definition

-- Drop table

-- DROP TABLE public.recipe;

CREATE TABLE public.recipe (
	id serial4 NOT NULL,
	compound_name varchar(100) NULL,
	compound_code varchar(100) NULL,
	stage_id int4 NULL,
	"data" jsonb NULL,
	prod_line_id int4 NULL,
	created_on timestamptz NULL,
	created_by varchar(50) NULL,
	modified_on timestamptz NULL,
	modified_by varchar(50) NULL,
	CONSTRAINT recipe_pkey PRIMARY KEY (id),
	CONSTRAINT fk_compound_name FOREIGN KEY (compound_name) REFERENCES public.compounds(compound_name),
	CONSTRAINT fk_stage_id FOREIGN KEY (stage_id) REFERENCES public.stage(id)
);

-- Permissions

ALTER TABLE public.recipe OWNER TO "vsys-kanban-user";
GRANT ALL ON TABLE public.recipe TO "vsys-kanban-user";


-- public.role_permissions definition

-- Drop table

-- DROP TABLE public.role_permissions;

CREATE TABLE public.role_permissions (
	role_id int4 NOT NULL,
	permission_id int4 NOT NULL,
	CONSTRAINT role_permissions_pkey PRIMARY KEY (role_id, permission_id),
	CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id),
	CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.user_roles(id)
);

-- Permissions

ALTER TABLE public.role_permissions OWNER TO "vsys-kanban-user";


-- public.user_to_vendor definition

-- Drop table

-- DROP TABLE public.user_to_vendor;

CREATE TABLE public.user_to_vendor (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	vendor_id int4 NOT NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_on timestamptz DEFAULT now() NULL,
	CONSTRAINT user_to_vendor_pkey PRIMARY KEY (id),
	CONSTRAINT fk_user_to_vendor_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE,
	CONSTRAINT fk_user_to_vendor_vendor FOREIGN KEY (vendor_id) REFERENCES public.vendors(id) ON DELETE CASCADE
);
CREATE INDEX idx_user_to_vendor_user_id ON public.user_to_vendor USING btree (user_id);
CREATE INDEX idx_user_to_vendor_vendor_id ON public.user_to_vendor USING btree (vendor_id);

-- Permissions

ALTER TABLE public.user_to_vendor OWNER TO "vsys-kanban-user";


-- public.usertorole definition

-- Drop table

-- DROP TABLE public.usertorole;

CREATE TABLE public.usertorole (
	id serial4 NOT NULL,
	userid int4 NULL,
	userroleid int4 NULL,
	createdon timestamptz DEFAULT now() NULL,
	createdby varchar(254) NULL,
	modifiedon timestamptz DEFAULT now() NULL,
	modifiedby varchar(254) NULL,
	CONSTRAINT usertorole_pkey PRIMARY KEY (id),
	CONSTRAINT usertorole_userid_fkey FOREIGN KEY (userid) REFERENCES public.users(id),
	CONSTRAINT usertorole_userroleid_fkey FOREIGN KEY (userroleid) REFERENCES public.user_roles(id)
);

-- Permissions

ALTER TABLE public.usertorole OWNER TO "vsys-kanban-user";
GRANT ALL ON TABLE public.usertorole TO "vsys-kanban-user";


-- public.kb_data definition

-- Drop table

-- DROP TABLE public.kb_data;

CREATE TABLE public.kb_data (
	id serial4 NOT NULL,
	compound_id int4 NULL,
	mfg_date_time timestamp NULL,
	demand_date_time timestamp NULL,
	exp_date timestamp NULL,
	cell_no varchar(50) NULL,
	"location" varchar(100) NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	kb_extension_id int4 NULL,
	no_of_lots int4 NULL,
	CONSTRAINT kb_data_pkey PRIMARY KEY (id),
	CONSTRAINT fk_kb_extension_id FOREIGN KEY (kb_extension_id) REFERENCES public.kb_extension(id),
	CONSTRAINT kb_data_compound_id_fkey FOREIGN KEY (compound_id) REFERENCES public.compounds(id)
);
CREATE INDEX idx_kb_data_compound_id ON public.kb_data USING btree (compound_id);

-- Permissions

ALTER TABLE public.kb_data OWNER TO "vsys-kanban-user";


-- public.kb_root definition

-- Drop table

-- DROP TABLE public.kb_root;

CREATE TABLE public.kb_root (
	id serial4 NOT NULL,
	running_no int4 NULL,
	initial_no int4 NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	kb_data_id int4 NULL,
	status varchar(50) NOT NULL,
	lot_no varchar(255) NULL,
	in_inventory bool DEFAULT false NULL,
	"comment" varchar(25) NULL,
	notes varchar(512) NULL,
	remark varchar(510) NULL,
	CONSTRAINT kb_root_pkey PRIMARY KEY (id),
	CONSTRAINT fk_kb_data_id FOREIGN KEY (kb_data_id) REFERENCES public.kb_data(id)
);

-- Permissions

ALTER TABLE public.kb_root OWNER TO "vsys-kanban-user";


-- public.kb_transaction definition

-- Drop table

-- DROP TABLE public.kb_transaction;

CREATE TABLE public.kb_transaction (
	id serial4 NOT NULL,
	prod_process_id int4 NOT NULL,
	status varchar(50) NULL,
	job_id int4 NULL,
	kb_root_id int4 NOT NULL,
	prod_process_line_id int4 NOT NULL,
	started_on timestamptz DEFAULT now() NULL,
	completed_on timestamptz DEFAULT now() NULL,
	created_by varchar(50) NULL,
	created_on timestamptz DEFAULT now() NULL,
	modified_by varchar(50) NULL,
	modified_on timestamptz DEFAULT now() NULL,
	CONSTRAINT kb_transaction_pkey PRIMARY KEY (id),
	CONSTRAINT unique_kbrootid_prodprocesslineid UNIQUE (kb_root_id, prod_process_line_id),
	CONSTRAINT unique_status_kbrootid UNIQUE (status, kb_root_id),
	CONSTRAINT kb_transaction_kb_root_id_fkey FOREIGN KEY (kb_root_id) REFERENCES public.kb_root(id),
	CONSTRAINT kb_transaction_prod_process_line_id_fkey FOREIGN KEY (prod_process_line_id) REFERENCES public.prod_process_line(id)
);
CREATE INDEX idx_kb_transaction_kb_root_id ON public.kb_transaction USING btree (kb_root_id);
CREATE INDEX idx_kb_transaction_prod_process_id ON public.kb_transaction USING btree (prod_process_id);
CREATE INDEX idx_kb_transaction_prod_process_line_id ON public.kb_transaction USING btree (prod_process_line_id);
CREATE INDEX idx_kb_transaction_status ON public.kb_transaction USING btree (status);

-- Permissions

ALTER TABLE public.kb_transaction OWNER TO "vsys-kanban-user";




-- Permissions;