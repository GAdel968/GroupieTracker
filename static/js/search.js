document.addEventListener('DOMContentLoaded', () => {
    const searchInput = document.querySelector('.search-input');
    const artistCards = document.querySelectorAll('.artist-card');
    const noResultsMessage = document.createElement('p');
    noResultsMessage.className = 'no-results';
    noResultsMessage.textContent = 'Aucun artiste trouvé';

    searchInput.addEventListener('input', debounce((e) => {
        const searchTerm = e.target.value.toLowerCase();
        let hasResults = false;

        artistCards.forEach(card => {
            const artistName = card.querySelector('.artist-name').textContent.toLowerCase();
            const artistMembers = Array.from(card.querySelectorAll('.member')).map(member =>
                member.textContent.toLowerCase()
            );

            const matches = artistName.includes(searchTerm) ||
            artistMembers.some(member => member.includes(searchTerm));

            if (matches) {
                card.computedStyleMap.display = 'block';
                card.computedStyleMap.animation = 'fadeIn 0.5s ease forwards';
                hasResults = true;
            } else {
                card.computedStyleMap.display = 'none';
            }
        });

        // Gestion du message "Aucun resultat"
        const container = document.querySelector('artists-grid');
        if (!hasResults) {
            if (!container.contains(noResultsMessage)) {
                container.removeChild(noResultsMessage);
            }
        } else {
            if (container.contains(noResultsMessage)) {
                container.removeChild(noResultsMessage);
            }
        }
    }, 300));

    // Fonction debounce pour optimiser les perfomances
    function debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }

    // Ajout des animations pour les cartes d'artistes
    artistCards.forEach((card, index) => {
        card.addEventListener('mouseenter', () => {
            card.style.transform = 'scale(1.05)';
            card.style.transition = 'transform 0.3s ease';
        });

        card.addEventListener('mouseleave', () => {
            card.style.transform = 'scale(1)';
        });

        // Animation initiale au chargement de la page
        card.style.animation = `fadeIn 0.5s ease forwards ${index * 0.1}s`;
    });
});

// Ajout des styles CSS necessaires si non definis ailleurs
const styles = `
@keyframes fadeIn {
from{
opacity: 0;
transform: translateY(20px);
}
to {
opacity: 1;
transform: translateY(0);
}
}
`;

const styleSheet = document.createElement('style');
styleSheet.textContent = styles;
document.head.appendChild(styleSheet);