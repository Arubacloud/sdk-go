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
  tutorialSidebar: [
    {
      type: 'doc',
      id: 'intro',
      label: 'Quick Start',
    },
    {
      type: 'doc',
      id: 'resources',
      label: 'Resources',
    },
    {
      type: 'doc',
      id: 'types',
      label: 'Types',
    },
    {
      type: 'doc',
      id: 'options',
      label: 'Options',
    },
    {
      type: 'doc',
      id: 'response-handling',
      label: 'Response Handling',
    },
    {
      type: 'doc',
      id: 'filters',
      label: 'Filters',
    },
  ],
};

module.exports = sidebars;

