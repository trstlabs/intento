module.exports = {
  theme: "cosmos",
  title: "Trustless Hub",

  head: [
    [
      "script",
      {
        async: true,
        src: "https://www.googletagmanager.com/gtag/js?id=G-NG42ZHDDXJ",
      },
    ],
    [
      "script",
      {},
      [
        "window.dataLayer = window.dataLayer || [];\nfunction gtag(){dataLayer.push(arguments);}\ngtag('js', new Date());\ngtag('config', 'G-NG42ZHDDXJ');",
      ],
    ],
  ],
  themeConfig: {
    logo: {
      // Image in ./vuepress/public/logo.svg
      src: "/logo.svg",
      // Image width relative to the sidebar
      width: "100%",
  
    },
    topbar: {
      banner: false,
    },
    header: {
      img: {
        // Image in ./vuepress/public/logo.svg
        src: "/logo.png",
        // Image width relative to the sidebar
        width: "75%",
      },
      title: "Documentation",
    },

    algolia: {
      id: "GB7BY4DOAY",
      key: "ed0f0d3998090c9dc5023ccffd9bb4dd",
      index: "trustlesshub",
    },
 
    sidebar: {
      
      auto: true,
      nav: [
        {
          title: "Resources",
          children: [ 
            {
              title: "CosmWasm Docs",
              path: "https://docs.cosmwasm.com/docs/",
            },
            {
              title: "Cosmos SDK Docs",
              path: "https://docs.cosmos.network",
            },
            {
              title: "Trustless Hub Github",
              path: "https://github.com/trstlabs",
            },
           /* {
              title: "DeFi Contract Bundle ",
              path: "https://github.com/trstlabs/dex-contracts",
            },*/
          ],
        },
      ],
    },
  

    custom: true,
    footer: {
      question: {
        text:
          "Chat with Trustless Hub and Cosmos SDK developers in <a href='https://discord.gg/7fwqwc3afK' target='_blank'>Discord</a>.",
      },
     
      textLink: {
        text: "Trustless Hub",
        url: "https://www.trustlesshub.com/",
      },
      services: [
       /*  {
          service: "medium",
          url: "https://trstlabs.medium.com/",
        }, */
        {
          service: "twitter",
          url: "https://twitter.com/trustlesshub",
        },
      
       /*  {
          service: "reddit",
          url: "https://reddit.com/r/trustlesshub",
        }, */
        {
          service: "discord",
          url: "https://discord.gg/vcExX9T",
        },
        {
          service: "youtube",
          url: "https://www.youtube.com/channel/UCia6iAp3VqLDeEQMoU9zt2w",
        },
      ],

      smallprint:
        "This website is maintained by TRST Labs. The contents and opinions of this website are those of TRST Labs.",
      links: [
        {
          title: "Documentation",
          children: [
            {
              title: "Cosmos SDK",
              url: "https://docs.cosmos.network",
            },
            {
              title: "Cosmos Hub",
              url: "https://hub.cosmos.network",
            },
            {
              title: "Tendermint Core",
              url: "https://docs.tendermint.com",
            },
          ],
        },
        {
          title: "Community",
          children: [
            {
              title: "Cosmos blog",
              url: "https://blog.cosmos.network",
            },
            {
              title: "Forum",
              url: "https://forum.cosmos.network",
            },
            {
              title: "Chat",
              url: "https://discord.gg/7fwqwc3afK",
            },
          ],
        },
        {
          title: "Contributing",
          children: [
            {
              title: "Contributing to the docs",
              url:
                "https://github.com/trstlabs/trst/blob/master/docs/README.md",
            },
            {
              title: "Source code on GitHub",
              url: "https://github.com/trstlabs/trst/",
            },
          ],
        },
      ],
      plugins: [
        [
          'sitemap',
          {
            hostname: 'https://docs.trustlesshub.com'
          }
        ]
      ]
    },
  },
};
