import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Easy to Use',
    Svg: require('@site/static/img/undraw_docusaurus_mountain.svg').default,
    description: (
      <>
        The SDK provides a simple, intuitive API that makes it easy to interact
        with the Aruba Cloud API. Get started in minutes with our quick start guide.
      </>
    ),
  },
  {
    title: 'Type Safe',
    Svg: require('@site/static/img/undraw_docusaurus_tree.svg').default,
    description: (
      <>
        Built with Go's strong typing system, the SDK provides compile-time
        safety and excellent IDE support for all API operations.
      </>
    ),
  },
  {
    title: 'Comprehensive',
    Svg: require('@site/static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
        Supports all Aruba Cloud services including Compute, Network, Storage,
        Database, Container, and more. Everything you need in one SDK.
      </>
    ),
  },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        {Svg && <Svg className={styles.featureSvg} role="img" />}
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}

