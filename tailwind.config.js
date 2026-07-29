/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          500: '#3b82f6',
          600: '#2563eb',
        },
        danger: '#ef4444',   // 可保留，也可直接用 red-500
        success: '#10b981',  // 可保留，也可直接用 emerald-500
        // 移除 purple 自定义，使用 Tailwind 默认 purple 调色板
        secondary: '#64748b',
        muted: '#94a3b8',
        glass: '#f8fafc',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      borderRadius: {
        'xl': '12px',
      },
      boxShadow: {
        'glass': '0 4px 12px rgba(0, 0, 0, 0.06)',
        'modal': '0 20px 40px rgba(0, 0, 0, 0.1)',
      },
      backdropBlur: {
        '4': '4px',
      },
    },
  },
  plugins: [],
}