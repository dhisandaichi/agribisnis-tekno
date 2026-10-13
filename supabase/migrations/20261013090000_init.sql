-- =========================================================================
-- SFA Agritec - skema PoC untuk Supabase (Auth + Postgres + Storage)
-- Menggantikan migrations/01_init.sql (skema lama backend Go).
-- =========================================================================

create extension if not exists postgis with schema extensions;

-- Schema privat: helper yang TIDAK diekspos lewat Data API
create schema if not exists private;
revoke all on schema private from public, anon;
grant usage on schema private to authenticated;

-- -------------------------------------------------------------------------
-- Profiles (1:1 dengan auth.users) + role
-- -------------------------------------------------------------------------
create table public.profiles (
    id         uuid primary key references auth.users (id) on delete cascade,
    email      text not null,
    full_name  text,
    role       text not null default 'farmer'
               check (role in ('sales', 'supervisor', 'farmer', 'distributor')),
    created_at timestamptz not null default now()
);

create or replace function private.app_role()
returns text
language sql
stable
security definer
set search_path = ''
as $$
    select role from public.profiles where id = (select auth.uid())
$$;
revoke all on function private.app_role() from public, anon;
grant execute on function private.app_role() to authenticated;

-- Profil dibuat otomatis saat signup. Role dari metadata signup dibatasi
-- ke 'sales' / 'farmer'; supervisor & distributor hanya diset oleh admin.
create or replace function private.handle_new_user()
returns trigger
language plpgsql
security definer
set search_path = ''
as $$
declare
    requested text := new.raw_user_meta_data ->> 'role';
begin
    insert into public.profiles (id, email, full_name, role)
    values (
        new.id,
        coalesce(new.email, ''),
        new.raw_user_meta_data ->> 'full_name',
        case when requested in ('sales', 'farmer') then requested else 'farmer' end
    );
    return new;
end;
$$;
revoke all on function private.handle_new_user() from public, anon, authenticated;

create trigger on_auth_user_created
    after insert on auth.users
    for each row execute function private.handle_new_user();

-- -------------------------------------------------------------------------
-- Data domain
-- -------------------------------------------------------------------------
create table public.farmers (
    id                 bigint generated always as identity primary key,
    name               text not null check (length(trim(name)) > 0),
    crop_type          text,
    land_area          numeric(10, 2) check (land_area >= 0),
    fertilizer_usage   numeric(10, 2),
    competitor_product text,
    latitude           double precision check (latitude between -90 and 90),
    longitude          double precision check (longitude between -180 and 180),
    location           extensions.geography(Point, 4326) generated always as (
        extensions.st_setsrid(extensions.st_makepoint(longitude, latitude), 4326)::extensions.geography
    ) stored,
    created_by         uuid default auth.uid() references public.profiles (id) on delete set null,
    created_at         timestamptz not null default now()
);

create table public.visits (
    id         bigint generated always as identity primary key,
    client_id  uuid not null unique, -- kunci idempoten untuk sinkronisasi offline
    sales_id   uuid not null default auth.uid() references public.profiles (id) on delete cascade,
    farmer_id  bigint not null references public.farmers (id) on delete cascade,
    latitude   double precision not null check (latitude between -90 and 90),
    longitude  double precision not null check (longitude between -180 and 180),
    location   extensions.geography(Point, 4326) generated always as (
        extensions.st_setsrid(extensions.st_makepoint(longitude, latitude), 4326)::extensions.geography
    ) stored,
    photo_path text, -- path objek di bucket 'visit-photos'
    notes      text,
    visit_time timestamptz not null default now(),
    created_at timestamptz not null default now()
);

create table public.distributors (
    id               bigint generated always as identity primary key,
    name             text not null,
    phone            text,
    credit_term_days int not null default 0 check (credit_term_days >= 0), -- 0 = COD
    latitude         double precision check (latitude between -90 and 90),
    longitude        double precision check (longitude between -180 and 180),
    location         extensions.geography(Point, 4326) generated always as (
        extensions.st_setsrid(extensions.st_makepoint(longitude, latitude), 4326)::extensions.geography
    ) stored,
    created_at       timestamptz not null default now()
);

create table public.orders (
    id             bigint generated always as identity primary key,
    distributor_id bigint not null references public.distributors (id) on delete cascade,
    supervisor_id  uuid default auth.uid() references public.profiles (id) on delete set null,
    quantity       int not null default 0 check (quantity >= 0),
    amount         numeric(14, 2) not null default 0 check (amount >= 0),
    order_date     date not null default current_date,
    due_date       date,
    payment_status text not null default 'unpaid' check (payment_status in ('paid', 'unpaid')),
    created_at     timestamptz not null default now()
);

create index farmers_created_by_idx  on public.farmers (created_by);
create index visits_sales_id_idx     on public.visits (sales_id);
create index visits_farmer_id_idx    on public.visits (farmer_id);
create index visits_visit_time_idx   on public.visits (visit_time desc);
create index visits_location_idx     on public.visits using gist (location);
create index orders_distributor_idx  on public.orders (distributor_id);
create index orders_supervisor_idx   on public.orders (supervisor_id);

-- -------------------------------------------------------------------------
-- Row Level Security
-- -------------------------------------------------------------------------
alter table public.profiles     enable row level security;
alter table public.farmers      enable row level security;
alter table public.visits       enable row level security;
alter table public.distributors enable row level security;
alter table public.orders       enable row level security;

