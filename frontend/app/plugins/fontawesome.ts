import { library, config } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faEdit, faTrash } from '@fortawesome/free-solid-svg-icons'

// Prevent Font Awesome from adding CSS since we did it manually
config.autoAddCss = false

// Add specific icons to the library
library.add(faEdit, faTrash)

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.component('font-awesome-icon', FontAwesomeIcon)
})
