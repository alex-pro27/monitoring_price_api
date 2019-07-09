--
-- PostgreSQL database dump
--

-- Dumped from database version 10.9 (Ubuntu 10.9-0ubuntu0.18.04.1)
-- Dumped by pg_dump version 10.9 (Ubuntu 10.9-0ubuntu0.18.04.1)

-- Started on 2019-07-09 16:01:41 +10

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 198 (class 1259 OID 32812)
-- Name: permissions; Type: TABLE; Schema: public; Owner: alex
--

CREATE TABLE public.permissions (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    view_id integer,
    access integer DEFAULT 3,
    view text
);


ALTER TABLE public.permissions OWNER TO alex;

--
-- TOC entry 199 (class 1259 OID 32816)
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: alex
--

CREATE SEQUENCE public.permissions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.permissions_id_seq OWNER TO alex;

--
-- TOC entry 3053 (class 0 OID 0)
-- Dependencies: 199
-- Name: permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: alex
--

ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;


--
-- TOC entry 200 (class 1259 OID 32818)
-- Name: roles; Type: TABLE; Schema: public; Owner: alex
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name character varying(255) NOT NULL
);


ALTER TABLE public.roles OWNER TO alex;

--
-- TOC entry 201 (class 1259 OID 32821)
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: alex
--

CREATE SEQUENCE public.roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.roles_id_seq OWNER TO alex;

--
-- TOC entry 3054 (class 0 OID 0)
-- Dependencies: 201
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: alex
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- TOC entry 202 (class 1259 OID 32823)
-- Name: roles_permissions; Type: TABLE; Schema: public; Owner: alex
--

CREATE TABLE public.roles_permissions (
    role_id integer NOT NULL,
    permission_id integer NOT NULL
);


ALTER TABLE public.roles_permissions OWNER TO alex;

--
-- TOC entry 203 (class 1259 OID 32826)
-- Name: views; Type: TABLE; Schema: public; Owner: alex
--

CREATE TABLE public.views (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name character varying(255),
    route_path character varying(255),
    content_type_id integer,
    icon text,
    view_type integer DEFAULT 0,
    parent_id integer,
    menu boolean DEFAULT false,
    position_menu integer DEFAULT 0
);


ALTER TABLE public.views OWNER TO alex;

--
-- TOC entry 204 (class 1259 OID 32834)
-- Name: views_id_seq; Type: SEQUENCE; Schema: public; Owner: alex
--

CREATE SEQUENCE public.views_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.views_id_seq OWNER TO alex;

--
-- TOC entry 3055 (class 0 OID 0)
-- Dependencies: 204
-- Name: views_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: alex
--

ALTER SEQUENCE public.views_id_seq OWNED BY public.views.id;


--
-- TOC entry 2898 (class 2604 OID 32837)
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);


--
-- TOC entry 2899 (class 2604 OID 32838)
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- TOC entry 2902 (class 2604 OID 32839)
-- Name: views id; Type: DEFAULT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.views ALTER COLUMN id SET DEFAULT nextval('public.views_id_seq'::regclass);


--
-- TOC entry 3041 (class 0 OID 32812)
-- Dependencies: 198
-- Data for Name: permissions; Type: TABLE DATA; Schema: public; Owner: alex
--

COPY public.permissions (id, created_at, updated_at, deleted_at, view_id, access, view) FROM stdin;
3	2019-04-26 21:21:37.835893+10	2019-04-30 21:46:49.697412+10	\N	7	3	\N
4	2019-04-27 03:13:34.72043+10	2019-04-30 23:18:55.370647+10	\N	8	5	\N
5	2019-04-30 19:40:38.817306+10	2019-04-30 23:18:55.373023+10	\N	9	3	\N
7	2019-05-01 01:40:17.531609+10	2019-05-01 01:58:11.470373+10	\N	11	3	\N
6	2019-04-30 19:40:56.917551+10	2019-05-01 02:03:26.564902+10	\N	10	5	\N
9	2019-05-01 02:27:00.09764+10	2019-05-01 02:27:00.09764+10	\N	15	3	\N
8	2019-05-01 01:58:22.720974+10	2019-05-03 15:41:50.030431+10	\N	12	5	\N
10	2019-05-03 15:47:20.095162+10	2019-05-03 15:47:20.095162+10	\N	16	3	\N
11	2019-05-06 17:54:52.061354+10	2019-05-06 17:54:52.061354+10	\N	17	3	\N
12	2019-05-06 17:55:19.773954+10	2019-05-06 17:55:19.773954+10	\N	18	3	\N
13	2019-05-07 00:39:25.154238+10	2019-05-07 00:39:25.154238+10	\N	19	5	\N
14	2019-05-07 00:49:17.912056+10	2019-05-07 00:49:17.912056+10	\N	20	3	\N
15	2019-05-07 00:49:36.328017+10	2019-05-07 00:49:36.328017+10	\N	23	5	\N
16	2019-05-07 00:49:51.835138+10	2019-05-07 00:49:51.835138+10	\N	21	3	\N
17	2019-05-07 00:50:09.112336+10	2019-05-07 00:50:09.112336+10	\N	22	5	\N
18	2019-05-07 01:05:37.311032+10	2019-05-07 01:05:37.311032+10	\N	26	3	\N
19	2019-05-07 01:05:48.024528+10	2019-05-07 01:05:48.024528+10	\N	27	5	\N
20	2019-05-07 01:18:24.80196+10	2019-05-07 01:18:24.80196+10	\N	24	3	\N
21	2019-05-07 01:18:52.514247+10	2019-05-07 01:18:52.514247+10	\N	25	5	\N
22	2019-05-07 01:23:34.84936+10	2019-05-07 01:23:34.84936+10	\N	28	3	\N
23	2019-05-07 20:32:39.744815+10	2019-05-07 20:32:39.744815+10	\N	29	3	\N
24	2019-05-07 20:33:03.811798+10	2019-05-07 20:33:03.811798+10	\N	30	5	\N
25	2019-05-07 20:36:29.116358+10	2019-05-07 20:36:29.116358+10	\N	31	3	\N
\.


