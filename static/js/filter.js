document.addEventListener('DOMContentLoaded', () => {
    const filterForm = document.querySelector('.filter-form');
    const artistsContainer = document.querySelector('.artists-grid');

    filterForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = new FormData(filterForm);
        const params = new URLSearchParams();

        // Ajout des parametres de filtrage 
        formData.forEach((value, key) => {
            if (value) params.append(key, value);
        });

        try {
            const response = await fetch(`/filter?${params.toString()}`);
            if (!response.ok) throw new Error('Erreur reseau');

            const data = await response.text();
            artistsContainer.innerHTML = data;

            // Animation des cartes d'artistes
            const cards = document.querySelectorAll('.artist-card');
            cards.forEach((card, index) => {
                cards.style.animation = `fadeIn 0.5s ease forwards ${index * 0.1}s`;
            });

        } catch (error) {
            console.error('Erreur:', error);
            artistsContainer.innerHTML = '<p class="error">Une erreur est survenue lors du filtrage</p>';
        }
    });

    // Reset du formulaire
    const resetButton = document.querySelector('.reset-button');
    if (resetButton) {
        resetButton.addEventListener('click', () => {
            filterForm.reset();
            filterForm.dispatchEvent(new Event('submit'));
        });
    }
});