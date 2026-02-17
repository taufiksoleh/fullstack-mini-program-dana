import { getListingsUseCase } from '../../useCases/getListingsUseCase';

Component({
  data: {
    listings: [],
    searchKeyword: '',
    layout: 'grid',
    page: 1,
    totalPages: 1,
    loading: false,
    hasMore: false,
    scrollHeight: 600,
    toolbarHidden: false,
    ptrVisible: false,
    ptrText: '',
  },
  didMount() {
    my.getSystemInfo({
      success: (res) => {
        const toolbarHeight = Math.round(res.windowWidth / 750 * 112);
        this._windowHeight = res.windowHeight;
        this._toolbarHeight = toolbarHeight;
        this.setData({ scrollHeight: res.windowHeight - 50 - toolbarHeight });
      },
    });
    this.fetchListings({ page: 1, keyword: '' });
  },
  detached() {
    clearTimeout(this._searchTimer);
  },
  methods: {
    fetchListings({ page, keyword }) {
      this.setData({ loading: true });
      getListingsUseCase({ page, keyword })
        .then(({ listings: newItems, page: currentPage, totalPages, hasMore }) => {
          const listings = currentPage === 1
            ? newItems
            : [...this.data.listings, ...newItems];
          this.setData({
            listings,
            page: currentPage,
            totalPages,
            hasMore,
            loading: false,
            ptrVisible: false,
            ptrText: '',
          });
        })
        .catch(() => {
          this.setData({ loading: false, ptrVisible: false, ptrText: '' });
          my.showToast({ content: 'Failed to load properties', type: 'fail' });
        });
    },
    doRefresh() {
      this._showToolbar();
      this.setData({ ptrVisible: true, ptrText: 'Refreshing...' });
      this.fetchListings({ page: 1, keyword: this.data.searchKeyword });
    },
    onTouchStart(e) {
      this._touchStartY = e.touches[0].pageY;
    },
    onTouchMove(e) {
      if ((this._lastScrollTop || 0) > 5 || !this._touchStartY) return;
      const dy = e.touches[0].pageY - this._touchStartY;
      if (dy <= 0) return;
      const text = dy > 60 ? 'Release to refresh' : 'Pull to refresh';
      if (text !== this.data.ptrText || !this.data.ptrVisible) {
        this.setData({ ptrVisible: true, ptrText: text });
      }
    },
    onTouchEnd(e) {
      if ((this._lastScrollTop || 0) > 5) {
        this._touchStartY = 0;
        return;
      }
      const dy = e.changedTouches[0].pageY - (this._touchStartY || 0);
      this._touchStartY = 0;
      if (dy > 60) {
        this.doRefresh();
      } else {
        this.setData({ ptrVisible: false, ptrText: '' });
      }
    },
    onSearchChange(value) {
      this.setData({ searchKeyword: value });
      this._showToolbar();
      clearTimeout(this._searchTimer);
      this._searchTimer = setTimeout(() => {
        this.fetchListings({ page: 1, keyword: value });
      }, 1000);
    },
    onScroll(e) {
      const scrollTop = e.detail.scrollTop;
      const last = this._lastScrollTop || 0;
      const delta = scrollTop - last;
      this._lastScrollTop = scrollTop;

      if (delta > 10 && !this.data.toolbarHidden) {
        this.setData({
          toolbarHidden: true,
          scrollHeight: this._windowHeight - 50,
        });
      } else if (delta < -10 && this.data.toolbarHidden) {
        this._showToolbar();
      }
    },
    _showToolbar() {
      this.setData({
        toolbarHidden: false,
        scrollHeight: (this._windowHeight || 600) - 50 - (this._toolbarHeight || 60),
      });
    },
    onScrollToLower() {
      if (!this.data.hasMore || this.data.loading) return;
      this.fetchListings({ page: this.data.page + 1, keyword: this.data.searchKeyword });
    },
    toggleLayout() {
      this.setData({ layout: this.data.layout === 'grid' ? 'list' : 'grid' });
    },
    handleCardTap(e) {
      const documentId = e.currentTarget.dataset.id;
      my.navigateTo({ url: `/pages/propertyDetail/propertyDetail?documentId=${documentId}` });
    },
  },
});
