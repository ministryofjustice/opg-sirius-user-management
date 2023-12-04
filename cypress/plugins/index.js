/// <reference types="cypress" />
// ***********************************************************
// This example plugins/index.js can be used to load plugins
//
// You can change the location of this file or turn off loading
// the plugins file with the 'pluginsFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/plugins-guide
// ***********************************************************

// This function is called when a project is opened or re-opened (e.g. due to
// the project's config changing)

const webpackPreprocessor = require("@cypress/webpack-preprocessor");

/**
 * @type {Cypress.PluginConfig}
 */
module.exports = (on, config) => {
  const options = {
    // see https://github.com/webpack/webpack.js.org/pull/3981/files
    webpackOptions: {
      resolve: {
        fallback: { "path": require.resolve("path-browserify") }
      },
      module: {
        rules: [
          {
            test: /\.m?js/,
            resolve: {
              fullySpecified: false
            }
          }
        ]
      }
    }
  };

  // `on` is used to hook into various events Cypress emits
  // `config` is the resolved Cypress config
  on("task", {
    log(message) {
        console.log(message)

        return null
    },
    table(message) {
        console.table(message)

        return null
    },
    failed: require("cypress-failed-log/src/failed")(),
  });

  on("file:preprocessor", webpackPreprocessor(options));

  on("before:browser:launch", (browser = {}, launchOptions) => {
    if (browser.name === "chrome") {
      launchOptions.args.push("--disable-dev-shm-usage");
    }
    return launchOptions;
  });
};
