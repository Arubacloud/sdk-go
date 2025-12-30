import {Redirect} from '@docusaurus/router';

export default function Home() {
  // Redirect root to intro page
  // Docusaurus will handle version routing automatically
  return <Redirect to="/sdk-go/intro" />;
}

