import { API_BASE_URL } from '../config/index';
import { request } from '../utils/request';

export function fetchListings({ page = 1, pageSize = 10, keyword = '' } = {}) {
  const params = `page=${page}&pageSize=${pageSize}${keyword ? `&keyword=${encodeURIComponent(keyword)}` : ''}`;
  return request(`${API_BASE_URL}/api/v1/properties?${params}`);
}

export function fetchDetail(documentId) {
  return request(`${API_BASE_URL}/api/v1/properties/doc/${documentId}`);
}
