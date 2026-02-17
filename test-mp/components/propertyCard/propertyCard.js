import { formatRupiah } from '../../utils/format';

Component({
  properties: {
    imageUrl: { type: String, value: '' },
    title: { type: String, value: '' },
    price: { type: String, value: '' },
    layout: { type: String, value: 'grid' },
  },
  data: {
    formattedPrice: '',
  },
  deriveDataFromProps(nextProps) {
    this.setData({ formattedPrice: formatRupiah(nextProps.price) });
  },
});
