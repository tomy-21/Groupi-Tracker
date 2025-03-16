# 🏎️ Groupie Tracker – Formula One API  

## 🚀 Présentation  
Groupie Tracker est une application web développée en **Golang**, permettant d'afficher des données sur la **Formule 1** grâce à l'**API Ergast Formula One**. L'application permet d'effectuer des recherches sur les pilotes, circuits et résultats des courses.  

### 📌 Fonctionnalités principales :  
✅ Recherche et filtres avancés sur les pilotes (4 filtres disponibles)  
✅ Affichage des circuits et résultats de courses par année  
✅ Système de favoris persistant avec stockage en JSON  
✅ Gestion des erreurs en cas d’indisponibilité de l’API  
✅ Interface simple et user-friendly  

---

## 📡 API utilisée  
L’application exploite l’**API Ergast Formula One**, qui fournit des données sur la Formule 1 depuis 1950.  

🔗 **Documentation API :** [Formula One API – Postman](https://www.postman.com/maintenance-astronomer-29796265/f1-api/collection/60v93se/formula-one-api)  

### 🔥 Endpoints exploités  
- **Pilotes** : [`/drivers`](http://ergast.com/api/f1/{{year}}/drivers) → Liste des pilotes par saison  
- **Circuits** : [`/circuits`](http://ergast.com/api/f1/{{year}}/circuits) → Liste des circuits par saison  
- **Résultats de courses** : [`/results`](http://ergast.com/api/f1/{{year}}/{{round}}/results) → Résultats d’un Grand Prix  

---

## 🛠️ Installation & Exécution  
> **Prérequis :** Golang installé sur votre machine  

### 1️⃣ Clonez le projet  
```sh
git clone https://github.com/tomy-21/Groupi-Tracker.git
cd Groupi-Tracker
2️⃣ Lancez l’application
sh
Copier
Modifier
go run main.go
💡 L’application tourne en local, ouvrez ensuite votre navigateur à l’adresse http://localhost:8080.

🎨 Pages disponibles
📍 Page d’accueil – Présentation du projet
📍 Page pilotes – Liste des pilotes avec recherche et filtres avancés
📍 Page circuits – Liste des circuits par année
📍 Page résultats – Résultats des courses
📍 Page favoris – Liste des favoris sauvegardés
📍 Page recherche – Recherche avancée
📍 Page à-propos – Détails sur le projet

🏆 Défis rencontrés
La première page des pilotes a été la plus complexe à réaliser. C’était la première implémentée et celle qui m’a demandé le plus de temps. Elle inclut 4 filtres, ce qui a rendu son développement plus difficile.
Gestion des erreurs API : L’application gère les erreurs en cas d’indisponibilité de l’API.
Stockage des favoris : J’ai utilisé des fichiers JSON pour enregistrer la liste des favoris, permettant ainsi une persistance des données.
📚 Ressources utilisées
🔹 Postman – Pour tester les endpoints et comprendre les réponses de l’API
🔹 ChatGPT – Pour obtenir des explications techniques et résoudre certains problèmes
🔹 YouTube & Forums – Pour trouver des tutoriels et bonnes pratiques en Golang

📝 Exemples de prompts utilisés sur ChatGPT :
1️⃣ "Comment parser un JSON en Golang ?"
2️⃣ "Comment mettre en place un système de filtres en Go ?"
3️⃣ "Comment gérer les erreurs d’une requête HTTP en Go ?"

📌 Auteur
👤 Nom : Tomy OUADHI
📂 GitHub : tomy-21

📅 Date limite du projet : 10/03/2025 avant 23h59

🔗 Lien du dépôt GitHub : Groupi-Tracker

💡 Améliorations futures possibles
✅ Ajouter un mode sombre pour améliorer l’UX
✅ Intégrer un graphique des performances des pilotes
✅ Héberger le projet en ligne pour un accès public

🔥 Merci d’avoir lu ce README ! 🚀