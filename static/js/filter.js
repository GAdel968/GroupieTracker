document.addEventListener('DOMContentLoaded', function() {
    const filterForm = document.getElementById('filterForm');
    const resetButton = document.getElementById('resetFilters');
    
    // Fonction pour charger dynamiquement les lieux de concert
    async function loadLocations() {
        try {
            const locationsFilter = document.querySelector('.locations-filter');
            if (!locationsFilter) return;
            
            // On récupère d'abord tous les lieux de l'API
            const response = await fetch('/api/locations');
            const data = await response.json();
            
            // On extrait tous les lieux uniques
            const allLocations = new Set();
            
            // Parcourir les données de localisation
            for (const location in data.locations) {
                // Extraire pays/ville
                const parts = location.split(', ');
                if (parts.length > 1) {
                    const country = parts[parts.length - 1];
                    allLocations.add(country);
                }
            }
            
            // On ajoute les cases à cocher pour chaque lieu
            const sortedLocations = Array.from(allLocations).sort();
            
            // On conserve les 5 premières cases à cocher et on supprime le reste
            const checkboxes = locationsFilter.querySelectorAll('label');
            for (let i = 5; i < checkboxes.length; i++) {
                locationsFilter.removeChild(checkboxes[i]);
            }
            
            // On ajoute les nouveaux lieux
            sortedLocations.forEach(location => {
                // On vérifie si la case existe déjà
                const existingCheckbox = locationsFilter.querySelector(`input[value="${location.toLowerCase()}"]`);
                if (!existingCheckbox) {
                    const label = document.createElement('label');
                    
                    const checkbox = document.createElement('input');
                    checkbox.type = 'checkbox';
                    checkbox.name = 'locations';
                    checkbox.value = location.toLowerCase();
                    
                    label.appendChild(checkbox);
                    label.appendChild(document.createTextNode(location));
                    
                    locationsFilter.appendChild(label);
                }
            });
        } catch (error) {
            console.error('Erreur lors du chargement des lieux:', error);
        }
    }
    
    // Charger les lieux au chargement de la page
    loadLocations();
    
    // Réinitialiser les filtres
    resetButton.addEventListener('click', function(e) {
        e.preventDefault();
        
        // Réinitialiser les champs de saisie
        const inputs = filterForm.querySelectorAll('input[type="number"], input[type="date"]');
        inputs.forEach(input => {
            input.value = '';
        });
        
        // Décocher toutes les cases
        const checkboxes = filterForm.querySelectorAll('input[type="checkbox"]');
        checkboxes.forEach(checkbox => {
            checkbox.checked = false;
        });
        
        // Soumettre le formulaire pour actualiser les résultats
        filterForm.submit();
    });
    
    // Ajouter une fonction d'aide pour les modèles Go
    if (typeof window.contains !== 'function') {
        window.contains = function(array, value) {
            return array && array.includes(value);
        };
    }
});