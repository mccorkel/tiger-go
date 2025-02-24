import sveltePreprocess from 'svelte-preprocess'
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte'

const config = {
  // Consult https://github.com/sveltejs/svelte-preprocess
  // for more information about preprocessors
  preprocess: [
    vitePreprocess(),
    sveltePreprocess({})
  ],
  vitePlugin: {
    compatibility: {
      componentApi: '4'
    }
  }
}

export default config
