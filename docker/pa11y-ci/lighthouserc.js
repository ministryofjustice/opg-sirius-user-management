module.exports = {
    ci: {
      collect: {
        url: ['http://app:8888/my-details'],
        settings: {
          extraHeaders: JSON.stringify({Cookie: 'XSRF-TOKEN=abcde; Other=other'}),
          chromeFlags: "--disable-gpu --no-sandbox",
        },
      },
      assert: {
        preset: "lighthouse:no-pwa",
      },
      upload: {
        target: 'temporary-public-storage',
      },
    },
  };
