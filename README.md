Welsh
==========

### Welcome to cheddar lovers !

# API Welsh
The requirements are in the go.mod file. In this test, sqlite was used as a database. A database is available with recipes and ingredients. Once compiled, the binary can run in a container with the help of a dockerfile.
The following command to linux allows to launch the Welsh API :
```
go build && ./welsh
```
The expected result :
```
Welsh
2023/04/19 16:11:09 Connected to DB
2023/04/19 16:11:09 Serving HTTP on port 8080
```
The http server listens on port 8080 and access to the routes is done by authentication except for the 2 routes index and createUser.

#  Routes
## Welcome
Request : curl -i http://localhost:8080/
Response :
```
HTTP/1.1 200 OK
Date: Wed, 19 Apr 2023 16:15:44 GMT
Content-Length: 28
Content-Type: text/plain; charset=utf-8

Welcome to cheddar lovers !!
```

### Admin – Create users
Request : curl -i -H "Content-Type: application/json" 
-X POST -d '{"username":"john","password":"test"}'
http://localhost:8080/api/admin/user

Response :
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 19 Apr 2023 19:55:23 GMT
Content-Length: 27

{"id":1,"username":"john"}
```

### Create ingredients
Request : curl -i -H "Content-Type: application/json"
-H "username: john" -H "password: test"
-X POST -d '{"name": "chou-fleur, cheddar"}'
http://localhost:8080/api/ingredient
Response :
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 19 Apr 2023 20:12:54 GMT
Content-Length: 57

[{"id":1,"name":"chou-fleur"},{"id":2,"name":"cheddar"}]
```

### List ingredients
Request : curl -i -H "Content-Type: application/json"
-H "username: john" -H "password: test"
http://localhost:8080/api/ingredient

Response :
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 19 Apr 2023 20:15:12 GMT
Content-Length: 57

[{"id":1,"name":"chou-fleur"},{"id":2,"name":"cheddar"}]
```

### Create recipe
Request : curl -i -H "Content-Type: application/json"
-H "username: john" -H "password: test"
-X POST -d '{"name": "welsh traditionnel","description": "Originaire du Pays de Galles, le welsh est désormais un plat typique proposé dans le Nord de la France. Cette recette traditionnelle est une variante du croque-monsieur, une tranche de pain garnie de jambon et cuite dans un plat à gratin, nappée d’une sauce au cheddar fondu , à la bière et à la moutarde. Le plus souvent, le welsh complet est surmonté d’un œuf au plat et proposé avec des frites. Facile à préparer et délicieusement gourmande, testez dès à présent cette recette hivernale et réconfortante par excellence ! ","instruction": "Étape 1 Préchauffez le four à 180°C (thermostat 6). Dans des plats à gratin individuels ou dans un grand plat, disposez les tranches de pain de campagne. Déposez une tranche de jambon blanc sur chaque morceau de pain. Étape 2 Râpez le cheddar et faites-le fondre dans une casserole à feu très doux, en remuant sans cesse à l’aide d’une cuillère en bois. Lorsque le fromage fondu nappe la cuillère, mouillez-le avec la bière blonde. Continuez de remuer jusqu’à obtenir une consistance homogène. Ajoutez alors la moutarde puis mélangez à nouveau. Nappez les tranches de pain couvertes de jambon de cette préparation au cheddar et à la bière. Étape 3 Enfournez le welsh complet pendant 10 minutes. Étape 4 Quelques minutes avant la cuisson, faites cuire les œufs au plat. À la sortie du four, déposez un œuf au plat sur le dessus de chaque welsh complet. Servez aussitôt avec un tour de moulin à poivre.","ingredient":[{"name": "cheddar"},{	"name": "oeufs"},{"name": "pain de campagne"},{"name": "jambon blanc"},{"name": "moutarde forte"},{"name": "bière blonde"}]}'
http://localhost:8080/api/recipe


Response :
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 19 Apr 2023 20:15:12 GMT
Content-Length: 1777

{
	"id": 1,
	"name": "welsh traditionnel",
	"description": "Originaire du Pays de Galles, le welsh est désormais un plat typique proposé dans le Nord de la France. Cette recette traditionnelle est une variante du croque-monsieur, une tranche de pain garnie de jambon et cuite dans un plat à gratin, nappée d’une sauce au cheddar fondu , à la bière et à la moutarde. Le plus souvent, le welsh complet est surmonté d’un œuf au plat et proposé avec des frites. Facile à préparer et délicieusement gourmande, testez dès à présent cette recette hivernale et réconfortante par excellence !",
	"instruction": "Étape 1 Préchauffez le four à 180°C (thermostat 6). Dans des plats à gratin individuels ou dans un grand plat, disposez les tranches de pain de campagne. Déposez une tranche de jambon blanc sur chaque morceau de pain. Étape 2 Râpez le cheddar et faites-le fondre dans une casserole à feu très doux, en remuant sans cesse à l’aide d’une cuillère en bois. Lorsque le fromage fondu nappe la cuillère, mouillez-le avec la bière blonde. Continuez de remuer jusqu’à obtenir une consistance homogène. Ajoutez alors la moutarde puis mélangez à nouveau. Nappez les tranches de pain couvertes de jambon de cette préparation au cheddar et à la bière. Étape 3 Enfournez le welsh complet pendant 10 minutes. Étape 4 Quelques minutes avant la cuisson, faites cuire les œufs au plat. À la sortie du four, déposez un œuf au plat sur le dessus de chaque welsh complet. Servez aussitôt avec un tour de moulin à poivre.",
	"Ingredient": [
		{
			"id": 2,
			"name": "cheddar"
		},
		{
			"id": 3,
			"name": "oeufs"
		},
		{
			"id": 4,
			"name": "pain de campagne"
		},
		{
			"id": 5,
			"name": "jambon blanc"
		},
		{
			"id": 6,
			"name": "moutarde forte"
		},
		{
			"id": 7,
			"name": "bière blonde"
		}
	]
}
```

