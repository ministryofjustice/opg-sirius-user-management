module.exports = {
    ci: {
      collect: {
        url: ['http://app:8888/my-details'],
        settings: {
          extraHeaders: JSON.stringify({Cookie: 'XSRF-TOKEN=abcde; Other=other'}),
          chromeFlags: "--disable-gpu --no-sandbox",
        },
      },
      upload: {
        target: 'temporary-public-storage',
      },
    },
  };
