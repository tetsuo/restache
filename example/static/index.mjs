import Base from './base.stache'
import exampleData from './data.json'
import { createRoot } from 'react-dom/client';

window.addEventListener("load", () => {
  const root = createRoot(document.getElementById('root'));
  root.render(Base(exampleData));
})

