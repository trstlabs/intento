// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

// const googleTrackingId = 'G-EB7MEE3TJ1';
const algoliaAppKey = "MLPM5572P7";
const algoliaAPIKey = "153018e543e060268ab4b74c1d7983dd";
const algoliaIndexName = "intento";

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "INTENTO",
  tagline: "",
  favicon: "img/favicon.ico",

  // Set the production url of your site here
  url: "https://docs.intento.zone",
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: "/",

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: "TRST Labs", // Usually your GitHub org/user name.
  projectName: "intento", // Usually your repo name.

  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "throw",
  trailingSlash: false,

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  scripts: [
    {
      src: "https://kit.fontawesome.com/401fb1e734.js",
      crossorigin: "anonymous",
    },
  ],

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          routeBasePath: "/",
          sidebarPath: require.resolve("./sidebars.js"),
          // lastVersion: lastVersion,
          versions: {
            current: {
              path: "/",
              banner: "none",
            },
          },
        },
        blog: false,

        // gtag: {
        //   trackingID: googleTrackingId,
        //   anonymizeIP: true,
        // },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: "img/web.png",
      docs: {
        sidebar: {
          autoCollapseCategories: true,
          hideable: true,
        },
      },
      navbar: {
        logo: {
          alt: "Intento Logo",
          src: "img/intento_text.png",
        },
        items: [
          {
            type: "doc",
            docId: "index",
            position: "left",
            label: "Home",
          },
          {
            to: "/getting-started",
            position: "left",
            label: "Getting Started",
            className: "navbar__link--getting-started",
          },
          {
            to: "/tutorials",
            position: "left",
            label: "Tutorials",
            className: "navbar__link--tutorials",
          },
          {
            to: "/concepts/intent",
            position: "left",
            label: "Concepts",
            className: "navbar__link--concepts",
          },
          {
            to: "/guides/portal/overview",
            position: "left",
            label: "Guides",
            className: "navbar__link--guides",
          },
          {
            to: "/reference/intent-engine",
            position: "left",
            label: "Reference",
            className: "navbar__link--reference",
          },
          // Right-aligned social links
          {
            type: "html",
            position: "right",
            value:
              '<a href="https://github.com/trstlabs/intento" target="_blank" rel="noopener noreferrer" class="navbar__item navbar__link"><i class="fa-fw fa-brands fa-github"></i></a>',
          },
          {
            type: "html",
            position: "right",
            value:
              '<a href="https://x.com/IntentoZone" target="_blank" rel="noopener noreferrer" class="navbar__item navbar__link"><i class="fa-fw fa-brands fa-x-twitter"></i></a>',
          },
          {
            type: "html",
            position: "right",
            value:
              '<a href="https://discord.gg/hsVf9sYyZW" target="_blank" rel="noopener noreferrer" class="navbar__item navbar__link"><i class="fa-fw fa-brands fa-discord"></i></a>',
          },
        ],
      },
      // announcementBar: {
      //   id: "support_us",
      //   content: "Mainnet live on Wednesday August 27!",
      //   backgroundColor: "#fafbfc",
      //   textColor: "#091E42",
      //   isCloseable: false,
      // },
      footer: {
        style: "dark",
        links: [
          {
            items: [
              {
                html: `<a href="https://intento.zone"><img src="/img/intento_text.png" alt="Intento Logo"></a>`,
              },
            ],
          },
          {
            title: "Ecosystem",
            items: [
              {
                label: "Intento Portal",
                href: "https://portal.intento.zone/",
              },
              {
                label: "tokenstream.fun",
                href: "https://tokenstream.fun/",
              },
              {
                label: "TRST Labs",
                href: "https://trstlabs.xyz/",
              },
            ],
          },
          {
            title: "Community",
            items: [
              {
                label: "Blog",
                href: "https://blog.intento.zone/",
              },
              {
                label: "X",
                href: "https://x.com/intentozone",
              },
              {
                label: "Discord",
                href: "https://discord.gg/hsVf9sYyZW",
              },
              /*               {
                label: "Forum",
                href: "https://forum.cosmos.network/",
              },
              {
                label: "Reddit",
                href: "https://reddit.com/r/cosmosnetwork",
              }, */
            ],
          },
        ],
        copyright: `This website is maintained by TRST Labs.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
        additionalLanguages: ["protobuf", "go-module"], // https://prismjs.com/#supported-languages
      },
      algolia: {
        appId: algoliaAppKey,
        apiKey: algoliaAPIKey,
        indexName: algoliaIndexName,
        contextualSearch: false,
      },
    }),
  themes: ["@you54f/theme-github-codeblock"],
  plugins: [
    async function myPlugin(context, options) {
      return {
        name: "docusaurus-tailwindcss",
        configurePostCss(postcssOptions) {
          postcssOptions.plugins.push(require("postcss-import"));
          postcssOptions.plugins.push(require("tailwindcss/nesting"));
          postcssOptions.plugins.push(require("tailwindcss"));
          postcssOptions.plugins.push(require("autoprefixer"));
          return postcssOptions;
        },
      };
    },
    [
      "@docusaurus/plugin-client-redirects",
      {
        fromExtensions: ["html"],
        toExtensions: ["html"],
        redirects: [],
      },
    ],
  ],
};

module.exports = config;
