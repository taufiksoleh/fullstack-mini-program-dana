Component({
  properties: {
    placeholder: { type: String, value: 'Search property...' },
  },
  methods: {
    handleChange(e) {
      this.props.onSearch && this.props.onSearch(e.detail.value);
    },
    handleConfirm(e) {
      this.props.onSearch && this.props.onSearch(e.detail.value);
    },
  },
});
