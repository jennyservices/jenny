/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

/* List of projects/orgs using your project for the users page */
const users = [
  {
    caption: 'Typeform',
    image:
      'https://d31kvrskfu54aq.cloudfront.net/dist/img/common/typeform_new_logo.png?v=1889',
    infoLink: 'https://www.typeform.com',
    pinned: true
  }
];

const siteConfig = {
  title: 'jenny' /* title for your website */,
  tagline: 'the generator',
  url: 'https://typeform.github.io' /* your website url */,
  baseUrl: '/jenny/' /* base url for your project */,
  organizationName: 'Typeform',
  projectName: 'jenny',
  headerLinks: [
    { doc: 'readme', label: 'Docs' },
    { doc: 'options', label: 'API' },
    { doc: 'tutorials', label: 'Tutorials' },
    { page: 'help', label: 'Help' },
    { blog: true, label: 'Blog' }
  ],
  users,
  /* path to images for header/footer */
  headerIcon: 'img/jenny.svg',
  footerIcon: 'img/jenny.svg',
  favicon: 'img/jenny.png',
  /* colors for website */
  colors: {
    primaryColor: '#262627',
    secondaryColor: '#f1ece3'
  },
  // This copyright info is used in /core/Footer.js and blog rss/atom feeds.
  copyright: 'Copyright © ' + new Date().getFullYear() + 'Typeform SL',
  // organizationName: 'deltice', // or set an env variable ORGANIZATION_NAME
  // projectName: 'test-site', // or set an env variable PROJECT_NAME
  highlight: {
    // Highlight.js theme to use for syntax highlighting in code blocks
    theme: 'default'
  },
  scripts: ['https://buttons.github.io/buttons.js'],
  // You may provide arbitrary config keys to be used as needed by your template.
  repoUrl: 'https://github.com/typeform/jenny'
};

module.exports = siteConfig;
