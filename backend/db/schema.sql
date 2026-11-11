-- =========================================================================
-- Agribisnis Tekno (SFA Agritec) - Skema database production (PostgreSQL)
-- Dipakai oleh backend Go. Aman dijalankan berulang (idempotent).
-- Jalankan: go run ./backend/cmd/server -migrate  (atau tempel di SQL editor)
-- =========================================================================

create table if not exists users (
  id bigserial primary key,
  email varchar(255) not null unique,
  password_hash varchar(255) not null,
  full_name varchar(255),
  phone varchar(50),
  role varchar(20) not null default 'farmer'
    check (role in ('sales', 'supervisor', 'farmer', 'distributor')),
  active boolean not null default true,
  created_at timestamptz not null default now()
);
-- kompatibel dengan skema lama (migrations/01_init.sql)
alter table users add column if not exists full_name varchar(255);
alter table users add column if not exists phone varchar(50);
alter table users add column if not exists active boolean not null default true;

create table if not exists farmers (
  id bigserial primary key,
  name varchar(255) not null,
  phone varchar(50),
  village varchar(255),
  crop_type varchar(100),
  land_area numeric(10,2),
  fertilizer_usage numeric(10,2),
  competitor_product varchar(255),
  latitude double precision,
  longitude double precision,
  created_by bigint references users(id) on delete set null,
  created_at timestamptz not null default now()
);
create index if not exists farmers_created_by_idx on farmers(created_by);

create table if not exists visits (
  id bigserial primary key,
  sales_id bigint not null references users(id) on delete cascade,
  farmer_id bigint not null references farmers(id) on delete cascade,
  has_photo boolean not null default false,
  notes text,
  latitude double precision,
  longitude double precision,
  visit_time timestamptz not null default now(),
  client_ref varchar(64) unique,
  created_at timestamptz not null default now()
);
create index if not exists visits_sales_id_idx on visits(sales_id);
create index if not exists visits_visit_time_idx on visits(visit_time desc);

-- Foto kunjungan (sudah dikompres di perangkat, ±100-300 KB)
create table if not exists visit_photos (
  visit_id bigint primary key references visits(id) on delete cascade,
  content_type varchar(50) not null default 'image/jpeg',
  data bytea not null,
  created_at timestamptz not null default now()
);

create table if not exists distributors (
  id bigserial primary key,
  name varchar(255) not null,
  owner_id bigint unique references users(id) on delete set null,
  address text,
  phone varchar(50),
  latitude double precision,
  longitude double precision,
  credit_limit_days int not null default 0, -- 0 = COD, 30 = tempo 30 hari
  created_at timestamptz not null default now()
);

create table if not exists orders (
  id bigserial primary key,
  distributor_id bigint not null references distributors(id) on delete cascade,
  created_by bigint references users(id) on delete set null,
  product varchar(255) not null,
  quantity int not null check (quantity > 0),
  total_amount numeric(14,2) not null default 0,
  order_date date not null default current_date,
  due_date date,
  payment_status varchar(20) not null default 'unpaid' check (payment_status in ('paid', 'unpaid')),
  created_at timestamptz not null default now()
);
create index if not exists orders_distributor_id_idx on orders(distributor_id);

create table if not exists commodity_prices (
  id bigserial primary key,
  name varchar(255) not null unique,
  price numeric(12,2) not null,
  unit varchar(20) not null default 'kg',
  trend varchar(10) not null default 'flat' check (trend in ('up', 'down', 'flat')),
  updated_at timestamptz not null default now()
);

create table if not exists education_contents (
  id bigserial primary key,
  title varchar(255) not null,
  description text,
  video_url text,
  thumbnail_url text,
  duration_minutes int,
  published boolean not null default true,
  created_at timestamptz not null default now()
);

-- Data awal konten (bukan data user)
insert into commodity_prices (name, price, unit, trend) values
  ('Gabah Kering Panen (GKP)', 6400, 'kg', 'up'),
  ('Bawang Merah', 24500, 'kg', 'down'),
  ('Jagung Pipil Kering', 5200, 'kg', 'flat')
on conflict (name) do nothing;

insert into education_contents (title, description, duration_minutes)
select * from (values
  ('Cara Aplikasi Pupuk Organik yang Efisien',
   'Panduan menekan biaya pupuk anorganik sekaligus memperkaya unsur hara lahan.', 3),
  ('Mengenali Hama Wereng dan Pengendaliannya',
   'Langkah pengamatan dini dan pengendalian hama wereng batang cokelat.', 5)
) as v(title, description, duration_minutes)
where not exists (select 1 from education_contents);

-- Jika database ini adalah Supabase: kunci tabel dari API publik (PostgREST),
-- karena semua akses lewat backend Go.
do $$
declare t text;
begin
  if exists (select 1 from pg_roles where rolname = 'anon') then
    foreach t in array array['users','farmers','visits','visit_photos','distributors','orders','commodity_prices','education_contents'] loop
      execute format('alter table %I enable row level security', t);
    end loop;
  end if;
end $$;
