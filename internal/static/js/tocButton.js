document.addEventListener('DOMContentLoaded', function () {
    const tocToggleButton = document.getElementById('toc-toggle');
    const tocContent = document.getElementById('toc-content');

    // Function to update the TOC state based on the stored preference
    function applyTocPreference() {
        const tocVisible = localStorage.getItem('tocVisible') === 'true';
        if (tocVisible) {
            tocContent.style.display = 'block';
        } else {
            tocContent.style.display = 'none';
        }
    }

    // Event listener for the TOC toggle button
    tocToggleButton.addEventListener('click', function () {
        const isTocVisible = tocContent.style.display !== 'none';
        tocContent.style.display = isTocVisible ? 'none' : 'block';
        localStorage.setItem('tocVisible', !isTocVisible);
    });

    // Apply the TOC preference on initial load
    applyTocPreference();
});