const toggleButton = document.getElementById('dark-mode-toggle');
const darkModeIcon = document.getElementById('dark-mode-icon');
const lightModeIcon = document.getElementById('light-mode-icon');
const darkModeText = document.getElementById('dark-mode-text');

function updateButtonText() {
    if (document.documentElement.classList.contains('dark')) {
        darkModeText.textContent = 'Light Mode';
        darkModeIcon.classList.remove('hidden');
        lightModeIcon.classList.add('hidden');
    } else {
        darkModeText.textContent = 'Dark Mode';
        darkModeIcon.classList.add('hidden');
        lightModeIcon.classList.remove('hidden');
    }
}

function applyDarkModePreference() {
    const darkModeEnabled = localStorage.getItem('darkMode') === 'true';
    if (darkModeEnabled) {
        document.documentElement.classList.add('dark');
    } else {
        document.documentElement.classList.remove('dark');
    }
    updateButtonText();
}

toggleButton.addEventListener('click', () => {
    document.documentElement.classList.toggle('dark');
    const darkModeEnabled = document.documentElement.classList.contains('dark');
    localStorage.setItem('darkMode', darkModeEnabled);
    updateButtonText();
});

applyDarkModePreference();