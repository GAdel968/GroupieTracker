document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('searchInput');
    const searchResults = document.getElementById('searchResults');
    
    // Fonction pour effectuer une recherche
    async function performSearch(query) {
        if (query.length < 2) {
            searchResults.style.display = 'none';
            return;
        }
        
        try {
            const response = await fetch(`/search?q=${encodeURIComponent(query)}`);
            const data = await response.json();
            
            searchResults.innerHTML = '';
            
            if (data.length === 0) {
                searchResults.innerHTML = '<div class="search-result-item">Aucun résultat trouvé</div>';
            } else {
                data.forEach(result => {
                    const resultItem = document.createElement('div');
                    resultItem.className = 'search-result-item';
                    
                    const valueSpan = document.createElement('span');
                    valueSpan.textContent = result.value;
                    
                    const typeSpan = document.createElement('span');
                    typeSpan.className = 'search-result-type';
                    typeSpan.textContent = result.type;
                    
                    resultItem.appendChild(valueSpan);
                    resultItem.appendChild(typeSpan);
                    
                    // Gestion du clic sur un résultat
                    resultItem.addEventListener('click', function() {
                        if (result.type === 'artiste/groupe') {
                            // Rediriger vers la page de l'artiste
                            window.location.href = `/artist/${result.id}`;
                        } else {
                            // Appliquer le filtre correspondant
                            let filterUrl = '/filter?';
                            
                            if (result.type === 'membre') {
                                filterUrl += `members=${encodeURIComponent(result.value)}`;
                            } else if (result.type === 'date de création') {
                                filterUrl += `creationDateMin=${result.value}&creationDateMax=${result.value}`;
                            } else if (result.type === 'date du premier album') {
                                // Convertir la date au format approprié
                                filterUrl += `firstAlbumMin=${result.value}&firstAlbumMax=${result.value}`;
                            } else if (result.type === 'lieu de concert') {
                                filterUrl += `locations=${encodeURIComponent(result.value)}`;
                            }
                            
                            window.location.href = filterUrl;
                        }
                    });
                    
                    searchResults.appendChild(resultItem);
                });
            }
            
            searchResults.style.display = 'block';
        } catch (error) {
            console.error('Erreur lors de la recherche:', error);
        }
    }
    
    // Événement sur l'input de recherche
    let debounceTimer;
    searchInput.addEventListener('input', function() {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            performSearch(this.value);
        }, 300);
    });
    
    // Cacher les résultats quand on clique ailleurs
    document.addEventListener('click', function(event) {
        if (!searchInput.contains(event.target) && !searchResults.contains(event.target)) {
            searchResults.style.display = 'none';
        }
    });
    
    // Gérer le focus sur l'input
    searchInput.addEventListener('focus', function() {
        if (this.value.length >= 2) {
            performSearch(this.value);
        }
    });
});