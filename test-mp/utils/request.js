export function request(url, { method = 'GET' } = {}) {
  return new Promise((resolve, reject) => {
    my.request({
      url,
      method,
      success: (res) => resolve(res),
      fail: (err) => reject(err),
    });
  });
}
