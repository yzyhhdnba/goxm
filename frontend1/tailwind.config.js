/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        ink: '#18191C',        // Bilibili main text
        paper: '#FFFFFF',      // White surfaces
        surface: '#F1F2F3',    // Light grey backgrounds
        accent: '#fb7299',     // Bilibili Pink
        hoverPink: '#ff85a2',  // Pink hover
        sea: '#00aeec',        // Bilibili Blue
        hoverBlue: '#00b5e5',  // Blue hover
        mist: '#E3E5E7',       // Bilibili Borders / dividers
        muted: '#9499A0',      // Muted text
      },
      boxShadow: {
        float: '0 2px 4px rgba(0,0,0,0.08)',
        bili: '0 0 0 1px #e3e5e7',
      },
    },
  },
  plugins: [],
};
