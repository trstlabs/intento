module.exports = {
  theme: "cosmos",
  title: "Trustless Hub",

  head: [
    [
      "script",
      {
        async: true,
        src: "https://www.googletagmanager.com/gtag/js?id=G-XL9GNV1KHW",
      },
    ],
    [
      "script",
      {},
      [
        "window.dataLayer = window.dataLayer || [];\nfunction gtag(){dataLayer.push(arguments);}\ngtag('js', new Date());\ngtag('config', 'G-XL9GNV1KHW');",
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
      id: "BH4D9OD16A",
      key: "d6908a9436133e03e9b0131bad808775",
      index: "docs-startport",
    },
 
    sidebar: {
      
      auto: true,
      nav: [
        {
          title: "Resources",
          children: [
            {
              title: "Trustless Hub on Github",
              path: "https://github.com/danieljdd/trst",
            },
            {
              title: "Cosmos SDK Docs",
              path: "https://docs.cosmos.network",
            },
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
        {
          service: "medium",
          url: "https://danieldijkstra.medium.com/",
        },
        {
          service: "twitter",
          url: "https://twitter.com/trustlesshub",
        },
      
        {
          service: "reddit",
          url: "https://reddit.com/r/trustlesshub",
        },
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
        "This website is maintained by Tendermint Inc. The contents and opinions of this website are those of Tendermint Inc.",
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
                "https://github.com/cosmos/cosmos-sdk/blob/master/docs/DOCS_README.md",
            },
            {
              title: "Source code on GitHub",
              url: "https://github.com/cosmos/cosmos-sdk/",
            },
          ],
        },
      ],
    },
  },
};
