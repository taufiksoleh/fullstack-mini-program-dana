export function formatRupiah(value) {
  const num = parseInt(String(value).replace(/\D/g, ''), 10);
  if (isNaN(num)) return '0';
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.');
}
