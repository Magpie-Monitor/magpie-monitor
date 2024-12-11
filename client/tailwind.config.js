import { nextui } from '@nextui-org/theme';

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx,mdx}',
    './node_modules/@nextui-org/theme/dist/components/(date-picker|button|ripple|spinner|calendar|date-input|form|popover).js',
  ],
  darkMode: 'class',
  plugins: [
    nextui({
      defaultTheme: 'magpie',
      layout: {
        radius: {
          small: '1rem',
          medium: '1rem',
          large: '1rem',
        },
      },

      themes: {
        magpie: {
          extend: 'dark',
          layout: {
            radius: {
              small: '1rem',
              medium: '1rem',
              large: '1rem',
            },
          },
          colors: {
            default: {
              50: '#122131',
            },
            primary: {
              50: '#5cd06080',
              DEFAULT: '#5cd060',
              foreground: '#122131',
            },
            background: '#07111b',

            content1: {
              DEFAULT: '#07111b',
            },
          },
        },
      },
    }),
  ],
};
