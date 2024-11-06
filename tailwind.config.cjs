/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./internal/templ/**/*.templ"],
    darkMode: 'selector',
    theme: {
        extend: {},
    },
    plugins: [
        require("@tailwindcss/typography"),
    ],
}
