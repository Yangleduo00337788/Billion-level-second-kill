/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module 'vfonts/FiraCode.css' {
  const content: string
  export default content
}

declare module 'highlight.js/lib/core' {
  import hljs from 'highlight.js'
  export default hljs
}
