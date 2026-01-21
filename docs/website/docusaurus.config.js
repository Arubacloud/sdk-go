// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer').themes.github;
const darkCodeTheme = require('prism-react-renderer').themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Aruba Cloud SDK for Go',
  tagline: 'Official Go SDK for the Aruba Cloud API',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://arubacloud.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployments, it is often '/<projectName>/'
  baseUrl: '/sdk-go/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'Arubacloud',
  projectName: 'sdk-go',

  onBrokenLinks: 'warn',
  markdown: {
    hooks: {
      onBrokenMarkdownLinks: 'ignore',
    },
  },

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to set "zh-Hans" here.
  i18n: {
    defaultLocale: 'en',
    locales: ['en', 'it'],
    localeConfigs: {
      en: {
        label: 'English',
        direction: 'ltr',
        htmlLang: 'en-US',
        calendar: 'gregory',
      },
      it: {
        label: 'Italiano',
        direction: 'ltr',
        htmlLang: 'it-IT',
        calendar: 'gregory',
      },
    },
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: 'https://github.com/Arubacloud/sdk-go/tree/main/docs/website/',
          routeBasePath: '/',
          // Enable versioning (disabled during PR checks via DISABLE_VERSIONING env var)
          // Docusaurus automatically uses first entry from versions.json as latest
          versions: process.env.DISABLE_VERSIONING === 'true' ? {} : {
            current: {
              label: 'Next',
              path: 'next',
            },
          },
          onlyIncludeVersions: process.env.DISABLE_VERSIONING === 'true' ? ['current'] : undefined,
          // Show last update time
          showLastUpdateTime: true,
          showLastUpdateAuthor: true,
        },
        blog: false,
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  plugins: [
    [
      require.resolve('@easyops-cn/docusaurus-search-local'),
      {
        // Options for the search plugin
        hashed: true,
        language: ['en', 'it'],
        highlightSearchTermsOnTargetPage: true,
        explicitSearchResultPath: true,
        // Fix for GitHub Pages deployment
        indexBlog: false,
        indexPages: false,
        docsRouteBasePath: '/',
        // Ensure search works with baseUrl
        removeDefaultStopWordFilter: false,
        removeDefaultStemmer: false,
      },
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/docusaurus-social-card.jpg',
      navbar: {
        title: 'Aruba Cloud SDK for Go',
        logo: {
          alt: 'Aruba Cloud Logo',
          src: 'img/logo-cloud.png',
        },
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'tutorialSidebar',
            position: 'left',
            label: 'Docs',
          },
          // Only show version dropdown when versioning is enabled
          ...(process.env.DISABLE_VERSIONING !== 'true' ? [{
            type: 'docsVersionDropdown',
            position: 'right',
          }] : []),
          {
            type: 'localeDropdown',
            position: 'right',
          },
          {
            href: 'https://api.arubacloud.com',
            position: 'right',
            label: 'API',
          },
          {
            href: 'https://github.com/Arubacloud/sdk-go',
            position: 'right',
            className: 'header-github-link',
            'aria-label': 'GitHub repository',
          },
        ],
        hideOnScroll: true,
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Community',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/Arubacloud/sdk-go',
              },
              {
                label: 'Issues',
                href: 'https://github.com/Arubacloud/sdk-go/issues',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'Aruba Cloud',
                href: 'https://www.arubacloud.com',
              },
              {
                label: 'Changelog',
                href: 'https://github.com/Arubacloud/sdk-go/releases',
              },
            ],
          },
        ],
        copyright: `Copyright © 2025 Aruba S.p.A. - via San Clemente, 53 - 24036 Ponte San Pietro (BG) P.IVA 01573850516 - C.F. 04552920482 - C.S. € 4.000.000,00 i.v. - Numero REA: BG – 434483 - All rights reserved`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
        additionalLanguages: ['go'],
      },
    }),
};

module.exports = config;

