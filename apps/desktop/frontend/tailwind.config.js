/******** Tailwind v4 uses CSS-first config via @config. Keeping JS for clarity. ********/
/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: [
    './index.html',
    './src/**/*.{ts,tsx}',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}

