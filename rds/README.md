Here are the `curl` commands for each API endpoint:  

### 1️⃣ **Get All Pokémon**  
```sh
curl -X GET http://localhost:8080/pokemon
```

### 2️⃣ **Get a Specific Pokémon**  
```sh
curl -X GET http://localhost:8080/pokemon/{id}
```
Replace `{id}` with the actual Pokémon ID, e.g.,  
```sh
curl -X GET http://localhost:8080/pokemon/1
```

### 3️⃣ **Create a New Pokémon**  
```sh
curl -X POST http://localhost:8080/pokemon \
     -H "Content-Type: application/json" \
     -d '{
          "name": "Pikachu",
          "type": "Electric",
          "hp": 35,
          "attack": 55,
          "defense": 40,
          "sp_attack": 50,
          "sp_defense": 50,
          "speed": 90
        }'
```

### 4️⃣ **Update an Existing Pokémon**  
```sh
curl -X PUT http://localhost:8080/pokemon/{id} \
     -H "Content-Type: application/json" \
     -d '{
          "name": "Raichu",
          "type": "Electric",
          "hp": 60,
          "attack": 90,
          "defense": 55,
          "sp_attack": 90,
          "sp_defense": 80,
          "speed": 110
        }'
```
Replace `{id}` with the actual Pokémon ID, e.g.,  
```sh
curl -X PUT http://localhost:8080/pokemon/1 \
     -H "Content-Type: application/json" \
     -d '{
          "name": "Raichu",
          "type": "Electric",
          "hp": 60,
          "attack": 90,
          "defense": 55,
          "sp_attack": 90,
          "sp_defense": 80,
          "speed": 110
        }'
```

### 5️⃣ **Delete a Pokémon**  
```sh
curl -X DELETE http://localhost:8080/pokemon/{id}
```
Replace `{id}` with the actual Pokémon ID, e.g.,  
```sh
curl -X DELETE http://localhost:8080/pokemon/1
```

These commands assume the API is running on `localhost:8080`. 🚀
