/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // Main documentation sidebar
  docsSidebar: [
    {
      type: 'doc',
      id: 'index',
      label: 'Introduction',
      className: 'sidebar-item-intro',
    },
    {
      type: 'category',
      label: 'Getting Started',
      items: [
        'getting-started/index',
        'getting-started/architecture',
        'getting-started/use-cases',
        'getting-started/differentiation',
        'getting-started/integration',
        'getting-started/into-token'
      ],
      className: 'sidebar-item-getting-started',
    },
    
    {
      type: 'category',
      label: 'Concepts',
      items: [
        'concepts/intent',
        'concepts/flow-patterns',
        'concepts/icq',
        'concepts/conditions',
        'concepts/fees',
      ],
      className: 'sidebar-item-concept',
    },
    {
      type: 'category',
      label: 'Guides',
      items: [
        {
          type: 'category',
          label: 'Intento Portal',
          items: [
            'guides/portal/overview',
            'guides/portal/submit-flows',
            'guides/portal/notifications',
          ],
        },
        {
          type: 'category',
          label: 'tokenstream.fun',
          items: ['guides/tokenstream/index'],
        },
        {
          type: 'category',
          label: 'CLI',
          items: ['guides/cli/using-the-cli'],
        },
        {
          type: 'category',
          label: 'TypeScript',
          items: ['guides/typescript/getting-started'],
        },
      ],
      className: 'sidebar-item-guide',
    },
    {
      type: 'category',
      label: 'Reference',
      items: [
        'reference/intent-engine/index',
        'reference/intent-engine/parameters',
        'reference/intent-engine/supported_types',
        'reference/intent-engine/authentication',
        'reference/api/querying',
        'reference/consensus',
      ],
      className: 'sidebar-item-reference',
    },
    {
      type: 'category',
      label: 'Tutorials',
      items: [
        'tutorials/index',
        {
          type: 'category',
          label: 'General',
          items: [
            'tutorials/cross-chain',
            'tutorials/auto-delegation',
            'tutorials/conditional-transfers',
            'tutorials/conditional-transfers-icq',
            'tutorials/trustless-agent',
          ],
        },
        {
          type: 'category',
          label: 'Case Studies',
          items: [
            'tutorials/case-study/auto-compound',
            'tutorials/case-study/rwa-integration',
          ],
        },
      ],
      className: 'sidebar-item-tutorial',
    },
    {
      type: 'category',
      label: 'Community Programs',
      items: [
        'community_programs/ambassador-program',
        'community_programs/delegation-program',
      ],
      className: 'sidebar-item-community',
    },
  ],
};

module.exports = sidebars;
