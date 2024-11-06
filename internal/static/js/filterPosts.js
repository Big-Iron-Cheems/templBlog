// Search bar functionality for filtering posts by title, tags, and keywords
function filterPosts(inputId, collectionId, itemClass) {
    const query = document.getElementById(inputId).value.trim().toLowerCase();
    const posts = document.querySelectorAll(`#${collectionId} .${itemClass}`);

    // Parse search terms
    const exactMatches = [...query.matchAll(/"([^"]+)"/g)].map(match => match[1]); // Extract phrases in quotes
    const tagMatches = [...query.matchAll(/#([\w-]+)/g)].map(match => `#${match[1]}`); // Extract tags
    const keywordMatches = query
        .replace(/"([^"]+)"/g, "") // Remove exact phrases
        .replace(/#([\w-]+)/g, "") // Remove tags
        .split(" ")                             // Split by whitespace
        .map(term => term.trim())
        .filter(term => term);                    // Extract remaining keywords, ignoring empty terms

    posts.forEach((post) => {
        const postTitle = post.querySelector(".post-title").textContent.toLowerCase();
        const postTags = Array.from(post.querySelectorAll(".post-tags")).map(tag => tag.textContent.toLowerCase());

        // Check for exact phrase matches
        const hasExactMatch = exactMatches.every(phrase => postTitle.includes(phrase));

        // Check for tag matches
        let hasTagMatch = tagMatches.every(tagQuery => postTags.some(tag => tag.startsWith(tagQuery)));

        // Check for keyword matches
        const hasKeywordMatch = keywordMatches.every(keyword => postTitle.includes(keyword));

        // Mixed condition: require all specified exact matches, tag conditions, and keywords
        const match = hasExactMatch && hasTagMatch && hasKeywordMatch;

        // Show or hide the post based on match result
        post.style.display = match ? "block" : "none";
    });
}

// Clicking a tag button adds the tag to the search input
document.addEventListener("DOMContentLoaded", () => {
    const tagButtons = document.querySelectorAll(".post-tags");
    const searchInput = document.getElementById("search");

    tagButtons.forEach(button => {
        button.addEventListener("click", () => {
            const tagText = button.textContent.trim();
            const currentQuery = searchInput.value.trim();
            if (!currentQuery.includes(tagText)) {
                searchInput.value = currentQuery ? `${currentQuery} ${tagText}` : tagText;
                filterPosts('search', 'postGrid', 'post');
            }
        });
    });
});