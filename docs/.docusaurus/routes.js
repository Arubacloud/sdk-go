import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
  {
    path: '/sdk-go/',
    component: ComponentCreator('/sdk-go/', '3fd'),
    routes: [
      {
        path: '/sdk-go/next',
        component: ComponentCreator('/sdk-go/next', '121'),
        routes: [
          {
            path: '/sdk-go/next',
            component: ComponentCreator('/sdk-go/next', 'ee4'),
            routes: [
              {
                path: '/sdk-go/next/',
                component: ComponentCreator('/sdk-go/next/', 'a05'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/sdk-go/next/filters',
                component: ComponentCreator('/sdk-go/next/filters', '7c1'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/sdk-go/next/options',
                component: ComponentCreator('/sdk-go/next/options', 'e46'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/sdk-go/next/resources',
                component: ComponentCreator('/sdk-go/next/resources', 'd3d'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/sdk-go/next/response-handling',
                component: ComponentCreator('/sdk-go/next/response-handling', '1e5'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/sdk-go/next/types',
                component: ComponentCreator('/sdk-go/next/types', '614'),
                exact: true,
                sidebar: "tutorialSidebar"
              }
            ]
          }
        ]
      },
      {
        path: '/sdk-go/',
        component: ComponentCreator('/sdk-go/', '5ee'),
        routes: [
          {
            path: '/sdk-go/',
            component: ComponentCreator('/sdk-go/', '486'),
            routes: [
              {
                path: '/sdk-go/filters',
                component: ComponentCreator('/sdk-go/filters', 'f8a'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/sdk-go/intro',
                component: ComponentCreator('/sdk-go/intro', 'e28'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/sdk-go/options',
                component: ComponentCreator('/sdk-go/options', '29b'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/sdk-go/resources',
                component: ComponentCreator('/sdk-go/resources', 'e80'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/sdk-go/response-handling',
                component: ComponentCreator('/sdk-go/response-handling', '1a1'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/sdk-go/types',
                component: ComponentCreator('/sdk-go/types', 'f61'),
                exact: true,
                sidebar: "tutorialSidebar"
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '*',
    component: ComponentCreator('*'),
  },
];
