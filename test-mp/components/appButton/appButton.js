Component({
  properties: {
    type: { type: String, value: 'primary' },
    disabled: { type: Boolean, value: false },
    loading: { type: Boolean, value: false },
    size: { type: String, value: 'large' },
  },
  methods: {
    handleTap(e) {
      this.triggerEvent('tap', e);
    },
  },
});
