import { fetchListings } from '../repositories/propertyRepository';

export function getListingsUseCase({ page = 1, keyword = '' } = {}) {
  return fetchListings({ page, pageSize: 10, keyword }).then((res) => {
    const listings = res.data.data;
    const { pagination } = res.data;
    return {
      listings,
      page: pagination.page,
      totalPages: pagination.totalPages,
      hasMore: pagination.page < pagination.totalPages,
    };
  });
}
