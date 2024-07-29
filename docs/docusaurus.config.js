// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

// const googleTrackingId = 'G-EB7MEE3TJ1';
const algoliaAppKey = "UTAVANG7NO";
const algoliaAPIKey = "e828406cee19753665694712b1bf7555";
const algoliaIndexName = "intento_docs";

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
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
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
      // announcementBar: {
      //   id: "support_us",
      //   content: "ddd",
      //   backgroundColor: "#fafbfc",
      //   textColor: "#091E42",
      //   isCloseable: false,
      // },
      navbar: {
        // title: 'INTENTO',
        logo: {
          alt: "Intento Logo",
          src: "img/intento_text.png",
        },
        items: [
          // {
          //   type: 'docsVersionDropdown',
          //   position: 'left',
          //   dropdownItemsAfter: [{to: '/versions', label: 'All versions'}],
          //   dropdownActiveClassDisabled: true,
          // },
        ],
      },
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
                label: "TriggerPortal",
                href: "https://triggerportal.zone/",
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
                href: "https://discord.gg/teRhcmwu",
              },
              /*               {
                label: "Forum",
                href: "https://forum.cosmos.network/",
              },
              {
                label: "Discord",
                href: "https://discord.gg/cosmosnetwork",
              },
              {
                label: "Reddit",
                href: "https://reddit.com/r/cosmosnetwork",
              }, */
            ],
          }
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