--
-- TOC entry 3043 (class 0 OID 32818)
-- Dependencies: 200
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: alex
--

COPY public.roles (id, created_at, updated_at, deleted_at, name) FROM stdin;
1	2019-04-17 21:33:56.305149+10	2019-05-07 20:56:48.249119+10	\N	Менеджер
2	2019-05-01 02:30:58.023948+10	2019-05-11 22:20:52.633821+10	\N	Админ
\.


--
-- TOC entry 3045 (class 0 OID 32823)
-- Dependencies: 202
-- Data for Name: roles_permissions; Type: TABLE DATA; Schema: public; Owner: alex
--

COPY public.roles_permissions (role_id, permission_id) FROM stdin;
1	5
1	6
1	4
1	7
1	8
2	3
2	4
2	5
2	6
2	7
2	8
2	10
1	10
2	11
2	12
1	11
1	12
1	13
2	13
1	14
1	15
1	16
1	17
2	14
2	15
2	16
2	17
1	18
1	19
2	18
2	19
1	20
1	21
2	20
2	21
1	22
2	22
2	23
2	24
2	25
\.


--
-- TOC entry 3046 (class 0 OID 32826)
-- Dependencies: 203
-- Data for Name: views; Type: TABLE DATA; Schema: public; Owner: alex
--

COPY public.views (id, created_at, updated_at, deleted_at, name, route_path, content_type_id, icon, view_type, parent_id, menu, position_menu) FROM stdin;
15	2019-05-01 02:25:40.145023+10	2019-05-07 00:19:34.688863+10	\N	Dashboard	/	\N	home	0	\N	t	0
28	2019-05-07 01:07:50.848408+10	2019-05-07 01:08:36.679191+10	\N	Фотографии	/photos	13	insert_photo	1	\N	f	0
31	2019-05-07 20:35:14.300985+10	2019-05-07 20:35:14.300985+10	\N	Разрешения	/permissions	2		1	\N	f	0
20	2019-05-07 00:44:22.462702+10	2019-05-21 12:48:04.413201+10	\N	Сегменты	/segments	4	list_alt	1	\N	t	7
23	2019-05-07 00:48:56.84317+10	2019-05-21 12:48:04.416561+10	\N	Редактировать сегмент	/:id	4		2	20	f	0
21	2019-05-07 00:45:53.481056+10	2019-05-21 12:48:14.609647+10	\N	Товары	/wares	5	shopping_basket	1	\N	t	8
22	2019-05-07 00:46:43.791149+10	2019-05-21 12:48:14.612619+10	\N	Редактировать товар	/:id	5		2	21	f	0
7	2019-04-26 21:15:25.654387+10	2019-05-21 12:32:26.003313+10	\N	Пользователи	/users	3	people	1	\N	t	1
8	2019-04-27 03:11:30.216069+10	2019-05-21 12:32:26.008558+10	\N	Новый пользователь	/:id	3		2	7	f	0
17	2019-05-06 17:51:51.900432+10	2019-05-21 12:48:24.299049+10	\N	Промонитореные товары	/complete-wares	11	done_all	1	\N	t	9
18	2019-05-06 17:54:04.175176+10	2019-05-21 12:48:24.302098+10	\N	Промонитореный товар	/:id	11	 	2	17	f	0
29	2019-05-07 20:31:25.645941+10	2019-05-21 12:48:31.475166+10	\N	Роли для администрирования	/roles	1	security	1	\N	t	10
30	2019-05-07 20:32:08.96122+10	2019-05-21 12:48:31.47856+10	\N	Редактировать роль	/:id	1		2	29	f	0
9	2019-04-30 19:38:12.679222+10	2019-05-21 12:41:04.159279+10	\N	Типы мониторинга	/monitoring-types	12	details	1	\N	t	2
10	2019-04-30 19:39:31.922393+10	2019-05-21 12:41:04.164552+10	\N	Редактировать тип мониторинга	/:id	12		2	9	f	0
24	2019-05-07 01:02:07.311235+10	2019-05-21 12:42:03.355796+10	\N	Группы мониторинга	/monitoring-groups	16	group_work	1	\N	t	5
25	2019-05-07 01:02:47.294613+10	2019-05-21 12:42:03.359641+10	\N	Редактировать группу мониторинга	/:id	16		2	24	f	0
26	2019-05-07 01:04:19.978113+10	2019-05-21 12:47:25.231872+10	\N	Магазины для мониторинга	/monitoring-shops	8	shopping_cart	1	\N	t	3
27	2019-05-07 01:05:18.209139+10	2019-05-21 12:47:25.235394+10	\N	Редактировать магазин для мониторинга	/:id	8		2	26	f	0
16	2019-05-03 15:47:03.854006+10	2019-05-21 12:47:46.086246+10	\N	Рабочая группа	/work-groups	6	work	1	\N	t	4
19	2019-05-07 00:38:46.817181+10	2019-05-21 12:47:46.089581+10	\N	Редактировать рабочую группу	/:id	6		2	16	f	0
11	2019-05-01 01:39:54.54909+10	2019-05-21 12:47:54.465223+10	\N	Периоды	/periods	10	date_range	1	\N	t	6
12	2019-05-01 01:59:48.56362+10	2019-05-21 12:47:54.46881+10	\N	Редактировать период	/:id	10		2	11	f	0
\.