### List Recipe
Request : curl -i -H "Content-Type: application/json"
-H "username: john" -H "password: test"
-d '{"ingredient": true}'
http://localhost:8080/api/recipe

Response :
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 19 Apr 2023 20:49:36 GMT
Content-Length: 1552

[
	{
		"id": 1,
		"name": "welsh traditionnel",
		"description": "Originaire du Pays de Galles, le welsh est désormais un plat typique proposé dans le Nord de la France. Cette recette traditionnelle est une variante du croque-monsieur, une tranche de pain garnie de jambon et cuite dans un plat à gratin, nappée d’une sauce au cheddar fondu , à la bière et à la moutarde. Le plus souvent, le welsh complet est surmonté d’un œuf au plat et proposé avec des frites. Facile à préparer et délicieusement gourmande, testez dès à présent cette recette hivernale et réconfortante par excellence !",
		"instruction": "Étape 1 Préchauffez le four à 180°C (thermostat 6). Dans des plats à gratin individuels ou dans un grand plat, disposez les tranches de pain de campagne. Déposez une tranche de jambon blanc sur chaque morceau de pain. Étape 2 Râpez le cheddar et faites-le fondre dans une casserole à feu très doux, en remuant sans cesse à l’aide d’une cuillère en bois. Lorsque le fromage fondu nappe la cuillère, mouillez-le avec la bière blonde. Continuez de remuer jusqu’à obtenir une consistance homogène. Ajoutez alors la moutarde puis mélangez à nouveau. Nappez les tranches de pain couvertes de jambon de cette préparation au cheddar et à la bière. Étape 3 Enfournez le welsh complet pendant 10 minutes. Étape 4 Quelques minutes avant la cuisson, faites cuire les œufs au plat. À la sortie du four, déposez un œuf au plat sur le dessus de chaque welsh complet. Servez aussitôt avec un tour de moulin à poivre.",
		"Ingredient": [
			{
				"id": 2,
				"name": "cheddar"
			},
			{
				"id": 3,
				"name": "oeufs"
			},
			{
				"id": 4,
				"name": "pain de campagne"
			},
			{
				"id": 5,
				"name": "jambon blanc"
			},
			{
				"id": 6,
				"name": "moutarde forte"
			},
			{
				"id": 7,
				"name": "bière blonde"
			}
		]
	}
]
```

### List recipe details
Request : curl -i -H "Content-Type: application/json"
-H "username: john" -H "password: test"
http://localhost:8080/api/recipe/1

Response :
```

Response :
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 19 Apr 2023 20:57:49 GMT
Content-Length: 1746

[
	{
		"id": 1,
		"name": "welsh traditionnel",
		"description": "….
```

### 2.8	Flag / Unflag favorirte recipes
=> Flag

Request : curl -i -H "Content-Type: application/json"
-H "username: john" -H "password: test"
-X POST -d '{"name": "welsh traditionnel","flag": true}'
http://localhost:8080/api/favorite

Response :
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 19 Apr 2023 21:09:32 GMT
Content-Length: 64

{"username":"john","recipe":["welsh traditionnel"],"flag":true}
```

=> Unflag

Request : curl -i -H "Content-Type: application/json"
-H "username: john" -H "password: test" \
-X POST -d '{"name": "welsh traditionnel","flag": false}'
http://localhost:8080/api/favorite

Response :
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 19 Apr 2023 21:07:50 GMT
Content-Length: 65

{"username":"john","recipe":["welsh traditionnel"],"flag":false}
```

### List favorite recipe
Request : curl -i -H "Content-Type: application/json"
-H "username: john" -H "password: test"
http://localhost:8080/api/favorite


Response :
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 19 Apr 2023 21:10:56 GMT
Content-Length: 1552

[{"id":1,"name":"welsh traditionnel","description" : "…
```

## Description of tests
### Test unit
The test only tests the handleIngredientCreate function :
```
go test -timeout 30s -run ^TestIngredientCreateUnit$ academy.go/welsh
```

### Test integration
The test allows for end-to-end testing of the Create Ingredient feature:
```
go test -timeout 30s -run ^TestIngredientCreateIntregation$ academy.go/welsh
```