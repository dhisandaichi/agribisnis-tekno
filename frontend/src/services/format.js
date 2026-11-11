export const rupiah = (n) =>
  new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(Number(n) || 0)

export const tanggal = (d) =>
  d ? new Date(d).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }) : '-'

export const waktu = (d) =>
  d ? new Date(d).toLocaleString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) : '-'

export const ROLE_LABELS = {
  sales: 'Sales Lapangan',
  supervisor: 'Supervisor',
  farmer: 'Petani',
  distributor: 'Distributor / Kios',
}

export const homeForRole = (role) =>
  ({ sales: '/sales', supervisor: '/supervisor', farmer: '/farmer', distributor: '/distributor' }[role] || '/login')
