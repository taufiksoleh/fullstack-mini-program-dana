import { fetchDetail } from '../repositories/propertyRepository';
import { formatRupiah } from '../utils/format';

export function getDetailUseCase(documentId) {
  return fetchDetail(documentId).then((res) => {
    const property = res.data.data;
    return {
      property,
      formattedPrice: formatRupiah(property.price),
    };
  });
}
