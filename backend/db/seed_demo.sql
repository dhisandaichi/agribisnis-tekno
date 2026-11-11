-- Akun demo, satu per peran, plus data contoh. Idempotent.
-- Password: supervisor@agritec.id / Supervisor@Agritec26, sales@agritec.id / Sales@Agritec26,
--           distributor@agritec.id / Distributor@Agritec26, petani@agritec.id / Petani@Agritec26
-- Ganti password setelah login pertama, atau nonaktifkan akun demo sebelum dipakai pelanggan.
insert into users (email, password_hash, full_name, phone, role) values
 ('supervisor@agritec.id', '$2a$10$1r.mEewO139nb8Do7.1jAO5NuNcH0HAYg0LjdNm0qvxE6p4zuqXke', 'Supervisor Demo', '081200000001', 'supervisor'),
 ('sales@agritec.id', '$2a$10$bav6k4i8h5B6DYloNsjpfuXHZQobg6meW1Udo2MFeOf8z0H0hB.Oa', 'Sales Demo', '081200000002', 'sales'),
 ('distributor@agritec.id', '$2a$10$c2/g9faJ2OVMuVmypeBo8.BqR0v2aqXNovN3TqfzRnOwf0ug445zK', 'Distributor Demo', '081200000003', 'distributor'),
 ('petani@agritec.id', '$2a$10$KP3yqN9tKKAn3peArLIW6OLzjmjDfOIHJkqlH7rUt5wK1jtI2i.8S', 'Petani Demo', '081200000004', 'farmer')
on conflict (email) do nothing;

insert into distributors (name, owner_id, address, phone, credit_limit_days)
select 'Kios Tani Makmur (Demo)', id, 'Jl. Raya Desa Sukamaju No. 1', '081200000003', 30
from users where email = 'distributor@agritec.id'
on conflict (owner_id) do nothing;

insert into orders (distributor_id, created_by, product, quantity, total_amount, order_date, due_date, payment_status)
select d.id, s.id, 'Pupuk Organik Agritec 25 kg', 20, 2500000, current_date - 10, current_date + 20, 'unpaid'
from distributors d, users s
where d.name = 'Kios Tani Makmur (Demo)' and s.email = 'supervisor@agritec.id'
  and not exists (select 1 from orders o where o.distributor_id = d.id);

insert into farmers (name, phone, village, crop_type, land_area, fertilizer_usage, competitor_product, latitude, longitude, created_by)
select 'Pak Budi (Demo)', '081300000001', 'Sukamaju, Cibadak', 'Padi', 1.5, 300, 'Pupuk Urea Subsidi', -6.8939, 106.7794, id
from users where email = 'sales@agritec.id'
  and not exists (select 1 from farmers where name = 'Pak Budi (Demo)');

insert into visits (sales_id, farmer_id, notes, latitude, longitude, client_ref)
select s.id, f.id, 'Kunjungan contoh: petani tertarik mencoba pupuk organik musim depan.', -6.8939, 106.7794, 'seed-demo-1'
from users s, farmers f
where s.email = 'sales@agritec.id' and f.name = 'Pak Budi (Demo)'
on conflict (client_ref) do nothing;