--
-- TOC entry 3056 (class 0 OID 0)
-- Dependencies: 199
-- Name: permissions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: alex
--

SELECT pg_catalog.setval('public.permissions_id_seq', 25, true);


--
-- TOC entry 3057 (class 0 OID 0)
-- Dependencies: 201
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: alex
--

SELECT pg_catalog.setval('public.roles_id_seq', 2, true);


--
-- TOC entry 3058 (class 0 OID 0)
-- Dependencies: 204
-- Name: views_id_seq; Type: SEQUENCE SET; Schema: public; Owner: alex
--

SELECT pg_catalog.setval('public.views_id_seq', 31, true);


--
-- TOC entry 2906 (class 2606 OID 32843)
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- TOC entry 2911 (class 2606 OID 32845)
-- Name: roles_permissions roles_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.roles_permissions
    ADD CONSTRAINT roles_permissions_pkey PRIMARY KEY (role_id, permission_id);


--
-- TOC entry 2909 (class 2606 OID 32847)
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- TOC entry 2914 (class 2606 OID 32849)
-- Name: views views_pkey; Type: CONSTRAINT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.views
    ADD CONSTRAINT views_pkey PRIMARY KEY (id);


--
-- TOC entry 2904 (class 1259 OID 32851)
-- Name: idx_permissions_deleted_at; Type: INDEX; Schema: public; Owner: alex
--

CREATE INDEX idx_permissions_deleted_at ON public.permissions USING btree (deleted_at);


--
-- TOC entry 2907 (class 1259 OID 32852)
-- Name: idx_roles_deleted_at; Type: INDEX; Schema: public; Owner: alex
--

CREATE INDEX idx_roles_deleted_at ON public.roles USING btree (deleted_at);


--
-- TOC entry 2912 (class 1259 OID 32853)
-- Name: idx_views_deleted_at; Type: INDEX; Schema: public; Owner: alex
--

CREATE INDEX idx_views_deleted_at ON public.views USING btree (deleted_at);


--
-- TOC entry 2915 (class 2606 OID 32854)
-- Name: permissions permissions_view_id_views_id_foreign; Type: FK CONSTRAINT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_view_id_views_id_foreign FOREIGN KEY (view_id) REFERENCES public.views(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 2916 (class 2606 OID 32859)
-- Name: roles_permissions roles_permissions_permission_id_permissions_id_foreign; Type: FK CONSTRAINT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.roles_permissions
    ADD CONSTRAINT roles_permissions_permission_id_permissions_id_foreign FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 2917 (class 2606 OID 32864)
-- Name: roles_permissions roles_permissions_role_id_roles_id_foreign; Type: FK CONSTRAINT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.roles_permissions
    ADD CONSTRAINT roles_permissions_role_id_roles_id_foreign FOREIGN KEY (role_id) REFERENCES public.roles(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 2918 (class 2606 OID 32869)
-- Name: views views_content_type_id_content_types_id_foreign; Type: FK CONSTRAINT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.views
    ADD CONSTRAINT views_content_type_id_content_types_id_foreign FOREIGN KEY (content_type_id) REFERENCES public.content_types(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- TOC entry 2919 (class 2606 OID 32874)
-- Name: views views_parent_id_views_id_foreign; Type: FK CONSTRAINT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.views
    ADD CONSTRAINT views_parent_id_views_id_foreign FOREIGN KEY (parent_id) REFERENCES public.views(id) ON UPDATE CASCADE ON DELETE RESTRICT;


-- Completed on 2019-07-09 16:01:41 +10

--
-- PostgreSQL database dump complete
--

