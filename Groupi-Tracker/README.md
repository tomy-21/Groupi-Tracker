🏎️ Groupie Tracker
Un site web permettant d'explorer les pilotes de F1, circuits, résultats de courses et de gérer ses pilotes favoris.


📖 Description
Groupie Tracker est une application web en Go qui permet d'explorer les pilotes de Formule 1, les circuits, les résultats des courses et d'ajouter des pilotes en favoris.

Les données sont récupérées dynamiquement via l'API Ergast et affichées sur une interface web utilisant Go, HTML, CSS et JavaScript.


🚀 Fonctionnalités
✔️ Affichage de la liste des pilotes de F1 par année
✔️ Affichage des circuits avec leurs informations
✔️ Affichage des résultats des courses par saison et manche
✔️ Ajout et gestion des pilotes favoris
✔️ Filtrage des pilotes par nationalité et numéro
✔️ Interface fluide et rapide grâce à Go et des templates HTML


⚙️ Installation
1️⃣ Prérequis
Go 1.18+ installé
Accès à Internet pour récupérer les données via l'API Ergast
2️⃣ Cloner le projet
sh
Copier
Modifier
git clone https://github.com/tonrepo/groupie-tracker.git
cd groupie-tracker
3️⃣ Lancer le serveur
sh
Copier
Modifier
go run main.go
Le serveur démarre sur : http://localhost:8080


🖥 Utilisation
Accéder à la page d'accueil → http://localhost:8080
Naviguer entre les pages :
/drivers → Liste des pilotes
/circuits → Liste des circuits
/results → Résultats des courses
/favorites → Pilotes favoris
Ajouter un pilote en favori en cliquant sur le ❤️
Voir les pilotes favoris dans la page Favoris



📂 Structure du projet
csharp
Copier
Modifier
groupie-tracker/
│── assets/                     # Fichiers statiques (CSS, JavaScript)
│   ├── styles.css
│── templates/                   # Templates HTML
│   ├── base.html
│   ├── drivers.html
│   ├── circuits.html
│   ├── results.html
│   ├── favorites.html
│── main.go                      # Serveur principal en Go
│── README.md                    # Documentation du projet


📄 Détails des pages
1️⃣ /drivers (Liste des pilotes)
Affichage des pilotes récupérés via l'API
Filtrage par nationalité et numéro
Bouton ❤️ Ajouter aux favoris
2️⃣ /circuits (Liste des circuits)
Affichage des circuits avec leurs détails
Recherche par année
3️⃣ /results (Résultats des courses)
Affichage des résultats des courses
Détails de chaque Grand Prix
4️⃣ /favorites (Pilotes favoris)
Liste des pilotes ajoutés en favoris
Possibilité de les retirer


🔮 Améliorations futures
✅ Persistance des favoris (localStorage ou base de données)
✅ Ajout des statistiques pilotes (pôles, victoires)
✅ Ajout d'un mode sombre
✅ Animations et transitions CSS


✍️ Auteur
👤 Ton Nom
📧 Contact : tomy.ouadhi@ynov.com
🔗 GitHub : https://github.com/tomy-21/Groupi-Tracker/blob/main/Groupi-Tracker/README.md