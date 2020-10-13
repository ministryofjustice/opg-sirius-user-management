module.exports = {
    ci: {
      collect: {
        url: ['http://app:8888/my-details'],
        settings: {
          chromeFlags: "--disable-gpu --no-sandbox",
        },
      },
      upload: {
        target: 'temporary-public-storage',
      },
    },
  };
