import Index from './index.stache'
import exampleData from './example-data.json'
import { createRoot } from 'react-dom/client';

window.addEventListener("load", () => {
  const root = createRoot(document.getElementById('root'));
  root.render(Index(exampleData));
})