revoke all on public.profiles, public.farmers, public.visits, public.distributors, public.orders from anon;

-- profiles: baca milik sendiri; supervisor baca semua. Tidak ada update dari
-- client sehingga user tidak bisa menaikkan role-nya sendiri.
create policy profiles_select on public.profiles
    for select to authenticated
    using (id = (select auth.uid()) or (select private.app_role()) = 'supervisor');

-- farmers: sales menambah & membaca data yang ia input; supervisor baca semua.
create policy farmers_insert_sales on public.farmers
    for insert to authenticated
    with check (created_by = (select auth.uid()) and (select private.app_role()) = 'sales');
create policy farmers_select on public.farmers
    for select to authenticated
    using (created_by = (select auth.uid()) or (select private.app_role()) = 'supervisor');

-- visits: sales menambah & membaca kunjungan sendiri; supervisor baca semua.
create policy visits_insert_sales on public.visits
    for insert to authenticated
    with check (sales_id = (select auth.uid()) and (select private.app_role()) = 'sales');
create policy visits_select on public.visits
    for select to authenticated
    using (sales_id = (select auth.uid()) or (select private.app_role()) = 'supervisor');

-- distributors & orders: hanya supervisor.
create policy distributors_supervisor on public.distributors
    for all to authenticated
    using ((select private.app_role()) = 'supervisor')
    with check ((select private.app_role()) = 'supervisor');
create policy orders_supervisor on public.orders
    for all to authenticated
    using ((select private.app_role()) = 'supervisor')
    with check ((select private.app_role()) = 'supervisor');

-- -------------------------------------------------------------------------
-- RPC & view (security invoker: RLS tetap berlaku)
-- -------------------------------------------------------------------------

-- Simpan kunjungan + data petani secara atomik dan idempoten (client_id).
create or replace function public.submit_visit(
    p_client_id   uuid,
    p_farmer_name text,
    p_crop_type   text,
    p_land_area   numeric,
    p_latitude    double precision,
    p_longitude   double precision,
    p_photo_path  text,
    p_notes       text default null,
    p_visit_time  timestamptz default now()
)
returns bigint
language plpgsql
security invoker
set search_path = ''
as $$
declare
    v_visit_id  bigint;
    v_farmer_id bigint;
begin
    select id into v_visit_id from public.visits where client_id = p_client_id;
    if found then
        return v_visit_id;
    end if;

    insert into public.farmers (name, crop_type, land_area, latitude, longitude)
    values (p_farmer_name, p_crop_type, p_land_area, p_latitude, p_longitude)
    returning id into v_farmer_id;

    insert into public.visits (client_id, farmer_id, latitude, longitude, photo_path, notes, visit_time)
    values (p_client_id, v_farmer_id, p_latitude, p_longitude, p_photo_path, p_notes, coalesce(p_visit_time, now()))
    returning id into v_visit_id;

    return v_visit_id;
end;
$$;
revoke all on function public.submit_visit(uuid, text, text, numeric, double precision, double precision, text, text, timestamptz) from public, anon;
grant execute on function public.submit_visit(uuid, text, text, numeric, double precision, double precision, text, text, timestamptz) to authenticated;

create or replace function public.supervisor_summary()
returns json
language sql
stable
security invoker
set search_path = ''
as $$
    select json_build_object(
        'total_farmers',        (select count(*) from public.farmers),
        'visits_7d',            (select count(*) from public.visits where visit_time >= now() - interval '7 days'),
        'overdue_distributors', (select count(distinct distributor_id) from public.orders
                                  where payment_status = 'unpaid' and due_date < current_date),
        'overdue_amount',       (select coalesce(sum(amount), 0) from public.orders
                                  where payment_status = 'unpaid' and due_date < current_date)
    )
$$;
revoke all on function public.supervisor_summary() from public, anon;
grant execute on function public.supervisor_summary() to authenticated;

create view public.distributor_balances
with (security_invoker = true) as
select
    d.id,
    d.name,
    d.credit_term_days,
    coalesce(sum(o.amount) filter (where o.payment_status = 'unpaid'), 0) as unpaid_amount,
    min(o.due_date) filter (where o.payment_status = 'unpaid')            as oldest_due_date,
    greatest(current_date - min(o.due_date) filter (where o.payment_status = 'unpaid'), 0) as days_overdue
from public.distributors d
left join public.orders o on o.distributor_id = d.id
group by d.id;
revoke all on public.distributor_balances from anon;

-- -------------------------------------------------------------------------
-- Storage: foto kunjungan (privat, path = {user_id}/{client_id}.jpg)
-- -------------------------------------------------------------------------
insert into storage.buckets (id, name, public, file_size_limit, allowed_mime_types)
values ('visit-photos', 'visit-photos', false, 1048576, array['image/jpeg', 'image/png', 'image/webp'])
on conflict (id) do nothing;

create policy visit_photos_insert_sales on storage.objects
    for insert to authenticated
    with check (
        bucket_id = 'visit-photos'
        and (storage.foldername(name))[1] = (select auth.uid())::text
        and (select private.app_role()) = 'sales'
    );

create policy visit_photos_select on storage.objects
    for select to authenticated
    using (
        bucket_id = 'visit-photos'
        and (
            (storage.foldername(name))[1] = (select auth.uid())::text
            or (select private.app_role()) = 'supervisor'
        )
    );
