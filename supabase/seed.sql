-- Data contoh untuk demo dashboard supervisor (PoC).
-- Sengaja TIDAK membuat user/password: buat akun lewat Supabase Dashboard
-- (Authentication > Users) lalu set role supervisor dengan:
--   update public.profiles set role = 'supervisor' where email = '<email>';

insert into public.distributors (name, phone, credit_term_days, latitude, longitude) values
    ('Kios Tani Makmur', '081200000001', 0,  -7.2575, 112.7521),
    ('UD Subur Jaya',    '081200000002', 30, -7.9666, 112.6326),
    ('Toko Pupuk Sejahtera', '081200000003', 30, -8.1724, 113.7008);

insert into public.orders (distributor_id, supervisor_id, quantity, amount, order_date, due_date, payment_status)
select d.id, null, o.quantity, o.amount, o.order_date, o.due_date, o.payment_status
from (values
    ('Kios Tani Makmur',     50, 7500000.00,  current_date - 10, current_date - 10, 'paid'),
    ('UD Subur Jaya',        80, 14000000.00, current_date - 35, current_date - 5,  'unpaid'),
    ('Toko Pupuk Sejahtera', 40, 6200000.00,  current_date - 12, current_date + 18, 'unpaid')
) as o(distributor_name, quantity, amount, order_date, due_date, payment_status)
join public.distributors d on d.name = o.distributor_name;
