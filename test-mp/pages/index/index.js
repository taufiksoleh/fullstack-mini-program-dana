Page({
  data: {
    tabs: [
      {
        icon: 'AppOutline',
        activeIcon: 'AppOutline',
        text: 'Home',
        title: 'About Property Finder',
      },
      {
        icon: 'SetOutline',
        activeIcon: 'SetOutline',
        text: 'Property List',
        title: 'Property Listing',
      },
    ],
    current: 0,
  },
  onLoad() {
    my.setNavigationBar({ title: this.data.tabs[0].title });
  },
  handleChange(index) {
    this.setData({ current: index });
    my.setNavigationBar({ title: this.data.tabs[index].title });
  },
});
