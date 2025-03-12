
    document.addEventListener("DOMContentLoaded", function () {
        const favoriteButtons = document.querySelectorAll(".favorite-btn");

        favoriteButtons.forEach(button => {
            const driverId = button.dataset.permanentNumber; // Récupère l'ID du pilote
            
            // Vérifie si le favori est déjà activé dans localStorage
            if (localStorage.getItem(`favorite-${driverId}`) === "true") {
                button.innerText = "❤️";
            }

            button.addEventListener("click", function () {
                if (this.innerText === "♡") {
                    this.innerText = "❤️"; 
                    localStorage.setItem(`favorite-${driverId}`, "true"); // Sauvegarde en favori
                } else {
                    this.innerText = "♡"; 
                    localStorage.removeItem(`favorite-${driverId}`); // Supprime des favoris
                }
            });
        });
    });



