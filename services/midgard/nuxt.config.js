export default {
    // Global page headers: https://go.nuxtjs.dev/config-head
    head: {
        title: 'Artchitect',
        htmlAttrs: {
            lang: 'en'
        },
        meta: [
            {charset: 'utf-8'},
            {name: 'viewport', content: 'width=device-width, initial-scale=1'},
            {hid: 'description', name: 'description', content: ''},
            {name: 'format-detection', content: 'telephone=no'}
        ],
        link: [{rel: 'icon', type: 'image/x-icon', href: '/icon64.png'}]
    },

    // Global CSS: https://go.nuxtjs.dev/config-css
    css: ['~/assets/style.scss', '@fortawesome/fontawesome-svg-core/styles.css'],

    // Plugins to run before rendering page: https://go.nuxtjs.dev/config-plugins
    plugins: ['~/plugins/axios.js', '~/plugins/fas.js'],

    // Auto import components: https://go.nuxtjs.dev/config-components
    components: true,

    // Modules for dev and build (recommended): https://go.nuxtjs.dev/config-modules
    buildModules: ['@nuxtjs/dotenv'],

    // Modules: https://go.nuxtjs.dev/config-modules
    modules: [// https://go.nuxtjs.dev/buefy
        'nuxt-buefy', // https://go.nuxtjs.dev/axios
        '@nuxtjs/axios', // https://i18n.nuxtjs.org/setup
        '@nuxtjs/i18n'],

    // Axios module configuration: https://go.nuxtjs.dev/config-axios
    axios: {
        baseURL: process.env.SERVER_API_URL,
        browserBaseURL: process.env.CLIENT_API_URL
    },

    i18n: {
        /* module options */
        locales: ["en", "ru"],
        defaultLocale: process.env.DEFAULT_LOCALE,
        strategy: 'prefix',
        vueI18nLoader: true,
        detectBrowserLanguage: {
            useCookie: true, cookieKey: 'i18n_redirected', redirectOn: 'root'  // recommended
        },
        vueI18n: {
            fallbackLocale: 'en'
        }
    }
}
