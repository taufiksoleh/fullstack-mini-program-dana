import { getDetailUseCase } from '../../useCases/getDetailUseCase';

Page({
  data: {
    property: null,
    formattedPrice: '',
    loading: true,
  },
  onLoad(query) {
    const { documentId } = query;
    getDetailUseCase(documentId)
      .then(({ property, formattedPrice }) => {
        this.setData({ property, formattedPrice, loading: false });
      })
      .catch(() => {
        this.setData({ loading: false });
        my.showToast({ content: 'Failed to load property', type: 'fail' });
      });
  },
  handleBook() {
    my.showToast({ content: 'Booking coming soon!', type: 'none' });
  },
});
