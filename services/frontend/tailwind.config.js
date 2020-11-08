/* eslint-disable */
const { colors } = require("tailwindcss/defaultTheme");

module.exports = {
  future: {
    removeDeprecatedGapUtilities: true,
    purgeLayersByDefault: true
  },
  purge: ["./src/**/*.html", "./src/**/*.vue"],
  theme: {
    extend: {
      colors: {
        primary: "#311b92",
        primarylight: "#6746c3",
        primarydark: "#000063",

        secondary: "#f57c00",
        secondarylight: "#ffad42",
        secondarydark: "#bb4d00",

        cancel: colors.red[500],
        accept: colors.green[500],
        info: colors.blue[500]
      }
    }
  },
  variants: {},
  plugins: []
};
